/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/asciifaceman/gomo/pkg/clients"
	"github.com/asciifaceman/gomo/pkg/clio"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
)

var pretty bool

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Do a single fetch and display",
	Long:  `Do a single fetch and display.`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := clients.NewOneShot(hostname, reqtimeout)
		if err != nil {
			fmt.Println(err)
			return
		}

		resp := c.Fetch()

		if pretty {
			p := clio.NewPrinter(40, 25, 2)
			p.PrintHeader(fmt.Sprintf("Gomo %s", version))

			if resp.Error != nil {
				p.PrintHeader("Failed to fetch data")
				p.PrintKVIndent("Error", resp.Error.Error())
				return
			}

			p.PrintKV("Online", resp.Body.ConnectionStatus[0].ConnectionStatus)
			p.PrintKV("IPV6", resp.Body.ApCfg[0].IPV6)
			p.PrintKV("Bytes Recv", fmt.Sprintf("%d (%.2fGB)", resp.StatCellular().BytesReceived, float64(resp.StatCellular().BytesReceived)*1e-9))
			p.PrintKV("Bytes Sent", fmt.Sprintf("%d (%.2fGB)", resp.StatCellular().BytesSent, float64(resp.StatCellular().BytesSent)*1e-9))
			p.PrintHeader("5G")
			p.PrintKVIndent("Band", resp.Stat5G().Band)
			p.PrintKVIndent("CellID", resp.Stat5G().PhysicalCellID)
			fmt.Println("")
			p.PrintKVIndent("SNR", resp.Stat5G().SNRCurrent)
			p.PrintKVIndent("RSRP", resp.Stat5G().RSRPCurrent)
			p.PrintKVIndent("RSRQ", resp.Stat5G().RSRQCurrent)
			p.PrintHeader("LTE")
			p.PrintKVIndent("Band", resp.StatLTE().Band)
			p.PrintKVIndent("CellID", resp.StatLTE().PhysicalCellID)
			fmt.Println("")
			p.PrintKVIndent("SNR", resp.StatLTE().SNRCurrent)
			p.PrintKVIndent("RSRP", resp.StatLTE().RSRPCurrent)
			p.PrintKVIndent("RSRQ", resp.StatLTE().RSRQCurrent)
			p.PrintKVIndent("RSSI", resp.StatLTE().RSSICurrent)
			p.PrintHeader("Ethernet")
			p.PrintKVIndent("Enabled", resp.StatEthernet().Enable)
			p.PrintKVIndent("Status", resp.StatEthernet().Status)
			fmt.Println("")
			p.PrintKVIndent("Bytes Recv", fmt.Sprintf("%d (%.2fGB)", resp.StatEthernet().Stat.BytesReceived, float64(resp.StatEthernet().Stat.BytesReceived)*1e-9))
			p.PrintKVIndent("Bytes Sent", fmt.Sprintf("%d (%.2fGB)", resp.StatEthernet().Stat.BytesSent, float64(resp.StatEthernet().Stat.BytesSent)*1e-9))
			p.PrintHeader("TODO")
			p.PrintKVIndent("Bring Back", "Ping Stats")

		} else {
			if resp.Error != nil {
				fmt.Println(err)
			}
			spew.Dump(resp)
		}

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
