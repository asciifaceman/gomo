package models

import "time"

type PingReportReturn struct {
	Error error
	Body  *PingReport
}

type PingReport struct {
	Hostname        string
	PacketsSent     int
	PacketLoss      float64
	AvgResponseTime time.Duration
}
