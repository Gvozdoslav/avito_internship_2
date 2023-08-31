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
	"path/filepath"
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

func (u *UserSegmentRepositoryImpl) RemoveExpiredUserSegments() error {

	tx, err := u.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(u.queries.SetExpiredUserSegmentsQuery, model.Expired, time.Now())
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
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

func (u *UserSegmentRepositoryImpl) GetUserSegmentsDataCsv(userId int, fromTime *time.Time, toTime *time.Time) (string, error) {

	rows, err := u.db.Query(u.queries.GetUserSegmentsByPeriod, userId, fromTime, toTime)
	if err != nil {
		return "", err
	}

	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	fileName := filepath.Join(dir, fmt.Sprintf("user_%d_report_%s__%s.csv",
		userId, fromTime.Format("2006-01-02"), toTime.Format("2006-01-02")))
	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()

	headers := []string{"user_id", "segment_slug", "operation", "time"}
	err = csvWriter.Write(headers)
	if err != nil {
		fmt.Println("Error writing CSV header:", err)
		return "", err
	}

	columns, err := rows.Columns()
	if err != nil {
		fmt.Println("Error writing CSV header:", err)
		return "", err
	}

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			return "", err
		}

		if values[3] != nil {
			var row []string
			row = append(row, fmt.Sprintf("%v", values[1]))
			row = append(row, fmt.Sprintf("%v", values[2]))
			row = append(row, "added")
			row = append(row, fmt.Sprintf("%v", values[3]))
			csvWriter.Write(row)
		}

		if values[4] != nil {
			var row []string
			row = append(row, fmt.Sprintf("%v", values[1]))
			row = append(row, fmt.Sprintf("%v", values[2]))
			row = append(row, "left")
			row = append(row, fmt.Sprintf("%v", values[4]))
			csvWriter.Write(row)
		}
	}

	return file.Name(), nil
}

func (u *UserSegmentRepositoryImpl) addUserSegment(userSegment *model.UserSegment, tx *sql.Tx) (int, error) {

	if u.isActiveUserSegmentExist(userSegment) {
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

func (u *UserSegmentRepositoryImpl) isActiveUserSegmentExist(userSegment *model.UserSegment) bool {

	var foundUserSegment model.UserSegment
	if err := u.db.Get(&foundUserSegment, u.queries.GetActiveUserSegmentByUserIdAndSlug,
		userSegment.UserId, userSegment.SegmentSlug, model.Active); err != nil {
		return false
	}

	return true
}
