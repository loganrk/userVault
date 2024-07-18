package api

import (
	"context"

	"mayilon/src/http/v1/request"
	"mayilon/src/http/v1/response"

	"mayilon/src/types"
	"net/http"
)

func (a *Api) UserLogin(w http.ResponseWriter, r *http.Request) {
	req := request.NewUserLogin()
	res := response.New()

	err := req.Parse(r)
	if err != nil {
		// TODO log
		res.SetError("invalid request parameters")
		res.Send(w)
		return
	}

	result := req.Validate()
	if result != "" {
		res.SetError(result)
		res.Send(w)
		return
	}

	ctx := context.Background()
	userId := a.Services.User.GetUserIdFromUsername(ctx, req.Username)
	if userId == 0 {
		//res.Status()
		res.SetError("username or password is incorrect")
		res.Send(w)
		return
	}

	attemptStatus := a.Services.User.CheckLoginAttempt(ctx, userId)
	if attemptStatus == types.LOGIN_ATTEMPT_MAX_REACHED {
		res.SetError("max login attempt reached. please try after sometime")
		res.Send(w)
		return

	} else if attemptStatus == types.LOGIN_ATTEMPT_FAILED {
		res.SetError("internal error")
		res.Send(w)
		return
	}

	userData := a.Services.User.GetUserDataFromUsernameAndPassword(ctx, req.Username, req.Password)
	if userData.Id == 0 {
		res.SetError("username or password is incorrect")
		res.Send(w)
		return
	}

	token, err := a.Authentication.CreateToken(userData.Id)
	if err != nil {
		res.SetError("internal error")
		res.Send(w)
		return
	}

	data := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}

	res.SetData(data)
	res.Send(w)
}

func (a *Api) UserRegister(w http.ResponseWriter, r *http.Request) {

	req := request.NewUserRegister()
	res := response.New()

	err := req.Parse(r)
	if err != nil {
		// TODO log
		res.SetError("invalid request parameters")
		res.Send(w)
		return
	}

	result := req.Validate()
	if result != "" {
		res.SetError(result)
		res.Send(w)
		return
	}

}

func (a *Api) UserForgotPassword(w http.ResponseWriter, r *http.Request) {

}

func (a *Api) UserResetPassword(w http.ResponseWriter, r *http.Request) {

}