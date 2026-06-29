package service

import (
	"context"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

type rateProtectionServiceAPIKeyRepo struct {
	APIKeyRepository
	created *APIKey
	current *APIKey
	updated *APIKey
}

func (r *rateProtectionServiceAPIKeyRepo) Create(_ context.Context, key *APIKey) error {
	clone := *key
	clone.ID = 101
	clone.Key = key.Key
	r.created = &clone
	return nil
}

func (r *rateProtectionServiceAPIKeyRepo) GetByID(_ context.Context, id int64) (*APIKey, error) {
	if r.current == nil || r.current.ID != id {
		return nil, ErrAPIKeyNotFound
	}
	clone := *r.current
	return &clone, nil
}

func (r *rateProtectionServiceAPIKeyRepo) ExistsByKey(_ context.Context, _ string) (bool, error) {
	return false, nil
}

func (r *rateProtectionServiceAPIKeyRepo) Update(_ context.Context, key *APIKey) error {
	clone := *key
	r.updated = &clone
	return nil
}

type rateProtectionServiceUserRepo struct {
	UserRepository
	user *User
}

func (r rateProtectionServiceUserRepo) GetByID(_ context.Context, id int64) (*User, error) {
	if r.user == nil || r.user.ID != id {
		return nil, ErrUserNotFound
	}
	clone := *r.user
	return &clone, nil
}

type rateProtectionServiceGroupRepo struct {
	GroupRepository
	group *Group
}

func (r rateProtectionServiceGroupRepo) GetByID(_ context.Context, id int64) (*Group, error) {
	if r.group == nil || r.group.ID != id {
		return nil, ErrGroupNotFound
	}
	clone := *r.group
	return &clone, nil
}

type rateProtectionServiceUserGroupRateRepo struct {
	UserGroupRateRepository
	rate *float64
}

func (r rateProtectionServiceUserGroupRateRepo) GetByUserAndGroup(context.Context, int64, int64) (*float64, error) {
	return r.rate, nil
}

func TestAPIKeyServiceCreateDefaultsRateProtectionOff(t *testing.T) {
	groupID := int64(12)
	userRate := 2.5
	apiKeyRepo := &rateProtectionServiceAPIKeyRepo{}
	svc := NewAPIKeyService(
		apiKeyRepo,
		rateProtectionServiceUserRepo{user: &User{ID: 7, Status: StatusActive}},
		rateProtectionServiceGroupRepo{group: &Group{ID: groupID, Status: StatusActive, RateMultiplier: 1.2}},
		nil,
		rateProtectionServiceUserGroupRateRepo{rate: &userRate},
		nil,
		&config.Config{Default: config.DefaultConfig{APIKeyPrefix: "sk-test-", RateMultiplier: 1.0}},
	)

	created, err := svc.Create(context.Background(), 7, CreateAPIKeyRequest{
		Name:    "protected",
		GroupID: &groupID,
	})

	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotNil(t, apiKeyRepo.created)
	require.False(t, apiKeyRepo.created.RateProtectionEnabled)
	require.Zero(t, apiKeyRepo.created.MaxRateMultiplier)
}

func TestAPIKeyServiceCreateDefaultsMaxMultiplierWhenRateProtectionEnabled(t *testing.T) {
	groupID := int64(12)
	userRate := 2.5
	apiKeyRepo := &rateProtectionServiceAPIKeyRepo{}
	svc := NewAPIKeyService(
		apiKeyRepo,
		rateProtectionServiceUserRepo{user: &User{ID: 7, Status: StatusActive}},
		rateProtectionServiceGroupRepo{group: &Group{ID: groupID, Status: StatusActive, RateMultiplier: 1.2}},
		nil,
		rateProtectionServiceUserGroupRateRepo{rate: &userRate},
		nil,
		&config.Config{Default: config.DefaultConfig{APIKeyPrefix: "sk-test-", RateMultiplier: 1.0}},
	)

	enabled := true
	created, err := svc.Create(context.Background(), 7, CreateAPIKeyRequest{
		Name:                  "protected",
		GroupID:               &groupID,
		RateProtectionEnabled: &enabled,
	})

	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotNil(t, apiKeyRepo.created)
	require.True(t, apiKeyRepo.created.RateProtectionEnabled)
	require.InDelta(t, userRate, apiKeyRepo.created.MaxRateMultiplier, 1e-12)
}

func TestAPIKeyServiceUpdatePersistsRateProtectionFields(t *testing.T) {
	groupID := int64(12)
	apiKeyRepo := &rateProtectionServiceAPIKeyRepo{
		current: &APIKey{
			ID:                    101,
			UserID:                7,
			Key:                   "sk-existing",
			Name:                  "existing",
			GroupID:               &groupID,
			Status:                StatusActive,
			RateProtectionEnabled: true,
			MaxRateMultiplier:     1.2,
		},
	}
	svc := NewAPIKeyService(apiKeyRepo, nil, nil, nil, nil, nil, &config.Config{})

	enabled := false
	maxMultiplier := 3.4
	updated, err := svc.Update(context.Background(), 101, 7, UpdateAPIKeyRequest{
		RateProtectionEnabled: &enabled,
		MaxRateMultiplier:     &maxMultiplier,
	})

	require.NoError(t, err)
	require.NotNil(t, updated)
	require.NotNil(t, apiKeyRepo.updated)
	require.False(t, apiKeyRepo.updated.RateProtectionEnabled)
	require.InDelta(t, maxMultiplier, apiKeyRepo.updated.MaxRateMultiplier, 1e-12)
}
