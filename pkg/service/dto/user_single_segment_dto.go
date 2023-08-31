package dto

import (
	"avito2/pkg/model"
)

type UserSingleSegmentDto struct {
	UserId  int        `json:"userId"`
	Segment SegmentDto `json:"segment"`
}

func NewUserSingleSegmentDto(userId int, userSegment *model.UserSegment) *UserSingleSegmentDto {
	return &UserSingleSegmentDto{
		UserId:  userId,
		Segment: *NewSegmentDtoFromModel(userSegment),
	}
}

func (userDto *UserSingleSegmentDto) ToUserSegment() *model.UserSegment {
	return &model.UserSegment{
		UserId:      userDto.UserId,
		SegmentSlug: userDto.Segment.Slug,
		CreateTime:  userDto.Segment.CreateTime,
		ExpireTime:  userDto.Segment.ExpireTime,
	}
}
