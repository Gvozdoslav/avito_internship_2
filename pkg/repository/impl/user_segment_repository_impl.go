package impl

import (
	"avito2/pkg/model"
	"avito2/pkg/repository/query"
	"database/sql"
	"encoding/csv"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go/types"
	"os"
	"time"
)

type UserSegmentRepositoryImpl struct {
	db      *sqlx.DB
	queries *query.UserSegmentQuery
}

func NewUserSegmentRepositoryImpl(db *sqlx.DB) *UserSegmentRepositoryImpl {
	return &UserSegmentRepositoryImpl{
		db:      db,
		queries: query.NewUserSegmentQuery(),
	}
}

func (u *UserSegmentRepositoryImpl) GetUserById(userId int) ([]*model.UserSegment, error) {

	var userSegments []*model.UserSegment

	if err := u.db.Select(&userSegments, u.queries.GetUserQuery, userId); err != nil {
		return nil, err
	}

	return userSegments, nil
}

func (u *UserSegmentRepositoryImpl) GetAllUsers() ([]*model.UserSegment, error) {

	var userSegments []*model.UserSegment

	if err := u.db.Select(&userSegments, u.queries.GetUsersQuery); err != nil {
		return nil, err
	}

	return userSegments, nil
}

func (u *UserSegmentRepositoryImpl) GetUserActiveSegments(userId int) ([]*model.UserSegment, error) {

	var userSegments []*model.UserSegment
	var now = time.Now()

	if err := u.db.Select(&userSegments, u.queries.GetUserActiveSegmentsQuery, userId, now); err != nil {
		return nil, err
	}

	return userSegments, nil
}

func (u *UserSegmentRepositoryImpl) AddUserToSegment(userSegment *model.UserSegment) ([]*model.UserSegment, error) {

	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}

	_, err = u.addUserSegment(userSegment, tx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	var createdUserSegment []*model.UserSegment
	if err := u.db.Select(&createdUserSegment, u.queries.GetUserQuery, userSegment.UserId); err != nil {
		return nil, err
	}

	return createdUserSegment, nil
}

func (u *UserSegmentRepositoryImpl) AddUserToSegments(userId int, userSegments []*model.UserSegment) ([]*model.UserSegment, error) {

	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}

	for _, userSegment := range userSegments {
		_, err = u.addUserSegment(userSegment, tx)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	var foundUserSegments []*model.UserSegment
	if err := u.db.Select(&foundUserSegments, u.queries.GetUserQuery, userId); err != nil {
		return nil, err
	}

	return foundUserSegments, nil
}

func (u *UserSegmentRepositoryImpl) RemoveUserFromSegment(userSegment *model.UserSegment) ([]*model.UserSegment, error) {

	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}

	if err = u.removeUserSegment(userSegment, tx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	var foundUserSegments []*model.UserSegment
	if err := u.db.Select(&foundUserSegments, u.queries.GetUserQuery, userSegment.UserId); err != nil {
		return nil, err
	}

	return foundUserSegments, nil
}

func (u *UserSegmentRepositoryImpl) RemoveUserFromSegments(userId int, userSegments []*model.UserSegment) ([]*model.UserSegment, error) {

	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}

	for _, userSegment := range userSegments {
		if err = u.removeUserSegment(userSegment, tx); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	var foundUserSegments []*model.UserSegment
	if err := u.db.Select(&foundUserSegments, u.queries.GetUserQuery, userId); err != nil {
		return nil, err
	}

	return foundUserSegments, nil
}

func (u *UserSegmentRepositoryImpl) UpdateUserSegments(userId int, userSegments []*model.UserSegment) ([]*model.UserSegment, error) {

	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(u.queries.RemoveAllUserSegmentsQuery, userId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, userSegment := range userSegments {
		_, err = u.addUserSegment(userSegment, tx)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	var foundUserSegments []*model.UserSegment
	if err := u.db.Select(&foundUserSegments, u.queries.GetUserQuery, userId); err != nil {
		return nil, err
	}

	return foundUserSegments, nil
}

func (u *UserSegmentRepositoryImpl) GetUserSegmentsDataCsv(userId int) (string, error) {

	rows, err := u.db.Query(u.queries.GetUserQuery, userId)
	if err != nil {
		return "", err
	}

	now := time.Now()
	fileName := fmt.Sprintf("user_%d_report_%s.csv", userId, now.Format("2006-01-02"))
	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()

	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}

	csvWriter.Write(columns)

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			return "", err
		}
		var row []string
		for _, value := range values {
			if value != nil {
				row = append(row, fmt.Sprintf("%v", value))
			} else {
				row = append(row, "")
			}
		}
		csvWriter.Write(row)
	}

	return file.Name(), nil
}

func (u *UserSegmentRepositoryImpl) addUserSegment(userSegment *model.UserSegment, tx *sql.Tx) (int, error) {

	if u.isUserSegmentExist(userSegment) {
		return -1, types.Error{Msg: "UserSegment already exists!"}
	}

	var id int
	if err := tx.QueryRow(u.queries.AddUserSegmentQuery, userSegment.UserId,
		userSegment.SegmentSlug, userSegment.AddTime, userSegment.ExpireTime).Scan(&id); err != nil {
		tx.Rollback()
		return -1, err
	}

	return id, nil
}

func (u *UserSegmentRepositoryImpl) removeUserSegment(userSegment *model.UserSegment, tx *sql.Tx) error {

	_, err := tx.Exec(u.queries.RemoveUserSegmentQuery, userSegment.UserId, userSegment.SegmentSlug)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (u *UserSegmentRepositoryImpl) isUserSegmentExist(userSegment *model.UserSegment) bool {

	var foundUserSegment model.UserSegment
	if err := u.db.Get(&foundUserSegment, u.queries.GetUserSegmentByUserIdAndSlug,
		userSegment.UserId, userSegment.SegmentSlug); err != nil {
		return false
	}

	return true
}
