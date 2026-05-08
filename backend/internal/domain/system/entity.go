package system

import (
	"strconv"
	"time"
)

const (
	ValueTypeString  = "string"
	ValueTypeNumber  = "number"
	ValueTypeBoolean = "boolean"
	ValueTypeJSON    = "json"

	StatusHealthy   = "healthy"
	StatusUnhealthy = "unhealthy"
)

type SiteSetting struct {
	ID          string
	Key         string
	Value       string
	ValueType   string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (s SiteSetting) PublicValue() any {
	switch s.ValueType {
	case ValueTypeBoolean:
		value, err := strconv.ParseBool(s.Value)
		if err != nil {
			return false
		}
		return value
	case ValueTypeNumber:
		value, err := strconv.ParseFloat(s.Value, 64)
		if err != nil {
			return float64(0)
		}
		return value
	default:
		return s.Value
	}
}

type HealthStatus struct {
	Status     string
	DatabaseOK bool
	RedisOK    bool
	Version    string
	CheckedAt  time.Time
}

func (h *HealthStatus) Normalize() {
	if h.DatabaseOK && h.RedisOK {
		h.Status = StatusHealthy
		return
	}
	h.Status = StatusUnhealthy
}
