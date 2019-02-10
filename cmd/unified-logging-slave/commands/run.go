/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package commands

import (
	"github.com/nalej/unified-logging/internal/app/slave"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var config = slave.Config{}

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
	runCmd.Flags().IntVar(&config.Port, "port", 8322, "Port for Unified Logging Slave gRPC API")
	runCmd.PersistentFlags().StringVar(&config.ElasticAddress, "elasticAddress", "localhost:9200",
		"ElasticSearch address (host:port)")
	rootCmd.AddCommand(runCmd)
}

func Run() {
	log.Info().Msg("Launching Unified Logging Slave service")

	server, err := slave.NewService(&config)
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
