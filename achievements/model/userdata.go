package model

import (
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/amidtime"
)

type UserData struct {
	PackPrice           usermodel.PackPrice
	CigaretteDayAmount  usermodel.CigaretteDayAmount
	CigarettePackAmount usermodel.CigarettePackAmount
	AbstinenceTime      amidtime.Timestamp
}

func NewUserData(
	packPrice usermodel.PackPrice,
	cigaretteDayAmount usermodel.CigaretteDayAmount,
	cigarettePackAmount usermodel.CigarettePackAmount,
	abstinenceTime amidtime.Timestamp,
) *UserData {
	return &UserData{
		PackPrice:           packPrice,
		CigaretteDayAmount:  cigaretteDayAmount,
		CigarettePackAmount: cigarettePackAmount,
		AbstinenceTime:      abstinenceTime,
	}
}
