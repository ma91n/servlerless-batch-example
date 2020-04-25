package unlink

import "time"

type Ncu struct {
	NcuID         string    `json:"ncu_id"`
	NcuType       string    `json:"ncu_type"`
	FirstUplinkAt time.Time `json:"first_uplink_at"`
}

type VMeter struct {
	VmeterID     string  `dynamodbav:"vmeter_id"`
	MeterID      string  `dynamodbav:"meter_id"`
	MeterSize    float64 `dynamodbav:"meter_size"`
	CurrentNcuID string  `dynamodbav:"current_ncu_id"`
}
