package user

import "github.com/Gabukuro/gymratz-api/internal/pkg/response"

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

func (r *RegisterUserRequest) Validate() *response.ErrorDetails {
	var errors response.ErrorDetails

	if r.Name == "" {
		errors = append(errors, ErrorNameIsRequired)
	}

	if r.Email == "" {
		errors = append(errors, ErrorEmailIsRequired)
	}

	if r.Password == "" {
		errors = append(errors, ErrorPasswordIsRequired)
	}

	if len(errors) == 0 {
		return nil
	}

	return &errors
}

func (l *LoginUserRequest) Validate() *response.ErrorDetails {
	var errors response.ErrorDetails

	if l.Email == "" {
		errors = append(errors, ErrorEmailIsRequired)
	}

	if l.Password == "" {
		errors = append(errors, ErrorPasswordIsRequired)
	}

	if len(errors) == 0 {
		return nil
	}

	return &errors
}
