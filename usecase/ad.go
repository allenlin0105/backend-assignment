package usecase

import (
	"context"
	"dcard-backend/domain"
	"fmt"
	"time"
)

var conditionToAnyValue = map[string]string{
	"gender":   "A",
	"country":  "AY",
	"platform": "any",
}

type adUsecase struct {
	adRepository   domain.AdRepository
	contextTimeout time.Duration
}

func NewAdUsecase(adRepository domain.AdRepository, timeout time.Duration) domain.AdUsecase {
	return &adUsecase{
		adRepository:   adRepository,
		contextTimeout: timeout,
	}
}

func changeTimeToUTF8(timeStr *string) error {
	t, err := time.Parse(time.RFC3339, *timeStr)
	if err != nil {
		return err
	}

	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		return err
	}
	*timeStr = t.In(loc).Format(time.RFC3339)
	return nil
}

func changeAgeIfZero(age *int, defaultAge int) {
	if *age == 0 {
		*age = defaultAge
	}
}

func changeSliceIfEmpty(inputSlice *[]string, anyValueStr string) {
	if len(*inputSlice) == 0 {
		*inputSlice = append(*inputSlice, anyValueStr)
	}
}

func (au *adUsecase) Create(c context.Context, ad *domain.Ad) error {
	ctx, cancel := context.WithTimeout(c, au.contextTimeout)
	defer cancel()

	if ad.Title == "" || ad.StartAt == "" || ad.EndAt == "" {
		return fmt.Errorf("title, startAt, and endAt should not be empty")
	}

	if ad.Condition == nil {
		ad.Condition = &domain.Condition{}
	}

	err := changeTimeToUTF8(&ad.StartAt)
	if err != nil {
		return err
	}

	err = changeTimeToUTF8(&ad.EndAt)
	if err != nil {
		return err
	}

	changeAgeIfZero(&ad.Condition.AgeStart, 1)
	changeAgeIfZero(&ad.Condition.AgeEnd, 100)

	changeSliceIfEmpty(&ad.Condition.Gender, conditionToAnyValue["gender"])
	changeSliceIfEmpty(&ad.Condition.Country, conditionToAnyValue["country"])
	changeSliceIfEmpty(&ad.Condition.Platform, conditionToAnyValue["platform"])

	err = au.adRepository.Create(ctx, ad)
	return err
}

func (au *adUsecase) GetByCondition(c context.Context, condition map[string][]string) ([]domain.Ad, error) {
	ctx, cancel := context.WithTimeout(c, au.contextTimeout)
	defer cancel()

	if _, ok := condition["offset"]; !ok {
		return nil, fmt.Errorf("condition should have offset provided")
	}

	if _, ok := condition["limit"]; !ok {
		condition["limit"] = []string{"5"}
	}

	for key, anyValue := range conditionToAnyValue {
		if _, ok := condition[key]; ok {
			condition[key] = append(condition[key], anyValue)
		}
	}

	ads, err := au.adRepository.GetByCondition(ctx, condition)
	if err != nil {
		return nil, err
	}

	for i, ad := range ads {
		loc, err := time.LoadLocation("Asia/Taipei")
		if err != nil {
			return nil, err
		}

		endAtTime, err := time.ParseInLocation("2006-01-02 15:04:05", ad.EndAt, loc)
		if err != nil {
			return nil, err
		}

		ads[i].EndAt = endAtTime.UTC().Format(time.RFC3339)
	}
	return ads, nil
}
