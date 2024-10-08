package user

import (
	"context"
	"mayilon/pkg/types"

	"gorm.io/gorm"
)

func (s *userStore) CreateActivation(ctx context.Context, tokenData types.UserActivationToken) (int, error) {
	result := s.db.GetDb().WithContext(ctx).Model(&types.UserActivationToken{}).Create(&tokenData)
	return tokenData.Id, result.Error
}

func (s *userStore) GetActivationByToken(ctx context.Context, token string) (types.UserActivationToken, error) {
	var tokenData types.UserActivationToken

	result := s.db.GetDb().WithContext(ctx).Model(&types.UserActivationToken{}).Select("id").Where("token = ?", token).First(&tokenData)
	if result.Error == gorm.ErrRecordNotFound {
		result.Error = nil
	}

	return tokenData, result.Error
}

func (s *userStore) UpdatedActivationtatus(ctx context.Context, id int, status int) error {
	result := s.db.GetDb().WithContext(ctx).Model(&types.UserActivationToken{}).Where("id = ?", id).Update("status", status)
	return result.Error

}

func (s *userStore) UpdateStatus(ctx context.Context, userid int, status int) error {
	result := s.db.GetDb().WithContext(ctx).Model(&types.UserActivationToken{}).Where("id = ?", userid).Update("status", status)
	return result.Error
}
