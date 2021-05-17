package grpc

import (
	monitoring "github.com/NataliaZabelina/monitoring/internal/app"
	l "github.com/NataliaZabelina/monitoring/internal/logger"
	grpcserver "github.com/NataliaZabelina/monitoring/internal/server"
	"github.com/NataliaZabelina/monitoring/internal/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ServerCmd = &cobra.Command{
		Use:     "grpc_server",
		Short:   "Start grpc server",
		Run:     startGrpcServer,
		Example: "monitoring grpc_server --config=configs/config.json",
	}

	port string
)

func init() {
	ServerCmd.Flags().StringVarP(&port, "port", "p", "50051", "server port")
	err := viper.BindPFlag("port", ServerCmd.Flags().Lookup("port"))
	if err != nil {
		l.Logger.Fatalf("Something is going wrong: %w", err)
	}
}

func startGrpcServer(cmd *cobra.Command, args []string) {
	logger, err := l.Init()
	if err != nil {
		l.Logger.Fatalf("Can't init logger: %w", err)
	}
	db := &storage.Db{}
	db.Init()
	monitoring := monitoring.New(db, logger)

	if err := grpcserver.Start(monitoring, logger); err != nil {
		l.Logger.Fatalf("Something is going wrong: %w", err)
	}
}
