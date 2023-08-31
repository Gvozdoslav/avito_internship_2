package model

import (
	"time"
)

type UserSegmentStatus int

const (
	Active UserSegmentStatus = iota
	Expired
)

type UserSegment struct {
	Id          int               `json:"id" db:"id"`
	UserId      int               `json:"userId" db:"user_id"`
	SegmentSlug string            `json:"segmentName" db:"segment_slug"`
	AddTime     *time.Time        `json:"addTime" db:"add_time"`
	ExpireTime  *time.Time        `json:"expireTime" db:"expire_time"`
	Status      UserSegmentStatus `json:"status" db:"status"`
}
