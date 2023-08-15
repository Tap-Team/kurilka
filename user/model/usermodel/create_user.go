package usermodel

import "github.com/Tap-Team/kurilka/pkg/validate"

type CreateUser struct {
	Name                Name                `json:"name"`
	CigaretteDayAmount  CigaretteDayAmount  `json:"cigaretteDayAmount"`
	CigarettePackAmount CigarettePackAmount `json:"cigarettePackAmount"`
	PackPrice           PackPrice           `json:"packPrice"`
}

func NewCreateUser(
	name string,
	cigaretteDayAmount uint8,
	cigarettePackAmount uint8,
	packPrice float32,
) *CreateUser {
	return &CreateUser{
		Name:                Name(name),
		CigaretteDayAmount:  CigaretteDayAmount(cigaretteDayAmount),
		CigarettePackAmount: CigarettePackAmount(cigarettePackAmount),
		PackPrice:           PackPrice(packPrice),
	}
}

func (u *CreateUser) ValidatableVariables() []validate.Validatable {
	return []validate.Validatable{u.Name, u.CigaretteDayAmount, u.CigarettePackAmount, u.PackPrice}
}
