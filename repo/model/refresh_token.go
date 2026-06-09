package model

import (
	"github.com/google/uuid"
	"time"
)

type RefreshToken struct {
	ID        uuid.UUID
	Token     string
	UserID    uuid.UUID
	ExpiresAt time.Time
	Revoked   bool
	CreatedAt time.Time
}
