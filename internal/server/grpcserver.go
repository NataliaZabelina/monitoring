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
	"github.com/NataliaZabelina/monitoring/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var grpcServer *grpc.Server

type GrpcServer struct {
	monitoring *monitoring.Monitoring
	logger     *zap.Logger
	api.UnimplementedMonitoringServer
}

func Start(monitoring *monitoring.Monitoring, log *zap.Logger) error {

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
		}

		address := "localhost:50051"
		listener, err := net.Listen("tcp", address)
		if err != nil {
			log.Fatal("Can not start listening")
			return
		}

		grpcServer = grpc.NewServer()
		api.RegisterMonitoringServer(grpcServer, srv)

		if err := grpcServer.Serve(listener); err != nil {
			logger.Logger.Fatalf("Can not accept incoming connection: %v", err)
			return
		}
	}()

	select {
	case <-interrupt:
		break
	case <-ctx.Done():
		break
	}

	logger.Logger.Info("Shutdown signal accepted")
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
			data.SystemValue = &api.SystemResponse{
				ResponseTime:     timestamppb.Now(),
				LoadAverageValue: 67.8}
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
