package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/predicate"
	"github.com/Wei-Shaw/sub2api/ent/redeemcode"
	"github.com/Wei-Shaw/sub2api/ent/user"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"

	entsql "entgo.io/ent/dialect/sql"
)

type redeemCodeRepository struct {
	client *dbent.Client
}

func NewRedeemCodeRepository(client *dbent.Client) service.RedeemCodeRepository {
	return &redeemCodeRepository{client: client}
}

func (r *redeemCodeRepository) Create(ctx context.Context, code *service.RedeemCode) error {
	client := clientFromContext(ctx, r.client)
	created, err := client.RedeemCode.Create().
		SetCode(code.Code).
		SetType(code.Type).
		SetValue(code.Value).
		SetStatus(code.Status).
		SetNotes(code.Notes).
		SetValidityDays(code.ValidityDays).
		SetNillableExpiresAt(code.ExpiresAt).
		SetNillableUsedBy(code.UsedBy).
		SetNillableUsedAt(code.UsedAt).
		SetNillableGroupID(code.GroupID).
		Save(ctx)
	if err == nil {
		code.ID = created.ID
		code.CreatedAt = created.CreatedAt
		err = r.updateGiftFields(ctx, code)
	}
	return err
}

func (r *redeemCodeRepository) CreateBatch(ctx context.Context, codes []service.RedeemCode) error {
	if len(codes) == 0 {
		return nil
	}

	client := clientFromContext(ctx, r.client)
	builders := make([]*dbent.RedeemCodeCreate, 0, len(codes))
	for i := range codes {
		c := &codes[i]
		b := client.RedeemCode.Create().
			SetCode(c.Code).
			SetType(c.Type).
			SetValue(c.Value).
			SetStatus(c.Status).
			SetNotes(c.Notes).
			SetValidityDays(c.ValidityDays).
			SetNillableExpiresAt(c.ExpiresAt).
			SetNillableUsedBy(c.UsedBy).
			SetNillableUsedAt(c.UsedAt).
			SetNillableGroupID(c.GroupID)
		builders = append(builders, b)
	}

	if err := client.RedeemCode.CreateBulk(builders...).Exec(ctx); err != nil {
		return err
	}
	for i := range codes {
		if codes[i].UsageLimit == 0 && codes[i].PerUserLimit == 0 && codes[i].GiftParentID == nil {
			continue
		}
		created, err := r.GetByCode(ctx, codes[i].Code)
		if err != nil {
			return err
		}
		codes[i].ID = created.ID
		codes[i].CreatedAt = created.CreatedAt
		if err := r.updateGiftFields(ctx, &codes[i]); err != nil {
			return err
		}
	}
	return nil
}

func (r *redeemCodeRepository) GetByID(ctx context.Context, id int64) (*service.RedeemCode, error) {
	client := clientFromContext(ctx, r.client)
	m, err := client.RedeemCode.Query().
		Where(redeemcode.IDEQ(id)).
		Only(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, service.ErrRedeemCodeNotFound
		}
		return nil, err
	}
	out := redeemCodeEntityToService(m)
	if err := r.loadGiftFields(ctx, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *redeemCodeRepository) GetByCode(ctx context.Context, code string) (*service.RedeemCode, error) {
	client := clientFromContext(ctx, r.client)
	m, err := client.RedeemCode.Query().
		Where(redeemcode.CodeEQ(code)).
		Only(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, service.ErrRedeemCodeNotFound
		}
		return nil, err
	}
	out := redeemCodeEntityToService(m)
	if err := r.loadGiftFields(ctx, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *redeemCodeRepository) GetByCodeForUpdate(ctx context.Context, code string) (*service.RedeemCode, error) {
	client := clientFromContext(ctx, r.client)
	m, err := client.RedeemCode.Query().
		Where(redeemcode.CodeEQ(code)).
		ForUpdate().
		Only(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, service.ErrRedeemCodeNotFound
		}
		return nil, err
	}
	out := redeemCodeEntityToService(m)
	if err := r.loadGiftFields(ctx, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *redeemCodeRepository) Delete(ctx context.Context, id int64) error {
	client := clientFromContext(ctx, r.client)
	_, err := client.RedeemCode.Delete().Where(redeemcode.IDEQ(id)).Exec(ctx)
	return err
}

func (r *redeemCodeRepository) DeleteGiftParent(ctx context.Context, id int64) error {
	if tx := dbent.TxFromContext(ctx); tx != nil {
		return r.deleteGiftParentWithClient(ctx, tx.Client(), id)
	}

	tx, err := r.client.Tx(ctx)
	if err != nil {
		return err
	}
	txCtx := dbent.NewTxContext(ctx, tx)
	defer func() { _ = tx.Rollback() }()

	if err := r.deleteGiftParentWithClient(txCtx, tx.Client(), id); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *redeemCodeRepository) deleteGiftParentWithClient(ctx context.Context, client *dbent.Client, id int64) error {
	parent, err := client.RedeemCode.Query().
		Where(redeemcode.IDEQ(id)).
		ForUpdate().
		Only(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return service.ErrRedeemCodeNotFound
		}
		return err
	}
	if parent.Type != service.RedeemTypeGiftBalance {
		return service.ErrRedeemCodeNotFound
	}

	used, err := r.countGiftChildrenByStatus(ctx, client, id, service.StatusUsed)
	if err != nil {
		return err
	}
	if used > 0 {
		return service.ErrRedeemCodeUsed
	}
	if _, err := client.RedeemCode.Delete().Where(rawRedeemCodeFieldEQ("gift_parent_id", id)).Exec(ctx); err != nil {
		return err
	}
	deleted, err := client.RedeemCode.Delete().Where(redeemcode.IDEQ(id)).Exec(ctx)
	if err != nil {
		return err
	}
	if deleted == 0 {
		return service.ErrRedeemCodeNotFound
	}
	return nil
}

func (r *redeemCodeRepository) DeleteGiftChild(ctx context.Context, id, parentID int64) error {
	if tx := dbent.TxFromContext(ctx); tx != nil {
		return r.deleteGiftChildWithClient(ctx, tx.Client(), id, parentID)
	}

	tx, err := r.client.Tx(ctx)
	if err != nil {
		return err
	}
	txCtx := dbent.NewTxContext(ctx, tx)
	defer func() { _ = tx.Rollback() }()

	if err := r.deleteGiftChildWithClient(txCtx, tx.Client(), id, parentID); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *redeemCodeRepository) deleteGiftChildWithClient(ctx context.Context, client *dbent.Client, id, parentID int64) error {
	parent, err := client.RedeemCode.Query().
		Where(redeemcode.IDEQ(parentID)).
		ForUpdate().
		Only(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return service.ErrRedeemCodeNotFound
		}
		return err
	}
	if parent.Type != service.RedeemTypeGiftBalance {
		return service.ErrRedeemCodeNotFound
	}

	deleted, err := client.RedeemCode.Delete().
		Where(
			redeemcode.IDEQ(id),
			rawRedeemCodeFieldEQ("gift_parent_id", parentID),
			redeemcode.StatusEQ(service.StatusUnused),
		).
		Exec(ctx)
	if err != nil {
		return err
	}
	if deleted == 0 {
		return service.ErrRedeemCodeUsed
	}
	return nil
}

func (r *redeemCodeRepository) List(ctx context.Context, params pagination.PaginationParams) ([]service.RedeemCode, *pagination.PaginationResult, error) {
	return r.ListWithFilters(ctx, params, "", "", "")
}

func (r *redeemCodeRepository) ListWithFilters(ctx context.Context, params pagination.PaginationParams, codeType, status, search string) ([]service.RedeemCode, *pagination.PaginationResult, error) {
	client := clientFromContext(ctx, r.client)
	q := client.RedeemCode.Query().Where(rawRedeemCodeFieldIsNil("gift_parent_id"))

	if codeType != "" {
		q = q.Where(redeemcode.TypeEQ(codeType))
	}
	if status != "" {
		now := time.Now()
		switch status {
		case service.StatusExpired:
			q = q.Where(redeemcode.Or(
				redeemcode.StatusEQ(service.StatusExpired),
				redeemcode.And(
					redeemcode.StatusEQ(service.StatusUnused),
					redeemcode.ExpiresAtNotNil(),
					redeemcode.ExpiresAtLTE(now),
				),
			))
		case service.StatusUnused:
			q = q.Where(
				redeemcode.StatusEQ(service.StatusUnused),
				redeemcode.Or(
					redeemcode.ExpiresAtIsNil(),
					redeemcode.ExpiresAtGT(now),
				),
			)
		default:
			q = q.Where(redeemcode.StatusEQ(status))
		}
	}
	if search != "" {
		q = q.Where(
			redeemcode.Or(
				redeemcode.CodeContainsFold(search),
				redeemcode.HasUserWith(user.EmailContainsFold(search)),
			),
		)
	}

	total, err := q.Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	codesQuery := q.
		WithUser().
		WithGroup().
		Offset(params.Offset()).
		Limit(params.Limit())
	for _, order := range redeemCodeListOrder(params) {
		codesQuery = codesQuery.Order(order)
	}

	codes, err := codesQuery.All(ctx)
	if err != nil {
		return nil, nil, err
	}

	outCodes := redeemCodeEntitiesToService(codes)
	if err := r.loadGiftFieldsForCodes(ctx, outCodes); err != nil {
		return nil, nil, err
	}

	return outCodes, paginationResultFromTotal(int64(total), params), nil
}

func redeemCodeListOrder(params pagination.PaginationParams) []func(*entsql.Selector) {
	sortBy := strings.ToLower(strings.TrimSpace(params.SortBy))
	sortOrder := params.NormalizedSortOrder(pagination.SortOrderDesc)

	var field string
	switch sortBy {
	case "type":
		field = redeemcode.FieldType
	case "value":
		field = redeemcode.FieldValue
	case "status":
		field = redeemcode.FieldStatus
	case "used_at":
		field = redeemcode.FieldUsedAt
	case "created_at":
		field = redeemcode.FieldCreatedAt
	case "expires_at":
		field = redeemcode.FieldExpiresAt
	case "code":
		field = redeemcode.FieldCode
	default:
		field = redeemcode.FieldID
	}

	if sortOrder == pagination.SortOrderAsc {
		return []func(*entsql.Selector){dbent.Asc(field), dbent.Asc(redeemcode.FieldID)}
	}
	return []func(*entsql.Selector){dbent.Desc(field), dbent.Desc(redeemcode.FieldID)}
}

func (r *redeemCodeRepository) Update(ctx context.Context, code *service.RedeemCode) error {
	client := clientFromContext(ctx, r.client)
	up := client.RedeemCode.UpdateOneID(code.ID).
		SetCode(code.Code).
		SetType(code.Type).
		SetValue(code.Value).
		SetStatus(code.Status).
		SetNotes(code.Notes).
		SetValidityDays(code.ValidityDays)

	if code.UsedBy != nil {
		up.SetUsedBy(*code.UsedBy)
	} else {
		up.ClearUsedBy()
	}
	if code.UsedAt != nil {
		up.SetUsedAt(*code.UsedAt)
	} else {
		up.ClearUsedAt()
	}
	if code.GroupID != nil {
		up.SetGroupID(*code.GroupID)
	} else {
		up.ClearGroupID()
	}
	if code.ExpiresAt != nil {
		up.SetExpiresAt(*code.ExpiresAt)
	} else {
		up.ClearExpiresAt()
	}

	updated, err := up.Save(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return service.ErrRedeemCodeNotFound
		}
		return err
	}
	code.CreatedAt = updated.CreatedAt
	return r.updateGiftFields(ctx, code)
}

func (r *redeemCodeRepository) BatchUpdate(ctx context.Context, ids []int64, fields service.RedeemCodeBatchUpdateFields) (int64, error) {
	uniqueIDs := make([]int64, 0, len(ids))
	seen := make(map[int64]struct{}, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		uniqueIDs = append(uniqueIDs, id)
	}
	if len(uniqueIDs) == 0 {
		return 0, nil
	}

	if tx := dbent.TxFromContext(ctx); tx != nil {
		return r.batchUpdate(ctx, tx.Client(), uniqueIDs, fields)
	}

	tx, err := r.client.Tx(ctx)
	if err != nil {
		return 0, err
	}
	txCtx := dbent.NewTxContext(ctx, tx)
	defer func() { _ = tx.Rollback() }()

	updated, err := r.batchUpdate(txCtx, tx.Client(), uniqueIDs, fields)
	if err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return updated, nil
}

func (r *redeemCodeRepository) batchUpdate(ctx context.Context, client *dbent.Client, ids []int64, fields service.RedeemCodeBatchUpdateFields) (int64, error) {
	existing, err := client.RedeemCode.Query().
		Where(redeemcode.IDIn(ids...)).
		All(ctx)
	if err != nil {
		return 0, err
	}
	if len(existing) != len(ids) {
		return 0, service.ErrRedeemCodeNotFound
	}
	if fields.TouchesUsedSensitiveFields() {
		for _, code := range existing {
			if code.Status == service.StatusUsed {
				return 0, service.ErrRedeemCodeUsed
			}
		}
	}

	up := client.RedeemCode.Update().Where(redeemcode.IDIn(ids...))
	if fields.Status != nil {
		up.SetStatus(*fields.Status)
	}
	if fields.Notes != nil {
		up.SetNotes(*fields.Notes)
	}
	if fields.ExpiresAt.Set {
		if fields.ExpiresAt.Value != nil {
			up.SetExpiresAt(*fields.ExpiresAt.Value)
		} else {
			up.ClearExpiresAt()
		}
	}
	if fields.GroupID.Set {
		if fields.GroupID.Value != nil {
			up.SetGroupID(*fields.GroupID.Value)
		} else {
			up.ClearGroupID()
		}
	}

	affected, err := up.Save(ctx)
	if err != nil {
		return 0, err
	}
	if affected != len(ids) {
		return 0, service.ErrRedeemCodeNotFound
	}
	return int64(affected), nil
}

func (r *redeemCodeRepository) Use(ctx context.Context, id, userID int64) error {
	now := time.Now()
	client := clientFromContext(ctx, r.client)
	affected, err := client.RedeemCode.Update().
		Where(redeemcode.IDEQ(id), redeemcode.StatusEQ(service.StatusUnused), rawRedeemCodeFieldIsNil("gift_parent_id")).
		SetStatus(service.StatusUsed).
		SetUsedBy(userID).
		SetUsedAt(now).
		Save(ctx)
	if err != nil {
		return err
	}
	if affected == 0 {
		return service.ErrRedeemCodeUsed
	}
	return nil
}

func (r *redeemCodeRepository) UseGiftChild(ctx context.Context, parentID, userID int64) (*service.RedeemCode, error) {
	client := clientFromContext(ctx, r.client)
	child, err := client.RedeemCode.Query().
		Where(
			rawRedeemCodeFieldEQ("gift_parent_id", parentID),
			redeemcode.TypeEQ(service.RedeemTypeBalance),
			redeemcode.StatusEQ(service.StatusUnused),
			redeemcode.Or(redeemcode.ExpiresAtIsNil(), redeemcode.ExpiresAtGT(time.Now())),
		).
		Order(dbent.Asc(redeemcode.FieldID)).
		ForUpdate().
		First(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, service.ErrRedeemCodeUsed
		}
		return nil, err
	}
	now := time.Now()
	affected, err := client.RedeemCode.Update().
		Where(redeemcode.IDEQ(child.ID), redeemcode.StatusEQ(service.StatusUnused), rawRedeemCodeFieldEQ("gift_parent_id", parentID)).
		SetStatus(service.StatusUsed).
		SetUsedBy(userID).
		SetUsedAt(now).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, service.ErrRedeemCodeUsed
	}
	return r.GetByID(ctx, child.ID)
}

func (r *redeemCodeRepository) ListGiftChildren(ctx context.Context, parentID int64, params pagination.PaginationParams) ([]service.RedeemCode, *pagination.PaginationResult, error) {
	client := clientFromContext(ctx, r.client)
	q := client.RedeemCode.Query().Where(rawRedeemCodeFieldEQ("gift_parent_id", parentID))

	total, err := q.Count(ctx)
	if err != nil {
		return nil, nil, err
	}
	codes, err := q.
		WithUser().
		Offset(params.Offset()).
		Limit(params.Limit()).
		Order(dbent.Asc(redeemcode.FieldID)).
		All(ctx)
	if err != nil {
		return nil, nil, err
	}
	out := redeemCodeEntitiesToService(codes)
	if err := r.loadGiftFieldsForCodes(ctx, out); err != nil {
		return nil, nil, err
	}
	return out, paginationResultFromTotal(int64(total), params), nil
}

func (r *redeemCodeRepository) CountGiftChildrenByStatus(ctx context.Context, parentID int64, status string) (int64, error) {
	return r.countGiftChildrenByStatus(ctx, clientFromContext(ctx, r.client), parentID, status)
}

func (r *redeemCodeRepository) countGiftChildrenByStatus(ctx context.Context, client *dbent.Client, parentID int64, status string) (int64, error) {
	q := client.RedeemCode.Query().Where(rawRedeemCodeFieldEQ("gift_parent_id", parentID))
	if status != "" {
		q = q.Where(redeemcode.StatusEQ(status))
		if status == service.StatusUnused {
			now := time.Now()
			q = q.Where(redeemcode.Or(
				redeemcode.ExpiresAtIsNil(),
				redeemcode.ExpiresAtGT(now),
			))
		}
	}
	n, err := q.Count(ctx)
	return int64(n), err
}

func (r *redeemCodeRepository) CountGiftChildrenByUser(ctx context.Context, parentID, userID int64) (int64, error) {
	client := clientFromContext(ctx, r.client)
	n, err := client.RedeemCode.Query().
		Where(rawRedeemCodeFieldEQ("gift_parent_id", parentID), redeemcode.UsedByEQ(userID)).
		Count(ctx)
	return int64(n), err
}

func (r *redeemCodeRepository) ListByUser(ctx context.Context, userID int64, limit int) ([]service.RedeemCode, error) {
	if limit <= 0 {
		limit = 10
	}

	client := clientFromContext(ctx, r.client)
	codes, err := r.listUserRedeemHistory(ctx, client, userID, pagination.PaginationParams{
		Page:     1,
		PageSize: limit,
	}, "")
	if err != nil {
		return nil, err
	}
	return codes, nil
}

// ListByUserPaginated returns paginated balance/concurrency history for a user.
// Supports optional type filter (e.g. "balance", "admin_balance", "concurrency", "admin_concurrency", "subscription").
func (r *redeemCodeRepository) ListByUserPaginated(ctx context.Context, userID int64, params pagination.PaginationParams, codeType string) ([]service.RedeemCode, *pagination.PaginationResult, error) {
	client := clientFromContext(ctx, r.client)
	total, err := r.countUserRedeemHistory(ctx, client, userID, codeType)
	if err != nil {
		return nil, nil, err
	}
	codes, err := r.listUserRedeemHistory(ctx, client, userID, params, codeType)
	if err != nil {
		return nil, nil, err
	}
	return codes, paginationResultFromTotal(total, params), nil
}

// SumPositiveBalanceByUser returns total recharged amount.
// Gift balance parent codes are audit containers; their used child balance rows
// represent the actual credit events and are counted instead.
func (r *redeemCodeRepository) SumPositiveBalanceByUser(ctx context.Context, userID int64) (float64, error) {
	var result []struct {
		Sum float64 `json:"sum"`
	}
	client := clientFromContext(ctx, r.client)
	err := client.RedeemCode.Query().
		Where(
			redeemcode.UsedByEQ(userID),
			redeemcode.ValueGT(0),
			redeemcode.Or(
				redeemcode.And(
					redeemcode.TypeIn(service.RedeemTypeBalance, service.AdjustmentTypeAdminBalance),
					rawRedeemCodeFieldIsNil("gift_parent_id"),
				),
				redeemcode.And(
					redeemcode.TypeEQ(service.RedeemTypeBalance),
					rawRedeemCodeFieldIsNotNil("gift_parent_id"),
				),
			),
		).
		Aggregate(dbent.As(dbent.Sum(redeemcode.FieldValue), "sum")).
		Scan(ctx, &result)
	if err != nil {
		return 0, err
	}
	if len(result) == 0 {
		return 0, nil
	}
	return result[0].Sum, nil
}

func (r *redeemCodeRepository) countUserRedeemHistory(ctx context.Context, client *dbent.Client, userID int64, codeType string) (int64, error) {
	rows, err := client.QueryContext(ctx, userRedeemHistoryCountSQL(codeType), userID)
	if err != nil {
		return 0, err
	}
	defer func() { _ = rows.Close() }()

	var total int64
	if rows.Next() {
		if err := rows.Scan(&total); err != nil {
			return 0, err
		}
	}
	return total, rows.Err()
}

func (r *redeemCodeRepository) listUserRedeemHistory(ctx context.Context, client *dbent.Client, userID int64, params pagination.PaginationParams, codeType string) ([]service.RedeemCode, error) {
	rows, err := client.QueryContext(ctx, userRedeemHistoryListSQL(codeType), userID, params.Offset(), params.Limit())
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	codes := make([]service.RedeemCode, 0, params.Limit())
	for rows.Next() {
		code, err := scanUserRedeemHistoryRow(rows)
		if err != nil {
			return nil, err
		}
		codes = append(codes, code)
	}
	return codes, rows.Err()
}

func userRedeemHistoryCountSQL(codeType string) string {
	return "SELECT COUNT(*) FROM (" + userRedeemHistoryBaseSQL(codeType) + ") AS history"
}

func userRedeemHistoryListSQL(codeType string) string {
	return userRedeemHistoryBaseSQL(codeType) + `
ORDER BY 7 DESC NULLS LAST, 1 DESC
OFFSET $2
LIMIT $3`
}

func userRedeemHistoryBaseSQL(codeType string) string {
	typePredicate := ""
	giftTypePredicate := ""
	if codeType != "" && !isRedeemHistoryType(codeType) {
		typePredicate = " AND false"
		giftTypePredicate = " AND false"
	} else if codeType != "" {
		typePredicate = " AND c.type = '" + codeType + "'"
		giftTypePredicate = " AND parent.type = '" + codeType + "'"
	}

	return `
SELECT c.id,
       c.code,
       c.type,
       c.value::double precision,
       c.status,
       c.used_by,
       c.used_at,
       c.notes,
       c.created_at,
       c.expires_at,
       c.group_id,
       c.validity_days,
       COALESCE(c.usage_limit, 1),
       COALESCE(c.per_user_limit, 1),
       c.gift_parent_id,
       g.id,
       g.name,
       g.description,
       g.platform,
       g.rate_multiplier,
       g.is_exclusive,
       g.status,
       g.subscription_type,
       g.created_at,
       g.updated_at
FROM redeem_codes c
LEFT JOIN groups g ON g.id = c.group_id
WHERE c.used_by = $1
  AND c.gift_parent_id IS NULL
  AND c.type <> '` + service.RedeemTypeGiftBalance + `'` + typePredicate + `
UNION ALL
SELECT child.id,
       parent.code,
       parent.type,
       child.value::double precision,
       child.status,
       child.used_by,
       child.used_at,
       parent.notes,
       parent.created_at,
       parent.expires_at,
       parent.group_id,
       parent.validity_days,
       COALESCE(parent.usage_limit, 1),
       COALESCE(parent.per_user_limit, 1),
       NULL::bigint,
       g.id,
       g.name,
       g.description,
       g.platform,
       g.rate_multiplier,
       g.is_exclusive,
       g.status,
       g.subscription_type,
       g.created_at,
       g.updated_at
FROM redeem_codes child
JOIN redeem_codes parent ON parent.id = child.gift_parent_id
LEFT JOIN groups g ON g.id = parent.group_id
WHERE child.used_by = $1
  AND parent.type = '` + service.RedeemTypeGiftBalance + `'` + giftTypePredicate
}

func isRedeemHistoryType(codeType string) bool {
	switch codeType {
	case service.RedeemTypeBalance,
		service.RedeemTypeGiftBalance,
		service.RedeemTypeConcurrency,
		service.RedeemTypeSubscription,
		service.RedeemTypeInvitation,
		service.AdjustmentTypeAdminBalance,
		service.AdjustmentTypeAdminConcurrency:
		return true
	default:
		return false
	}
}

func scanUserRedeemHistoryRow(rows *sql.Rows) (service.RedeemCode, error) {
	var code service.RedeemCode
	var usedBy sql.NullInt64
	var usedAt sql.NullTime
	var notes sql.NullString
	var createdAt time.Time
	var expiresAt sql.NullTime
	var groupID sql.NullInt64
	var giftParentID sql.NullInt64
	var historyGroupID sql.NullInt64
	var groupName sql.NullString
	var group service.Group
	var groupDescription sql.NullString
	var groupPlatform sql.NullString
	var groupRate sql.NullFloat64
	var groupExclusive sql.NullBool
	var groupStatus sql.NullString
	var groupSubscriptionType sql.NullString
	var groupCreatedAt sql.NullTime
	var groupUpdatedAt sql.NullTime

	if err := rows.Scan(
		&code.ID,
		&code.Code,
		&code.Type,
		&code.Value,
		&code.Status,
		&usedBy,
		&usedAt,
		&notes,
		&createdAt,
		&expiresAt,
		&groupID,
		&code.ValidityDays,
		&code.UsageLimit,
		&code.PerUserLimit,
		&giftParentID,
		&historyGroupID,
		&groupName,
		&groupDescription,
		&groupPlatform,
		&groupRate,
		&groupExclusive,
		&groupStatus,
		&groupSubscriptionType,
		&groupCreatedAt,
		&groupUpdatedAt,
	); err != nil {
		return service.RedeemCode{}, fmt.Errorf("scan user redeem history: %w", err)
	}

	if usedBy.Valid {
		code.UsedBy = &usedBy.Int64
	}
	if usedAt.Valid {
		code.UsedAt = &usedAt.Time
	}
	if notes.Valid {
		code.Notes = notes.String
	}
	code.CreatedAt = createdAt
	if expiresAt.Valid {
		code.ExpiresAt = &expiresAt.Time
	}
	if groupID.Valid {
		code.GroupID = &groupID.Int64
	}
	if giftParentID.Valid {
		code.GiftParentID = &giftParentID.Int64
	}
	if code.UsageLimit <= 0 {
		code.UsageLimit = 1
	}
	if code.PerUserLimit <= 0 {
		code.PerUserLimit = 1
	}
	if historyGroupID.Valid {
		group.ID = historyGroupID.Int64
		group.Name = nullableStringValue(groupName)
		group.Description = nullableStringValue(groupDescription)
		group.Platform = nullableStringValue(groupPlatform)
		group.RateMultiplier = nullableFloat64Value(groupRate)
		group.IsExclusive = nullableBoolValue(groupExclusive)
		group.Status = nullableStringValue(groupStatus)
		group.SubscriptionType = nullableStringValue(groupSubscriptionType)
		group.CreatedAt = nullableTimeValue(groupCreatedAt)
		group.UpdatedAt = nullableTimeValue(groupUpdatedAt)
		group.Hydrated = true
		code.Group = &group
	}
	return code, nil
}

func nullableStringValue(v sql.NullString) string {
	if !v.Valid {
		return ""
	}
	return v.String
}

func nullableFloat64Value(v sql.NullFloat64) float64 {
	if !v.Valid {
		return 0
	}
	return v.Float64
}

func nullableBoolValue(v sql.NullBool) bool {
	if !v.Valid {
		return false
	}
	return v.Bool
}

func nullableTimeValue(v sql.NullTime) time.Time {
	if !v.Valid {
		return time.Time{}
	}
	return v.Time
}

func redeemCodeEntityToService(m *dbent.RedeemCode) *service.RedeemCode {
	if m == nil {
		return nil
	}
	out := &service.RedeemCode{
		ID:           m.ID,
		Code:         m.Code,
		Type:         m.Type,
		Value:        m.Value,
		Status:       m.Status,
		UsedBy:       m.UsedBy,
		UsedAt:       m.UsedAt,
		Notes:        derefString(m.Notes),
		CreatedAt:    m.CreatedAt,
		ExpiresAt:    m.ExpiresAt,
		GroupID:      m.GroupID,
		ValidityDays: m.ValidityDays,
	}
	if m.Edges.User != nil {
		out.User = userEntityToService(m.Edges.User)
	}
	if m.Edges.Group != nil {
		out.Group = groupEntityToService(m.Edges.Group)
	}
	return out
}

func redeemCodeEntitiesToService(models []*dbent.RedeemCode) []service.RedeemCode {
	out := make([]service.RedeemCode, 0, len(models))
	for i := range models {
		if s := redeemCodeEntityToService(models[i]); s != nil {
			out = append(out, *s)
		}
	}
	return out
}

func rawRedeemCodeFieldEQ(field string, value any) predicate.RedeemCode {
	return predicate.RedeemCode(entsql.FieldEQ(field, value))
}

func rawRedeemCodeFieldIsNil(field string) predicate.RedeemCode {
	return predicate.RedeemCode(entsql.FieldIsNull(field))
}

func rawRedeemCodeFieldIsNotNil(field string) predicate.RedeemCode {
	return predicate.RedeemCode(entsql.FieldNotNull(field))
}

func (r *redeemCodeRepository) updateGiftFields(ctx context.Context, code *service.RedeemCode) error {
	if code == nil || code.ID == 0 {
		return nil
	}
	usageLimit := code.UsageLimit
	if usageLimit <= 0 {
		usageLimit = 1
	}
	perUserLimit := code.PerUserLimit
	if perUserLimit <= 0 {
		perUserLimit = 1
	}

	client := clientFromContext(ctx, r.client)
	_, err := client.ExecContext(
		ctx,
		`UPDATE redeem_codes
		 SET usage_limit = $1, per_user_limit = $2, gift_parent_id = $3
		 WHERE id = $4`,
		usageLimit,
		perUserLimit,
		code.GiftParentID,
		code.ID,
	)
	if err != nil {
		return err
	}
	code.UsageLimit = usageLimit
	code.PerUserLimit = perUserLimit
	return nil
}

func (r *redeemCodeRepository) loadGiftFields(ctx context.Context, code *service.RedeemCode) error {
	if code == nil || code.ID == 0 {
		return nil
	}
	client := clientFromContext(ctx, r.client)
	rows, err := client.QueryContext(
		ctx,
		`SELECT usage_limit, per_user_limit, gift_parent_id
		 FROM redeem_codes
		 WHERE id = $1`,
		code.ID,
	)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	if !rows.Next() {
		return service.ErrRedeemCodeNotFound
	}
	if err := scanGiftFields(rows, code); err != nil {
		return err
	}
	return rows.Err()
}

func (r *redeemCodeRepository) loadGiftFieldsForCodes(ctx context.Context, codes []service.RedeemCode) error {
	for i := range codes {
		if err := r.loadGiftFields(ctx, &codes[i]); err != nil {
			return err
		}
	}
	return nil
}

func scanGiftFields(rows *sql.Rows, code *service.RedeemCode) error {
	var usageLimit int
	var perUserLimit int
	var parentID sql.NullInt64
	if err := rows.Scan(&usageLimit, &perUserLimit, &parentID); err != nil {
		return fmt.Errorf("scan gift redeem fields: %w", err)
	}
	if usageLimit <= 0 {
		usageLimit = 1
	}
	if perUserLimit <= 0 {
		perUserLimit = 1
	}
	code.UsageLimit = usageLimit
	code.PerUserLimit = perUserLimit
	if parentID.Valid {
		id := parentID.Int64
		code.GiftParentID = &id
	} else {
		code.GiftParentID = nil
	}
	return nil
}
