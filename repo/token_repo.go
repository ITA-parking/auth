package repo

import (
	"auth-service/repo/model"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
)

func SaveRefreshToken(ctx context.Context, token model.RefreshToken) error {
	sql, args, err := psql.
		Insert("refresh_token_tbl").
		Columns("id", "token", "user_id", "expires_at", "revoked").
		Values(
			token.ID,
			token.Token,
			token.UserID,
			token.ExpiresAt,
			token.Revoked,
		).
		ToSql()
	if err != nil {
		return err
	}

	_, err = Pool.Exec(ctx, sql, args...)
	return err
}

func FindRefreshToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	sql, args, err := psql.
		Select("id", "token", "user_id", "expires_at", "revoked", "created_at").
		From("refresh_token_tbl").
		Where("token = ?", token).
		ToSql()
	if err != nil {
		return nil, err
	}

	var rt model.RefreshToken
	err = Pool.QueryRow(ctx, sql, args...).Scan(
		&rt.ID,
		&rt.Token,
		&rt.UserID,
		&rt.ExpiresAt,
		&rt.Revoked,
		&rt.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &rt, nil
}

func RevokeRefreshToken(ctx context.Context, token string) error {
	sql, args, err := psql.
		Update("refresh_token_tbl").
		Set("revoked", true).
		Where("token = ?", token).
		ToSql()
	if err != nil {
		return err
	}

	_, err = Pool.Exec(ctx, sql, args...)
	return err
}
