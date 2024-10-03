package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/niccolot/Chirpy/internal/database"
)


type User struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string `json:"email"`
	HashedPassword string `json:"hashed_password"`
	IsChirpyred bool `json:"is_chirpy_red"`
}

func (u *User) mapUser(user *database.User) {
	u.Id = user.ID
	u.CreatedAt = user.CreatedAt
	u.UpdatedAt = user.UpdatedAt
	u.Email = user.Email
	u.HashedPassword = user.HashedPassword
	u.IsChirpyred = user.IsChirpyRed.Bool
}

type TemplateData struct {
	FileserverHits int32
}

type Chirp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body     string `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}

func (c *Chirp) mapChirp(chirp *database.Chirp) {
	c.Id = chirp.ID
	c.CreatedAt = chirp.CreatedAt
	c.UpdatedAt = chirp.UpdatedAt
	c.Body = chirp.Body
	c.UserId = chirp.UserID
}

