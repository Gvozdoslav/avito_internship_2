package service

import (
	"avito2/pkg/model"
	"avito2/pkg/repository"
	"avito2/pkg/service/dto"
	"avito2/pkg/service/impl"
)

type UserService interface {
	GetUserById(id int) (*dto.UserDto, error)
	GetUsers() ([]*dto.UserDto, error)
	GetUserActiveSegments(id int) (*dto.UserDto, error)
	CreateUser(id int) (*dto.UserDto, error)
	AddUserToSegment(userSingleSegmentDto *dto.UserSingleSegmentDto) (*dto.UserDto, error)
	AddUserToSegments(userDto *dto.UserDto) (*dto.UserDto, error)
	RemoveUserFromSegment(userSingleSegmentDto *dto.UserSingleSegmentDto) (*dto.UserDto, error)
	RemoveUserFromSegments(userDto *dto.UserDto) (*dto.UserDto, error)
	UpdateUserSegments(userDto *dto.UserDto) (*dto.UserDto, error)
	GetUserSegmentsDataCsvUrl(userId int) (string, error)
	DeleteUser(id int) error
}

type SegmentService interface {
	GetSegment(slug string) (*model.Segment, error)
	GetAllSegments() ([]*model.Segment, error)
	CreateSegment(slug string) (*model.Segment, error)
	DeleteSegment(slug string) error
}

type Services struct {
	SegmentService
	UserService
}

func NewServices(repositories *repository.Repositories) *Services {
	return &Services{
		SegmentService: impl.NewSegmentServiceImpl(&repositories.SegmentRepository),
		UserService: impl.NewUserServiceImpl(&repositories.UserRepository,
			&repositories.UserSegmentRepository),
	}
}
