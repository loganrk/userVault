package handler

import (
	"context"
	request "mayilon/pkg/http/v1/request/user"
	"mayilon/pkg/http/v1/response"
	"mayilon/pkg/types"
	"net/http"
	"time"
)

func (h *Handler) UserActivation(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	req := request.NewUserActivation()
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

	tokenData := h.Services.User.GetUserActivationByToken(ctx, req.Token)
	if tokenData.Id == 0 {
		res.SetError("invalid token")
		res.Send(w)
		return
	}

	if tokenData.Status != types.USER_ACTIVATION_TOKEN_STATUS_ACTIVE {
		res.SetError("activation token already used")
		res.Send(w)
		return
	}

	if tokenData.ExpiredAt.Before(time.Now()) {
		res.SetError("activation link expired")
		res.Send(w)
		return
	}

	userData := h.Services.User.GetUserByUserid(ctx, tokenData.UserId)

	if userData.Status != types.USER_STATUS_PENDING {
		if userData.Status == types.USER_STATUS_ACTIVE {
			res.SetError("your account is already activated")
		} else if userData.Status == types.USER_STATUS_INACTIVE {
			res.SetError("your account is currently inactive")
		} else {
			res.SetError("your account has been banned")
		}

		res.Send(w)
		return
	}
	h.Services.User.UpdatedActivationtatus(ctx, tokenData.Id, types.USER_ACTIVATION_TOKEN_STATUS_INACTIVE)

	success := h.Services.User.UpdateStatus(ctx, userData.Id, types.USER_STATUS_ACTIVE)

	if !success {

		res.SetError("internal server error")
		res.Send(w)
		return
	}

	resData := "account has been activated successfully"
	res.SetData(resData)
	res.Send(w)
}
