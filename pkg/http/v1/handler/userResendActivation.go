package handler

import (
	"context"
	request "mayilon/pkg/http/v1/request/user"
	"mayilon/pkg/http/v1/response"
	"mayilon/pkg/types"
	"net/http"
)

func (h *Handler) UserResendActivation(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	req := request.NewUserResendActivation()
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
		res.SetError("username is incorrect")
		res.Send(w)
		return
	}

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

	tokenId, activationToken := h.Services.User.CreateActivationToken(ctx, userData.Id)
	if tokenId != 0 && activationToken != "" {
		activationLink := h.Services.User.GetActivationLink(tokenId, activationToken)
		if activationLink != "" {
			template := h.Services.User.GetActivationEmailTemplate(ctx, userData.Name, activationLink)
			if template != "" {
				emailStatus := h.Services.User.SendActivation(ctx, userData.Username, template)
				if emailStatus == types.EMAIL_STATUS_SUCCESS {
					resData := "please check your email for activate account"
					res.SetData(resData)
					res.Send(w)
					return
				}
			}
		}

	}
	res.SetError("internal server error")
	res.Send(w)
}
