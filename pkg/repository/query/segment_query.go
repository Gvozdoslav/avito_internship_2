package query

import "fmt"

type SegmentQuery struct {
	GetSegmentBySlugQuery       string
	GetAllSegmentsQuery         string
	GetPercentOfUsersQuery      string
	AddUserToSegmentQuery       string
	CreateSegmentQuery          string
	DeleteSegmentQuery          string
	DeleteUserSegmentsBySegment string
}

func NewSegmentQuery() *SegmentQuery {
	return &SegmentQuery{
		GetSegmentBySlugQuery:       fmt.Sprintf("select * from segment where slug ilike ($1)"),
		GetAllSegmentsQuery:         fmt.Sprintf("select * from segment"),
		GetPercentOfUsersQuery:      fmt.Sprintf("select id from users limit (select count(*) * $1 / 100 from users)"),
		AddUserToSegmentQuery:       fmt.Sprintf("insert into user_segment (user_id, segment_slug, add_time) values ($1, $2, $3)"),
		CreateSegmentQuery:          fmt.Sprintf("insert into segment (slug) values ($1) returning slug"),
		DeleteSegmentQuery:          fmt.Sprintf("delete from segment where slug ilike $1"),
		DeleteUserSegmentsBySegment: fmt.Sprintf("delete from user_segment where segment_slug ilike $1"),
	}
}
