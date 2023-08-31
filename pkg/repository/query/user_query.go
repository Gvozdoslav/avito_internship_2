package query

import "fmt"

type UserQuery struct {
	GetUserQuery              string
	CreateUserQuery           string
	DeleteUserQuery           string
	DeleteSegmentsByUserQuery string
}

func NewUserQueries() *UserQuery {
	return &UserQuery{
		GetUserQuery:              fmt.Sprintf("select * from users where id = $1"),
		CreateUserQuery:           fmt.Sprintf("insert into users (id) values ($1) returning id"),
		DeleteUserQuery:           fmt.Sprintf("delete from users where id = $1"),
		DeleteSegmentsByUserQuery: fmt.Sprintf("delete from user_segment where user_id = $1"),
	}
}
