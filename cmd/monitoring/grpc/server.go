package grpc

import (
	monitoring "github.com/NataliaZabelina/monitoring/internal/app"
	"github.com/NataliaZabelina/monitoring/internal/config"
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

	cfgFilePath string
)

func init() {
	ServerCmd.Flags().StringVarP(&cfgFilePath, "config", "c", "", "path to configuration file")
	err := viper.BindPFlag("config", ServerCmd.Flags().Lookup("config"))
	if err != nil {
		l.Logger.Fatalf("Something is going wrong: %w", err)
	}
}

func startGrpcServer(cmd *cobra.Command, args []string) {
	cfg, err := config.LoadConfig(cfgFilePath)
	if err != nil {
		l.Logger.Fatalf("Can't load config: %w", err)
	}

	logger, err := l.Init(cfg.Log)
	if err != nil {
		l.Logger.Fatalf("Can't init logger: %w", err)
	}

	db := &storage.DB{}
	db.Init()
	monitoring := monitoring.New(db, logger, cfg)

	if err := grpcserver.Start(db, monitoring, logger, cfg); err != nil {
		l.Logger.Fatalf("Something is going wrong: %w", err)
	}
}
