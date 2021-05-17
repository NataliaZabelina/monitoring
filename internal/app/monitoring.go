package monitoring

import (
	"context"
	"time"

	"github.com/NataliaZabelina/monitoring/internal/app/command"
	"github.com/NataliaZabelina/monitoring/internal/storage"
	"go.uber.org/zap"
)

type Monitoring struct {
	db     *storage.Db
	logger *zap.SugaredLogger
}

func New(db *storage.Db, logger *zap.SugaredLogger) *Monitoring {
	return &Monitoring{db: db, logger: logger}
}

func (m *Monitoring) Run(ctx context.Context) error {
	timoutCollection := 1
	m.logger.Info("Start data collection.")

	go func() {
		err := m.newWorker(ctx, timoutCollection, "LoadAvg", command.GetLoadAvg)
		if err != nil {
			m.logger.Error("Can't start loadAvgWorker.", zap.Error(err))
			return
		}
	}()

	go func() {
		err := m.newWorker(ctx, timoutCollection, "CpuLoad", command.GetCpu)
		if err != nil {
			m.logger.Error("Cannot start workerLoadCPU", zap.Error(err))
			return
		}
	}()

	return nil
}

func (m *Monitoring) newWorker(ctx context.Context, timout int, name string, getStat func(*storage.Db) (string, error)) error {
	m.logger.Infof("Starting collection %s.", name)
	for {
		d := time.Duration(int64(time.Second) * int64(timout))

		select {
		case <-time.After(d):

			res, err := getStat(m.db)
			if err != nil {
				m.logger.Error(name, zap.Error(err))
				return err
			}

			m.logger.Debug(name, zap.String(name, res))

		case <-ctx.Done():
			m.logger.Infof("Data collection for %s completed. Context closed.", name)
			return nil
		}
	}
}
