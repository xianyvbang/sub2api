package admin

import (
	"net/http"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

func TestUpdateSettingsEnablesModelMarketplaceAfterItWasDisabled(t *testing.T) {
	h, repo := newStepUpSwitchTestHandler(t, map[string]string{
		service.SettingKeyModelMarketplaceEnabled:       "false",
		service.SettingKeyModelMarketplaceRequiresLogin: "true",
	})

	rec := doUpdateSettings(t, h, map[string]any{
		"model_marketplace_enabled":        true,
		"model_marketplace_requires_login": false,
	}, nil)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "true", repo.values[service.SettingKeyModelMarketplaceEnabled])
	require.Equal(t, "false", repo.values[service.SettingKeyModelMarketplaceRequiresLogin])
	require.Contains(t, rec.Body.String(), `"model_marketplace_enabled":true`)
	require.Contains(t, rec.Body.String(), `"model_marketplace_requires_login":false`)
}

func TestUpdateSettingsPreservesModelMarketplaceFlagsWhenOmitted(t *testing.T) {
	h, repo := newStepUpSwitchTestHandler(t, map[string]string{
		service.SettingKeyModelMarketplaceEnabled:       "true",
		service.SettingKeyModelMarketplaceRequiresLogin: "false",
	})

	rec := doUpdateSettings(t, h, map[string]any{
		"registration_enabled": true,
	}, nil)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "true", repo.values[service.SettingKeyModelMarketplaceEnabled])
	require.Equal(t, "false", repo.values[service.SettingKeyModelMarketplaceRequiresLogin])
	require.Contains(t, rec.Body.String(), `"model_marketplace_enabled":true`)
	require.Contains(t, rec.Body.String(), `"model_marketplace_requires_login":false`)
}
