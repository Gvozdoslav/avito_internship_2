package impl

import (
	"avito2/pkg/model"
	"avito2/pkg/repository"
	"avito2/pkg/service/dto"
)

type UserServiceImpl struct {
	userRepository        *repository.UserRepository
	userSegmentRepository *repository.UserSegmentRepository
}

func NewUserServiceImpl(userRepository *repository.UserRepository,
	userSegmentRepository *repository.UserSegmentRepository) *UserServiceImpl {

	return &UserServiceImpl{
		userRepository:        userRepository,
		userSegmentRepository: userSegmentRepository,
	}
}

func (u *UserServiceImpl) GetUserById(id int) (*dto.UserDto, error) {

	userSegments, err := (*u.userSegmentRepository).GetUserById(id)
	if err != nil {
		return nil, err
	}

	userDto := dto.NewUserDto(id, userSegments)
	return userDto, nil
}

func (u *UserServiceImpl) GetUsers() ([]*dto.UserDto, error) {

	usersSegments, err := (*u.userSegmentRepository).GetAllUsers()
	if err != nil {
		return nil, err
	}

	return getUserSegmentsDtos(usersSegments)
}

func (u *UserServiceImpl) GetUserActiveSegments(id int) (*dto.UserDto, error) {

	userSegments, err := (*u.userSegmentRepository).GetUserActiveSegments(id)
	if err != nil {
		return nil, err
	}

	userDto := dto.NewUserDto(id, userSegments)
	return userDto, err
}

func (u *UserServiceImpl) CreateUser(id int) (*dto.UserDto, error) {

	user, err := (*u.userRepository).CreateUser(id)
	if err != nil {
		return nil, err
	}

	return dto.NewUserDto(user.Id, []*model.UserSegment{}), nil
}

func (u *UserServiceImpl) AddUserToSegment(userSingleSegmentDto *dto.UserSingleSegmentDto) (*dto.UserDto, error) {

	userSegments, err := (*u.userSegmentRepository).AddUserToSegment(userSingleSegmentDto.ToUserSegment())
	if err != nil {
		return nil, err
	}

	userDto := dto.NewUserDto(userSingleSegmentDto.UserId, userSegments)
	return userDto, nil
}

func (u *UserServiceImpl) AddUserToSegments(userDto *dto.UserDto) (*dto.UserDto, error) {

	userSegments, err := (*u.userSegmentRepository).AddUserToSegments(userDto.UserId, userDto.ToUserSegments())
	if err != nil {
		return nil, err
	}

	return dto.NewUserDto(userDto.UserId, userSegments), nil
}

func (u *UserServiceImpl) RemoveUserFromSegment(userSingleSegmentDto *dto.UserSingleSegmentDto) (*dto.UserDto, error) {

	userSegments, err := (*u.userSegmentRepository).RemoveUserFromSegment(userSingleSegmentDto.ToUserSegment())
	if err != nil {
		return nil, err
	}

	userDto := dto.NewUserDto(userSingleSegmentDto.UserId, userSegments)
	return userDto, nil
}

func (u *UserServiceImpl) RemoveUserFromSegments(userDto *dto.UserDto) (*dto.UserDto, error) {

	userSegments, err := (*u.userSegmentRepository).RemoveUserFromSegments(userDto.UserId, userDto.ToUserSegments())
	if err != nil {
		return nil, err
	}

	return dto.NewUserDto(userDto.UserId, userSegments), nil
}

func (u *UserServiceImpl) UpdateUserSegments(userDto *dto.UserDto) (*dto.UserDto, error) {

	userSegments, err := (*u.userSegmentRepository).UpdateUserSegments(userDto.UserId, userDto.ToUserSegments())
	if err != nil {
		return nil, err
	}

	return dto.NewUserDto(userDto.UserId, userSegments), nil
}

func (u *UserServiceImpl) DeleteUser(id int) error {

	if err := (*u.userRepository).DeleteUser(id); err != nil {
		return err
	}

	return nil
}

func (u *UserServiceImpl) GetUserSegmentsDataCsvUrl(userCsvDto *dto.UserSegmentsCsvDto) (string, error) {

	csvUrl, err := (*u.userSegmentRepository).GetUserSegmentsDataCsv(userCsvDto.UserId,
		userCsvDto.FromTime, userCsvDto.ToTime)
	if err != nil {
		return "", err
	}

	return csvUrl, nil
}

func getUserSegmentsDtos(usersSegments []*model.UserSegment) ([]*dto.UserDto, error) {
	var usersMap = make(map[int][]*model.UserSegment)
	for _, userSegment := range usersSegments {
		if _, elemExist := usersMap[userSegment.UserId]; !elemExist {
			usersMap[userSegment.UserId] = []*model.UserSegment{userSegment}
			continue
		}

		usersMap[userSegment.UserId] = append(usersMap[userSegment.UserId], userSegment)
	}

	var userSegmentDtos []*dto.UserDto
	for userId, userSegment := range usersMap {
		userSegmentDtos = append(userSegmentDtos, dto.NewUserDto(userId, userSegment))
	}

	return userSegmentDtos, nil
}
