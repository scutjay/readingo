package model

import "time"

type OperateRedisReq struct {
	Action string `json:"action" binding:"required"`
	Key    string `json:"key"`
	Field  string `json:"field"`
	Value  string `json:"value"`
	Score  string `json:"score"`
	Min    string `json:"min"`
	Max    string `json:"max"`
	TTL    string `json:"ttl"`

	Host  string `json:"host" binding:"required"`
	Index string `json:"index" binding:"required"`
	Type  string `json:"type" binding:"required"`
}

type GetDBTreeResp struct {
	LastUpdateTime time.Time  `json:"lastUpdateTime"`
	Hosts          []HostInfo `json:"hosts"`
}

type HostInfo struct {
	Host       string          `json:"host"'`
	Type       string          `json:"type"`
	Partitions []PartitionInfo `json:"partitions"`
}

type PartitionInfo struct {
	Index  int   `json:"index"`
	DBSize int64 `json:"DBSize"`
}

type GetSupportedActionResp struct {
	Actions []Action `json:"actions"`
}

type Action struct {
	Action         string   `json:"action"`
	RequiredParams []string `json:"requiredParams"`
	Tips           string   `json:"tips"`
}
