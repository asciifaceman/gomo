/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/asciifaceman/gomo/pkg/gotmo"
	"github.com/spf13/cobra"
)

var pretty bool

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Do a single fetch and display",
	Long:  `Do a single fetch and display.`,
	Run: func(cmd *cobra.Command, args []string) {
		g, err := gotmo.NewGotmo(hostname, reqtimeout, pingtargets, pingWorkerCount)
		if err != nil {
			fmt.Printf("Failed to setup: %v", err)
		}

		g.Printer.PrintHeader(fmt.Sprintf("GOMO %s", version))

		g.CLIEntry(pretty)
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
	showCmd.PersistentFlags().BoolVar(&pretty, "pretty", false, "Print a prettified table layout instead of raw data")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
