package sqlite

import (
	"context"
	"database/sql"
	"time"

	"school-management/internal/domain/identity"
	"school-management/internal/domain/shared"
)

type RecoveryCodeRepository struct {
	q Querier
}

func (r *RecoveryCodeRepository) ReplaceForUser(ctx context.Context, userID shared.ID, codes []identity.RecoveryCode, at time.Time) error {
	if err := r.q.ReplaceUserRecoveryCodes(ctx, userID.String()); err != nil {
		return shared.Internal(err)
	}
	createdAt := at.UTC().Format(time.RFC3339Nano)
	for _, code := range codes {
		err := r.q.InsertRecoveryCode(ctx, InsertRecoveryCodeParams{
			ID:        code.ID.String(),
			UserID:    userID.String(),
			CodeHash:  code.CodeHash,
			UsedAt:    sql.NullString{},
			CreatedAt: createdAt,
		})
		if err != nil {
			return shared.Internal(err)
		}
	}
	return nil
}

func (r *RecoveryCodeRepository) ListUnused(ctx context.Context, userID shared.ID) ([]identity.RecoveryCode, error) {
	rows, err := r.q.ListUnusedRecoveryCodes(ctx, userID.String())
	if err != nil {
		return nil, shared.Internal(err)
	}
	codes := make([]identity.RecoveryCode, len(rows))
	for i, row := range rows {
		codes[i] = identity.RecoveryCode{
			ID:       shared.ID(row.ID),
			UserID:   shared.ID(row.UserID),
			CodeHash: row.CodeHash,
		}
	}
	return codes, nil
}

func (r *RecoveryCodeRepository) MarkUsed(ctx context.Context, id shared.ID, at time.Time) error {
	return r.q.MarkRecoveryCodeUsed(ctx, MarkRecoveryCodeUsedParams{
		ID:     id.String(),
		UsedAt: sql.NullString{String: at.UTC().Format(time.RFC3339Nano), Valid: true},
	})
}

