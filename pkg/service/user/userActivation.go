package user

import (
	"context"
	"mayilon/pkg/types"
	"mayilon/pkg/utils"
	"strconv"
	"strings"
	"time"
)

const (
	USER_ACTIVATION_TOKEN_ID_MACRO = "{{tokenId}}"
	USER_ACTIVATION_TOKEN_MACRO    = "{{token}}"
	USER_ACTIVATION_LINK_MACRO     = "{{link}}"
	USER_ACTIVATION_NAME_MACRO     = "{{name}}"
	USER_ACTIVATION_APP_NAME_MACRO = "{{appName}}"
)

func (u *userService) CreateActivationToken(ctx context.Context, userid int) (int, string) {
	activationToken := utils.GenerateRandomString(25)

	tokenData, err := u.store.GetActivationByToken(ctx, activationToken)
	if err != nil {
		return 0, ""
	}

	if tokenData.Id != 0 {
		return u.CreateActivationToken(ctx, userid)
	}

	tokenData = types.UserActivationToken{
		UserId:    userid,
		Token:     activationToken,
		Status:    types.USER_ACTIVATION_TOKEN_STATUS_ACTIVE,
		ExpiredAt: time.Now().Add(time.Duration(u.activationLinkExpiry) * time.Second),
	}

	tokenId, err := u.store.CreateActivation(ctx, tokenData)
	if err != nil {
		return 0, ""
	}
	return tokenId, activationToken
}

func (u *userService) GetActivationLink(tokenId int, token string) string {
	activationLink := u.activationLink
	return u.activationLinkMacroReplacement(activationLink, tokenId, token)

}

func (u *userService) activationLinkMacroReplacement(activationLink string, tokenId int, token string) string {
	s := strings.NewReplacer(
		USER_ACTIVATION_TOKEN_ID_MACRO, strconv.Itoa(tokenId),
		USER_ACTIVATION_TOKEN_MACRO, token)

	return s.Replace(activationLink)

}

func (u *userService) GetActivationEmailTemplate(ctx context.Context, name string, activationLink string) string {
	templatePath := u.conf.activationTemplatePath
	template, err := utils.FindFileContent(templatePath)
	if err != nil {
		return ""
	}
	return u.activationTemplateMacroReplacement(template, name, activationLink)
}

func (u *userService) activationTemplateMacroReplacement(template string, name string, activationLink string) string {
	s := strings.NewReplacer(
		USER_ACTIVATION_APP_NAME_MACRO, u.conf.appName,
		USER_ACTIVATION_NAME_MACRO, name,
		USER_ACTIVATION_LINK_MACRO, activationLink)

	return s.Replace(template)

}

func (u *userService) SendActivation(ctx context.Context, email string, template string) int {

	return types.EMAIL_STATUS_FAILED
}

func (u *userService) GetUserActivationByToken(ctx context.Context, token string) types.UserActivationToken {
	tokenData, err := u.store.GetActivationByToken(ctx, token)
	if err != nil {
		return types.UserActivationToken{}
	}

	return tokenData

}

func (u *userService) UpdatedActivationtatus(ctx context.Context, id int, status int) {
	err := u.store.UpdatedActivationtatus(ctx, id, status)
	if err != nil {

	}
}

func (u *userService) UpdateStatus(ctx context.Context, userid int, status int) bool {

	err := u.store.UpdateStatus(ctx, userid, status)
	if err != nil {
		return false
	}

	return true
}