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
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"
	"win-a-thon/database"
	"win-a-thon/models"
	"win-a-thon/token"
)

func TestNotifyParticipant(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	//Authorization
	r.Use(func(c *gin.Context) {
		c.Set("authorization_payload", &token.Payload{
			Username:  "admin",
			IssuedAt:  time.Now(),
			ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
		})
	})

	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	ctx.Request.Header.Add("Accept", "application/json")

	//Parameters
	ctx.Params = []gin.Param{
		{
			Key:   "username",
			Value: "admin",
		},
		{
			Key:   "hackathon_id",
			Value: "1",
		},
	}

	//JSON params
	jsonParams := `{"subject": "modi15lakhdo", "message": "ye gormint bik chuki hai"}`

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	//GetProfileByUsername
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

	//GetProfileByUsername
	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, user.Username, user.FullName, user.HashedPassword, user.Email, user.LinkedIn, user.GitHub, user.WebLink, user.Organisation, user.IsAdmin)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("admin").
		WillReturnRows(userRows)

	var created_at = time.Now()

	hackathonMockRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_time", "ending_time", "res" +
		"ult_time", "organisation_name", "organisation_id", "description", "admin_approved"}).AddRow("1", created_at, created_at, created_at, "Test Hackathon", created_at, created_at, created_at, "Winathon", "1", "It is a test Hackathon", true)

	//HackathonFromHackathonID
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs("1").WillReturnRows(hackathonMockRows)

	//GetProfileByUsername
	userRows1 := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, user.Username, user.FullName, user.HashedPassword, user.Email, user.LinkedIn, user.GitHub, user.WebLink, user.Organisation, user.IsAdmin)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("admin").
		WillReturnRows(userRows1)

	//Notify participant
	r.POST("/hackathons/:hackathon_id/notify/:username", NotifyParticipant)
	req, _ := http.NewRequest("POST", "/hackathons/1/notify/admin", strings.NewReader(jsonParams))
	r.ServeHTTP(w, req)

	//NotifyParticipant(ctx)
	assert.EqualValues(t, http.StatusInternalServerError, w.Code)

}

func TestNotifyParticipantGetProfileByUsernameStatusBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	//Authorization
	r.Use(func(c *gin.Context) {
		c.Set("authorization_payload", &token.Payload{
			Username:  "admin",
			IssuedAt:  time.Now(),
			ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
		})
	})

	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	ctx.Request.Header.Add("Accept", "application/json")

	//Parameters
	ctx.Params = []gin.Param{
		{
			Key:   "username",
			Value: "admin",
		},
		{
			Key:   "hackathon_id",
			Value: "1",
		},
	}

	//JSON params
	jsonParams := `{"subject": "modi15lakhdo", "message": "ye gormint bik chuki hai"}`

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
		WillReturnError(errors.New("bad request"))

	//Notify participant
	r.POST("/hackathons/:hackathon_id/notify/:username", NotifyParticipant)
	req, _ := http.NewRequest("POST", "/hackathons/1/notify/admin", strings.NewReader(jsonParams))
	r.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusBadRequest, w.Code)

}

func TestNotifyParticipantHackathonFromHackathonIDStatusNotFound(t *testing.T) {

	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	//Authorization
	r.Use(func(c *gin.Context) {
		c.Set("authorization_payload", &token.Payload{
			Username:  "admin",
			IssuedAt:  time.Now(),
			ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
		})
	})

	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	ctx.Request.Header.Add("Accept", "application/json")

	//Parameters
	ctx.Params = []gin.Param{
		{
			Key:   "username",
			Value: "admin",
		},
		{
			Key:   "hackathon_id",
			Value: "1",
		},
	}

	//JSON params
	jsonParams := `{"subject": "modi15lakhdo", "message": "ye gormint bik chuki hai"}`

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
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs("2").WillReturnError(errors.New("Status not found"))

	//Notify participant
	r.POST("/hackathons/:hackathon_id/notify/:username", NotifyParticipant)
	req, _ := http.NewRequest("POST", "/hackathons/1/notify/admin", strings.NewReader(jsonParams))
	r.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusNotFound, w.Code)

}

func TestNotifyParticipantHackathonFromHackathonIDLoggedInUserNotCorrect(t *testing.T) {

	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	//Authorization
	r.Use(func(c *gin.Context) {
		c.Set("authorization_payload", &token.Payload{
			Username:  "admin",
			IssuedAt:  time.Now(),
			ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
		})
	})

	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	ctx.Request.Header.Add("Accept", "application/json")

	//Parameters
	ctx.Params = []gin.Param{
		{
			Key:   "username",
			Value: "admin",
		},
		{
			Key:   "hackathon_id",
			Value: "1",
		},
	}

	//JSON params
	jsonParams := `{"subject": "modi15lakhdo", "message": "ye gormint bik chuki hai"}`

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

	date := time.Date(2021, 8, 2, 0, 0, 0, 0, time.UTC)
	hackathonMockRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_" +
		"time", "ending_time", "result_time", "organisation_name", "organiser_id", "description", "admin_" +
		"approved"}).AddRow("1", date, date, nil, "Test Hackathon", date, date, date, "Winathon", 3, "It is a test"+
		" Hackathon", true)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs("1").WillReturnRows(hackathonMockRow)

	//Notify participant
	r.POST("/hackathons/:hackathon_id/notify/:username", NotifyParticipant)
	req, _ := http.NewRequest("POST", "/hackathons/1/notify/admin", strings.NewReader(jsonParams))
	r.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusUnauthorized, w.Code)

}

func TestNotifyParticipantGetProfileByUsernameStatusBadRequest2(t *testing.T) {

	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	//Authorization
	r.Use(func(c *gin.Context) {
		c.Set("authorization_payload", &token.Payload{
			Username:  "admin",
			IssuedAt:  time.Now(),
			ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
		})
	})

	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	ctx.Request.Header.Add("Accept", "application/json")

	//Parameters
	ctx.Params = []gin.Param{
		{
			Key:   "username",
			Value: "admin",
		},
		{
			Key:   "hackathon_id",
			Value: "1",
		},
	}

	//JSON params
	jsonParams := `{"subject": "modi15lakhdo", "message": "ye gormint bik chuki hai"}`

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

	date := time.Date(2021, 8, 2, 0, 0, 0, 0, time.UTC)
	hackathonMockRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_" +
		"time", "ending_time", "result_time", "organisation_name", "organiser_id", "description", "admin_" +
		"approved"}).AddRow("1", date, date, nil, "Test Hackathon", date, date, date, "Winathon", 0, "It is a test"+
		" Hackathon", true)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs("1").WillReturnRows(hackathonMockRow)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("admin").
		WillReturnError(errors.New("bad request"))

	//Notify participant
	r.POST("/hackathons/:hackathon_id/notify/:username", NotifyParticipant)
	req, _ := http.NewRequest("POST", "/hackathons/1/notify/admin", strings.NewReader(jsonParams))
	r.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusBadRequest, w.Code)

}

//Not completed yet
func TestNotifyParticipantFailure(t *testing.T) {

	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	//Authorization
	r.Use(func(c *gin.Context) {
		c.Set("authorization_payload", &token.Payload{
			Username:  "admin",
			IssuedAt:  time.Now(),
			ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
		})
	})

	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	ctx.Request.Header.Add("Accept", "application/json")

	//Parameters
	ctx.Params = []gin.Param{
		{
			Key:   "username",
			Value: "admin",
		},
		{
			Key:   "hackathon_id",
			Value: "1",
		},
	}

	//JSON params
	jsonParams := `{"subject": "modi15lakhdo", "message": "ye gormint bik chuki hai"}`

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	//GetProfileByUsername
	user := &models.User{
		Username:       "admin",
		FullName:       "abcdef ghijkl",
		HashedPassword: "fjaskfjaslfj",
		Email:          "abc.defghi",
		LinkedIn:       "",
		GitHub:         "",
		WebLink:        "",
		Organisation:   "Razorpay",
		IsAdmin:        true,
	}

	//GetProfileByUsername
	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, user.Username, user.FullName, user.HashedPassword, user.Email, user.LinkedIn, user.GitHub, user.WebLink, user.Organisation, user.IsAdmin)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("admin").
		WillReturnRows(userRows)

	var created_at = time.Now()

	hackathonMockRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_time", "ending_time", "res" +
		"ult_time", "organisation_name", "organisation_id", "description", "admin_approved"}).AddRow("1", created_at, created_at, created_at, "Test Hackathon", created_at, created_at, created_at, "Winathon", "1", "It is a test Hackathon", true)

	//HackathonFromHackathonID
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs("1").WillReturnRows(hackathonMockRows)

	//GetProfileByUsername
	userRows1 := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, user.Username, user.FullName, user.HashedPassword, user.Email, user.LinkedIn, user.GitHub, user.WebLink, user.Organisation, user.IsAdmin)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("admin").
		WillReturnRows(userRows1)

	//Notify participant
	r.POST("/hackathons/:hackathon_id/notify/:username", NotifyParticipant)
	req, _ := http.NewRequest("POST", "/hackathons/1/notify/admin", strings.NewReader(jsonParams))
	r.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusInternalServerError, w.Code)

}

func TestNotifyAllSuccess(t *testing.T) {

	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	//Authorization
	r.Use(func(c *gin.Context) {
		c.Set("authorization_payload", &token.Payload{
			Username:  "admin",
			IssuedAt:  time.Now(),
			ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
		})
	})

	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	ctx.Request.Header.Add("Accept", "application/json")

	//Parameters
	ctx.Params = []gin.Param{
		{
			Key:   "username",
			Value: "admin",
		},
		{
			Key:   "hackathon_id",
			Value: "1",
		},
	}

	//JSON params
	jsonParams := `{"subject": "modi15lakhdo", "message": "ye gormint bik chuki hai"}`

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

	//GetProfileByUsername
	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, user.Username, user.FullName, user.HashedPassword, user.Email, user.LinkedIn, user.GitHub, user.WebLink, user.Organisation, user.IsAdmin)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("admin").
		WillReturnRows(userRows)

	var created_at = time.Now()

	hackathonMockRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_time", "ending_time", "res" +
		"ult_time", "organisation_name", "organisation_id", "description", "admin_approved"}).AddRow("1", created_at, created_at, created_at, "Test Hackathon", created_at, created_at, created_at, "Winathon", "1", "It is a test Hackathon", true)

	//HackathonFromHackathonID
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs("1").WillReturnRows(hackathonMockRows)

	//GetSubmissions
	hackathonMockRowsSubmissions := sqlmock.NewRows([]string{"hackathon_id", "user_id", "demo_url", "code_url", "score"}).
		AddRow(1, 1, "abc", "abc", 10).
		AddRow(2, 2, "def", "def", 20).
		AddRow(1, 3, "ijk", "ijk", 30)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `participants` WHERE hackathon_id = ? ")).WithArgs("1").WillReturnRows(hackathonMockRowsSubmissions)

	//UserFromUserID
	var user_id = 1

	userRows2 := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, "name", "fullname", "password", "email", "abc", "abc", "abc", "abc", false)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE ID = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).WithArgs(user_id).WillReturnRows(userRows2)

	//GetProfileByUsername
	userRows1 := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, user.Username, user.FullName, user.HashedPassword, user.Email, user.LinkedIn, user.GitHub, user.WebLink, user.Organisation, user.IsAdmin)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("admin").
		WillReturnRows(userRows1)

	//Notify All
	r.POST("/hackathons/:hackathon_id/notify", NotifyAll)
	req, _ := http.NewRequest("POST", "/hackathons/1/notify", strings.NewReader(jsonParams))
	r.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusOK, w.Code)

}

func TestNotifyAllGetProfileByUsernameError(t *testing.T) {

	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	//Authorization
	r.Use(func(c *gin.Context) {
		c.Set("authorization_payload", &token.Payload{
			Username:  "admin",
			IssuedAt:  time.Now(),
			ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
		})
	})

	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	ctx.Request.Header.Add("Accept", "application/json")

	//Parameters
	ctx.Params = []gin.Param{
		{
			Key:   "username",
			Value: "admin",
		},
		{
			Key:   "hackathon_id",
			Value: "1",
		},
	}

	//JSON params
	jsonParams := `{"subject": "modi15lakhdo", "message": "ye gormint bik chuki hai"}`

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
		WillReturnError(errors.New("bad request"))

	//Notify All
	r.POST("/hackathons/:hackathon_id/notify", NotifyAll)
	req, _ := http.NewRequest("POST", "/hackathons/1/notify", strings.NewReader(jsonParams))
	r.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusBadRequest, w.Code)

}

func TestNotifyAllHackathonFromHackathonIDError(t *testing.T) {

	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	//Authorization
	r.Use(func(c *gin.Context) {
		c.Set("authorization_payload", &token.Payload{
			Username:  "admin",
			IssuedAt:  time.Now(),
			ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
		})
	})

	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	ctx.Request.Header.Add("Accept", "application/json")

	//Parameters
	ctx.Params = []gin.Param{
		{
			Key:   "username",
			Value: "admin",
		},
		{
			Key:   "hackathon_id",
			Value: "1",
		},
	}

	//JSON params
	jsonParams := `{"subject": "modi15lakhdo", "message": "ye gormint bik chuki hai"}`

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
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs("2").WillReturnError(errors.New("Status not found"))

	//Notify All
	r.POST("/hackathons/:hackathon_id/notify", NotifyAll)
	req, _ := http.NewRequest("POST", "/hackathons/1/notify", strings.NewReader(jsonParams))
	r.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusNotFound, w.Code)

}

func TestNotifyAllHackathonFromHackathonIDLoggedInUserNotCorrect(t *testing.T) {

	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	//Authorization
	r.Use(func(c *gin.Context) {
		c.Set("authorization_payload", &token.Payload{
			Username:  "admin",
			IssuedAt:  time.Now(),
			ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
		})
	})

	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	ctx.Request.Header.Add("Accept", "application/json")

	//Parameters
	ctx.Params = []gin.Param{
		{
			Key:   "username",
			Value: "admin",
		},
		{
			Key:   "hackathon_id",
			Value: "1",
		},
	}

	//JSON params
	jsonParams := `{"subject": "modi15lakhdo", "message": "ye gormint bik chuki hai"}`

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

	date := time.Date(2021, 8, 2, 0, 0, 0, 0, time.UTC)
	hackathonMockRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_" +
		"time", "ending_time", "result_time", "organisation_name", "organiser_id", "description", "admin_" +
		"approved"}).AddRow("1", date, date, nil, "Test Hackathon", date, date, date, "Winathon", 3, "It is a test"+
		" Hackathon", true)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs("1").WillReturnRows(hackathonMockRow)

	//Notify All
	r.POST("/hackathons/:hackathon_id/notify", NotifyAll)
	req, _ := http.NewRequest("POST", "/hackathons/1/notify", strings.NewReader(jsonParams))
	r.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusUnauthorized, w.Code)

}

func TestNotifyAllSuccessUserFromUserIDError(t *testing.T) {

	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	//Authorization
	r.Use(func(c *gin.Context) {
		c.Set("authorization_payload", &token.Payload{
			Username:  "admin",
			IssuedAt:  time.Now(),
			ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
		})
	})

	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	ctx.Request.Header.Add("Accept", "application/json")

	//Parameters
	ctx.Params = []gin.Param{
		{
			Key:   "username",
			Value: "admin",
		},
		{
			Key:   "hackathon_id",
			Value: "1",
		},
	}

	//JSON params
	jsonParams := `{"subject": "modi15lakhdo", "message": "ye gormint bik chuki hai"}`

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

	//GetProfileByUsername
	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, user.Username, user.FullName, user.HashedPassword, user.Email, user.LinkedIn, user.GitHub, user.WebLink, user.Organisation, user.IsAdmin)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("admin").
		WillReturnRows(userRows)

	var created_at = time.Now()

	hackathonMockRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_time", "ending_time", "res" +
		"ult_time", "organisation_name", "organisation_id", "description", "admin_approved"}).AddRow("1", created_at, created_at, created_at, "Test Hackathon", created_at, created_at, created_at, "Winathon", "1", "It is a test Hackathon", true)

	//HackathonFromHackathonID
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs("1").WillReturnRows(hackathonMockRows)

	//GetSubmissions
	hackathonMockRowsSubmissions := sqlmock.NewRows([]string{"hackathon_id", "user_id", "demo_url", "code_url", "score"}).
		AddRow(1, 1, "abc", "abc", 10).
		AddRow(2, 2, "def", "def", 20).
		AddRow(1, 3, "ijk", "ijk", 30)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `participants` WHERE hackathon_id = ? ")).WithArgs("1").WillReturnRows(hackathonMockRowsSubmissions)

	var user_id = 3
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE ID = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).WithArgs(user_id).WillReturnError(errors.New("internal server error"))

	//Notify All
	r.POST("/hackathons/:hackathon_id/notify", NotifyAll)
	req, _ := http.NewRequest("POST", "/hackathons/1/notify", strings.NewReader(jsonParams))
	r.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusInternalServerError, w.Code)

}

func TestNotifyAllSuccessError(t *testing.T) {

	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	//Authorization
	r.Use(func(c *gin.Context) {
		c.Set("authorization_payload", &token.Payload{
			Username:  "admin",
			IssuedAt:  time.Now(),
			ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
		})
	})

	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	ctx.Request.Header.Add("Accept", "application/json")

	//Parameters
	ctx.Params = []gin.Param{
		{
			Key:   "username",
			Value: "admin",
		},
		{
			Key:   "hackathon_id",
			Value: "1",
		},
	}

	//JSON params
	jsonParams := `{"subject": "modi15lakhdo", "message": "ye gormint bik chuki hai"}`

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
		Email:          "abefgh@gmail.com",
		LinkedIn:       "",
		GitHub:         "",
		WebLink:        "",
		Organisation:   "Razorpay",
		IsAdmin:        true,
	}

	//GetProfileByUsername
	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, user.Username, user.FullName, user.HashedPassword, user.Email, user.LinkedIn, user.GitHub, user.WebLink, user.Organisation, user.IsAdmin)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("admin").
		WillReturnRows(userRows)

	var created_at = time.Now()

	hackathonMockRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_time", "ending_time", "res" +
		"ult_time", "organisation_name", "organisation_id", "description", "admin_approved"}).AddRow("1", created_at, created_at, created_at, "Test Hackathon", created_at, created_at, created_at, "Winathon", "1", "It is a test Hackathon", true)

	//HackathonFromHackathonID
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs("1").WillReturnRows(hackathonMockRows)

	//GetSubmissions
	hackathonMockRowsSubmissions := sqlmock.NewRows([]string{"hackathon_id", "user_id", "demo_url", "code_url", "score"}).
		AddRow(1, 1, "abc", "abc", 10).
		AddRow(2, 2, "def", "def", 20).
		AddRow(1, 3, "ijk", "ijk", 30)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `participants` WHERE hackathon_id = ? ")).WithArgs("4").WillReturnRows(hackathonMockRowsSubmissions)

	//UserFromUserID
	var user_id = 1

	userRows2 := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, "name", "fullname", "password", "email", "abc", "abc", "abc", "abc", false)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE ID = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).WithArgs(user_id).WillReturnRows(userRows2)

	//GetProfileByUsername
	userRows1 := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, user.Username, user.FullName, user.HashedPassword, user.Email, user.LinkedIn, user.GitHub, user.WebLink, user.Organisation, user.IsAdmin)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("admin").
		WillReturnRows(userRows1)

	//Notify All
	r.POST("/hackathons/:hackathon_id/notify", NotifyAll)
	req, _ := http.NewRequest("POST", "/hackathons/1/notify", strings.NewReader(jsonParams))
	r.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusInternalServerError, w.Code)

}
