package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAPIKeyRateMultiplierProtectionExceeded(t *testing.T) {
	groupID := int64(10)
	userRate := 1.5

	tests := []struct {
		name    string
		key     *APIKey
		wantHit bool
	}{
		{
			name: "blocks when effective user group rate exceeds max",
			key: &APIKey{
				GroupID:                 &groupID,
				Group:                   &Group{ID: groupID, RateMultiplier: 1.0},
				UserGroupRateMultiplier: &userRate,
				RateProtectionEnabled:   true,
				MaxRateMultiplier:       1.2,
			},
			wantHit: true,
		},
		{
			name: "allows when effective rate equals max",
			key: &APIKey{
				GroupID:               &groupID,
				Group:                 &Group{ID: groupID, RateMultiplier: 1.2},
				RateProtectionEnabled: true,
				MaxRateMultiplier:     1.2,
			},
			wantHit: false,
		},
		{
			name: "allows when protection disabled",
			key: &APIKey{
				Group:                 &Group{RateMultiplier: 5},
				RateProtectionEnabled: false,
				MaxRateMultiplier:     1,
			},
			wantHit: false,
		},
		{
			name: "allows when max multiplier unset",
			key: &APIKey{
				Group:                 &Group{RateMultiplier: 5},
				RateProtectionEnabled: true,
				MaxRateMultiplier:     0,
			},
			wantHit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantHit, tt.key.RateMultiplierProtectionExceeded())
		})
	}
}
