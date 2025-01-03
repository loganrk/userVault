package config

type User interface {
	GetMaxLoginAttempt() int
	GetLoginAttemptSessionPeriod() int
	GetPasswordHashCost() int
	GetActivationLink() string
	GetActivationEmailTemplate() string
	GetActivationLinkExpiry() int
	GetPasswordResetLink() string
	GetPasswordResetTemplate() string
	GetPasswordResetLinkExpiry() int

	GetRefreshTokenEnabled() bool
	GetRefreshTokenRotationEnabled() bool
	GetRefreshTokenExpiry() int

	GetAccessTokenExpiry() int
}

func (u user) GetMaxLoginAttempt() int {

	return u.MaxLoginAttempt
}

func (u user) GetLoginAttemptSessionPeriod() int {

	return u.LoginAttemptSessionPeriod
}
func (u user) GetPasswordHashCost() int {

	return u.PasswordHashCost
}

func (u user) GetActivationLink() string {
	return u.Activation.Link
}

func (u user) GetActivationEmailTemplate() string {
	return u.Activation.TemplatePath
}

func (u user) GetActivationLinkExpiry() int {
	return u.Activation.LinkExpiry
}

func (u user) GetPasswordResetLink() string {
	return u.PasswordReset.Link
}

func (u user) GetPasswordResetTemplate() string {
	return u.PasswordReset.TemplatePath
}

func (u user) GetPasswordResetLinkExpiry() int {
	return u.PasswordReset.LinkExpiry
}

func (u user) GetRefreshTokenEnabled() bool {
	return u.RefreshToken.Enabled
}

func (u user) GetRefreshTokenRotationEnabled() bool {
	return u.RefreshToken.Rotation
}

func (u user) GetRefreshTokenExpiry() int {
	return u.RefreshToken.Expiry
}

func (u user) GetAccessTokenExpiry() int {
	return u.AccessToken.Expiry
}
