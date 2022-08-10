package repo

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"regexp"
	"testing"
	"time"
	"win-a-thon/database"
	"win-a-thon/models"
)

func TestGetLeaderboard(t *testing.T) {
	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      driver,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	var participants []models.Participant

	hackathonMockRows := sqlmock.NewRows([]string{"hackathon_id", "user_id", "demo_url", "code_url", "score"}).
		AddRow(1, 1, "abc", "abc", 10).
		AddRow(2, 2, "def", "def", 20).
		AddRow(1, 3, "ijk", "ijk", 30)

	fmt.Println(hackathonMockRows, "test line")
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `participants` WHERE hackathon_id = ?  ORDER BY score desc")).WithArgs("1").WillReturnRows(hackathonMockRows)

	err = GetLeaderboard(&participants, "1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(participants)
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestGetSubmissions(t *testing.T) {
	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      driver,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	var participants []models.Participant

	hackathonMockRows := sqlmock.NewRows([]string{"hackathon_id", "user_id", "demo_url", "code_url", "score"}).
		AddRow(1, 1, "abc", "abc", 10).
		AddRow(2, 2, "def", "def", 20).
		AddRow(1, 3, "ijk", "ijk", 30)

	fmt.Println(hackathonMockRows, "test line")
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `participants` WHERE hackathon_id = ? ")).WithArgs("1").WillReturnRows(hackathonMockRows)

	err = GetSubmissions(&participants, "1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(participants)
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestGetWinners(t *testing.T) {
	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      driver,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	var winners []models.Participant

	participantMockRows := sqlmock.NewRows([]string{"hackathon_id", "user_id", "demo_url", "code_url", "score"}).
		AddRow(1, 1, "abc", "abc", 10).
		AddRow(2, 2, "def", "def", 20).
		AddRow(1, 3, "ijk", "ijk", 30)

	fmt.Println(participantMockRows, "test line")
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `participants` WHERE hackathon_id = ?  ORDER BY score desc LIMIT 3")).WithArgs("1").WillReturnRows(participantMockRows)

	err = GetWinners(&winners, "1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(winners)
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestHackathonFromHackathonID(t *testing.T) {
	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      driver,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	var hackathon models.Hackathon
	var hackathon_id = "1"
	var created_at = time.Now()

	hackathonMockRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_time", "ending_time", "res" +
		"ult_time", "organisation_name", "organisation_id", "description", "admin_approved"}).AddRow("1", created_at, created_at, created_at, "Test Hackathon", created_at, created_at, created_at, "Winathon", "1", "It is a test Hackathon", true)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs(hackathon_id).WillReturnRows(hackathonMockRows)

	hackathon, err = HackathonFromHackathonID(hackathon_id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(hackathon)
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestUserFromUserID(t *testing.T) {
	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      driver,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	var user models.User
	var user_id = 1

	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, "name", "fullname", "password", "email", "abc", "abc", "abc", "abc", false)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE ID = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).WithArgs(user_id).WillReturnRows(userRows)
	user, err = UserFromUserID(user_id)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestJudgeSubmisson(t *testing.T) {
	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      driver,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	// First Query
	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, "name", "fullname", "password", "email", "abc", "abc", "abc", "abc", false)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).WithArgs("name").WillReturnRows(userRows)

	// Second Query
	participantMockRows := sqlmock.NewRows([]string{"hackathon_id", "user_id", "demo_url", "code_url", "score"}).
		AddRow(1, 1, "abc", "abc", 10)
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `participants` WHERE hackathon_id = ? AND user_id = ? ORDER BY `participants`.`hackathon_id` LIMIT 1")).WithArgs().WillReturnRows(participantMockRows)

	// Third Query
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		"UPDATE `participants` SET `demo_url`=?,`code_url`=?,`score`=? WHERE `hackathon_id` = ? AND `user_id` = ?")).WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	var participant models.Participant
	err = JudgeSubmisson(&participant, "1", "name")

	if err != nil {
		log.Fatal(err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}
