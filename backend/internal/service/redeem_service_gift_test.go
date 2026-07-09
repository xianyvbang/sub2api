package service_test

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/enttest"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "modernc.org/sqlite"
)

type giftRedeemRepoStub struct {
	service.RedeemCodeRepository

	nextID int64
	codes  map[int64]*service.RedeemCode
	byCode map[string]int64
}

func newGiftRedeemRepoStub() *giftRedeemRepoStub {
	return &giftRedeemRepoStub{
		nextID: 1,
		codes:  make(map[int64]*service.RedeemCode),
		byCode: make(map[string]int64),
	}
}

func (r *giftRedeemRepoStub) Create(_ context.Context, code *service.RedeemCode) error {
	if code.ID == 0 {
		code.ID = r.nextID
		r.nextID++
	}
	clone := *code
	r.codes[clone.ID] = &clone
	r.byCode[clone.Code] = clone.ID
	return nil
}

func (r *giftRedeemRepoStub) GetByID(_ context.Context, id int64) (*service.RedeemCode, error) {
	code, ok := r.codes[id]
	if !ok {
		return nil, service.ErrRedeemCodeNotFound
	}
	clone := *code
	return &clone, nil
}

func (r *giftRedeemRepoStub) GetByCode(ctx context.Context, code string) (*service.RedeemCode, error) {
	id, ok := r.byCode[code]
	if !ok {
		return nil, service.ErrRedeemCodeNotFound
	}
	return r.GetByID(ctx, id)
}

func (r *giftRedeemRepoStub) GetByCodeForUpdate(ctx context.Context, code string) (*service.RedeemCode, error) {
	return r.GetByCode(ctx, code)
}

func (r *giftRedeemRepoStub) Update(_ context.Context, code *service.RedeemCode) error {
	if _, ok := r.codes[code.ID]; !ok {
		return service.ErrRedeemCodeNotFound
	}
	clone := *code
	r.codes[clone.ID] = &clone
	r.byCode[clone.Code] = clone.ID
	return nil
}

func (r *giftRedeemRepoStub) UseGiftChild(_ context.Context, parentID, userID int64) (*service.RedeemCode, error) {
	for _, code := range r.codes {
		if code.GiftParentID == nil || *code.GiftParentID != parentID || code.Status != service.StatusUnused {
			continue
		}
		now := time.Now()
		code.Status = service.StatusUsed
		code.UsedBy = &userID
		code.UsedAt = &now
		clone := *code
		return &clone, nil
	}
	return nil, service.ErrRedeemCodeUsed
}

func (r *giftRedeemRepoStub) CountGiftChildrenByStatus(_ context.Context, parentID int64, status string) (int64, error) {
	var count int64
	for _, code := range r.codes {
		if code.GiftParentID != nil && *code.GiftParentID == parentID && code.Status == status {
			count++
		}
	}
	return count, nil
}

func (r *giftRedeemRepoStub) CountGiftChildrenByUser(_ context.Context, parentID, userID int64) (int64, error) {
	var count int64
	for _, code := range r.codes {
		if code.GiftParentID != nil && *code.GiftParentID == parentID && code.UsedBy != nil && *code.UsedBy == userID {
			count++
		}
	}
	return count, nil
}

type giftUserRepoStub struct {
	service.UserRepository

	users map[int64]*service.User
}

func newGiftUserRepoStub(users ...*service.User) *giftUserRepoStub {
	repo := &giftUserRepoStub{users: make(map[int64]*service.User)}
	for _, user := range users {
		clone := *user
		repo.users[clone.ID] = &clone
	}
	return repo
}

func (r *giftUserRepoStub) GetByID(_ context.Context, id int64) (*service.User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, service.ErrUserNotFound
	}
	clone := *user
	return &clone, nil
}

func (r *giftUserRepoStub) UpdateBalance(_ context.Context, id int64, amount float64) error {
	user, ok := r.users[id]
	if !ok {
		return service.ErrUserNotFound
	}
	user.Balance += amount
	return nil
}

func newRedeemServiceTestEnt(t *testing.T) *dbent.Client {
	t.Helper()

	dbName := strings.NewReplacer("/", "_", "\\", "_", " ", "_").Replace(t.Name())
	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", dbName))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	db.SetMaxOpenConns(1)

	drv := entsql.OpenDB(dialect.SQLite, db)
	client := enttest.NewClient(t, enttest.WithOptions(dbent.Driver(drv)))
	t.Cleanup(func() { _ = client.Close() })

	return client
}

func TestRedeemGiftBalanceParentCodeCreditsUser(t *testing.T) {
	ctx := context.Background()
	client := newRedeemServiceTestEnt(t)
	redeemRepo := newGiftRedeemRepoStub()

	user := &service.User{
		ID:          1,
		Email:       "gift-parent@example.com",
		Role:        service.RoleUser,
		Status:      service.StatusActive,
		Balance:     2.5,
		Concurrency: 5,
	}
	userRepo := newGiftUserRepoStub(user)

	parent := &service.RedeemCode{
		Code:         "GIFT-PARENT-1",
		Type:         service.RedeemTypeGiftBalance,
		Status:       service.StatusUnused,
		UsageLimit:   1,
		PerUserLimit: 1,
	}
	require.NoError(t, redeemRepo.Create(ctx, parent))

	parentID := parent.ID
	child := &service.RedeemCode{
		Code:         "GIFT-CHILD-1",
		Type:         service.RedeemTypeBalance,
		Value:        7.25,
		Status:       service.StatusUnused,
		GiftParentID: &parentID,
	}
	require.NoError(t, redeemRepo.Create(ctx, child))

	redeemService := service.NewRedeemService(redeemRepo, userRepo, nil, nil, nil, client, nil, nil)

	got, err := redeemService.Redeem(ctx, user.ID, parent.Code)

	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, parent.ID, got.ID)
	require.Equal(t, service.RedeemTypeGiftBalance, got.Type)
	require.Equal(t, service.StatusUsed, got.Status)

	updatedUser, err := userRepo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	require.InDelta(t, 9.75, updatedUser.Balance, 0.000001)

	usedChild, err := redeemRepo.GetByID(ctx, child.ID)
	require.NoError(t, err)
	require.Equal(t, service.StatusUsed, usedChild.Status)
	require.NotNil(t, usedChild.UsedBy)
	require.Equal(t, user.ID, *usedChild.UsedBy)
}
