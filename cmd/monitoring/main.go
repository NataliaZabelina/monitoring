package main

import (
	"github.com/NataliaZabelina/monitoring/cmd/monitoring/grpc"
	l "github.com/NataliaZabelina/monitoring/internal/logger"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "monitoring",
	Short: "System monitoring",
}

func init() {
	rootCmd.AddCommand(grpc.ServerCmd, grpc.GrpcClientCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		l.Logger.Fatal(err)
	}
}
