package model

import (
	"time"
)

type UserSegment struct {
	Id          int        `json:"id" db:"id"`
	UserId      int        `json:"userId" db:"user_id"`
	SegmentSlug string     `json:"segmentName" db:"segment_slug"`
	CreateTime  *time.Time `json:"createTime" db:"create_time"`
	ExpireTime  *time.Time `json:"expireTime" db:"expire_time"`
}
