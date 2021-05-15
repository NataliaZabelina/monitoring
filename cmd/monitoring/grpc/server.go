package grpc

import (
	monitoring "github.com/NataliaZabelina/monitoring/internal/app"
	l "github.com/NataliaZabelina/monitoring/internal/logger"
	grpcserver "github.com/NataliaZabelina/monitoring/internal/server"
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
		l.Logger.Fatalf("Something is going wrong: %v", err)
	}
}

func startGrpcServer(cmd *cobra.Command, args []string) {
	monitoring := monitoring.New("Some test data")

	logger, err := l.Init()
	if err != nil {
		l.Logger.Fatalf("Can't init logger: %v", err)
	}

	if err := grpcserver.Start(monitoring, logger); err != nil {
		l.Logger.Fatalf("Something is going wrong: %v", err)
	}
}
