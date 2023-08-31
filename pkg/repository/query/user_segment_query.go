package query

import "fmt"

type UserSegmentQuery struct {
	GetUserQuery                        string
	GetUsersQuery                       string
	GetUserSegmentByIdQuery             string
	GetUserSegmentsByPeriod             string
	GetActiveUserSegmentByUserIdAndSlug string
	GetUserActiveSegmentsQuery          string
	AddUserSegmentQuery                 string
	RemoveUserSegmentQuery              string
	RemoveAllUserSegmentsQuery          string
	SetExpiredUserSegmentsQuery         string
}

func NewUserSegmentQuery() *UserSegmentQuery {
	return &UserSegmentQuery{
		GetUserQuery:                        fmt.Sprintf("select * from user_segment where user_id = $1"),
		GetUsersQuery:                       fmt.Sprintf("select * from user_segment"),
		GetUserSegmentByIdQuery:             fmt.Sprintf("select * from user_segment where id = $1"),
		GetActiveUserSegmentByUserIdAndSlug: fmt.Sprintf("select * from user_segment where user_id = $1 and segment_slug = $2 and status = $3"),
		GetUserActiveSegmentsQuery:          fmt.Sprintf("select * from user_segment where user_id = $1 and expire_time > $2"),
		GetUserSegmentsByPeriod:             fmt.Sprintf("select * from user_segment where user_id = $1 and add_time > $2 and add_time < $3"),
		AddUserSegmentQuery:                 fmt.Sprintf("insert into user_segment (user_id, segment_slug, add_time, expire_time) values ($1, $2, $3, $4) returning id"),
		RemoveUserSegmentQuery:              fmt.Sprintf("delete from user_segment where user_id = $1 and segment_slug ilike $2"),
		RemoveAllUserSegmentsQuery:          fmt.Sprintf("delete from user_segment where user_id = $1"),
		SetExpiredUserSegmentsQuery:         fmt.Sprintf("update user_segment set status = $1 where expire_time < $2"),
	}
}
