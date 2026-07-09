package service_test

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/enttest"
	"github.com/Wei-Shaw/sub2api/internal/repository"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "modernc.org/sqlite"
)

func newRedeemServiceTestEnt(t *testing.T) (*sql.DB, *dbent.Client) {
	t.Helper()

	dbName := strings.NewReplacer("/", "_", "\\", "_", " ", "_").Replace(t.Name())
	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", dbName))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	db.SetMaxOpenConns(1)

	drv := entsql.OpenDB(dialect.SQLite, db)
	client := enttest.NewClient(t, enttest.WithOptions(dbent.Driver(drv)))
	t.Cleanup(func() { _ = client.Close() })

	return db, client
}

func TestRedeemGiftBalanceParentCodeCreditsUser(t *testing.T) {
	ctx := context.Background()
	sqlDB, client := newRedeemServiceTestEnt(t)
	redeemRepo := repository.NewRedeemCodeRepository(client)
	userRepo := repository.NewUserRepository(client, sqlDB)

	user, err := client.User.Create().
		SetEmail("gift-parent@example.com").
		SetPasswordHash("hash").
		SetRole(service.RoleUser).
		SetStatus(service.StatusActive).
		SetBalance(2.5).
		SetConcurrency(5).
		Save(ctx)
	require.NoError(t, err)

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
