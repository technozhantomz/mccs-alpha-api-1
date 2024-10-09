package logic

import (
	"time"

	"github.com/technoshantoms/mccs-alpha-api/internal/app/repository/mongo"
	"github.com/technoshantoms/mccs-alpha-api/internal/app/types"
	"github.com/spf13/viper"
)

type lostpassword struct{}

var Lostpassword = &lostpassword{}

func (s *lostpassword) Create(l *types.LostPassword) error {
	err := mongo.LostPassword.Create(l)
	if err != nil {
		return err
	}
	return nil
}

func (s *lostpassword) FindByToken(token string) (*types.LostPassword, error) {
	lostPassword, err := mongo.LostPassword.FindByToken(token)
	if err != nil {
		return nil, err
	}
	return lostPassword, nil
}

func (s *lostpassword) FindByEmail(email string) (*types.LostPassword, error) {
	lostPassword, err := mongo.LostPassword.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return lostPassword, nil
}

func (s *lostpassword) SetTokenUsed(token string) error {
	err := mongo.LostPassword.SetTokenUsed(token)
	if err != nil {
		return err
	}
	return nil
}

func (s *lostpassword) IsTokenValid(l *types.LostPassword) bool {
	if time.Now().Sub(l.CreatedAt).Seconds() >= viper.GetFloat64("reset_password_timeout") || l.TokenUsed == true {
		return false
	}
	return true
}

func (s *lostpassword) IsTokenInvalid(l *types.LostPassword) bool {
	if time.Now().Sub(l.CreatedAt).Seconds() >= viper.GetFloat64("reset_password_timeout") || l.TokenUsed == true {
		return true
	}
	return false
}
