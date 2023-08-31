package impl

import (
	"avito2/pkg/model"
	"avito2/pkg/repository/query"
	"github.com/jmoiron/sqlx"
)

type UserRepositoryImpl struct {
	db      *sqlx.DB
	queries *query.UserQuery
}

func NewUserRepositoryImpl(db *sqlx.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		db:      db,
		queries: query.NewUserQueries(),
	}
}

func (u *UserRepositoryImpl) CreateUser(id int) (*model.User, error) {

	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}

	var createdUserId int
	if err = tx.QueryRow(u.queries.CreateUserQuery, id).Scan(&createdUserId); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	var createdUser model.User
	if err := u.db.Get(&createdUser, u.queries.GetUserQuery, createdUserId); err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func (u *UserRepositoryImpl) DeleteUser(id int) error {

	tx, err := u.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(u.queries.DeleteSegmentsByUserQuery, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(u.queries.DeleteUserQuery, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
