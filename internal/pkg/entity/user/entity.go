package user

type (
	RegisterUserRequest struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	LoginUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
)

func (r *RegisterUserRequest) Validate() bool {
	if r.Name == "" {
		return false
	}

	if r.Email == "" {
		return false
	}

	if r.Password == "" {
		return false
	}

	return true
}

func (l *LoginUserRequest) Validate() bool {
	if l.Email == "" {
		return false
	}

	if l.Password == "" {
		return false
	}

	return true
}
