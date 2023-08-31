package dto

import (
	"avito2/pkg/model"
	"time"
)

type SegmentDto struct {
	Slug       string     `json:"slug"`
	CreateTime *time.Time `json:"createTime"`
	ExpireTime *time.Time `json:"expireTime"`
}

func NewSegmentDtoFromModel(userSegment *model.UserSegment) *SegmentDto {
	return &SegmentDto{
		Slug:       userSegment.SegmentSlug,
		CreateTime: userSegment.CreateTime,
		ExpireTime: userSegment.ExpireTime,
	}
}
