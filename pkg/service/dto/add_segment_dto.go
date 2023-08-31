package dto

type AddSegmentDto struct {
	Slug    string `json:"slug"`
	Percent *int   `json:"percent"`
}

func NewAddSegmentDto(slug string, percent *int) *AddSegmentDto {
	return &AddSegmentDto{
		Slug:    slug,
		Percent: percent,
	}
}
