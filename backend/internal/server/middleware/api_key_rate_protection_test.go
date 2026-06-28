package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type rateProtectionAPIKeyRepo struct {
	key *service.APIKey
}

func (r rateProtectionAPIKeyRepo) Create(context.Context, *service.APIKey) error { return nil }
func (r rateProtectionAPIKeyRepo) GetByID(context.Context, int64) (*service.APIKey, error) {
	return nil, errors.New("not implemented")
}
func (r rateProtectionAPIKeyRepo) GetKeyAndOwnerID(context.Context, int64) (string, int64, error) {
	return "", 0, errors.New("not implemented")
}
func (r rateProtectionAPIKeyRepo) GetByKey(_ context.Context, key string) (*service.APIKey, error) {
	if r.key == nil || r.key.Key != key {
		return nil, service.ErrAPIKeyNotFound
	}
	clone := *r.key
	return &clone, nil
}
func (r rateProtectionAPIKeyRepo) GetByKeyForAuth(ctx context.Context, key string) (*service.APIKey, error) {
	return r.GetByKey(ctx, key)
}
func (r rateProtectionAPIKeyRepo) Update(context.Context, *service.APIKey) error { return nil }
func (r rateProtectionAPIKeyRepo) Delete(context.Context, int64) error           { return nil }
func (r rateProtectionAPIKeyRepo) DeleteWithAudit(context.Context, int64) error  { return nil }
func (r rateProtectionAPIKeyRepo) ListByUserID(context.Context, int64, pagination.PaginationParams, service.APIKeyListFilters) ([]service.APIKey, *pagination.PaginationResult, error) {
	return nil, nil, errors.New("not implemented")
}
func (r rateProtectionAPIKeyRepo) VerifyOwnership(context.Context, int64, []int64) ([]int64, error) {
	return nil, errors.New("not implemented")
}
func (r rateProtectionAPIKeyRepo) CountByUserID(context.Context, int64) (int64, error) {
	return 0, errors.New("not implemented")
}
func (r rateProtectionAPIKeyRepo) ExistsByKey(context.Context, string) (bool, error) {
	return false, nil
}
func (r rateProtectionAPIKeyRepo) ListByGroupID(context.Context, int64, pagination.PaginationParams) ([]service.APIKey, *pagination.PaginationResult, error) {
	return nil, nil, errors.New("not implemented")
}
func (r rateProtectionAPIKeyRepo) SearchAPIKeys(context.Context, int64, string, int) ([]service.APIKey, error) {
	return nil, errors.New("not implemented")
}
func (r rateProtectionAPIKeyRepo) ClearGroupIDByGroupID(context.Context, int64) (int64, error) {
	return 0, nil
}
func (r rateProtectionAPIKeyRepo) UpdateGroupIDByUserAndGroup(context.Context, int64, int64, int64) (int64, error) {
	return 0, nil
}
func (r rateProtectionAPIKeyRepo) CountByGroupID(context.Context, int64) (int64, error) {
	return 0, nil
}
func (r rateProtectionAPIKeyRepo) ListKeysByUserID(context.Context, int64) ([]string, error) {
	return nil, nil
}
func (r rateProtectionAPIKeyRepo) ListKeysByGroupID(context.Context, int64) ([]string, error) {
	return nil, nil
}
func (r rateProtectionAPIKeyRepo) IncrementQuotaUsed(context.Context, int64, float64) (float64, error) {
	return 0, nil
}
func (r rateProtectionAPIKeyRepo) UpdateLastUsed(context.Context, int64, time.Time) error {
	return nil
}
func (r rateProtectionAPIKeyRepo) IncrementRateLimitUsage(context.Context, int64, float64) error {
	return nil
}
func (r rateProtectionAPIKeyRepo) ResetRateLimitWindows(context.Context, int64) error {
	return nil
}
func (r rateProtectionAPIKeyRepo) GetRateLimitData(context.Context, int64) (*service.APIKeyRateLimitData, error) {
	return &service.APIKeyRateLimitData{}, nil
}

func newRateProtectionAPIKey(rate, max float64, enabled bool) *service.APIKey {
	groupID := int64(88)
	return &service.APIKey{
		ID:                    100,
		UserID:                7,
		Key:                   "rate-key",
		Status:                service.StatusActive,
		GroupID:               &groupID,
		RateProtectionEnabled: enabled,
		MaxRateMultiplier:     max,
		User: &service.User{
			ID:          7,
			Role:        service.RoleUser,
			Status:      service.StatusActive,
			Balance:     10,
			Concurrency: 3,
		},
		Group: &service.Group{
			ID:               groupID,
			Name:             "g",
			Status:           service.StatusActive,
			Platform:         service.PlatformOpenAI,
			SubscriptionType: service.SubscriptionTypeStandard,
			RateMultiplier:   rate,
			Hydrated:         true,
		},
	}
}

func withUserGroupRateMultiplier(key *service.APIKey, rate float64) *service.APIKey {
	key.UserGroupRateMultiplier = &rate
	return key
}

func TestAPIKeyAuthRateProtection(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		key        *service.APIKey
		path       string
		wantStatus int
		wantCode   string
	}{
		{
			name:       "blocks when effective rate exceeds max",
			key:        newRateProtectionAPIKey(1.5, 1.2, true),
			path:       "/t",
			wantStatus: http.StatusForbidden,
			wantCode:   "API_KEY_RATE_MULTIPLIER_EXCEEDED",
		},
		{
			name:       "blocks when user-specific rate exceeds max",
			key:        withUserGroupRateMultiplier(newRateProtectionAPIKey(1.0, 1.2, true), 1.5),
			path:       "/t",
			wantStatus: http.StatusForbidden,
			wantCode:   "API_KEY_RATE_MULTIPLIER_EXCEEDED",
		},
		{
			name:       "allows when user-specific rate is within max even if group default exceeds",
			key:        withUserGroupRateMultiplier(newRateProtectionAPIKey(1.5, 1.2, true), 1.0),
			path:       "/t",
			wantStatus: http.StatusOK,
		},
		{
			name:       "allows when effective rate equals max",
			key:        newRateProtectionAPIKey(1.2, 1.2, true),
			path:       "/t",
			wantStatus: http.StatusOK,
		},
		{
			name:       "allows when protection disabled",
			key:        newRateProtectionAPIKey(1.5, 1.2, false),
			path:       "/t",
			wantStatus: http.StatusOK,
		},
		{
			name:       "skips protection on usage endpoint",
			key:        newRateProtectionAPIKey(1.5, 1.2, true),
			path:       "/v1/usage",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{RunMode: config.RunModeStandard}
			apiKeyService := service.NewAPIKeyService(rateProtectionAPIKeyRepo{key: tt.key}, nil, nil, nil, nil, nil, cfg)
			router := gin.New()
			var businessLimitedReason string
			router.Use(func(c *gin.Context) {
				c.Next()
				if v, ok := c.Get(service.OpsClientBusinessLimitedReasonKey); ok {
					businessLimitedReason, _ = v.(string)
				}
			})
			router.Use(gin.HandlerFunc(NewAPIKeyAuthMiddleware(apiKeyService, nil, cfg)))
			router.GET("/t", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })
			router.GET("/v1/usage", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			req.Header.Set("x-api-key", tt.key.Key)
			router.ServeHTTP(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantCode != "" {
				var resp ErrorResponse
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				require.Equal(t, tt.wantCode, resp.Code)
				require.Equal(t, service.APIKeyRateMultiplierExceededMessage, resp.Message)
				require.Equal(t, service.OpsClientBusinessLimitedReasonAPIKeyRateMultiplier, businessLimitedReason)
			}
		})
	}
}

func TestAPIKeyAuthGoogleRateProtection(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &config.Config{RunMode: config.RunModeSimple}
	apiKey := withUserGroupRateMultiplier(newRateProtectionAPIKey(1.0, 1.2, true), 1.5)
	apiKeyService := service.NewAPIKeyService(rateProtectionAPIKeyRepo{key: apiKey}, nil, nil, nil, nil, nil, cfg)
	router := gin.New()
	router.Use(APIKeyAuthWithSubscriptionGoogle(apiKeyService, nil, cfg))
	router.GET("/v1beta/test", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1beta/test", nil)
	req.Header.Set("x-goog-api-key", apiKey.Key)
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusForbidden, rec.Code)
	var resp struct {
		Error struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Status  string `json:"status"`
		} `json:"error"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, http.StatusForbidden, resp.Error.Code)
	require.Equal(t, service.APIKeyRateMultiplierExceededMessage, resp.Error.Message)
	require.Equal(t, "PERMISSION_DENIED", resp.Error.Status)
}
