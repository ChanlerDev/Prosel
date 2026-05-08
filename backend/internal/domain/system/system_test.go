package system

import (
	"testing"
	"time"
)

func TestSiteSettingPublicValueParsesBooleanAndNumber(t *testing.T) {
	tests := []struct {
		name    string
		setting SiteSetting
		want    any
	}{
		{
			name:    "boolean true",
			setting: SiteSetting{Key: "analytics_enabled", Value: "true", ValueType: ValueTypeBoolean},
			want:    true,
		},
		{
			name:    "number",
			setting: SiteSetting{Key: "posts_per_page", Value: "10", ValueType: ValueTypeNumber},
			want:    float64(10),
		},
		{
			name:    "string",
			setting: SiteSetting{Key: "site_name", Value: "Prosel", ValueType: ValueTypeString},
			want:    "Prosel",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.setting.PublicValue()
			if got != tt.want {
				t.Fatalf("PublicValue() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestHealthStatusReportsUnhealthyWhenDependencyFails(t *testing.T) {
	status := HealthStatus{
		Status:     StatusHealthy,
		DatabaseOK: true,
		RedisOK:    false,
		Version:    "test",
		CheckedAt:  time.Now(),
	}

	status.Normalize()

	if status.Status != StatusUnhealthy {
		t.Fatalf("Status = %q, want %q", status.Status, StatusUnhealthy)
	}
}
