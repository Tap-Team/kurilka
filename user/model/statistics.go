package model

import "fmt"

type IntUserStatistics struct {
	Day   int `json:"day"`
	Week  int `json:"week"`
	Month int `json:"month"`
	Year  int `json:"year"`
}

func NewIntUserStatistics(day int) IntUserStatistics {
	return IntUserStatistics{
		Day:   day,
		Week:  day * 7,
		Month: day * 30,
		Year:  day * 365,
	}
}

type FloatStatisticsUnit float64

func (u FloatStatisticsUnit) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`%.2f`, u)), nil
}

type FloatUserStatistics struct {
	Day   FloatStatisticsUnit `json:"day"`
	Week  FloatStatisticsUnit `json:"week"`
	Month FloatStatisticsUnit `json:"month"`
	Year  FloatStatisticsUnit `json:"year"`
}

func NewFloatUserStatisctics(day float64) FloatUserStatistics {
	return FloatUserStatistics{
		Day:   FloatStatisticsUnit(day),
		Week:  FloatStatisticsUnit(day * 7),
		Month: FloatStatisticsUnit(day * 30),
		Year:  FloatStatisticsUnit(day * 365),
	}
}
