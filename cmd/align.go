/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/asciifaceman/gomo/pkg/alignui"
	"github.com/spf13/cobra"
)

var (
	pollFrequency = 1
	silentCellID  = false
)

// alignCmd represents the align command
var alignCmd = &cobra.Command{
	Use:   "align",
	Short: "Continuously fetch data and display timeseries CLI charts",
	Long: `Continuously fetch data and display timeseries CLI charts.
Useful for aligning antennas.`,
	Run: func(cmd *cobra.Command, args []string) {
		a, err := alignui.New(hostname, reqtimeout, pollFrequency, silentCellID)
		if err != nil {
			a.Close()
			fmt.Println(err)
			return
		}
		err = a.Run()
		if err != nil {
			fmt.Println(err)
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(alignCmd)

	alignCmd.PersistentFlags().IntVarP(&pollFrequency, "poll", "x", pollFrequency, "How often to fetch data and redraw")
	alignCmd.PersistentFlags().BoolVarP(&silentCellID, "silent", "z", silentCellID, "Silence cell ID for screenshots (avoid leaking location data unintentionally)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// alignCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// alignCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
