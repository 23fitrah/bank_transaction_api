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
