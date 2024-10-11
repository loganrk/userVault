package user

import (
	"encoding/json"
	"errors"
	"mayilon/internal/adapters/handler/http/v1/request"
	"mayilon/internal/utils"
	"net/http"
)

func NewUserResetPassword() *userResetPassword {
	return &userResetPassword{}
}

func (u *userResetPassword) Parse(r *http.Request) error {
	if r.Method == http.MethodPost {
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(u)
		if err != nil {
			return err
		}
	} else {
		u.Password = r.URL.Query().Get("password")
		u.Token = r.URL.Query().Get("token")
	}

	return nil
}

func (u *userResetPassword) Validate() error {
	if u.Token == "" {

		return errors.New("invalid token")
	}

	if u.Password == "" {
		return errors.New("invalid password")
	}

	if !request.PasswordRegex.MatchString(u.Password) {

		return errors.New("password must be between 8 and 12 characters long")
	}

	if !utils.HasDigit(u.Password) {

		return errors.New("password must contain at least one digit")
	}

	if !utils.HasLowercase(u.Password) {

		return errors.New("password must contain at least one lowercase letter")
	}

	if !utils.HasUppercase(u.Password) {

		return errors.New("password must contain at least one uppercase letter")
	}

	if !utils.HasSpecialChar(u.Password) {

		return errors.New("password must contain at least one special character")
	}

	return nil
}