package model

import (
	"time"

	"github.com/Tap-Team/kurilka/internal/model/usermodel"
)

type UserData struct {
	PackPrice           usermodel.PackPrice
	CigaretteDayAmount  usermodel.CigaretteDayAmount
	CigarettePackAmount usermodel.CigarettePackAmount
	AbstinenceTime      time.Time
}

func NewUserData(
	packPrice usermodel.PackPrice,
	cigaretteDayAmount usermodel.CigaretteDayAmount,
	cigarettePackAmount usermodel.CigarettePackAmount,
	abstinenceTime time.Time,
) *UserData {
	return &UserData{
		PackPrice:           packPrice,
		CigaretteDayAmount:  cigaretteDayAmount,
		CigarettePackAmount: cigarettePackAmount,
		AbstinenceTime:      abstinenceTime,
	}
}
