package impl

import (
	"avito2/pkg/model"
	"avito2/pkg/repository"
)

type SegmentServiceImpl struct {
	segmentRepository *repository.SegmentRepository
}

func NewSegmentServiceImpl(segmentRepository *repository.SegmentRepository) *SegmentServiceImpl {
	return &SegmentServiceImpl{
		segmentRepository: segmentRepository,
	}
}

func (s *SegmentServiceImpl) GetSegment(slug string) (*model.Segment, error) {

	segment, err := (*s.segmentRepository).GetSegmentBySlug(slug)
	if err != nil {
		return nil, err
	}

	return segment, nil
}

func (s *SegmentServiceImpl) GetAllSegments() ([]*model.Segment, error) {

	segments, err := (*s.segmentRepository).GetAllSegments()
	if err != nil {
		return nil, err
	}

	return segments, nil
}

func (s *SegmentServiceImpl) CreateSegment(slug string) (*model.Segment, error) {

	segment, err := (*s.segmentRepository).CreateSegment(slug)
	if err != nil {
		return nil, err
	}

	return segment, err
}

func (s *SegmentServiceImpl) DeleteSegment(slug string) error {

	if err := (*s.segmentRepository).DeleteSegment(slug); err != nil {
		return err
	}

	return nil
}
