package repo

import (
	"database/sql/driver"
	"fmt"
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
	"win-a-thon/utils"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestCreateUser(t *testing.T) {
	db, mock, err = sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	user := &models.User{
		Username:       "abcdef",
		FullName:       "abcdef ghijkl",
		HashedPassword: "fjaskfjaslfj",
		Email:          "abc.defghi@gmail.com",
		LinkedIn:       "",
		GitHub:         "",
		WebLink:        "",
		Organisation:   "Razorpay",
		IsAdmin:        false,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`username`,`full_name`,`hashed_password`,`email`,`linked_in`,`git_hub`,`web_link`,`organisation`,`is_admin`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)")).
		WithArgs(AnyTime{}, AnyTime{}, nil, user.Username, user.FullName, user.HashedPassword, user.Email, user.LinkedIn, user.GitHub, user.WebLink, user.Organisation, user.IsAdmin).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := CreateUser(user)

	require.NoError(t, err)
}

func TestGetProfileByUsername(t *testing.T) {
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	user := &models.User{
		Username:       "abcdef",
		FullName:       "abcdef ghijkl",
		HashedPassword: "fjaskfjaslfj",
		Email:          "abc.defghi@gmail.com",
		LinkedIn:       "",
		GitHub:         "",
		WebLink:        "",
		Organisation:   "Razorpay",
		IsAdmin:        false,
	}

	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, user.Username, user.FullName, user.HashedPassword, user.Email, user.LinkedIn, user.GitHub, user.WebLink, user.Organisation, user.IsAdmin)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("abcdef").
		WillReturnRows(userRows)

	reqUser, err := GetProfileByUsername("abcdef")

	fmt.Println(reqUser)
	require.NoError(t, err)
}

func TestUpdateUserProfile(t *testing.T) {
	db, mock, err = sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	user := &models.User{
		Username:       "abcdef",
		FullName:       "abcdef ghijkl",
		HashedPassword: "fjaskfjaslfj",
		Email:          "abc.defghi@gmail.com",
		LinkedIn:       "abc",
		GitHub:         "abc",
		WebLink:        "abc",
		Organisation:   "Razorpay",
		IsAdmin:        false,
	}

	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, user.Username, user.FullName, user.HashedPassword, user.Email, user.LinkedIn, user.GitHub, user.WebLink, user.Organisation, user.IsAdmin)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("abcdef").
		WillReturnRows(userRows)

	req := utils.UpdateProfileRequest{
		Email:        "abc.defghi@gmail.com",
		LinkedIn:     "abc",
		GitHub:       "abc",
		WebLink:      "abc",
		Organisation: "Razorpay",
	}

	reqUser, err := UpdateUserProfile("abcdef", req)

	fmt.Println(reqUser)
	require.NoError(t, err)
}

func TestDeleteUserByUsername(t *testing.T) {
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	user := &models.User{
		Username:       "abcdef",
		FullName:       "abcdef ghijkl",
		HashedPassword: "fjaskfjaslfj",
		Email:          "abc.defghi@gmail.com",
		LinkedIn:       "",
		GitHub:         "",
		WebLink:        "",
		Organisation:   "Razorpay",
		IsAdmin:        false,
	}

	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, user.Username, user.FullName, user.HashedPassword, user.Email, user.LinkedIn, user.GitHub, user.WebLink, user.Organisation, user.IsAdmin)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `deleted_at`=? WHERE username = ? AND `users`.`deleted_at` IS NULL")).WithArgs(AnyTime{}, "abcdef").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	fmt.Println(userRows)
	err := DeleteUserByUsername("abcdef")
	require.NoError(t, err)
}

func TestShowUserHackathons(t *testing.T) {
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	database.DB, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	userhackathonMockRows := sqlmock.NewRows([]string{"hackathon_id", "user_id", "demo_url", "code_url", "score"}).AddRow("1", "1", "", "", "70")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `participants` WHERE user_id = ?")).WithArgs(1).WillReturnRows(userhackathonMockRows)
	var p []models.Participant
	err := ShowUserHackathons(&p, 1)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(hackathon)
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

