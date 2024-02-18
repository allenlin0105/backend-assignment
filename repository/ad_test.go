package repository_test

import (
	"context"
	"dcard-backend/domain"
	"dcard-backend/repository"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

const (
	query_ads         = "INSERT INTO ads (title, start_at, end_at, age_start, age_end) VALUES (?, ?, ?, ?, ?)"
	query_ad_gender   = "INSERT INTO ad_gender (ad_id, gender_id) VALUES (?, (SELECT id FROM genders WHERE gender = ?))"
	query_ad_country  = "INSERT INTO ad_country (ad_id, country_id) VALUES (?, (SELECT id FROM countries WHERE country = ?))"
	query_ad_platform = "INSERT INTO ad_platform (ad_id, platform_id) VALUES (?, (SELECT id FROM platforms WHERE platform = ?))"
)

var mockAd = domain.Ad{
	Title:   "AD 0",
	StartAt: "2024-01-01 00:00:00",
	EndAt:   "2025-01-01 00:00:00",
	Condition: &domain.Condition{
		AgeStart: 10,
		AgeEnd:   20,
		Gender:   []string{"M", "F"},
		Country:  []string{"TW", "JP"},
		Platform: []string{"web", "ios"},
	},
}

func TestCreate_SuccessInsertOnAllTables_NoError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()

	prep := mock.ExpectPrepare(query_ads)
	prep.ExpectExec().
		WithArgs(mockAd.Title, mockAd.StartAt, mockAd.EndAt, mockAd.Condition.AgeStart, mockAd.Condition.AgeEnd).
		WillReturnResult(sqlmock.NewResult(1, 1))

	for i, gender := range mockAd.Condition.Gender {
		prep = mock.ExpectPrepare(query_ad_gender)
		prep.ExpectExec().
			WithArgs(1, gender).
			WillReturnResult(sqlmock.NewResult(int64(i+1), 1))
	}

	for i, country := range mockAd.Condition.Country {
		prep = mock.ExpectPrepare(query_ad_country)
		prep.ExpectExec().
			WithArgs(1, country).
			WillReturnResult(sqlmock.NewResult(int64(i+1), 1))
	}

	for i, platform := range mockAd.Condition.Platform {
		prep = mock.ExpectPrepare(query_ad_platform)
		prep.ExpectExec().
			WithArgs(1, platform).
			WillReturnResult(sqlmock.NewResult(int64(i+1), 1))
	}

	mock.ExpectCommit()

	testAr := repository.NewAdRepository(db)
	err = testAr.Create(context.Background(), &mockAd)
	assert.NoError(t, err, "Create function should return with no error")
}

func TestCreate_FailOnFirstInsert_ShouldRollbackOnError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	prep := mock.ExpectPrepare(query_ads)
	prep.ExpectExec().
		WithArgs(mockAd.Title, mockAd.StartAt, mockAd.EndAt, mockAd.Condition.AgeStart, mockAd.Condition.AgeEnd).
		WillReturnError(fmt.Errorf("Error"))
	mock.ExpectRollback()

	testAr := repository.NewAdRepository(db)
	err = testAr.Create(context.Background(), &mockAd)
	assert.Error(t, err, "If inserting ads fail, it should return error")
}

func TestCreate_FailOnSecondInsert_ShouldRollbackOnError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()

	prep := mock.ExpectPrepare(query_ads)
	prep.ExpectExec().
		WithArgs(mockAd.Title, mockAd.StartAt, mockAd.EndAt, mockAd.Condition.AgeStart, mockAd.Condition.AgeEnd).
		WillReturnResult(sqlmock.NewResult(1, 1))

	prep = mock.ExpectPrepare(query_ad_gender)
	prep.ExpectExec().
		WithArgs(1, mockAd.Condition.Gender[0]).
		WillReturnError(fmt.Errorf("Error"))

	mock.ExpectRollback()

	testAr := repository.NewAdRepository(db)
	err = testAr.Create(context.Background(), &mockAd)
	assert.Error(t, err, "If inserting ad_gender fail, it should return error")
}

func TestCreate_FailOnThirdInsert_ShouldRollbackOnError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()

	prep := mock.ExpectPrepare(query_ads)
	prep.ExpectExec().
		WithArgs(mockAd.Title, mockAd.StartAt, mockAd.EndAt, mockAd.Condition.AgeStart, mockAd.Condition.AgeEnd).
		WillReturnResult(sqlmock.NewResult(1, 1))

	for i, gender := range mockAd.Condition.Gender {
		prep = mock.ExpectPrepare(query_ad_gender)
		prep.ExpectExec().WithArgs(1, gender).WillReturnResult(sqlmock.NewResult(int64(i+1), 1))
	}

	prep = mock.ExpectPrepare(query_ad_country)
	prep.ExpectExec().
		WithArgs(1, mockAd.Condition.Country[0]).
		WillReturnError(fmt.Errorf("Error"))

	mock.ExpectRollback()

	testAr := repository.NewAdRepository(db)
	err = testAr.Create(context.Background(), &mockAd)
	assert.Error(t, err, "If inserting ad_country fail, it should return error")
}

func TestCreate_FailOnFourthInsert_ShouldRollbackOnError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()

	prep := mock.ExpectPrepare(query_ads)
	prep.ExpectExec().
		WithArgs(mockAd.Title, mockAd.StartAt, mockAd.EndAt, mockAd.Condition.AgeStart, mockAd.Condition.AgeEnd).
		WillReturnResult(sqlmock.NewResult(1, 1))

	for i, gender := range mockAd.Condition.Gender {
		prep = mock.ExpectPrepare(query_ad_gender)
		prep.ExpectExec().
			WithArgs(1, gender).
			WillReturnResult(sqlmock.NewResult(int64(i+1), 1))
	}

	for i, country := range mockAd.Condition.Country {
		prep = mock.ExpectPrepare(query_ad_country)
		prep.ExpectExec().
			WithArgs(1, country).
			WillReturnResult(sqlmock.NewResult(int64(i+1), 1))
	}

	prep = mock.ExpectPrepare(query_ad_platform)
	prep.ExpectExec().
		WithArgs(1, mockAd.Condition.Platform[0]).
		WillReturnError(fmt.Errorf("Error"))

	mock.ExpectRollback()

	testAr := repository.NewAdRepository(db)
	err = testAr.Create(context.Background(), &mockAd)
	assert.Error(t, err, "If inserting ad_platform fail, it should return error")
}

func TestGetByCondition_SuccessWithAllConditionsProvided_AdsReturnWithNoError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	query := "SELECT ads.title, ads.end_at FROM ads "
	query += "INNER JOIN ad_gender ON ads.id = ad_gender.ad_id INNER JOIN genders ON genders.id = ad_gender.gender_id "
	query += "INNER JOIN ad_country ON ads.id = ad_country.ad_id INNER JOIN countries ON countries.id = ad_country.country_id "
	query += "INNER JOIN ad_platform ON ads.id = ad_platform.ad_id INNER JOIN platforms ON platforms.id = ad_platform.platform_id "
	query += "WHERE genders.gender IN (?,?) AND countries.country IN (?,?) AND platforms.platform IN (?,?) AND "
	query += "ads.age_start <= ? AND ads.age_end >= ? AND ads.start_at <= NOW() AND ads.end_at >= NOW() "
	query += "ORDER BY ads.end_at ASC LIMIT ? OFFSET ?"

	mockRows := sqlmock.NewRows([]string{"title", "end_at"}).AddRow(mockAd.Title, mockAd.EndAt)

	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().
		WithArgs("M", "A", "TW", "AY", "web", "any", "14", "14", "10", "0").
		WillReturnRows(mockRows)

	testAr := repository.NewAdRepository(db)
	condition := map[string][]string{
		"gender":   {"M", "A"},
		"country":  {"TW", "AY"},
		"platform": {"web", "any"},
		"age":      {"14"},
		"limit":    {"10"},
		"offset":   {"0"},
	}
	ads, err := testAr.GetByCondition(context.Background(), condition)

	assert.NoError(t, err)
	if assert.NotNil(t, ads) {
		assert.Equal(t, len(ads), 1)
		assert.Equal(t, ads[0].Title, mockAd.Title)
		assert.Equal(t, ads[0].EndAt, mockAd.EndAt)
	}
}

func TestGetByCondition_SuccessWithNoAgeCondition_AdsReturnWithNoError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	query := "SELECT ads.title, ads.end_at FROM ads "
	query += "INNER JOIN ad_gender ON ads.id = ad_gender.ad_id INNER JOIN genders ON genders.id = ad_gender.gender_id "
	query += "INNER JOIN ad_country ON ads.id = ad_country.ad_id INNER JOIN countries ON countries.id = ad_country.country_id "
	query += "INNER JOIN ad_platform ON ads.id = ad_platform.ad_id INNER JOIN platforms ON platforms.id = ad_platform.platform_id "
	query += "WHERE genders.gender IN (?,?) AND countries.country IN (?,?) AND platforms.platform IN (?,?) AND "
	query += "ads.start_at <= NOW() AND ads.end_at >= NOW() "
	query += "ORDER BY ads.end_at ASC LIMIT ? OFFSET ?"

	mockRows := sqlmock.NewRows([]string{"title", "end_at"}).AddRow(mockAd.Title, mockAd.EndAt)

	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().
		WithArgs("M", "A", "TW", "AY", "web", "any", "10", "0").
		WillReturnRows(mockRows)

	testAr := repository.NewAdRepository(db)
	condition := map[string][]string{
		"gender":   {"M", "A"},
		"country":  {"TW", "AY"},
		"platform": {"web", "any"},
		"limit":    {"10"},
		"offset":   {"0"},
	}
	ads, err := testAr.GetByCondition(context.Background(), condition)

	assert.NoError(t, err)
	if assert.NotNil(t, ads) {
		assert.Equal(t, len(ads), 1)
		assert.Equal(t, ads[0].Title, mockAd.Title)
		assert.Equal(t, ads[0].EndAt, mockAd.EndAt)
	}
}

func TestGetByCondition_SuccessWithNoGenderCondition_AdsReturnWithNoError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	query := "SELECT ads.title, ads.end_at FROM ads "
	query += "INNER JOIN ad_country ON ads.id = ad_country.ad_id INNER JOIN countries ON countries.id = ad_country.country_id "
	query += "INNER JOIN ad_platform ON ads.id = ad_platform.ad_id INNER JOIN platforms ON platforms.id = ad_platform.platform_id "
	query += "WHERE countries.country IN (?,?) AND platforms.platform IN (?,?) AND "
	query += "ads.age_start <= ? AND ads.age_end >= ? AND ads.start_at <= NOW() AND ads.end_at >= NOW() "
	query += "ORDER BY ads.end_at ASC LIMIT ? OFFSET ?"

	mockRows := sqlmock.NewRows([]string{"title", "end_at"}).AddRow(mockAd.Title, mockAd.EndAt)

	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().
		WithArgs("TW", "AY", "web", "any", "14", "14", "10", "0").
		WillReturnRows(mockRows)

	testAr := repository.NewAdRepository(db)
	condition := map[string][]string{
		"country":  {"TW", "AY"},
		"platform": {"web", "any"},
		"age":      {"14"},
		"limit":    {"10"},
		"offset":   {"0"},
	}
	ads, err := testAr.GetByCondition(context.Background(), condition)

	assert.NoError(t, err)
	if assert.NotNil(t, ads) {
		assert.Equal(t, len(ads), 1)
		assert.Equal(t, ads[0].Title, mockAd.Title)
		assert.Equal(t, ads[0].EndAt, mockAd.EndAt)
	}
}

func TestGetByCondition_SuccessWithNoCountryCondition_AdsReturnWithNoError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	query := "SELECT ads.title, ads.end_at FROM ads "
	query += "INNER JOIN ad_gender ON ads.id = ad_gender.ad_id INNER JOIN genders ON genders.id = ad_gender.gender_id "
	query += "INNER JOIN ad_platform ON ads.id = ad_platform.ad_id INNER JOIN platforms ON platforms.id = ad_platform.platform_id "
	query += "WHERE genders.gender IN (?,?) AND platforms.platform IN (?,?) AND "
	query += "ads.age_start <= ? AND ads.age_end >= ? AND ads.start_at <= NOW() AND ads.end_at >= NOW() "
	query += "ORDER BY ads.end_at ASC LIMIT ? OFFSET ?"

	mockRows := sqlmock.NewRows([]string{"title", "end_at"}).AddRow(mockAd.Title, mockAd.EndAt)

	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().
		WithArgs("M", "A", "web", "any", "14", "14", "10", "0").
		WillReturnRows(mockRows)

	testAr := repository.NewAdRepository(db)
	condition := map[string][]string{
		"gender":   {"M", "A"},
		"platform": {"web", "any"},
		"age":      {"14"},
		"limit":    {"10"},
		"offset":   {"0"},
	}
	ads, err := testAr.GetByCondition(context.Background(), condition)

	assert.NoError(t, err)
	if assert.NotNil(t, ads) {
		assert.Equal(t, len(ads), 1)
		assert.Equal(t, ads[0].Title, mockAd.Title)
		assert.Equal(t, ads[0].EndAt, mockAd.EndAt)
	}
}

func TestGetByCondition_SuccessWithNoPlatformCondition_AdsReturnWithNoError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	query := "SELECT ads.title, ads.end_at FROM ads "
	query += "INNER JOIN ad_gender ON ads.id = ad_gender.ad_id INNER JOIN genders ON genders.id = ad_gender.gender_id "
	query += "INNER JOIN ad_country ON ads.id = ad_country.ad_id INNER JOIN countries ON countries.id = ad_country.country_id "
	query += "WHERE genders.gender IN (?,?) AND countries.country IN (?,?) AND "
	query += "ads.age_start <= ? AND ads.age_end >= ? AND ads.start_at <= NOW() AND ads.end_at >= NOW() "
	query += "ORDER BY ads.end_at ASC LIMIT ? OFFSET ?"

	mockRows := sqlmock.NewRows([]string{"title", "end_at"}).AddRow(mockAd.Title, mockAd.EndAt)

	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().
		WithArgs("M", "A", "TW", "AY", "14", "14", "10", "0").
		WillReturnRows(mockRows)

	testAr := repository.NewAdRepository(db)
	condition := map[string][]string{
		"gender":  {"M", "A"},
		"country": {"TW", "AY"},
		"age":     {"14"},
		"limit":   {"10"},
		"offset":  {"0"},
	}
	ads, err := testAr.GetByCondition(context.Background(), condition)

	assert.NoError(t, err)
	if assert.NotNil(t, ads) {
		assert.Equal(t, len(ads), 1)
		assert.Equal(t, ads[0].Title, mockAd.Title)
		assert.Equal(t, ads[0].EndAt, mockAd.EndAt)
	}
}
