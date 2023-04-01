package models

import "github.com/asciifaceman/gomo/pkg/helpers"

/*
	RSRP dBm
	SNR dB
	RSRQ dB
	RSSI dBm
*/

const (
	SNR_LOWER_BOUND  = -20
	SNR_UPPER_BOUND  = 20
	RSRP_LOWER_BOUND = -115
	RSRP_UPPER_BOUND = -80
	RSRQ_LOWER_BOUND = -20
	RSRQ_UPPER_BOUND = -10
)

type FastmileReturn struct {
	Error error
	Body  *FastmileRadioStatus
}

type FastmileRadioStatus struct {
	ConnectionStatus []*ConnectionStatus `json:"connection_status"`
	ApCfg            []*ApnCfg           `json:"apn_cfg"`
	CellularStats    []*CellularStats    `json:"cellular_status"`
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
	SNRCurrent               int    `json:"SNRCurrent"`
	RSRPCurrent              int    `json:"RSRPCurrent"`
	RSRPStrengthIndexCurrent int    `json:"RSRPStrengthIndexCurrent"`
	PhysicalCellID           string `json:"PhysicalCellID"`
	RSRQCurrent              int    `json:"RSRQCurrent"`
	DownlinkNRARFCN          int    `json:"Downlink_NR_ARFCN"`
	SignalStrengthLevel      int    `json:"SignalStrengthLevel"`
	Band                     string `json:"Band"`
}

func (c *Cell5GStat) SNRQuality(min float64, max float64) float64 {
	if c.SNRCurrent < SNR_LOWER_BOUND {
		return min
	}
	if c.SNRCurrent > SNR_UPPER_BOUND {
		return max
	}
	return helpers.NumMap(float64(c.SNRCurrent), float64(SNR_LOWER_BOUND), float64(SNR_UPPER_BOUND), float64(min), float64(max))

}

func (c *Cell5GStat) RSRPQuality(min float64, max float64) float64 {
	if c.RSRPCurrent < RSRP_LOWER_BOUND {
		return min
	}
	if c.RSRPCurrent > RSRP_UPPER_BOUND {
		return max
	}
	return helpers.NumMap(float64(c.RSRPCurrent), float64(RSRP_LOWER_BOUND), float64(RSRP_UPPER_BOUND), float64(min), float64(max))
}

func (c *Cell5GStat) RSRQQuality(min float64, max float64) float64 {
	if c.RSRQCurrent < RSRQ_LOWER_BOUND {
		return min
	}
	if c.RSRQCurrent > RSRQ_UPPER_BOUND {
		return max
	}
	return helpers.NumMap(float64(c.RSRQCurrent), float64(RSRQ_LOWER_BOUND), float64(RSRQ_UPPER_BOUND), min, max)
}

type CellLTEStats struct {
	Stat *CellLTEStat `json:"stat"`
}

type CellLTEStat struct {
	RSSICurrent              int    `json:"RSSICurrent"`
	SNRCurrent               int    `json:"SNRCurrent"`
	RSRPCurrent              int    `json:"RSRPCurrent"`
	RSRPStrengthIndexCurrent int    `json:"RSRPStrengthIndexCurrent"`
	PhysicalCellID           string `json:"PhysicalCellID"`
	RSRQCurrent              int    `json:"RSRQCurrent"`
	DownlinkEarfcn           int    `json:"DownlinkEarfcn"`
	SignalStrengthLevel      int    `json:"SignalStrengthLevel"`
	Band                     string `json:"Band"`
}

func (c *CellLTEStat) SNRQuality(min float64, max float64) float64 {
	if c.SNRCurrent < SNR_LOWER_BOUND {
		return min
	}
	if c.SNRCurrent > SNR_UPPER_BOUND {
		return max
	}
	return helpers.NumMap(float64(c.SNRCurrent), float64(SNR_LOWER_BOUND), float64(SNR_UPPER_BOUND), float64(min), float64(max))

}

func (c *CellLTEStat) RSRPQuality(min float64, max float64) float64 {
	if c.RSRPCurrent < RSRP_LOWER_BOUND {
		return min
	}
	if c.RSRPCurrent > RSRP_UPPER_BOUND {
		return max
	}
	return helpers.NumMap(float64(c.RSRPCurrent), float64(RSRP_LOWER_BOUND), float64(RSRP_UPPER_BOUND), float64(min), float64(max))
}

func (c *CellLTEStat) RSRQQuality(min float64, max float64) float64 {
	if c.RSRQCurrent < RSRQ_LOWER_BOUND {
		return min
	}
	if c.RSRQCurrent > RSRQ_UPPER_BOUND {
		return max
	}
	return helpers.NumMap(float64(c.RSRQCurrent), float64(RSRQ_LOWER_BOUND), float64(RSRQ_UPPER_BOUND), min, max)
}
