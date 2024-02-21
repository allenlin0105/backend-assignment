package domain

import "context"

type Ad struct {
	Title     string     `json:"title" binding:"required"`
	StartAt   string     `json:"startAt,omitempty" binding:"required"`
	EndAt     string     `json:"endAt" binding:"required"`
	Condition *Condition `json:"condition,omitempty"`
}

type Condition struct {
	AgeStart int      `json:"ageStart"`
	AgeEnd   int      `json:"ageEnd"`
	Gender   []string `json:"gender"`
	Country  []string `json:"country"`
	Platform []string `json:"platform"`
}

type AdRepository interface {
	Create(c context.Context, ad *Ad) error
	GetByCondition(c context.Context, condition map[string][]string) ([]Ad, error)
}

type AdUsecase interface {
	Create(c context.Context, ad *Ad) error
	GetByCondition(c context.Context, condition map[string][]string) ([]Ad, error)
}
