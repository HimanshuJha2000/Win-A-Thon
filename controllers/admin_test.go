package controllers

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"
	"win-a-thon/database"
	"win-a-thon/models"
	"win-a-thon/token"
)

func TestGetAdminApprovalSuccessForTrue(t *testing.T) {

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Keys = make(map[string]interface{})
	ctx.Keys["authorization_payload"] = &token.Payload{
		Username:  "admin",
		IssuedAt:  time.Now(),
		ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
	}

	ctx.Params = []gin.Param{
		{
			Key:   "hackathon_id",
			Value: "1",
		},
		{
			Key:   "value",
			Value: "1",
		},
	}
	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	user := &models.User{
		Username:       "admin",
		FullName:       "abcdef ghijkl",
		HashedPassword: "fjaskfjaslfj",
		Email:          "abc.defghi@gmail.com",
		LinkedIn:       "",
		GitHub:         "",
		WebLink:        "",
		Organisation:   "Razorpay",
		IsAdmin:        true,
	}

	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, user.Username, user.FullName, user.HashedPassword, user.Email, user.LinkedIn, user.GitHub, user.WebLink, user.Organisation, user.IsAdmin)

	var user_id = 1
	userRows1 := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, "name", "fullname", "password", "email", "abc", "abc", "abc", "abc", true)

	date := time.Date(2021, 8, 2, 0, 0, 0, 0, time.UTC)
	hackathonMockRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_" +
		"time", "ending_time", "result_time", "organisation_name", "organiser_id", "description", "admin_" +
		"approved"}).AddRow("1", date, date, nil, "Test Hackathon", date, date, date, "Winathon", 1, "It is a test"+
		" Hackathon", true)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("admin").
		WillReturnRows(userRows)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs("1").WillReturnRows(hackathonMockRow)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE ID = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).WithArgs(user_id).WillReturnRows(userRows1)

	date3 := time.Date(2021, 8, 2, 0, 0, 0, 0, time.UTC)

	hackathonMockRow3 := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_" +
		"time", "ending_time", "result_time", "organisation_name", "organiser_id", "description", "admin_" +
		"approved"}).AddRow(1, date3, date3, nil, "Test Hackathon", date3, date3, date3, "Winathon", 1, "It is a test"+
		" Hackathon", false)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE id = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs("1").WillReturnRows(hackathonMockRow3)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `hackathons` SET `created_at`=?,`updated_at`=?,`deleted_at`=?,`title`=?,`starting_time`=?,`ending_time`=?,`result_time`=?,`organisation_name`=?,`organiser_id`=?,`description`=?,`admin_approved`=? WHERE `id` = ?")).WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	GetAdminApproval(ctx)
	assert.EqualValues(t, http.StatusOK, w.Code)
}

func TestGetAdminApprovalSuccessForFalse(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Keys = make(map[string]interface{})
	ctx.Keys["authorization_payload"] = &token.Payload{
		Username:  "admin",
		IssuedAt:  time.Now(),
		ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
	}

	ctx.Params = []gin.Param{
		{
			Key:   "hackathon_id",
			Value: "1",
		},
		{
			Key:   "value",
			Value: "0",
		},
	}
	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	user := &models.User{
		Username:       "admin",
		FullName:       "abcdef ghijkl",
		HashedPassword: "fjaskfjaslfj",
		Email:          "abc.defghi@gmail.com",
		LinkedIn:       "",
		GitHub:         "",
		WebLink:        "",
		Organisation:   "Razorpay",
		IsAdmin:        true,
	}

	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, user.Username, user.FullName, user.HashedPassword, user.Email, user.LinkedIn, user.GitHub, user.WebLink, user.Organisation, user.IsAdmin)

	var user_id = 1
	userRows1 := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, "name", "fullname", "password", "email", "abc", "abc", "abc", "abc", true)

	date := time.Date(2021, 8, 2, 0, 0, 0, 0, time.UTC)
	hackathonMockRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_" +
		"time", "ending_time", "result_time", "organisation_name", "organiser_id", "description", "admin_" +
		"approved"}).AddRow("1", date, date, nil, "Test Hackathon", date, date, date, "Winathon", 1, "It is a test"+
		" Hackathon", true)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("admin").
		WillReturnRows(userRows)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs("1").WillReturnRows(hackathonMockRow)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE ID = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).WithArgs(user_id).WillReturnRows(userRows1)

	date3 := time.Date(2021, 8, 2, 0, 0, 0, 0, time.UTC)

	hackathonMockRow3 := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_" +
		"time", "ending_time", "result_time", "organisation_name", "organiser_id", "description", "admin_" +
		"approved"}).AddRow(1, date3, date3, nil, "Test Hackathon", date3, date3, date3, "Winathon", 1, "It is a test"+
		" Hackathon", false)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE id = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs("1").WillReturnRows(hackathonMockRow3)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `hackathons` SET `created_at`=?,`updated_at`=?,`deleted_at`=?,`title`=?,`starting_time`=?,`ending_time`=?,`result_time`=?,`organisation_name`=?,`organiser_id`=?,`description`=?,`admin_approved`=? WHERE `id` = ?")).WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `hackathons` SET `deleted_at`=? WHERE id = ? AND `hackathons`.`deleted_at` IS NULL")).WithArgs(AnyTime{}, "1").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	GetAdminApproval(ctx)
	assert.EqualValues(t, http.StatusOK, w.Code)
}

func TestGetAdminApprovalIsAdminFalse(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Keys = make(map[string]interface{})
	ctx.Keys["authorization_payload"] = &token.Payload{
		Username:  "admin",
		IssuedAt:  time.Now(),
		ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
	}

	ctx.Params = []gin.Param{
		{
			Key:   "hackathon_id",
			Value: "1",
		},
		{
			Key:   "value",
			Value: "0",
		},
	}
	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	user := &models.User{
		Username:       "admin",
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

	var user_id = 1
	userRows1 := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, "name", "fullname", "password", "email", "abc", "abc", "abc", "abc", false)

	date := time.Date(2021, 8, 2, 0, 0, 0, 0, time.UTC)
	hackathonMockRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_" +
		"time", "ending_time", "result_time", "organisation_name", "organiser_id", "description", "admin_" +
		"approved"}).AddRow("1", date, date, nil, "Test Hackathon", date, date, date, "Winathon", 1, "It is a test"+
		" Hackathon", true)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("admin").
		WillReturnRows(userRows)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs("1").WillReturnRows(hackathonMockRow)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE ID = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).WithArgs(user_id).WillReturnRows(userRows1)

	date3 := time.Date(2021, 8, 2, 0, 0, 0, 0, time.UTC)

	hackathonMockRow3 := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_" +
		"time", "ending_time", "result_time", "organisation_name", "organiser_id", "description", "admin_" +
		"approved"}).AddRow(1, date3, date3, nil, "Test Hackathon", date3, date3, date3, "Winathon", 1, "It is a test"+
		" Hackathon", false)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE id = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs("1").WillReturnRows(hackathonMockRow3)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `hackathons` SET `created_at`=?,`updated_at`=?,`deleted_at`=?,`title`=?,`starting_time`=?,`ending_time`=?,`result_time`=?,`organisation_name`=?,`organiser_id`=?,`description`=?,`admin_approved`=? WHERE `id` = ?")).WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	GetAdminApproval(ctx)
	assert.EqualValues(t, http.StatusOK, w.Code)

}

func TestGetAdminApprovalRecordNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Keys = make(map[string]interface{})
	ctx.Keys["authorization_payload"] = &token.Payload{
		Username:  "admin",
		IssuedAt:  time.Now(),
		ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
	}

	ctx.Params = []gin.Param{
		{
			Key:   "hackathon_id",
			Value: "1",
		},
		{
			Key:   "value",
			Value: "0",
		},
	}

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	//GetProfileByUsername

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("admin").
		WillReturnError(errors.New("record not found"))

	GetAdminApproval(ctx)
	assert.EqualValues(t, http.StatusNotFound, w.Code)
}

func TestGetAdminApprovalProfileByUsernameInternalServerError(t *testing.T) {

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Keys = make(map[string]interface{})
	ctx.Keys["authorization_payload"] = &token.Payload{
		Username:  "admin",
		IssuedAt:  time.Now(),
		ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
	}

	ctx.Params = []gin.Param{
		{
			Key:   "hackathon_id",
			Value: "1",
		},
		{
			Key:   "value",
			Value: "0",
		},
	}

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	//GetProfileByUsername

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("admin").
		WillReturnError(errors.New("server error"))

	GetAdminApproval(ctx)
	assert.EqualValues(t, http.StatusInternalServerError, w.Code)
}

func TestGetAdminApprovalHackathonFromHackathonIDError(t *testing.T) {

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Keys = make(map[string]interface{})
	ctx.Keys["authorization_payload"] = &token.Payload{
		Username:  "admin",
		IssuedAt:  time.Now(),
		ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
	}

	ctx.Params = []gin.Param{
		{
			Key:   "hackathon_id",
			Value: "1",
		},
		{
			Key:   "value",
			Value: "0",
		},
	}

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	user := &models.User{
		Username:       "admin",
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
		WithArgs("admin").
		WillReturnRows(userRows)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs("2").WillReturnError(errors.New("server error"))

	GetAdminApproval(ctx)
	assert.EqualValues(t, http.StatusInternalServerError, w.Code)
}

func TestGetAdminApprovalUserFromUserIDError(t *testing.T) {

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Keys = make(map[string]interface{})
	ctx.Keys["authorization_payload"] = &token.Payload{
		Username:  "admin",
		IssuedAt:  time.Now(),
		ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
	}

	ctx.Params = []gin.Param{
		{
			Key:   "hackathon_id",
			Value: "1",
		},
		{
			Key:   "value",
			Value: "0",
		},
	}

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	user := &models.User{
		Username:       "admin",
		FullName:       "abcdef ghijkl",
		HashedPassword: "fjaskfjaslfj",
		Email:          "abc.defghi@gmail.com",
		LinkedIn:       "",
		GitHub:         "",
		WebLink:        "",
		Organisation:   "Razorpay",
		IsAdmin:        true,
	}

	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, user.Username, user.FullName, user.HashedPassword, user.Email, user.LinkedIn, user.GitHub, user.WebLink, user.Organisation, user.IsAdmin)

	date := time.Date(2021, 8, 2, 0, 0, 0, 0, time.UTC)
	hackathonMockRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_" +
		"time", "ending_time", "result_time", "organisation_name", "organiser_id", "description", "admin_" +
		"approved"}).AddRow("1", date, date, nil, "Test Hackathon", date, date, date, "Winathon", 1, "It is a test"+
		" Hackathon", true)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("admin").
		WillReturnRows(userRows)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs("1").WillReturnRows(hackathonMockRow)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE ID = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).WithArgs("4").WillReturnError(errors.New("server error"))

	GetAdminApproval(ctx)
	assert.EqualValues(t, http.StatusInternalServerError, w.Code)
}

func TestGetAdminApprovalSuccessForError(t *testing.T) {

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Keys = make(map[string]interface{})
	ctx.Keys["authorization_payload"] = &token.Payload{
		Username:  "admin",
		IssuedAt:  time.Now(),
		ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
	}

	ctx.Params = []gin.Param{
		{
			Key:   "hackathon_id",
			Value: "1",
		},
		{
			Key:   "value",
			Value: "1",
		},
	}
	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	user := &models.User{
		Username:       "admin",
		FullName:       "abcdef ghijkl",
		HashedPassword: "fjaskfjaslfj",
		Email:          "abc.defghi@gmail.com",
		LinkedIn:       "",
		GitHub:         "",
		WebLink:        "",
		Organisation:   "Razorpay",
		IsAdmin:        true,
	}

	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, user.Username, user.FullName, user.HashedPassword, user.Email, user.LinkedIn, user.GitHub, user.WebLink, user.Organisation, user.IsAdmin)

	var user_id = 1
	userRows1 := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, "name", "fullname", "password", "email", "abc", "abc", "abc", "abc", true)

	date := time.Date(2021, 8, 2, 0, 0, 0, 0, time.UTC)
	hackathonMockRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_" +
		"time", "ending_time", "result_time", "organisation_name", "organiser_id", "description", "admin_" +
		"approved"}).AddRow("1", date, date, nil, "Test Hackathon", date, date, date, "Winathon", 1, "It is a test"+
		" Hackathon", true)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("admin").
		WillReturnRows(userRows)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs("1").WillReturnRows(hackathonMockRow)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE ID = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).WithArgs(user_id).WillReturnRows(userRows1)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE id = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs("4").WillReturnError(errors.New("server error"))

	GetAdminApproval(ctx)
	assert.EqualValues(t, http.StatusInternalServerError, w.Code)
}
