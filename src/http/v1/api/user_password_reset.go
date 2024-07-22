package api

import (
	"context"
	"mayilon/src/http/v1/request"
	"mayilon/src/http/v1/response"
	"mayilon/src/types"
	"net/http"
	"time"
)

func (a *Api) UserPasswordReset(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	req := request.NewUserResetPassword()
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

	tokenData := a.Services.User.GetPasswordResetDataByToken(ctx, req.Token)
	if tokenData.Id != 0 {
		res.SetError("invalid token")
		res.Send(w)
		return
	}
	if tokenData.ExpiredAt.Before(time.Now()) {
		res.SetError("activation link expired")
		res.Send(w)
		return
	}

	userData := a.Services.User.GetUserByUserid(ctx, tokenData.UserId)

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

	success := a.Services.User.UpdatePassword(ctx, userData.Id, req.Password, userData.Salt)

	if !success {

		res.SetError("internal server error")
		res.Send(w)
		return
	}

	resData := "password has been reset successfully"
	res.SetData(resData)
	res.Send(w)

}
