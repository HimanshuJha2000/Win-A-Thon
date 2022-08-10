package repo

import (
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"regexp"
	"testing"
	"time"
	"win-a-thon/database"
)

func TestAdminAppoved(t *testing.T) {
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	date := time.Date(2021, 8, 2, 0, 0, 0, 0, time.UTC)

	hackathonMockRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_" +
		"time", "ending_time", "result_time", "organisation_name", "organiser_id", "description", "admin_" +
		"approved"}).AddRow(1, date, date, nil, "Test Hackathon", date, date, date, "Winathon", 1, "It is a test"+
		" Hackathon", false)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE id = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs("1").WillReturnRows(hackathonMockRow)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `hackathons` SET `created_at`=?,`updated_at`=?,`deleted_at`=?,`title`=?,`starting_time`=?,`ending_time`=?,`result_time`=?,`organisation_name`=?,`organiser_id`=?,`description`=?,`admin_approved`=? WHERE `id` = ?")).WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := AdminAppoved("1", "1")
	if err != nil {
		log.Fatal(err)
	}
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}
