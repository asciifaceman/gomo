/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/asciifaceman/gomo/pkg/daemon"
	"github.com/spf13/cobra"
)

var (
	serverPort = 2112
)

// daemonCmd represents the daemon command
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Daemonized Gomo which will continuously run",
	Long: `Daemonized Gomo which will continuously run and insert
discovered metrics into prometheus time series for graphing and
historical analysis.`,
	Run: func(cmd *cobra.Command, args []string) {
		d, err := daemon.New(hostname, serverPort, reqtimeout)
		if err != nil {
			fmt.Printf("Failed to setup daemon: %v\n", err)
			return
		}

		err = d.Run()
		if err != nil {
			d.Logger.Errorw("Runtime error", "error", err)
			return
		}

		d.Logger.Info("Done. exiting.")

	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)

	daemonCmd.PersistentFlags().IntVarP(&serverPort, "port", "m", serverPort, "Port to bind metrics webserver to")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// daemonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// daemonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
