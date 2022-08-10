package repo

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"regexp"
	"testing"
	"time"
	"win-a-thon/database"
	"win-a-thon/models"
)

func TestCountParticipants(t *testing.T) {
	driver, mock, err := sqlmock.New()
	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT count(*) FROM `participants` WHERE hackathon_id = ?")).WithArgs("1").WillReturnRows(sqlmock.NewRows(
		[]string{"count(*)"}).AddRow(1))

	count, err := CountParticipants("1")
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Error("Wrong result")
	}
}

func TestIsParticipant(t *testing.T) {
	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      driver,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	participantMockRows := sqlmock.NewRows([]string{"hackathon_id", "user_id", "demo_url", "code_url", "score"}).AddRow(1, "1", "", "", "70")
	t.Run("no errors", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `participants` WHERE hackathon_id = ? AND user_id = ? ORDER BY `participants`.`hackathon_id` LIMIT 1")).WithArgs("1", 1).WillReturnRows(participantMockRows)
		_, err = IsParticipant(1, "1")
		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("Failed to meet expectations, got error: %v", err)
		}
	})
	t.Run("errors", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `participants` WHERE hackathon_id = ? AND user_id = ? ORDER BY `participants`.`hackathon_id` LIMIT 1")).WithArgs("1", 2).WillReturnError(gorm.ErrRecordNotFound)
		_, err = IsParticipant(2, "1")
		err = mock.ExpectationsWereMet()
	})
}

func TestIsOrganiser(t *testing.T) {
	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      driver,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	var created_at = time.Now()

	hackathonMockRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_time", "ending_time", "result_time", "organisation_name", "organiser_id", "description", "admin_approved"}).AddRow(1, created_at, created_at, sql.NullString{String: "", Valid: false}, "Test Hackathon", created_at, created_at, created_at, "Winathon", "1", "It is a test Hackathon", "1")

	t.Run("No error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `hackathons` WHERE id = ? AND `hackathons`.`deleted_at` IS NULL")).WithArgs("1").WillReturnRows(hackathonMockRows)
		var flag bool
		flag, err = IsOrganiser(1, "1")
		assert.Equal(t, true, flag)
		err = mock.ExpectationsWereMet()
	})
	t.Run("Error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `hackathons` WHERE id = ? AND `hackathons`.`deleted_at` IS NULL")).WithArgs("1").WillReturnRows(hackathonMockRows)
		var flag bool
		flag, err = IsOrganiser(2, "1")
		assert.Equal(t, false, flag)
		err = mock.ExpectationsWereMet()
	})
}

func TestGetSubmission(t *testing.T) {
	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      driver,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	var p models.Participant
	participantMockRows := sqlmock.NewRows([]string{"hackathon_id", "user_id", "demo_url", "code_url", "score"}).AddRow(1, "1", "", "", "70")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `participants` WHERE hackathon_id = ? AND user_id = ? ORDER BY `participants`.`hackathon_id` LIMIT 1")).WithArgs(1, 1).WillReturnRows(participantMockRows)
	err = GetSubmission(&p, 1, 1)
	if err != nil {
		log.Fatal(err)
	}
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}
