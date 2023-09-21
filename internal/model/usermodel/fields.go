package usermodel

import (
	"fmt"

	"github.com/Tap-Team/kurilka/internal/errorutils/privacysettingerror"
	"github.com/Tap-Team/kurilka/internal/errorutils/subscriptiontypeerror"
	"github.com/Tap-Team/kurilka/internal/errorutils/triggererror"
	"github.com/Tap-Team/kurilka/pkg/validate"
)

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

func (st SubscriptionType) Validate() error {
	for _, t := range []SubscriptionType{NONE, TRIAL, BASIC} {
		if t == st {
			return nil
		}
	}
	return subscriptiontypeerror.ExceptionSubscriptionTypeNotExists()
}

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
	return "Стоимость пачки"
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

type PrivacySetting string

const (
	STATISTICS_MONEY        PrivacySetting = "STATISTICS_MONEY"
	STATISTICS_CIGARETTE    PrivacySetting = "STATISTICS_CIGARETTE"
	STATISTICS_LIFE         PrivacySetting = "STATISTICS_LIFE"
	STATISTICS_TIME         PrivacySetting = "STATISTICS_TIME"
	ACHIEVEMENTS_DURATION   PrivacySetting = "ACHIEVEMENTS_DURATION"
	ACHIEVEMENTS_HEALTH     PrivacySetting = "ACHIEVEMENTS_HEALTH"
	ACHIEVEMENTS_WELL_BEING PrivacySetting = "ACHIEVEMENTS_WELL_BEING"
	ACHIEVEMENTS_SAVING     PrivacySetting = "ACHIEVEMENTS_SAVING"
	ACHIEVEMENTS_CIGARETTE  PrivacySetting = "ACHIEVEMENTS_CIGARETTE"
)

func (p PrivacySetting) Validate() error {
	for _, setting := range []PrivacySetting{
		STATISTICS_MONEY,
		STATISTICS_CIGARETTE,
		STATISTICS_LIFE,
		STATISTICS_TIME,
		ACHIEVEMENTS_DURATION,
		ACHIEVEMENTS_HEALTH,
		ACHIEVEMENTS_WELL_BEING,
		ACHIEVEMENTS_SAVING,
		ACHIEVEMENTS_CIGARETTE,
	} {
		if setting == p {
			return nil
		}
	}
	return privacysettingerror.ExceptionPrivacySettingNotExist()
}

func (p *PrivacySetting) Scan(src any) error {
	switch src := src.(type) {
	case []byte:
		*p = PrivacySetting(src)
	}
	return nil
}

type Trigger string

const (
	THANK_YOU          Trigger = "THANK_YOU"
	SUPPORT_CIGGARETTE Trigger = "SUPPORT_CIGGARETTE"
	SUPPORT_HEALTH     Trigger = "SUPPORT_HEALTH"
	SUPPORT_TRIAL      Trigger = "SUPPORT_TRIAL"
)

func (t Trigger) Validate() error {
	for _, tr := range [4]Trigger{
		THANK_YOU,
		SUPPORT_CIGGARETTE,
		SUPPORT_HEALTH,
		SUPPORT_TRIAL,
	} {
		if tr == t {
			return nil
		}
	}
	return triggererror.ExceptionTriggerNotExist()
}

func (t *Trigger) Scan(src any) error {
	switch src := src.(type) {
	case string:
		*t = Trigger(src)
	case []byte:
		*t = Trigger(src)
	}
	return nil
}

type Money float64

func (m Money) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`%.2f`, m)), nil
}
