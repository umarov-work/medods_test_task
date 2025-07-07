package dto

import "github.com/google/uuid"

type UserIDResponse struct {
	UserID uuid.UUID `json:"user_id"`
}
