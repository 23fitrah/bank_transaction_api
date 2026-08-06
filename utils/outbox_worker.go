package utils

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"transaction_api/config"

	"github.com/bytedance/sonic"
)

type OutboxWorker struct {
	db       *sql.DB
	rabbitmq *config.RabbitMQ
	cfg      *config.Config
}

func NewOutboxWorker(db *sql.DB, rabbit *config.RabbitMQ, cfg *config.Config) *OutboxWorker {
	return &OutboxWorker{
		db:       db,
		rabbitmq: rabbit,
		cfg:      cfg,
	}
}

type OutboxRow struct {
	ID         int64
	Exchange   string
	RoutingKey string
	PAYLOAD    string
}

func (w *OutboxWorker) Start(ctx context.Context) {
	go func() {
		log.Println("Outbox worker start...")

		for {
			/*select {
			case <-ctx.Done():
				log.Println("Worker stopped:", ctx.Err())
				return // keluar dari loop & goroutine
			default:
				w.processOne(ctx)
				time.Sleep(2 * time.Second)
			}*/
			w.processOne(ctx)

		}
	}()
}

func (w *OutboxWorker) processOne(ctx context.Context) {
	var row OutboxRow
	log.Println("Worker standby")

	err := w.db.QueryRowContext(ctx, `
		WITH cte AS (
			SELECT TOP 1 *
			FROM MESSAGE_OUTBOX WITH (ROWLOCK, READPAST)
			WHERE STATUS='PENDING'
			AND RETRY_COUNT < 10
			ORDER BY ID
		)
		UPDATE cte
		SET STATUS='PROCESSING'
		OUTPUT inserted.ID, inserted.EXCHANGE, inserted.ROUTING_KEY, inserted.PAYLOAD	`).Scan(&row.ID, &row.Exchange, &row.RoutingKey, &row.PAYLOAD)

	if err != nil || row.ID == 0 {
		log.Println("Worker standby : %s", err)

		return
	}
	log.Printf("payload : %s", string(row.PAYLOAD))
	var payload interface{}
	if err := sonic.Unmarshal([]byte(row.PAYLOAD), &payload); err != nil {
		w.db.ExecContext(ctx, `
			UPDATE MESSAGE_OUTBOX
			SET STATUS='FAILED', LAST_ERROR=@p1
			WHERE ID=@p2
		`, sql.Named("p1", err.Error()), sql.Named("p2", row.ID))
		log.Printf("%s", err.Error())
		return
	}

	err2 := Publish(w.rabbitmq, w.cfg.RabbitMQ.Exchange, w.cfg.RabbitMQ.Routing, w.cfg.RabbitMQ.Queue, payload)
	if err2 != nil {
		log.Printf("error rabbit : %s", err.Error())
		_, err = w.db.ExecContext(ctx, `
			UPDATE MESSAGE_OUTBOX
			SET STATUS='PENDING',
				RETRY_COUNT = RETRY_COUNT + 1,
				LAST_ERROR = @p1
			WHERE ID=@p2
		`, sql.Named("p1", err2.Error()), sql.Named("p2", row.ID))
		log.Println("Worker publish error : %w", err2.Error())
		return
	}

	_, err = w.db.ExecContext(ctx, `
		UPDATE MESSAGE_OUTBOX
		SET STATUS='SENT',
			SENT_AT=GETDATE()
		WHERE ID=@p1
	`, sql.Named("p1", row.ID))
	if err != nil {
		log.Printf("error update sent : %s", err)
	}

	fmt.Println("outbox sent:", row.ID)
}
