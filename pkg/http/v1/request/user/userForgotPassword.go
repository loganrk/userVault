package user

import (
	"encoding/json"
	"errors"
	"mayilon/pkg/http/v1/request"
	"net/http"
)

func NewUserForgotPassword() *userForgotPassword {
	return &userForgotPassword{}
}

func (u *userForgotPassword) Parse(r *http.Request) error {
	if r.Method == http.MethodPost {
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(u)
		if err != nil {
			return err
		}
	} else {
		u.Username = r.URL.Query().Get("username")
	}

	return nil
}

func (u *userForgotPassword) Validate() error {
	if !request.EmailRegex.MatchString(u.Username) {

		return errors.New("invalid username")
	}
	return nil
}
