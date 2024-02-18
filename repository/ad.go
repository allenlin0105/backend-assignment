package repository

import (
	"context"
	"database/sql"
	"dcard-backend/domain"
	"fmt"
	"strings"
)

type adRepository struct {
	database *sql.DB
}

func NewAdRepository(db *sql.DB) domain.AdRepository {
	return &adRepository{
		database: db,
	}
}

func prepareAndCreate(c context.Context, tx *sql.Tx, command string, args ...interface{}) (sql.Result, error) {
	stmt, err := tx.Prepare(command)
	if err != nil {
		return nil, err
	}

	result, err := stmt.Exec(args...)
	return result, err
}

func (ar *adRepository) Create(c context.Context, ad *domain.Ad) error {
	tx, err := ar.database.Begin()
	if err != nil {
		return err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	command := "INSERT INTO ads (title, start_at, end_at, age_start, age_end) VALUES (?, ?, ?, ?, ?)"
	result, err := prepareAndCreate(c, tx, command, ad.Title, ad.StartAt, ad.EndAt, ad.Condition.AgeStart, ad.Condition.AgeEnd)
	if err != nil {
		fmt.Println("Error created when inserting into ads:", err.Error())
		return err
	}

	adId, err := result.LastInsertId()
	if err != nil {
		return err
	}

	command = "INSERT INTO ad_gender (ad_id, gender_id) VALUES (?, (SELECT id FROM genders WHERE gender = ?))"
	for _, gender := range ad.Condition.Gender {
		if _, err = prepareAndCreate(c, tx, command, adId, gender); err != nil {
			fmt.Println("Error created when inserting into ad_gender:", err.Error())
			return err
		}
	}

	command = "INSERT INTO ad_country (ad_id, country_id) VALUES (?, (SELECT id FROM countries WHERE country = ?))"
	for _, country := range ad.Condition.Country {
		if _, err = prepareAndCreate(c, tx, command, adId, country); err != nil {
			fmt.Println("Error created when inserting into ad_country:", err.Error())
			return err
		}
	}

	command = "INSERT INTO ad_platform (ad_id, platform_id) VALUES (?, (SELECT id FROM platforms WHERE platform = ?))"
	for _, platform := range ad.Condition.Platform {
		if _, err = prepareAndCreate(c, tx, command, adId, platform); err != nil {
			fmt.Println("Error created when inserting into ad_platform:", err.Error())
			return err
		}
	}

	return err
}

func repeatQuestionMarks(length int) string {
	return "?" + strings.Repeat(",?", length-1)
}

func stringSliceToGenericSlice(slice []string) []interface{} {
	genericSlice := make([]interface{}, len(slice))
	for i, v := range slice {
		genericSlice[i] = v
	}
	return genericSlice
}

func (ar *adRepository) GetByCondition(c context.Context, condition map[string][]string) ([]domain.Ad, error) {
	var args []string
	innerJoinCommands, whereCommands := []string{}, []string{}

	// Gender condition
	if values, ok := condition["gender"]; ok {
		innerJoinCommands = append(innerJoinCommands,
			"INNER JOIN ad_gender ON ads.id = ad_gender.ad_id INNER JOIN genders ON genders.id = ad_gender.gender_id")
		whereCommands = append(whereCommands, "genders.gender IN ("+repeatQuestionMarks(len(values))+")")
		args = append(args, values...)
	}

	// Country condition
	if values, ok := condition["country"]; ok {
		innerJoinCommands = append(innerJoinCommands,
			"INNER JOIN ad_country ON ads.id = ad_country.ad_id INNER JOIN countries ON countries.id = ad_country.country_id")
		whereCommands = append(whereCommands, "countries.country IN ("+repeatQuestionMarks(len(values))+")")
		args = append(args, values...)
	}

	// Platform condition
	if values, ok := condition["platform"]; ok {
		innerJoinCommands = append(innerJoinCommands,
			"INNER JOIN ad_platform ON ads.id = ad_platform.ad_id INNER JOIN platforms ON platforms.id = ad_platform.platform_id")
		whereCommands = append(whereCommands, "platforms.platform IN ("+repeatQuestionMarks(len(values))+")")
		args = append(args, values...)
	}

	// Age condition
	if values, ok := condition["age"]; ok {
		whereCommands = append(whereCommands, "ads.age_start <= ? AND ads.age_end >= ?")
		args = append(args, values[0], values[0])
	}

	// Time condition
	whereCommands = append(whereCommands, "ads.start_at <= NOW() AND ads.end_at >= NOW()")

	// Set limit and offset
	args = append(args, condition["limit"][0], condition["offset"][0])

	command := "SELECT ads.title, ads.end_at FROM ads "
	command += strings.Join(innerJoinCommands, " ") + " "
	command += "WHERE " + strings.Join(whereCommands, " AND ") + " "
	command += "ORDER BY ads.end_at ASC LIMIT ? OFFSET ?"

	stmt, err := ar.database.PrepareContext(c, command)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(stringSliceToGenericSlice(args)...)
	if err != nil {
		return nil, err
	}

	var ads []domain.Ad
	for rows.Next() {
		var ad domain.Ad
		if err := rows.Scan(&ad.Title, &ad.EndAt); err != nil {
			return nil, err
		}
		ads = append(ads, ad)
	}
	err = rows.Err()
	fmt.Println(err)

	return ads, err
}
