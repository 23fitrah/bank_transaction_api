package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"transaction_api/model/log_monitoring"
)

type LogMonitoringRepo struct {
	db *sql.DB
}

func NewLogMonitoringRepository(db *sql.DB) *LogMonitoringRepo {
	return &LogMonitoringRepo{db: db}
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
