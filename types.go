package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/niccolot/Chirpy/internal/database"
)


type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string `json:"email"`
}

func (u *User) mapUser(user *database.User) {
	u.ID = user.ID
	u.CreatedAt = user.CreatedAt
	u.UpdatedAt = user.UpdatedAt
	u.Email = user.Email
}

type TemplateData struct {
	FileserverHits int32
}

