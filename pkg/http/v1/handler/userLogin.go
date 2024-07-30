package handler

import (
	"context"

	request "mayilon/pkg/http/v1/request/user"
	"mayilon/pkg/http/v1/response"

	"mayilon/pkg/types"
	"net/http"
)

func (h *Handler) UserLogin(w http.ResponseWriter, r *http.Request) {
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

	userData := h.Services.User.GetUserByUsername(ctx, req.Username)
	if userData.Id == 0 {
		//res.Status()
		res.SetError("username or password is incorrect")
		res.Send(w)
		return
	}

	attemptStatus := h.Services.User.CheckLoginFailedAttempt(ctx, userData.Id)
	if attemptStatus == types.LOGIN_ATTEMPT_MAX_REACHED {
		res.SetError("max login attempt reached. please try after sometime")
		res.Send(w)
		return

	} else if attemptStatus == types.LOGIN_ATTEMPT_FAILED {
		res.SetError("internal server error")
		res.Send(w)
		return
	}

	passwordMatch := h.Services.User.CheckPassword(ctx, req.Password, userData.Password, userData.Salt)
	if !passwordMatch {
		loginAttempId := h.Services.User.CreateLoginAttempt(ctx, userData.Id, false)
		if loginAttempId == 0 {
			res.SetError("internal server error")
			res.Send(w)
			return
		}

		res.SetError("username or password is incorrect")
		res.Send(w)
		return
	} else {
		loginAttempId := h.Services.User.CreateLoginAttempt(ctx, userData.Id, true)

		if loginAttempId == 0 {
			res.SetError("internal server error")
			res.Send(w)
			return
		}
	}
	userData = h.Services.User.GetUserByUserid(ctx, userData.Id)

	if userData.Status != types.USER_STATUS_ACTIVE {
		if userData.Status == types.USER_STATUS_INACTIVE {
			res.SetError("your account is currently inactive")
		} else if userData.Status == types.USER_STATUS_PENDING {
			res.SetError("your account verification is pending")
		} else {
			res.SetError("your account has been banned")
		}

		res.Send(w)
		return
	}

	token, err := h.Authentication.CreateToken(userData.Id)

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