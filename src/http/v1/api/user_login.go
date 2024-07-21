package api

import (
	"context"

	"mayilon/src/http/v1/request"
	"mayilon/src/http/v1/response"

	"mayilon/src/types"
	"net/http"
)

func (a *Api) UserLogin(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

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

	userData := a.Services.User.GetUserByUsername(ctx, req.Username)
	if userData.Id == 0 {
		//res.Status()
		res.SetError("username or password is incorrect")
		res.Send(w)
		return
	}

	attemptStatus := a.Services.User.CheckLoginAttempt(ctx, userData.Id)
	if attemptStatus == types.LOGIN_ATTEMPT_MAX_REACHED {
		res.SetError("max login attempt reached. please try after sometime")
		res.Send(w)
		return

	} else if attemptStatus == types.LOGIN_ATTEMPT_FAILED {
		res.SetError("internal server error")
		res.Send(w)
		return
	}

	userData = a.Services.User.GetUserByUseridAndPassword(ctx, userData.Id, req.Password, userData.Salt)
	if userData.Id == 0 {
		loginAttempId := a.Services.User.CreateLoginAttempt(ctx, userData.Id, false)
		if loginAttempId == 0 {
			res.SetError("internal server error")
			res.Send(w)
			return
		}

		res.SetError("username or password is incorrect")
		res.Send(w)
		return
	} else {
		loginAttempId := a.Services.User.CreateLoginAttempt(ctx, userData.Id, true)

		if loginAttempId == 0 {
			res.SetError("internal server error")
			res.Send(w)
			return
		}
	}

	if userData.State != types.USER_STATUS_ACTIVE {
		if userData.State == types.USER_STATUS_INACTIVE {
			res.SetError("your account is currently inactive")
		} else if userData.State == types.USER_STATUS_PENDING {
			res.SetError("your account verification is pending")
		} else {
			res.SetError("your account has been banned")
		}

		res.Send(w)
		return
	}

	token, err := a.Authentication.CreateToken(userData.Id)
	if err != nil {
		res.SetError("internal server error")
		res.Send(w)
		return
	}

	resData := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}

	res.SetData(resData)
	res.Send(w)
}