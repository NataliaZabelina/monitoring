package monitoring

import (
	"context"
	"time"

	"github.com/NataliaZabelina/monitoring/internal/app/command"
	"github.com/NataliaZabelina/monitoring/internal/config"
	"github.com/NataliaZabelina/monitoring/internal/storage"
	"go.uber.org/zap"
)

type Monitoring struct {
	db     *storage.DB
	logger *zap.SugaredLogger
	cfg    *config.Config
}

func New(db *storage.DB, logger *zap.SugaredLogger, cfg *config.Config) *Monitoring {
	return &Monitoring{db: db, logger: logger, cfg: cfg}
}

func (m *Monitoring) Run(ctx context.Context) error {
	timoutCollection := m.cfg.Collector.Timeout
	m.logger.Info("Start data collection.")

	if m.cfg.Collector.Statistics.LoadSystem {
		go func() {
			err := m.newWorker(ctx, timoutCollection, "LoadAvg", command.GetLoadAvg)
			if err != nil {
				m.logger.Error("Can't start loadAvgWorker.", zap.Error(err))
				return
			}
		}()
	}

	if m.cfg.Collector.Statistics.LoadCPU {
		go func() {
			err := m.newWorker(ctx, timoutCollection, "CpuLoad", command.GetCPU)
			if err != nil {
				m.logger.Error("Cannot start workerLoadCPU", zap.Error(err))
				return
			}
		}()
	}

	if m.cfg.Collector.Statistics.LoadDisk {
		go func() {
			err := m.newWorker(ctx, timoutCollection, "DislIO", command.GetDiskIO)
			if err != nil {
				m.logger.Error("Cannot start workerDiskIO", zap.Error(err))
				return
			}
		}()
	}

	return nil
}

func (m *Monitoring) newWorker(ctx context.Context, timout int, name string, getStat func(*storage.DB) (string, error)) error {
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
