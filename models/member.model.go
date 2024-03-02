package models

import (
	"time"

	"github.com/google/uuid"
)

type Member struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid; default:uuid_generate_v4(); primary_key"`
	Name      string    `json:"name" gorm:"type:varchar(255); not null"`
	Email     string    `json:"email" gorm:"uniqueIndex; not null"`
	Handphone string    `json:"handphone" gorm:"uniqueIndex"`
	Password  string    `json:"password" gorm:"not null"`
	Photo     string    `json:"photo" gorm:"type:varchar(255)"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type CreateMemberRequest struct {
	Name      string    `json:"name" binding:"required"`
	Email     string    `json:"email" binding:"required,email"`
	Handphone string    `json:"handphone,omitempty" binding:"numeric"`
	Password  string    `json:"password" binding:"required,min=8"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
