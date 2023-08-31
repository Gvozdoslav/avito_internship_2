package query

import "fmt"

type UserSegmentQuery struct {
	GetUserQuery                  string
	GetUsersQuery                 string
	GetUserSegmentByIdQuery       string
	GetUserSegmentByUserIdAndSlug string
	GetUserActiveSegmentsQuery    string
	AddUserSegmentQuery           string
	RemoveUserSegmentQuery        string
	RemoveAllUserSegmentsQuery    string
}

func NewUserSegmentQuery() *UserSegmentQuery {
	return &UserSegmentQuery{
		GetUserQuery:                  fmt.Sprintf("select * from user_segment where user_id = $1"),
		GetUsersQuery:                 fmt.Sprintf("select * from user_segment"),
		GetUserSegmentByIdQuery:       fmt.Sprintf("select * from user_segment where id = $1"),
		GetUserSegmentByUserIdAndSlug: fmt.Sprintf("select * from user_segment where user_id = $1 and segment_slug = $2"),
		GetUserActiveSegmentsQuery:    fmt.Sprintf("select * from user_segment where user_id = $1 and expire_time > $2"),
		AddUserSegmentQuery:           fmt.Sprintf("insert into user_segment (user_id, segment_slug, create_time, expire_time) values ($1, $2, $3, $4) returning id"),
		RemoveUserSegmentQuery:        fmt.Sprintf("delete from user_segment where user_id = $1 and segment_slug ilike $2"),
		RemoveAllUserSegmentsQuery:    fmt.Sprintf("delete from user_segment where user_id = $1"),
	}
}
