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
	"time"

	"github.com/asciifaceman/gomo/pkg/gotmo"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/spf13/cobra"
)

const (
	scaleFactor = 4
)

// alignCmd represents the align command
var alignCmd = &cobra.Command{
	Use:   "align",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		g, err := gotmo.NewGotmo(hostname, reqtimeout, pingtargets, pingWorkerCount)
		if err != nil {
			fmt.Printf("Failed to setup: %v", err)
		}

		if err := ui.Init(); err != nil {
			fmt.Printf("Failed to initialize UI: %v\n", err)
			return
		}
		defer ui.Close()

		p := widgets.NewParagraph()
		p.Title = "Gomo Alignment"
		p.Text = "Press q TO QUIT"
		p.SetRect(0, 0, 50, 5)
		p.TextStyle.Fg = ui.ColorWhite
		p.BorderStyle.Fg = ui.ColorCyan

		lc5G := widgets.NewPlot()
		lc5G.Title = " 5G "
		lc5G.Data = make([][]float64, 3)
		lc5G.SetRect(0, 5, 75, 25)
		lc5G.AxesColor = ui.ColorWhite
		lc5G.LineColors[0] = ui.ColorRed
		lc5G.LineColors[1] = ui.ColorBlue
		lc5G.LineColors[2] = ui.ColorYellow
		lc5G.Marker = widgets.MarkerDot

		p5GSNRKey := widgets.NewParagraph()
		p5GSNRKey.Text = "SNR"
		p5GSNRKey.TextStyle.Fg = ui.ColorRed
		p5GSNRKey.SetRect(75, 5, 100, 8)

		p5GRSRPKey := widgets.NewParagraph()
		p5GRSRPKey.Text = "RSRP"
		p5GRSRPKey.TextStyle.Fg = ui.ColorBlue
		p5GRSRPKey.SetRect(75, 8, 100, 11)

		p5GRSRQKey := widgets.NewParagraph()
		p5GRSRQKey.Text = "RSRQ"
		p5GRSRQKey.TextStyle.Fg = ui.ColorYellow
		p5GRSRQKey.SetRect(75, 11, 100, 14)

		lcLTE := widgets.NewPlot()
		lcLTE.Title = " LTE "
		lcLTE.Data = make([][]float64, 3)
		lcLTE.SetRect(0, 25, 75, 45)
		lcLTE.AxesColor = ui.ColorWhite
		lcLTE.LineColors[0] = ui.ColorRed
		lcLTE.LineColors[1] = ui.ColorBlue
		lcLTE.LineColors[2] = ui.ColorYellow
		lcLTE.Marker = widgets.MarkerDot

		pLTESNRKey := widgets.NewParagraph()
		pLTESNRKey.Text = "SNR"
		pLTESNRKey.TextStyle.Fg = ui.ColorRed
		pLTESNRKey.SetRect(75, 25, 100, 28)

		pLTERSRPKey := widgets.NewParagraph()
		pLTERSRPKey.Text = "RSRP"
		pLTERSRPKey.TextStyle.Fg = ui.ColorBlue
		pLTERSRPKey.SetRect(75, 28, 100, 31)

		pLTERSRQKey := widgets.NewParagraph()
		pLTERSRQKey.Text = "RSRQ"
		pLTERSRQKey.TextStyle.Fg = ui.ColorYellow
		pLTERSRQKey.SetRect(75, 31, 100, 34)

		slice5GSNR := []float64{}
		highest5GSNR := float64(0)
		slice5GRSRP := []float64{}
		highest5GRSRP := float64(0)
		slice5GRSRQ := []float64{}
		highestRSRQ := float64(0)

		sliceLTESNR := []float64{}
		highestLTESNR := float64(0)
		sliceLTERSRP := []float64{}
		highestLTERSRP := float64(0)
		sliceLTERSRQ := []float64{}
		highestLTERSRQ := float64(0)

		draw := func(count int) {
			radioData := g.AlignEntry()

			// 5G
			this5GSNR := radioData.Cell5GStats[0].Stat.SNRQuality(0, 1)

			if this5GSNR > highest5GSNR {
				highest5GSNR = this5GSNR
			}
			slice5GSNR = append(slice5GSNR, this5GSNR)
			if len(slice5GSNR) > 65 {
				slice5GSNR = slice5GSNR[1:]
			}

			this5GRSRP := radioData.Cell5GStats[0].Stat.RSRPQuality(0, 1)
			if this5GRSRP > highest5GRSRP {
				highest5GRSRP = this5GRSRP
			}
			slice5GRSRP = append(slice5GRSRP, this5GRSRP)
			if len(slice5GRSRP) > 65 {
				slice5GRSRP = slice5GRSRP[1:]
			}

			this5GRSRQ := radioData.Cell5GStats[0].Stat.RSRQQuality(0, 1)
			if this5GRSRQ > highestRSRQ {
				highestRSRQ = this5GRSRQ
			}
			slice5GRSRQ = append(slice5GRSRQ, this5GRSRQ)
			if len(slice5GRSRQ) > 65 {
				slice5GRSRQ = slice5GRSRQ[1:]
			}

			lc5G.Data[0] = slice5GSNR
			lc5G.Data[1] = slice5GRSRP
			lc5G.Data[2] = slice5GRSRQ

			p5GSNRKey.Text = fmt.Sprintf(" SNR (peak: %f)", highest5GSNR)
			p5GRSRPKey.Text = fmt.Sprintf(" RSRP (peak: %f)", highest5GRSRP)
			p5GRSRQKey.Text = fmt.Sprintf(" RSRQ (peak: %f)", highestRSRQ)

			// LTE

			thisLTESNR := radioData.CellLTEStats[0].Stat.SNRQuality(0, 1)
			thisLTERSRP := radioData.CellLTEStats[0].Stat.RSRPQuality(0, 1)
			thisLTERSRQ := radioData.CellLTEStats[0].Stat.RSRQQuality(0, 1)

			if thisLTESNR > highestLTESNR {
				highestLTESNR = thisLTESNR
			}
			sliceLTESNR = append(sliceLTESNR, thisLTESNR)
			if len(sliceLTESNR) > 65 {
				sliceLTESNR = sliceLTESNR[1:]
			}

			if thisLTERSRP > highestLTERSRP {
				highestLTERSRP = thisLTERSRP
			}
			sliceLTERSRP = append(sliceLTERSRP, thisLTERSRP)
			if len(sliceLTERSRP) > 65 {
				sliceLTERSRP = sliceLTERSRP[1:]
			}

			if thisLTERSRQ > highestLTERSRQ {
				highestLTERSRQ = thisLTERSRQ
			}
			sliceLTERSRQ = append(sliceLTERSRQ, thisLTERSRQ)
			if len(sliceLTERSRQ) > 65 {
				sliceLTERSRQ = sliceLTERSRQ[1:]
			}

			lcLTE.Data[0] = sliceLTESNR
			lcLTE.Data[1] = sliceLTERSRP
			lcLTE.Data[2] = sliceLTERSRQ

			pLTESNRKey.Text = fmt.Sprintf(" SNR (peak: %f)", highestLTESNR)
			pLTERSRPKey.Text = fmt.Sprintf(" RSRP (peak: %f)", highestLTERSRP)
			pLTERSRQKey.Text = fmt.Sprintf(" RSRQ (peak: %f)", highestRSRQ)

			ui.Render(p, lc5G, lcLTE, p5GSNRKey, p5GRSRPKey, p5GRSRQKey, pLTESNRKey, pLTERSRPKey, pLTERSRQKey)
			//ui.Render(p, lc5GSNR, lc5GRSRP, lc5GRSRQ, lcLTESNR)
		}

		uiEvents := ui.PollEvents()
		ticker := time.NewTicker(time.Second).C
		tickerCount := 1
		draw(tickerCount)

		for {
			select {
			case e := <-uiEvents:
				switch e.ID {
				case "q", "<C-c>":
					return
				}
			case <-ticker:
				tickerCount++
				if tickerCount > 100 {
					tickerCount = 1
				}
				draw(tickerCount)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(alignCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// alignCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// alignCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
