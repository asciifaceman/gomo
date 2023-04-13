package models

import (
	"strconv"

	"github.com/asciifaceman/gomo/pkg/helpers"
	"github.com/asciifaceman/gomo/pkg/radiofreq"
)

/*
	RSRP dBm
	SNR dB
	RSRQ dB
	RSSI dBm

	Trashcan 1.2201.00.0328
*/

const (
	SNR_UNIT                 = "dB"
	RSRQ_UNIT                = "dB"
	RSRP_UNIT                = "dBm"
	RSSI_UNIT                = "dBm"
	SNR_LOWER_BOUND  float64 = -20
	SNR_UPPER_BOUND  float64 = 20
	RSRP_LOWER_BOUND float64 = -115
	RSRP_UPPER_BOUND float64 = -80
	RSRQ_LOWER_BOUND float64 = -20
	RSRQ_UPPER_BOUND float64 = -10
)

type FastmileReturn struct {
	Error error
	Body  *FastmileRadioStatus
}

// StatLTE returns the attached LTE stats
func (f *FastmileReturn) StatLTE() *CellLTEStat {
	return f.Body.CellLTEStats[0].Stat
}

// Stat5G returns the attached 5G stats
func (f *FastmileReturn) Stat5G() *Cell5GStat {
	return f.Body.Cell5GStats[0].Stat
}

// StatCellular returns the attached Cellular stats
func (f *FastmileReturn) StatCellular() *CellularStats {
	return f.Body.CellularStats[0]
}

func (f *FastmileReturn) Status() float64 {
	return float64(f.Body.ConnectionStatus[0].ConnectionStatus)
}

func (f *FastmileReturn) BytesSent() float64 {
	return float64(f.Body.CellularStats[0].BytesSent)
}

func (f *FastmileReturn) BytesRecv() float64 {
	return float64(f.Body.CellularStats[0].BytesReceived)
}

func (f *FastmileReturn) StatEthernet() *EthernetStats {
	return f.Body.EthernetStats[0]
}

type FastmileRadioStatus struct {
	ConnectionStatus []*ConnectionStatus `json:"connection_status"`
	ApCfg            []*ApnCfg           `json:"apn_cfg"`
	CellularStats    []*CellularStats    `json:"cellular_stats"`
	EthernetStats    []*EthernetStats    `json:"ethernet_stats"`
	CellCAStats      []*CellCAStats      `json:"cell_CA_stats_cfg"`
	Cell5GStats      []*Cell5GStats      `json:"cell_5G_stats_cfg"`
	CellLTEStats     []*CellLTEStats     `json:"cell_LTE_stats_cfg"`
}

type ConnectionStatus struct {
	ConnectionStatus int `json:"ConnectionStatus"`
}

type ApnCfg struct {
	OID             int    `json:"oid"`
	Enable          int    `json:"Enable"`
	APN             string `json:"APN"`
	ServiceType     string `json:"X_ALU_COM_ServiceType"`
	ConnectionState int    `json:"X_ALU_COM_ConnectionState"`
	IPV4            string `json:"X_ALU_COM_IPAddressV4"`
	IPV6            string `json:"X_ALU_COM_IPAddressV6"`
}

type CellularStats struct {
	BytesReceived int
	BytesSent     int
}

type EthernetStats struct {
	Enable int                `json:"Enable"`
	Status string             `json:"Status"`
	Stat   *EthernetStatsStat `json:"stat"`
}

type EthernetStatsStat struct {
	BytesSent       int `json:"BytesSent"`
	BytesReceived   int `json:"BytesReceived"`
	PacketsSent     int `json:"PacketsSent"`
	PacketsReceived int `json:"PacketsReceived"`
}

type CellCAStats struct {
	DLCarrierAggregationNumberOfEntries int              `json:"X_ALU_COM_DLCarrierAggregationNumberOfEntries"`
	ULCarrierAggregationNumberOfEntries int              `json:"X_ALU_COM_ULCarrierAggregationNumberOfEntries"`
	Ca4GDL                              *CellCA4GDLStat0 `json:"ca4GDL"`
}

type CellCA4GDLStat0 struct {
	Ca4GDL0 *CellCA4GDLStat `json:"0"`
}

type CellCA4GDLStat struct {
	PhysicalCellID int    `json:"PhysicalCellID"`
	ScellBand      string `json:"ScellBand"`
	ScellChannel   int    `json:"ScellChannel"`
}

type Cell5GStats struct {
	Stat *Cell5GStat `json:"stat"`
}

type Cell5GStat struct {
	SNRCurrent               float64 `json:"SNRCurrent"`
	RSRPCurrent              float64 `json:"RSRPCurrent"`
	RSRPStrengthIndexCurrent float64 `json:"RSRPStrengthIndexCurrent"`
	PhysicalCellID           string  `json:"PhysicalCellID"`
	RSRQCurrent              float64 `json:"RSRQCurrent"`
	DownlinkNRARFCN          float64 `json:"Downlink_NR_ARFCN"`
	SignalStrengthLevel      float64 `json:"SignalStrengthLevel"`
	Band                     string  `json:"Band"`
}

func (c *Cell5GStat) Band64() float64 {
	return radiofreq.BandMap.FrequencyFromShortname(c.Band)
}

func (c *Cell5GStat) ID() float64 {
	id, err := strconv.ParseFloat(c.PhysicalCellID, 64)
	if err != nil {
		id = 0
	}
	return id
}

func (c *Cell5GStat) SNRQuality(min float64, max float64) float64 {
	if c.SNRCurrent < SNR_LOWER_BOUND {
		return min
	}
	if c.SNRCurrent > SNR_UPPER_BOUND {
		return max
	}
	return helpers.NumReMap(c.SNRCurrent, SNR_LOWER_BOUND, SNR_UPPER_BOUND, min, max)

}

func (c *Cell5GStat) RSRPQuality(min float64, max float64) float64 {
	if c.RSRPCurrent < RSRP_LOWER_BOUND {
		return min
	}
	if c.RSRPCurrent > RSRP_UPPER_BOUND {
		return max
	}
	return helpers.NumReMap(c.RSRPCurrent, RSRP_LOWER_BOUND, RSRP_UPPER_BOUND, min, max)
}

func (c *Cell5GStat) RSRQQuality(min float64, max float64) float64 {
	if c.RSRQCurrent < RSRQ_LOWER_BOUND {
		return min
	}
	if c.RSRQCurrent > RSRQ_UPPER_BOUND {
		return max
	}
	return helpers.NumReMap(c.RSRQCurrent, RSRQ_LOWER_BOUND, RSRQ_UPPER_BOUND, min, max)
}

type CellLTEStats struct {
	Stat *CellLTEStat `json:"stat"`
}

type CellLTEStat struct {
	RSSICurrent              float64 `json:"RSSICurrent"`
	SNRCurrent               float64 `json:"SNRCurrent"`
	RSRPCurrent              float64 `json:"RSRPCurrent"`
	RSRPStrengthIndexCurrent float64 `json:"RSRPStrengthIndexCurrent"`
	PhysicalCellID           string  `json:"PhysicalCellID"`
	RSRQCurrent              float64 `json:"RSRQCurrent"`
	DownlinkEarfcn           float64 `json:"DownlinkEarfcn"`
	SignalStrengthLevel      float64 `json:"SignalStrengthLevel"`
	Band                     string  `json:"Band"`
}

func (c *CellLTEStat) Band64() float64 {
	return radiofreq.BandMap.FrequencyFromShortname(c.Band)
}

func (c *CellLTEStat) ID() float64 {
	id, err := strconv.ParseFloat(c.PhysicalCellID, 64)
	if err != nil {
		id = 0
	}
	return id
}

func (c *CellLTEStat) SNRQuality(min float64, max float64) float64 {
	if c.SNRCurrent < SNR_LOWER_BOUND {
		return min
	}
	if c.SNRCurrent > SNR_UPPER_BOUND {
		return max
	}
	return helpers.NumReMap(c.SNRCurrent, SNR_LOWER_BOUND, SNR_UPPER_BOUND, min, max)

}

func (c *CellLTEStat) RSRPQuality(min float64, max float64) float64 {
	if c.RSRPCurrent < RSRP_LOWER_BOUND {
		return min
	}
	if c.RSRPCurrent > RSRP_UPPER_BOUND {
		return max
	}
	return helpers.NumReMap(c.RSRPCurrent, RSRP_LOWER_BOUND, RSRP_UPPER_BOUND, min, max)
}

func (c *CellLTEStat) RSRQQuality(min float64, max float64) float64 {
	if c.RSRQCurrent < RSRQ_LOWER_BOUND {
		return min
	}
	if c.RSRQCurrent > RSRQ_UPPER_BOUND {
		return max
	}
	return helpers.NumReMap(c.RSRQCurrent, RSRQ_LOWER_BOUND, RSRQ_UPPER_BOUND, min, max)
}
