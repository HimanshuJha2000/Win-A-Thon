package controllers

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	mysql2 "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"
	"win-a-thon/database"
	"win-a-thon/models"
	"win-a-thon/token"
	"win-a-thon/utils"
)

var mock sqlmock.Sqlmock
var err error
var db *sql.DB

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestCreateUser(t *testing.T) {

	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	ctx.Request.Header.Add("Accept", "application/json")

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`username`,`full_name`,`hashed_password`,`email`,`linked_in`,`git_hub`,`web_link`,`organisation`,`is_admin`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)")).
		WithArgs(AnyTime{}, AnyTime{}, nil, "bar", "abc def", sqlmock.AnyArg(), sqlmock.AnyArg(), "", "", "", "", false).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	jsonParams := `{"username": "bar", "password": "abcdefghi",
		"email": "abcdef.ghijkl@gmail.com", "full_name": "abc def"}`

	r.POST("/users", CreateUser)
	req, _ := http.NewRequest("POST", "/users", strings.NewReader(jsonParams))
	r.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusOK, w.Code)
}

func TestCreateUserBadRequest(t *testing.T) {

	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	ctx.Request.Header.Add("Accept", "application/json")

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`username`,`full_name`,`hashed_password`,`email`,`linked_in`,`git_hub`,`web_link`,`organisation`,`is_admin`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)")).
		WithArgs(AnyTime{}, AnyTime{}, nil, "bar", "abc def", sqlmock.AnyArg(), sqlmock.AnyArg(), "", "", "", "", false).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	jsonParams := `{"username": "bar", "password": "abc",
		"email": "abcdef.ghijkl@gmail.com", "full_name": "abc def"}`

	r.POST("/users", CreateUser)
	req, _ := http.NewRequest("POST", "/users", strings.NewReader(jsonParams))
	r.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusBadRequest, w.Code)
}

func TestCreateUserDupicateError(t *testing.T) {

	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	ctx.Request.Header.Add("Accept", "application/json")

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`username`,`full_name`,`hashed_password`,`email`,`linked_in`,`git_hub`,`web_link`,`organisation`,`is_admin`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)")).
		WithArgs(AnyTime{}, AnyTime{}, nil, "bar", "abc def", sqlmock.AnyArg(), sqlmock.AnyArg(), "", "", "", "", false).
		WillReturnError(&mysql2.MySQLError{Number: utils.DuplicateRecordErrorCode})
	mock.ExpectRollback()

	jsonParams := `{"username": "bar", "password": "abcdefghi",
		"email": "abcdef.ghijkl@gmail.com", "full_name": "abc def"}`

	r.POST("/users", CreateUser)
	req, _ := http.NewRequest("POST", "/users", strings.NewReader(jsonParams))
	r.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusBadRequest, w.Code)
}

func TestCreateUserInternalServerError(t *testing.T) {

	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	ctx.Request.Header.Add("Accept", "application/json")

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`username`,`full_name`,`hashed_password`,`email`,`linked_in`,`git_hub`,`web_link`,`organisation`,`is_admin`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)")).
		WithArgs(AnyTime{}, AnyTime{}, nil, "bar", "abc def", sqlmock.AnyArg(), sqlmock.AnyArg(), "", "", "", "", false).
		WillReturnError(&mysql2.MySQLError{Number: 1061})
	mock.ExpectRollback()

	jsonParams := `{"username": "bar", "password": "abcdefghi",
		"email": "abcdef.ghijkl@gmail.com", "full_name": "abc def"}`

	r.POST("/users", CreateUser)
	req, _ := http.NewRequest("POST", "/users", strings.NewReader(jsonParams))
	r.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusInternalServerError, w.Code)
}

func TestViewUserProfile(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	ctx.Request.Header.Add("Accept", "application/json")
	ctx.Keys = make(map[string]interface{})
	ctx.Keys["authorization_payload"] = &token.Payload{
		Username:  "admin",
		IssuedAt:  time.Now(),
		ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
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
		WithArgs("admin").
		WillReturnRows(userRows)

	ViewUserProfile(ctx)
	assert.EqualValues(t, http.StatusOK, w.Code)
}

func TestViewUserProfileInternalServerError(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	ctx.Request.Header.Add("Accept", "application/json")
	ctx.Keys = make(map[string]interface{})
	ctx.Keys["authorization_payload"] = &token.Payload{
		Username:  "admin",
		IssuedAt:  time.Now(),
		ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
	}

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("admin").
		WillReturnError(errors.New("internal server error"))

	ViewUserProfile(ctx)
	assert.EqualValues(t, http.StatusInternalServerError, w.Code)
}

func TestLoginUser(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	ctx.Request.Header.Add("Accept", "application/json")

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

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

	jsonParams := `{"username": "bar", "password": "abcdefghi"}`

	hashedPassword, err := utils.HashPassword("abcdefghi")

	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, user.Username, user.FullName, hashedPassword, user.Email, user.LinkedIn, user.GitHub, user.WebLink, user.Organisation, user.IsAdmin)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("bar").
		WillReturnRows(userRows)

	r.POST("/users", LoginUser)
	req, _ := http.NewRequest("POST", "/users", strings.NewReader(jsonParams))
	r.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusOK, w.Code)
}

func TestLoginUserNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	ctx.Request.Header.Add("Accept", "application/json")

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	jsonParams := `{"username": "bar", "password": "abcdefghi"}`

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("bar").
		WillReturnError(gorm.ErrRecordNotFound)

	r.POST("/users", LoginUser)
	req, _ := http.NewRequest("POST", "/users", strings.NewReader(jsonParams))
	r.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusNotFound, w.Code)
}

func TestLoginUserBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	ctx.Request.Header.Add("Accept", "application/json")

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

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

	jsonParams := `{"password": "abcdefghi"}`

	hashedPassword, err := utils.HashPassword("abcdefghi")

	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, user.Username, user.FullName, hashedPassword, user.Email, user.LinkedIn, user.GitHub, user.WebLink, user.Organisation, user.IsAdmin)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("bar").
		WillReturnRows(userRows)

	r.POST("/users", LoginUser)
	req, _ := http.NewRequest("POST", "/users", strings.NewReader(jsonParams))
	r.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusBadRequest, w.Code)
}

func TestDeleteUserProfile(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	ctx.Request.Header.Add("Accept", "application/json")
	ctx.Keys = make(map[string]interface{})
	ctx.Keys["authorization_payload"] = &token.Payload{
		Username:  "admin",
		IssuedAt:  time.Now(),
		ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
	}

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `deleted_at`=? WHERE username = ? AND `users`.`deleted_at` IS NULL")).WithArgs(AnyTime{}, "admin").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	DeleteUserProfile(ctx)
	assert.EqualValues(t, http.StatusOK, w.Code)
}

func TestUpdateUserProfileBadRequestCase(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	ctx.Request.Header.Add("Accept", "application/json")
	ctx.Keys = make(map[string]interface{})
	ctx.Keys["authorization_payload"] = &token.Payload{
		Username:  "admin",
		IssuedAt:  time.Now(),
		ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
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

	//jsonParams :=  `{"Email": "bar", "LinkedIn": "abcdefghi",
	//	"Github": "abcdef.ghijkl@gmail.com", "weblink": "abc def", "organisation": "Razorpay"}`

	UpdateUserProfile(ctx)

	assert.EqualValues(t, http.StatusBadRequest, w.Code)
}

func TestUpdateUserProfileCorrectRequestCase(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	ctx.Request.Header.Add("Accept", "application/json")
	//ctx.Keys = make(map[string]interface{})
	//ctx.Keys["authorization_payload"] = &token.Payload{
	//	Username: "admin",
	//	IssuedAt: time.Now(),
	//	ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
	//}

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

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
		WithArgs("admin").
		WillReturnRows(userRows)

	//mock.ExpectExec(regexp.QuoteMeta(
	//	))

	jsonParams := `{"Email": "bar", "LinkedIn": "abcdefghi",
		"Github": "abcdef.ghijkl@gmail.com", "weblink": "abc def", "organisation": "Razorpay"}`

	//ctx.Params = []gin.Param{{Key: "email", Value: "bar"}}

	//UpdateUserProfile(ctx)

	r.Use(func(c *gin.Context) {
		c.Set("authorization_payload", &token.Payload{
			Username:  "admin",
			IssuedAt:  time.Now(),
			ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
		})
	})
	r.PATCH("/users", UpdateUserProfile)
	req, _ := http.NewRequest("PATCH", "/users", strings.NewReader(jsonParams))
	r.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusOK, w.Code)
}

func TestCreateAdminUser(t *testing.T) {

	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	ctx.Request.Header.Add("Accept", "application/json")

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`username`,`full_name`,`hashed_password`,`email`,`linked_in`,`git_hub`,`web_link`,`organisation`,`is_admin`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)")).
		WithArgs(AnyTime{}, AnyTime{}, nil, "abdfbar", "abc def", sqlmock.AnyArg(), sqlmock.AnyArg(), "", "", "", "", true).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	jsonParams := `{"username": "abdfbar", "password": "!@~#$%^&*",
		"email": "abcdef.ghijkl@gmail.com", "full_name": "abc def", "linked_in": "",
		"github": "", "web_link": "", "organisation": ""}`

	r.Use(func(c *gin.Context) {
		c.Set("authorization_payload", &token.Payload{
			Username:  "admin",
			IssuedAt:  time.Now(),
			ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
		})
	})
	r.POST("/users/makeAdmin", CreateAdminUser)
	req, _ := http.NewRequest("POST", "/users/makeAdmin", strings.NewReader(jsonParams))
	r.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusOK, w.Code)
}

func TestCreateAdminUserDuplicateError(t *testing.T) {

	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	ctx.Request.Header.Add("Accept", "application/json")

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`username`,`full_name`,`hashed_password`,`email`,`linked_in`,`git_hub`,`web_link`,`organisation`,`is_admin`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)")).
		WithArgs(AnyTime{}, AnyTime{}, nil, "abdfbar", "abc def", sqlmock.AnyArg(), sqlmock.AnyArg(), "", "", "", "", true).
		WillReturnError(&mysql2.MySQLError{Number: utils.DuplicateRecordErrorCode})
	mock.ExpectRollback()

	jsonParams := `{"username": "abdfbar", "password": "!@~#$%^&*",
		"email": "abcdef.ghijkl@gmail.com", "full_name": "abc def", "linked_in": "",
		"github": "", "web_link": "", "organisation": ""}`

	r.Use(func(c *gin.Context) {
		c.Set("authorization_payload", &token.Payload{
			Username:  "admin",
			IssuedAt:  time.Now(),
			ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
		})
	})
	r.POST("/users/makeAdmin", CreateAdminUser)
	req, _ := http.NewRequest("POST", "/users/makeAdmin", strings.NewReader(jsonParams))
	r.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusBadRequest, w.Code)
}

func TestCreateAdminUserInternalServerError(t *testing.T) {

	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	ctx.Request.Header.Add("Accept", "application/json")

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`username`,`full_name`,`hashed_password`,`email`,`linked_in`,`git_hub`,`web_link`,`organisation`,`is_admin`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)")).
		WithArgs(AnyTime{}, AnyTime{}, nil, "abdfbar", "abc def", sqlmock.AnyArg(), sqlmock.AnyArg(), "", "", "", "", true).
		WillReturnError(&mysql2.MySQLError{Number: 1061})
	mock.ExpectRollback()

	jsonParams := `{"username": "abdfbar", "password": "!@~#$%^&*",
		"email": "abcdef.ghijkl@gmail.com", "full_name": "abc def", "linked_in": "",
		"github": "", "web_link": "", "organisation": ""}`

	r.Use(func(c *gin.Context) {
		c.Set("authorization_payload", &token.Payload{
			Username:  "admin",
			IssuedAt:  time.Now(),
			ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
		})
	})
	r.POST("/users/makeAdmin", CreateAdminUser)
	req, _ := http.NewRequest("POST", "/users/makeAdmin", strings.NewReader(jsonParams))
	r.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusInternalServerError, w.Code)
}

func TestCreateAdminUserBadRequest(t *testing.T) {

	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	ctx.Request.Header.Add("Accept", "application/json")

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`username`,`full_name`,`hashed_password`,`email`,`linked_in`,`git_hub`,`web_link`,`organisation`,`is_admin`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)")).
		WithArgs(AnyTime{}, AnyTime{}, nil, "abdfbar", "abc def", sqlmock.AnyArg(), sqlmock.AnyArg(), "", "", "", "", true).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	jsonParams := `{"username": "abdfbar", "password": "abcdefghi",
		"email": "abcdef.ghijkl@gmail.com", "full_name": "abc def", "linked_in": "",
		"github": "", "web_link": "", "organisation": "", "is_admin": "true"}`

	r.Use(func(c *gin.Context) {
		c.Set("authorization_payload", &token.Payload{
			Username:  "admin",
			IssuedAt:  time.Now(),
			ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
		})
	})
	r.POST("/users/makeAdmin", CreateAdminUser)
	req, _ := http.NewRequest("POST", "/users/makeAdmin", strings.NewReader(jsonParams))
	r.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusBadRequest, w.Code)
}

func TestShowUserHackathons(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Keys = make(map[string]interface{})
	ctx.Keys["authorization_payload"] = &token.Payload{
		Username:  "admin",
		IssuedAt:  time.Now(),
		ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
	}
	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	ctx.Request.Header.Add("Accept", "application/json")

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

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

	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).AddRow(time.Now(), time.Now(), nil, user.Username, user.FullName, user.HashedPassword, user.Email, user.LinkedIn, user.GitHub, user.WebLink, user.Organisation, user.IsAdmin)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).WithArgs("admin").WillReturnRows(userRows)

	userhackathonMockRows := sqlmock.NewRows([]string{"hackathon_id", "user_id", "demo_url", "code_url", "score"}).AddRow("1", "0", "", "", "70")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `participants` WHERE user_id = ?")).WithArgs(0).WillReturnRows(userhackathonMockRows)
	created_at := time.Now()
	hackathonMockRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_time", "ending_time", "result_time", "organisation_name", "organiser_id", "description", "admin_approved"}).AddRow("1", created_at, created_at, created_at, "Test Hackathon", created_at, created_at, created_at, "Winathon", "1", "It is a test Hackathon", "1")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `hackathons` WHERE (id = ? AND admin_approved = ?)")).WithArgs("1", true).WillReturnRows(hackathonMockRows)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `participants` WHERE hackathon_id = ?")).WithArgs("1").WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))
	ShowUserHackathons(ctx)
	assert.EqualValues(t, http.StatusOK, w.Code)
}
