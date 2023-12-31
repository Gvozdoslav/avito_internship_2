package repository

import (
	"avito2/pkg/model"
	"avito2/pkg/repository/impl"
	"github.com/jmoiron/sqlx"
	"time"
)

type UserRepository interface {
	CreateUser(id int) (*model.User, error)
	DeleteUser(id int) error
}

type SegmentRepository interface {
	GetSegmentBySlug(slug string) (*model.Segment, error)
	GetAllSegments() ([]*model.Segment, error)
	CreateSegment(slug string, percent *int) (*model.Segment, error)
	DeleteSegment(slug string) error
}

type UserSegmentRepository interface {
	GetUserById(userId int) ([]*model.UserSegment, error)
	GetAllUsers() ([]*model.UserSegment, error)
	GetUserActiveSegments(userId int) ([]*model.UserSegment, error)
	GetUserSegmentsDataCsv(userId int, fromTime *time.Time, toTime *time.Time) (string, error)
	AddUserToSegment(userSegment *model.UserSegment) ([]*model.UserSegment, error)
	AddUserToSegments(userId int, userSegment []*model.UserSegment) ([]*model.UserSegment, error)
	RemoveUserFromSegment(userSegment *model.UserSegment) ([]*model.UserSegment, error)
	RemoveUserFromSegments(userId int, userSegments []*model.UserSegment) ([]*model.UserSegment, error)
	RemoveExpiredUserSegments() error
	UpdateUserSegments(userId int, userSegment []*model.UserSegment) ([]*model.UserSegment, error)
}

type Repositories struct {
	UserRepository
	SegmentRepository
	UserSegmentRepository
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		UserRepository:        impl.NewUserRepositoryImpl(db),
		SegmentRepository:     impl.NewSegmentRepositoryImpl(db),
		UserSegmentRepository: impl.NewUserSegmentRepositoryImpl(db),
	}
}
