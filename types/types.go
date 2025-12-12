package types

import "github.com/google/uuid"

type Id uuid.UUID

type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Biography string `json:"bio"`
	ID string `json:"id"`
}