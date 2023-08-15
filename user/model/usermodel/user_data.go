package usermodel

type UserData struct {
	Name                Name                `json:"name"`
	CigaretteDayAmount  CigaretteDayAmount  `json:"cigaretteDayAmount"`
	CigarettePackAmount CigarettePackAmount `json:"cigarettePackAmount"`
	PackPrice           PackPrice           `json:"packPrice"`
	Level               LevelInfo           `json:"level"`
	Subscription        Subscription        `json:"subscription"`
}

func NewUserData(
	name string,
	cigaretteDayAmount uint8,
	cigarettePackAmount uint8,
	packPrice float32,
	level LevelInfo,
	subscription Subscription,
) UserData {
	return UserData{
		Name:                Name(name),
		CigaretteDayAmount:  CigaretteDayAmount(cigaretteDayAmount),
		CigarettePackAmount: CigarettePackAmount(cigarettePackAmount),
		PackPrice:           PackPrice(packPrice),
		Level:               level,
		Subscription:        subscription,
	}
}
