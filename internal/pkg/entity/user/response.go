package user

type (
	RegisterUserResponse struct {
		Message string `json:"message"`
	}

	LoginUserResponse struct {
		Token *string `json:"token"`
	}
)
