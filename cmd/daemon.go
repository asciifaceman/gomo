/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/asciifaceman/gomo/pkg/gotmo"
	"github.com/spf13/cobra"
)

// daemonCmd represents the daemon command
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Daemonized Gomo which will continuously run",
	Long: `Daemonized Gomo which will continuously run and insert
discovered metrics into prometheus time series for graphing and
historical analysis.`,
	Run: func(cmd *cobra.Command, args []string) {
		g, err := gotmo.NewGotmo(hostname, reqtimeout, pingtargets, pingWorkerCount)
		if err != nil {
			fmt.Printf("Failed to setup: %v", err)
		}
		err = g.Daemon()
		if err != nil {
			g.Logger.Errorw("Exited with error", "error", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// daemonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// daemonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
