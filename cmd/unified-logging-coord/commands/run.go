/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package commands

import (
	"github.com/nalej/unified-logging/internal/app/coord"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var config = coord.Config{}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Launch the server API",
	Long:  `Launch the server API`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		Run()
	},
}

func init() {
	runCmd.Flags().IntVar(&config.Port, "port", 8323, "Port for Unified Logging Coordinator gRPC API")
	runCmd.PersistentFlags().StringVar(&config.SystemModelAddress, "systemModelAddress", "localhost:8800", "System Model address (host:port)")
	runCmd.PersistentFlags().StringVar(&config.AppClusterPrefix, "appClusterPrefix", "appcluster", "Prefix for application cluster hostnames")
	runCmd.PersistentFlags().IntVar(&config.AppClusterPort, "appClusterPort", 443, "Port used by app-cluster-api")
	runCmd.PersistentFlags().BoolVar(&config.UseTLS, "useTLS", true, "Use TLS to connect to application cluster")
	runCmd.PersistentFlags().BoolVar(&config.Insecure, "insecure", false, "Don't validate TLS certificates")
	runCmd.PersistentFlags().StringVar(&config.CACert, "caCert", "", "Alternative certificate file to use for validation")
	rootCmd.AddCommand(runCmd)
}

func Run() {
	log.Info().Msg("Launching Unified Logging Coordinator service")

	server, err := coord.NewService(&config)
	if err != nil {
		log.Fatal().Str("err", err.DebugReport()).Err(err)
		panic(err.Error())
	}

	err = server.Run()
	if err != nil {
		log.Fatal().Str("err", err.DebugReport()).Err(err)
		panic(err.Error())
	}
}