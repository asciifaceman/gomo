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
	"os"

	"github.com/asciifaceman/gomo/pkg/tmo"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string
var hostname string
var pretty bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gomo",
	Short: "Fetch and log trashcan data",
	Long:  `Fetch and log trashcan data to a long term data store for analyzing`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		c, err := tmo.NewClient(hostname)
		if err != nil {
			fmt.Println(err)
			return
		}

		data, err := c.FetchRadioStatus()
		if err != nil {
			fmt.Println(err)
			return
		}

		if pretty {
			fmt.Printf("Stats for %s\n", hostname)
			fmt.Printf("Connected: %d\n", data.ApCfg[0].ConnectionState)
			fmt.Printf("IP6: %s\n", data.ApCfg[0].IPV6)
			fmt.Println("== LTE =======================")
			fmt.Printf("  RSSI: %5d\n", data.CellLTEStats[0].Stat.RSSICurrent)
			fmt.Printf("  SNR: %6d\n", data.CellLTEStats[0].Stat.SNRCurrent)
			fmt.Printf("  RSRP: %5d\n", data.CellLTEStats[0].Stat.RSRPCurrent)
			fmt.Printf("  RSRQ: %5d\n", data.CellLTEStats[0].Stat.RSRQCurrent)
			fmt.Printf("  Band: %5s\n", data.CellLTEStats[0].Stat.Band)
			fmt.Printf("  CellID: %2s\n", data.CellLTEStats[0].Stat.PhysicalCellID)
			fmt.Println("== 5G =======================")
			fmt.Printf("  SNR: %6d\n", data.Cell5GStats[0].Stat.SNRCurrent)
			fmt.Printf("  RSRP: %5d\n", data.Cell5GStats[0].Stat.RSRPCurrent)
			fmt.Printf("  RSRQ: %5d\n", data.Cell5GStats[0].Stat.RSRQCurrent)
			fmt.Printf("  Band: %5s\n", data.Cell5GStats[0].Stat.Band)
			fmt.Printf("  CellID: %2s\n", data.Cell5GStats[0].Stat.PhysicalCellID)

		} else {
			spew.Dump(data)
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gomo.yaml)")
	rootCmd.PersistentFlags().StringVarP(&hostname, "hostname", "u", "http://192.168.12.1", "hostname of your tmobile trashcan")
	rootCmd.PersistentFlags().BoolVar(&pretty, "pretty", false, "Whether to just print radio data or not")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".gomo" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".gomo")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
