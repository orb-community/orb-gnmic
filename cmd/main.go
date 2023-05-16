package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/orb-community/orb-gnmic/config"
	"github.com/orb-community/orb-gnmic/gnmic"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Debug         bool
	SelfTelemetry bool
	ServerHost    string
	ServerPort    uint64
)

func Run(cmd *cobra.Command, args []string) {

	initConfig()

	// configuration
	var config config.Config

	err := viper.Unmarshal(&config)
	if err != nil {
		cobra.CheckErr(fmt.Errorf("orb-gnmic start up error (config): %w", err))
		os.Exit(1)
	}

	// logger
	var logger *zap.Logger
	atomicLevel := zap.NewAtomicLevel()
	if Debug {
		atomicLevel.SetLevel(zap.DebugLevel)
	} else {
		atomicLevel.SetLevel(zap.InfoLevel)
	}
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		os.Stdout,
		atomicLevel,
	)
	logger = zap.New(core, zap.AddCaller())
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	// new otlpinf
	a, err := gnmic.New(logger, &config)
	if err != nil {
		logger.Error("orb-gnmic start up error", zap.Error(err))
		os.Exit(1)
	}

	// handle signals
	done := make(chan bool, 1)
	rootCtx, cancelFunc := context.WithCancel(context.WithValue(context.Background(), "routine", "mainRoutine"))

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		for {
			select {
			case <-sigs:
				logger.Warn("stop signal received, stopping orb-gnmic")
				a.Stop(rootCtx)
				cancelFunc()
			case <-rootCtx.Done():
				logger.Warn("mainRoutine context cancelled")
				done <- true
				return
			}
		}
	}()

	// start otlpinf
	err = a.Start(rootCtx, cancelFunc)
	if err != nil {
		logger.Error("orb-gnmic startup error", zap.Error(err))
		os.Exit(1)
	}

	<-done
}

func initConfig() {
	v := viper.New()
	v.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)
	// note: viper seems to require a default (or a BindEnv) to be overridden by environment variables
	v.SetDefault("orb_gnmic_debug", Debug)
	v.SetDefault("orb_gnmic_server_host", ServerHost)
	v.SetDefault("orb_gnmic_server_port", ServerPort)
	cobra.CheckErr(viper.MergeConfigMap(v.AllSettings()))
}

func main() {

	rootCmd := &cobra.Command{
		Use: "orb-gnmic",
	}

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run orb-gnmic",
		Long:  `Run orb-gnmic`,
		Run:   Run,
	}

	runCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "Enable verbose (debug level) output")
	runCmd.PersistentFlags().StringVarP(&ServerHost, "server_host", "a", "localhost", "Define REST Host")
	runCmd.PersistentFlags().Uint64VarP(&ServerPort, "server_port", "p", 10337, "Define REST Port")

	rootCmd.AddCommand(runCmd)
	rootCmd.Execute()
}
