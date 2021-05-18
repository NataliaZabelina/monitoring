package grpcserver

import (
	"context"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/NataliaZabelina/monitoring/api"
	monitoring "github.com/NataliaZabelina/monitoring/internal/app"
	"github.com/NataliaZabelina/monitoring/internal/config"
	"github.com/NataliaZabelina/monitoring/internal/logger"
	"github.com/NataliaZabelina/monitoring/internal/storage"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var grpcServer *grpc.Server

type GrpcServer struct {
	api.UnimplementedMonitoringServer
	logger     *zap.SugaredLogger
	monitoring *monitoring.Monitoring
	db         *storage.DB
	cfg        *config.Config
}

func Start(db *storage.DB, monitoring *monitoring.Monitoring, log *zap.SugaredLogger, cfg *config.Config) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := monitoring.Run(ctx)
		if err != nil {
			log.Fatal("Can't start monitoring")
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		srv := &GrpcServer{
			monitoring: monitoring,
			logger:     log,
			db:         db,
			cfg:        cfg,
		}

		address := cfg.Host + ":" + cfg.Port
		listener, err := net.Listen("tcp", address)
		if err != nil {
			log.Fatal("Can not start listening")
			return
		}

		grpcServer = grpc.NewServer()
		api.RegisterMonitoringServer(grpcServer, srv)

		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Can not accept incoming connection: %w", err)
			return
		}
	}()

	select {
	case <-interrupt:
		break
	case <-ctx.Done():
		break
	}

	log.Info("Shutdown signal accepted")
	cancel()
	if grpcServer != nil {
		grpcServer.GracefulStop()
	}

	wg.Wait()
	return nil
}

func (s *GrpcServer) GetInfo(req *api.Request, stream api.Monitoring_GetInfoServer) error {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)
	ctx := stream.Context()

	logger.Logger.Info("Client connected with params: ", zap.Int32("timeout", req.Time), zap.Int32("period", req.Period))

	for {
		d := time.Duration(int64(time.Second) * int64(req.Time))
		select {
		case <-time.After(d):
			data := api.Result{}
			timestamp := timestamppb.Now()
			data.SystemValue = &api.SystemResponse{
				ResponseTime:     timestamp,
				LoadAverageValue: float64((s.db.SystemTable.GetAverage(req.Period)).LoadAvg),
			}
			cpu := s.db.CPUTable.GetAverage(req.Period)
			data.CpuValue = &api.CPUResponse{
				ResponseTime: timestamp,
				UserMode:     float64(cpu.UserMode),
				SystemMode:   float64(cpu.SystemMode),
				Idle:         float64(cpu.Idle),
			}
			err := stream.Send(&data)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			logger.Logger.Errorf("Stop streaming to client: $v", ctx.Err())
			return ctx.Err()

		case <-stopChan:
			logger.Logger.Info("Close streaming")
			return nil
		}
	}
}
