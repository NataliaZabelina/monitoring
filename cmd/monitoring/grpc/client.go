package grpc

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/NataliaZabelina/monitoring/api"
	l "github.com/NataliaZabelina/monitoring/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var (
	address string
	timeout int32
	period  int32

	loadAvg   bool
	cpu       bool
	disk      bool
	toptalker bool
	netstat   bool
	all       bool
)

var (
	GrpcClientCmd = &cobra.Command{
		Use:     "grpc_client",
		Short:   "Start grpc client",
		Run:     startGrpcClient,
		Example: "monitoring grpc_client --address=':50051'",
	}
)

func init() {
	GrpcClientCmd.Flags().StringVarP(&address, "address", "", "localhost:50051", "host:port to establish connection")
	GrpcClientCmd.Flags().Int32VarP(&timeout, "timeout", "", 5, "time parameter in sec to define how often to collect statistics")
	GrpcClientCmd.Flags().Int32VarP(&period, "period", "", 15, "time period in sec for collecting statistics")
	GrpcClientCmd.Flags().BoolVarP(&loadAvg, "load", "s", false, "collect system load averge statistics")
	GrpcClientCmd.Flags().BoolVarP(&cpu, "cpu", "c", false, "collect CPU statistics")
	GrpcClientCmd.Flags().BoolVarP(&disk, "disk", "d", false, "collect Disk statistics")
	GrpcClientCmd.Flags().BoolVarP(&toptalker, "toptalk", "t", false, "collect top talkers statistics")
	GrpcClientCmd.Flags().BoolVarP(&netstat, "netstat", "n", false, "collect network statistics")
	GrpcClientCmd.Flags().BoolVarP(&all, "all", "a", true, "collect all statistics")
	err := viper.BindPFlags(GrpcClientCmd.Flags())
	if err != nil {
		l.Logger.Fatal(err)
	}
	viper.AutomaticEnv()
	address = viper.GetString("address")
	timeout = viper.GetInt32("timeout")
	period = viper.GetInt32("period")
}

func startGrpcClient(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		l.Logger.Fatal(err)
	}
	defer conn.Close()

	l.Logger.Infof("Connection established")
	client := api.NewMonitoringClient(conn)

	if !loadAvg && !cpu && !disk && !toptalker && !netstat {
		loadAvg = true
		cpu = true
		disk = true
		loadAvg = true
		netstat = true
	}

	collectStatistics(client, timeout, period)
}

func collectStatistics(client api.MonitoringClient, timeout int32, period int32) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := &api.Request{Time: timeout, Period: period}
	stream, err := client.GetInfo(ctx, req)
	if err != nil {
		l.Logger.Fatalf("Can't get statistics for client: ", err)
	}

	for {
		msg, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			l.Logger.Info("IO EOF: %w", err)
			return
		}
		if err != nil {
			l.Logger.Fatalf("Can't read from stream: %w", err)
		}

		l.Logger.Info("Collected System Information")
		l.Logger.Info(msg)
	}
}
