package repo

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"regexp"
	"testing"
	"time"
	"win-a-thon/database"
	"win-a-thon/models"
)

var mock sqlmock.Sqlmock
var err error
var db *sql.DB

func TestListHackathons(t *testing.T) {
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	var created_at = time.Now()

	hackathonMockRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_time", "ending_time", "result_time", "organisation_name", "organiser_id", "description", "admin_approved"}).AddRow("1", created_at, created_at, created_at, "Test Hackathon", created_at, created_at, created_at, "Winathon", "1", "It is a test Hackathon", "1")

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE admin_approved = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY starting_time desc")).WithArgs(true).WillReturnRows(hackathonMockRows)

	var hackathons []models.Hackathon

	err := ListHackathons(&hackathons, true)
	if err != nil {
		log.Fatal(err)
	}
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestCreateHackathon(t *testing.T) {
	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	hackathon := models.Hackathon{
		Title:            "Test",
		OrganisationName: "Winathon",
		OrganiserID:      1,
		Description:      "Lord of the Mysteries",
		AdminApproved:    false,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `hackathons` (`created_at`,`updated_at`,`deleted_at`,`title`,`starting_time`,`ending_time`,`result_time`,`organisation_name`,`organiser_id`,`description`,`admin_approved`) VALUES (?,?,?,?,?,?,?,?,?,?,?)")).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = CreateHackathon(&hackathon)

	require.NoError(t, err)
}

func TestCreateParticipant(t *testing.T) {
	driver, mock, err := sqlmock.New()
	database.DB, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      driver,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	var participant models.Participant
	participant.UserId = 1
	participant.HackathonId = 1
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `participants` (`hackathon_id`,`user_id`,`demo_url`,`code_url`,`score`) VALUES (?,?,?,?,?)")).WithArgs(participant.HackathonId, participant.UserId, participant.DemoUrl, participant.CodeUrl, participant.Score).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	if err = CreateParticipant(&participant); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}
func TestViewHackathonDetails(t *testing.T) {
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}
	var created_at = time.Now()
	hackathonMockRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_time", "ending_time", "result_time", "organisation_name", "organiser_id", "description", "admin_approved"}).AddRow("1", created_at, created_at, created_at, "Test Hackathon", created_at, created_at, created_at, "Winathon", "1", "It is a test Hackathon", "1")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `hackathons` WHERE (id = ? AND admin_approved = ?)")).WithArgs("1", true).WillReturnRows(hackathonMockRows)
	var hackathons models.Hackathon

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `participants` WHERE hackathon_id = ?")).WithArgs("1").WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))
	_, err = CountParticipants("1")
	err, _ := ViewHackathonDetails(&hackathons, "1")

	if err != nil {
		log.Fatal(err)
	}
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestListOrganisedHackathons(t *testing.T) {
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	date := time.Date(2021, 8, 4, 0, 0, 0, 0, time.UTC)

	hackathonMockRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_" +
		"time", "ending_time", "result_time", "organisation_name", "organiser_id", "description", "admin_" +
		"approved"}).AddRow(1, date, date, nil, "Test Hackathon", date, date, date, "Winathon", 1, "It is a test"+
		" Hackathon", true)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE (organiser_id = ?) AND `hackathons`.`deleted_at` IS NULL ORDER BY starting_time desc")).
		WithArgs(1).
		WillReturnRows(hackathonMockRow)

	var hackathons []models.Hackathon
	err := ListOrganisedHackathons(&hackathons, 1)
	if err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}
