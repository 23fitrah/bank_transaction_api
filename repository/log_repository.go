package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
	"transaction_api/constants"
	"transaction_api/model/log_monitoring"

	"github.com/olivere/elastic/v7"
	"github.com/redis/go-redis/v9"
)

type LogMonitoringRepo struct {
	db  *sql.DB
	rdb *redis.Client
	es  *elastic.Client
}

func NewLogMonitoringRepository(db *sql.DB, rdb *redis.Client, es *elastic.Client) *LogMonitoringRepo {
	return &LogMonitoringRepo{db: db, rdb: rdb, es: es}
}

func (r *LogMonitoringRepo) InsertLogMonitoring(ctx context.Context, data log_monitoring.Log_monitoring) (int64, error) {

	const query = `
		INSERT INTO [h_audit_trail] (
			MENU,
			ACTION,
			OLD_VALUE,
			NEW_VALUE,
			RESPONSE_MSG,
			CHANGED_BY,
			CHANGED_DATE,
			IP_CLIENT,
			USER_AGENT
		)
		OUTPUT INSERTED.ID_AUDIT
		VALUES (
			@p1,            
			@p2,            
			@p3,            
			@p4,            
			@p5,
			@p6,
			@p7,
			@p8,
			@p9                
		);
	`

	var insertedID int64

	err := r.db.QueryRowContext(ctx, query,
		sql.Named("p1", data.MENU),
		sql.Named("p2", data.ACTION),
		sql.Named("p3", data.OLD_VALUE),
		sql.Named("p4", data.NEW_VALUE),
		sql.Named("p5", data.RESPONSE_MSG),
		sql.Named("p6", data.CHANGED_BY),
		sql.Named("p7", data.CHANGED_DATE),
		sql.Named("p8", data.IP_CLIENT),
		sql.Named("p9", data.USER_AGENT),
	).Scan(&insertedID)

	if err != nil {

		return 0, fmt.Errorf("error inserting into h_audit_trail: %w", err)
	}
	log.Printf("Successfully inserted log monitoring with ID: %v", insertedID)
	return insertedID, nil

}

func (r *LogMonitoringRepo) GetAllLogRepository(ctx context.Context, limit, offset int, menu, dateFrom, dateTo string) ([]*log_monitoring.Log_monitoring, int64, error) {
	key := fmt.Sprintf(
		constants.LOGMONITORING+":menu=%s:limit=%d:offset=%d:from=%s:to=%s",
		menu,
		limit,
		offset,
		dateFrom,
		dateTo,
	)
	return GetOrSet(ctx, r.rdb, key, 10*time.Minute, func() ([]*log_monitoring.Log_monitoring, int64, error) {
		return r.FetchLogFromDB(ctx, limit, offset, menu, dateFrom, dateTo)
	})
}

func (r *LogMonitoringRepo) FetchLogFromDB(ctx context.Context, limit, offset int, menu, dateFrom, dateTo string) ([]*log_monitoring.Log_monitoring, int64, error) {

	query := `
		SELECT 
		    id_audit,
			menu, 
			action, 
			old_value, 
			new_value, 
			response_msg, 
			changed_by, 
			changed_date, 
			ip_client, 
			user_agent
		FROM 
			h_audit_trail WITH (NOLOCK)
		WHERE 
			1=1 `

	args := []interface{}{}
	paramIdx := 1

	if menu != "" {
		query = query + " AND menu = @p" + strconv.Itoa(paramIdx)
		args = append(args, menu)
		paramIdx++
	}

	if dateFrom != "" && dateTo != "" {
		query = query + " AND changed_date BETWEEN @p" + strconv.Itoa(paramIdx) + " AND @p" + strconv.Itoa(paramIdx+1)
		args = append(args, dateFrom, dateTo)
		paramIdx += 2
	}

	var total int64
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM (`+query+`) AS tb`, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count logs: %w", err)
	}

	query += " ORDER BY id_audit DESC"

	if limit != -1 {
		query += fmt.Sprintf("  OFFSET @p%d ROWS FETCH NEXT @p%d ROWS ONLY", paramIdx, paramIdx+1)
	}
	args = append(args, offset)
	args = append(args, limit)

	sqlRows, err := r.db.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, 0, fmt.Errorf("query get all log failed: %w", err)
	}

	defer sqlRows.Close()

	results := []*log_monitoring.Log_monitoring{}
	var tt time.Time

	for sqlRows.Next() {
		var trx log_monitoring.Log_monitoring
		err := sqlRows.Scan(
			&trx.ID_AUDIT,
			&trx.MENU,
			&trx.ACTION,
			&trx.OLD_VALUE,
			&trx.NEW_VALUE,
			&trx.RESPONSE_MSG,
			&trx.CHANGED_BY,
			&tt,
			&trx.IP_CLIENT,
			&trx.USER_AGENT)
		if err != nil {
			return nil, 0, fmt.Errorf("scan transaction row failed: %w", err)
		}
		trx.CHANGED_DATE = tt.Format("2006-01-02 15:04:05")
		results = append(results, &trx)
	}

	return results, total, nil
}

func (r *LogMonitoringRepo) FetchLogFromElasticSearch(ctx context.Context, limit, offset int, menu, dateFrom, dateTo string) ([]*log_monitoring.Log_monitoring, int64, error) {
	//query := elastic.NewMatchQuery("message", menu)

	if menu == "" {
		menu = "*"
	}
	result, err := r.es.Search().
		Index("transaction_api-"+menu).
		//Query(query).
		From(offset).
		Size(limit).
		Sort("@timestamp", false). // false = descending (terbaru duluan)
		Do(ctx)                    // <-- WAJIB, ini yang benar-benar eksekusi request
	if err != nil {
		return nil, 0, err
	}

	logs := make([]*log_monitoring.Log_monitoring, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		log.Printf("Source: %s\n", string(hit.Source))

		var l *log_monitoring.Log_monitoring
		if err := json.Unmarshal(hit.Source, &l); err != nil {
			continue
		}
		// hit.Id is a string; convert to int64 to match ID_AUDIT type
		if id, err := strconv.ParseInt(hit.Id, 10, 64); err == nil {
			l.ID_AUDIT = id
		} else {
			// skip entries with non-numeric IDs
			continue
		}

		logs = append(logs, l)
	}

	return logs, result.TotalHits(), nil
}
