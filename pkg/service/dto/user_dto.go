package dto

import (
	"avito2/pkg/model"
)

type UserDto struct {
	UserId   int          `json:"userId"`
	Segments []SegmentDto `json:"segments"`
}

func NewUserDto(userId int, userSegments []*model.UserSegment) *UserDto {
	return &UserDto{
		UserId:   userId,
		Segments: getSegments(userSegments),
	}
}

func (userDto *UserDto) ToUserSegments() []*model.UserSegment {
	var userSegments []*model.UserSegment
	for _, segment := range userDto.Segments {
		userSegment := userSegmentFromUserDto(userDto.UserId, &segment)
		userSegments = append(userSegments, userSegment)
	}

	return userSegments
}

func getSegments(userSegments []*model.UserSegment) []SegmentDto {
	var userSegmentDtos = make([]SegmentDto, len(userSegments))
	for index, userSegment := range userSegments {
		userSegmentDtos[index] = *NewSegmentDtoFromModel(userSegment)
	}
	return userSegmentDtos
}

func userSegmentFromUserDto(userId int, segment *SegmentDto) *model.UserSegment {
	return &model.UserSegment{
		UserId:      userId,
		SegmentSlug: segment.Slug,
		AddTime:     segment.AddTime,
		ExpireTime:  segment.ExpireTime,
	}
}
