package impl

import (
	"avito2/pkg/model"
	"avito2/pkg/repository/query"
	"github.com/jmoiron/sqlx"
)

type SegmentRepositoryImpl struct {
	db      *sqlx.DB
	queries *query.SegmentQuery
}

func NewSegmentRepositoryImpl(db *sqlx.DB) *SegmentRepositoryImpl {
	return &SegmentRepositoryImpl{
		db:      db,
		queries: query.NewSegmentQuery(),
	}
}

func (s *SegmentRepositoryImpl) GetSegmentBySlug(slug string) (*model.Segment, error) {

	var segment model.Segment
	if err := s.db.Get(&segment, s.queries.GetSegmentBySlugQuery, slug); err != nil {
		return nil, err
	}

	return &segment, nil
}

func (s *SegmentRepositoryImpl) GetAllSegments() ([]*model.Segment, error) {

	var segments []*model.Segment
	if err := s.db.Select(&segments, s.queries.GetAllSegmentsQuery); err != nil {
		return nil, err
	}

	return segments, nil
}

func (s *SegmentRepositoryImpl) CreateSegment(slug string) (*model.Segment, error) {

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	var createdSlug string
	if err = tx.QueryRow(s.queries.CreateSegmentQuery, slug).Scan(&createdSlug); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	var createdSegment model.Segment
	if err = s.db.Get(&createdSegment, s.queries.GetSegmentBySlugQuery, slug); err != nil {
		tx.Rollback()
		return nil, err
	}

	return &createdSegment, nil
}

func (s *SegmentRepositoryImpl) DeleteSegment(slug string) error {

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(s.queries.DeleteUserSegmentsBySegment, slug)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(s.queries.DeleteSegmentQuery, slug)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
