package service

import (
	"context"
	"fmt"
	"transaction_api/model/log_monitoring"
	"transaction_api/repository"
)

type LogMonitoringService struct {
	Repo *repository.LogMonitoringRepo
}

func NewLogMonitoringService(repo *repository.LogMonitoringRepo) *LogMonitoringService {
	return &LogMonitoringService{
		Repo: repo,
	}
}

func (s *LogMonitoringService) InsertLogMonitoringService(ctx context.Context, log log_monitoring.Log_monitoring) (int64, error) {
	insertID, err := s.Repo.InsertLogMonitoring(ctx, log)
	if err != nil {
		return 0, fmt.Errorf("insert log monitoring: %w", err)
	}
	return insertID, err
}

func (s *LogMonitoringService) GetAllLogService(ctx context.Context, limit, offset int, menu, dateFrom, dateTo, method string) ([]*log_monitoring.Log_monitoring, int64, error) {
	var list []*log_monitoring.Log_monitoring
	var total int64
	var err error

	if method == "elastic" {
		list, total, err = s.Repo.FetchLogFromElasticSearch(ctx, limit, offset, menu, dateFrom, dateTo)

	} else {
		list, total, err = s.Repo.GetAllLogRepository(ctx, limit, offset, menu, dateFrom, dateTo)

	}
	if err != nil {

		return nil, 0, fmt.Errorf("service get all logs: %w", err)
	}
	return list, total, nil
}
