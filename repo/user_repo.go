package repo

import (
	"auth-service/repo/model"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
)

func FindUserByUsername(ctx context.Context, username string) (*model.User, error) {
	sql, args, err := psql.
		Select("id", "username", "email", "password", "created_at").
		From("user_tbl").
		Where("username = ?", username).
		ToSql()
	if err != nil {
		return nil, err
	}

	var u model.User
	err = Pool.QueryRow(ctx, sql, args...).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // not found
		}
		return nil, err
	}

	return &u, nil
}

func FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	sql, args, err := psql.
		Select("id", "username", "email", "password", "created_at").
		From("user_tbl").
		Where("email = ?", email).
		ToSql()
	if err != nil {
		return nil, err
	}

	var u model.User
	err = Pool.QueryRow(ctx, sql, args...).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

func SaveUser(ctx context.Context, user model.User) error {
	tx, poolErr := Pool.Begin(ctx)
	if poolErr != nil {
		return poolErr
	}
	defer tx.Rollback(ctx)

	// --- insert user ---
	userSQL, userArgs, usrPsqlErr := psql.
		Insert("user_tbl").
		Columns("id", "username", "email", "password").
		Values(
			user.ID,
			user.Username,
			user.Email,
			user.Password,
		).
		ToSql()
	if usrPsqlErr != nil {
		return usrPsqlErr
	}

	_, saveUsrErr := tx.Exec(ctx, userSQL, userArgs...)
	if saveUsrErr != nil {
		return saveUsrErr
	}

	return tx.Commit(ctx)
}

func FindUserByID(ctx context.Context, id string) (*model.User, error) {
	sql, args, err := psql.
		Select("id", "username", "email", "password", "created_at").
		From("user_tbl").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return nil, err
	}

	var u model.User
	err = Pool.QueryRow(ctx, sql, args...).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // user not found
		}
		return nil, err
	}

	return &u, nil
}
