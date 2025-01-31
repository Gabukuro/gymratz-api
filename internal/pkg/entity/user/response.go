package user

import "github.com/google/uuid"

type (
	RegisterUserResponse struct {
		Message string `json:"message"`
	}

	LoginUserResponse struct {
		Token *string `json:"token"`
	}

	GetUserProfileResponse struct {
		ID    uuid.UUID `json:"id"`
		Name  string    `json:"name"`
		Email string    `json:"email"`
	}
)
