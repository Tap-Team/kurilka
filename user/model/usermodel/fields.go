package usermodel

import "github.com/Tap-Team/kurilka/pkg/validate"

type Rank string

const (
	Noob Rank = "Новичок"
)

type Level uint8

const (
	One Level = iota + 1
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
)

type SubscriptionType string

const (
	NONE  SubscriptionType = "NONE"
	TRIAL SubscriptionType = "TRIAL"
	BASIC SubscriptionType = "BASIC"
)

type Name string

func (n Name) String() string {
	return string(n)
}

func (n Name) Min() int {
	return 1
}

func (n Name) Max() int {
	return 15
}

func (n Name) Name() string {
	return "Имя"
}

func (n Name) Validate() error {
	return validate.StringValidate(n)
}

type CigarettePackAmount uint8

func (p CigarettePackAmount) Name() string {
	return "Количество сигарет в пачке"
}

func (p CigarettePackAmount) Int() int64 {
	return int64(p)
}

func (p CigarettePackAmount) Min() int {
	return 1
}

func (p CigarettePackAmount) Max() int {
	return 99
}

func (p CigarettePackAmount) Validate() error {
	return validate.IntValidate(p)
}

type CigaretteDayAmount uint8

func (d CigaretteDayAmount) Min() int {
	return 1
}

func (d CigaretteDayAmount) Max() int {
	return 99
}

func (d CigaretteDayAmount) Int() int64 {
	return int64(d)
}

func (d CigaretteDayAmount) Name() string {
	return "Количество сигарет в день"
}

func (d CigaretteDayAmount) Validate() error {
	return validate.IntValidate(d)
}

type PackPrice float32

func (p PackPrice) Name() string {
	return ""
}

func (p PackPrice) Int() int64 {
	return int64(p)
}

func (p PackPrice) Min() int {
	return 100
}

func (p PackPrice) Max() int {
	return 5000
}

func (p PackPrice) Validate() error {
	return validate.IntValidate(p)
}
