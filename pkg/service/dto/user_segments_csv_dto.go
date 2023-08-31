package dto

import "time"

type UserSegmentsCsvDto struct {
	UserId   int        `json:"userId"`
	FromTime *time.Time `json:"fromTime"`
	ToTime   *time.Time `json:"toTime"`
}

func NewUserSegmentsCsvDto(userId int, fromTime *time.Time, toTime *time.Time) *UserSegmentsCsvDto {
	return &UserSegmentsCsvDto{
		UserId:   userId,
		FromTime: fromTime,
		ToTime:   toTime,
	}
}
