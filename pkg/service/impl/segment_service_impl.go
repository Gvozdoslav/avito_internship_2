package impl

import (
	"avito2/pkg/model"
	"avito2/pkg/repository"
	"avito2/pkg/service/dto"
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

func (s *SegmentServiceImpl) CreateSegment(segmentDto *dto.AddSegmentDto) (*dto.AddSegmentDto, error) {

	_, err := (*s.segmentRepository).CreateSegment(segmentDto.Slug, segmentDto.Percent)
	if err != nil {
		return nil, err
	}

	return segmentDto, err
}

func (s *SegmentServiceImpl) DeleteSegment(slug string) error {

	if err := (*s.segmentRepository).DeleteSegment(slug); err != nil {
		return err
	}

	return nil
}
