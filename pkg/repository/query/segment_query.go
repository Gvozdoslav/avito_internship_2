package query

import "fmt"

type SegmentQuery struct {
	GetSegmentBySlugQuery       string
	GetAllSegmentsQuery         string
	CreateSegmentQuery          string
	DeleteSegmentQuery          string
	DeleteUserSegmentsBySegment string
}

func NewSegmentQuery() *SegmentQuery {
	return &SegmentQuery{
		GetSegmentBySlugQuery:       fmt.Sprintf("select * from segment where slug ilike ($1)"),
		GetAllSegmentsQuery:         fmt.Sprintf("select * from segment"),
		CreateSegmentQuery:          fmt.Sprintf("insert into segment (slug) values ($1) returning slug"),
		DeleteSegmentQuery:          fmt.Sprintf("delete from segment where slug ilike $1"),
		DeleteUserSegmentsBySegment: fmt.Sprintf("delete from user_segment where segment_slug ilike $1"),
	}
}
