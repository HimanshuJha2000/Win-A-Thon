package controllers

import (
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	mySQL "github.com/go-sql-driver/mysql"
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

func TestListHackathons(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	date := time.Date(2021, 8, 2, 0, 0, 0, 0, time.UTC)

	hackathonMockRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_" +
		"time", "ending_time", "result_time", "organisation_name", "organiser_id", "description", "admin_" +
		"approved"}).AddRow(1, date, date, nil, "Test Hackathon", date, date, date, "Winathon", 1, "It is a test"+
		" Hackathon", true)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE admin_approved = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY starting_" +
			"time desc")).WithArgs(true).WillReturnRows(hackathonMockRow)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT count(*) FROM `participants` WHERE hackathon_id = ?")).WithArgs("1").WillReturnRows(sqlmock.NewRows(
		[]string{"count(*)"}).AddRow(1))

	ListHackathons(ctx)
	if ctx.Writer.Status() != http.StatusOK {
		t.Fatal("Wrong error code")
	}

	expected := "{\"hackathons\":[{\"id\":1,\"title\":\"\",\"starting_time\":\"2021-08-02T00:00:00Z\",\"ending_time\":\"2021-08-02T00:00:00Z\",\"participants\":1}],\"status\":\"successful\"}"

	ctx.Writer.Flush()

	assert.Equal(t, w.Body.String(), expected)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestListHackathonsRepoError1(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE admin_approved = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY starting_" +
			"time desc")).WithArgs(true).WillReturnError(errors.New("Custom Repo Error"))

	ListHackathons(ctx)
	if ctx.Writer.Status() != http.StatusInternalServerError {
		t.Fatal("Wrong error code")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestListHackathonsEmptyList(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE admin_approved = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY starting_" +
			"time desc")).WithArgs(true).WillReturnRows(sqlmock.NewRows(nil))

	ListHackathons(ctx)
	if ctx.Writer.Status() != http.StatusOK {
		t.Fatal("Wrong error code")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestListHackathonsRepoError2(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	date := time.Date(2021, 8, 2, 0, 0, 0, 0, time.UTC)

	hackathonMockRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_" +
		"time", "ending_time", "result_time", "organisation_name", "organiser_id", "description", "admin_" +
		"approved"}).AddRow(1, date, date, nil, "Test Hackathon", date, date, date, "Winathon", 1, "It is a test"+
		" Hackathon", true)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE admin_approved = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY starting_" +
			"time desc")).WithArgs(true).WillReturnRows(hackathonMockRow)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT count(*) FROM `participants` WHERE hackathon_id = ?")).WithArgs("1").WillReturnError(errors.New("Custom error"))

	ListHackathons(ctx)
	if ctx.Writer.Status() != http.StatusInternalServerError {
		t.Fatal("Wrong error code")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestListUnapprovedHackathons(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
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

	// For repo.GetUserByUsername
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

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("admin").
		WillReturnRows(userRows)

	// For repo.ListHackathons
	date := time.Date(2021, 8, 2, 0, 0, 0, 0, time.UTC)

	hackathonMockRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_" +
		"time", "ending_time", "result_time", "organisation_name", "organiser_id", "description", "admin_" +
		"approved"}).AddRow(1, date, date, nil, "Test Hackathon", date, date, date, "Winathon", int(user.ID), "It is a test"+
		" Hackathon", false)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE admin_approved = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY starting_" +
			"time desc")).WithArgs(false).WillReturnRows(hackathonMockRow)

	ListUnapprovedHackathons(ctx)
	if ctx.Writer.Status() != http.StatusOK {
		t.Fatal("Wrong error code")
	}

	expected := "{\"hackathons\":[{\"id\":1,\"title\":\"\",\"starting_time\":\"2021-08-02T00:00:00Z\",\"ending_time\":\"2021-08-02T00:00:00Z\",\"result_time\":\"2021-08-02T00:00:00Z\",\"organisation_name\":\"Winathon\",\"description\":\"It is a test Hackathon\"}],\"status\":\"successful\"}"

	ctx.Writer.Flush()

	fmt.Println(w.Body.String())
	assert.Equal(t, w.Body.String(), expected)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestListUnapprovedHackathonsRepoError1(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
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

	// For repo.GetUserByUsername
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("admin").
		WillReturnError(errors.New("Custom Errror"))

	ListUnapprovedHackathons(ctx)
	if ctx.Writer.Status() != http.StatusInternalServerError {
		t.Fatal("Wrong error code")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestListUnapprovedHackathonsRepoError2(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
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

	// For repo.GetUserByUsername
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

	ListUnapprovedHackathons(ctx)
	if ctx.Writer.Status() != http.StatusBadRequest {
		t.Fatal("Wrong error code")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestListUnapprovedHackathonsRepoError3(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
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

	// For repo.GetUserByUsername
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

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("admin").
		WillReturnRows(userRows)

	// For repo.ListHackathons
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE admin_approved = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY starting_" +
			"time desc")).WithArgs(false).WillReturnError(errors.New("Custom Error"))

	ListUnapprovedHackathons(ctx)
	if ctx.Writer.Status() != http.StatusInternalServerError {
		t.Fatal("Wrong error code")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestListUnapprovedHackathonsRepoError4(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
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

	// For repo.GetUserByUsername
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

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("admin").
		WillReturnRows(userRows)

	// For repo.ListHackathons
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE admin_approved = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY starting_" +
			"time desc")).WithArgs(false).WillReturnRows(sqlmock.NewRows(nil))

	ListUnapprovedHackathons(ctx)
	if ctx.Writer.Status() != http.StatusOK {
		t.Fatal("Wrong error code")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestCreateHackathon(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	// For repo.GetUserByUsername
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
	fmt.Println(user.ID)
	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, user.Username, user.FullName, user.HashedPassword, user.Email, user.LinkedIn, user.GitHub, user.WebLink, user.Organisation, user.IsAdmin)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("admin").
		WillReturnRows(userRows)

	// For repo.CreateHackathon
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `hackathons` (`created_at`,`updated_at`,`deleted_at`,`title`,`starting_time`,`ending_time`,`result_time`,`organisation_name`,`organiser_id`,`description`,`admin_approved`) VALUES (?,?,?,?,?,?,?,?,?,?,?)")).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	r.Use(func(c *gin.Context) {
		c.Set("authorization_payload", &token.Payload{
			Username:  "admin",
			IssuedAt:  time.Now(),
			ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
		})
	})

	jsonStr := `{
		"title" : "debangan1",
		"starting_time" : "2021-08-02T00:00:00+05:30",
		"ending_time" : "2021-08-02T01:00:00+05:30",
		"result_time" : "2021-08-02T02:00:00+05:30",
		"description" : "ABC",
		"organisation_name" : "ABCD"
	}`
	r.POST("/hackathons", CreateHackathon)
	req, _ := http.NewRequest("POST", "/hackathons", strings.NewReader(jsonStr))
	r.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusOK, w.Code)
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestCreateHackathonBadInput(t *testing.T) {
	inputStrings := []string{`{
			"title" : "",
			"starting_time" : "2021-08-02T00:00:00+05:30",
			"ending_time" : "2021-08-02T01:00:00+05:30",
			"result_time" : "2021-08-02T02:00:00+05:30",
			"description" : "ABC",
			"organisation_name" : "ABCD"
		}`, `{
			"title" : "debangan1",
			"starting_time" : "2021-08-02T00:00:00+05:30",
			"ending_time" : "2021-08-02T01:00:00+05:30",
			"result_time" : "2021-08-02T02:00:00+05:30",
			"description" : "",
			"organisation_name" : "ABCD"
		}`, `{
			"title" : "debangan1",
			"starting_time" : "2021-08-02T00:00:00+05:30",
			"ending_time" : "2021-07-02T00:00:00+05:30",
			"result_time" : "2021-08-02T00:00:00+05:30",
			"description" : "ABC",
			"organisation_name" : "ABCD"
		}`, `{
			"title" : "debangan1",
			"starting_time" : "2021-08-02T00:00:00+05:30",
			"ending_time" : "2021-08-02T00:00:00+05:30",
			"result_time" : "2021-07-02T00:00:00+05:30",
			"description" : "ABC",
			"organisation_name" : "ABCD"
		}`, `{
			"title" : "debangan1",
			"starting_time" : "2021-08-02T00:00:00+05:30",
			"ending_time" : "2021-08-02T00:00:00+05:30",
			"result_time" : "2021-08-02T01:00:00+05:30",
			"description" : "ABC",
			"organisation_name" : "ABCD"
		}`, `{
			"title" : "debangan1",
			"starting_time" : "2021-08-02T00:00:00+05:30",
			"ending_time" : "2021-08-02T01:00:00+05:30",
			"result_time" : "2021-08-02T01:00:00+05:30",
			"description" : "ABC",
			"organisation_name" : "ABCD"
		}`}
	for _, input := range inputStrings {
		t.Run("Test", func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, r := gin.CreateTestContext(w)
			ctx.Request = &http.Request{
				URL:    &url.URL{},
				Header: make(http.Header),
			}

			driver, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening the database connection", err)
			}

			// For repo.GetUserByUsername
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
			fmt.Println(user.ID)
			userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
				AddRow(time.Now(), time.Now(), nil, user.Username, user.FullName, user.HashedPassword, user.Email, user.LinkedIn, user.GitHub, user.WebLink, user.Organisation, user.IsAdmin)

			mock.ExpectQuery(regexp.QuoteMeta(
				"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
				WithArgs("admin").
				WillReturnRows(userRows)

			// Testing
			r.Use(func(c *gin.Context) {
				c.Set("authorization_payload", &token.Payload{
					Username:  "admin",
					IssuedAt:  time.Now(),
					ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
				})
			})

			r.POST("/hackathons", CreateHackathon)
			req, _ := http.NewRequest("POST", "/hackathons", strings.NewReader(input))
			r.ServeHTTP(w, req)

			assert.EqualValues(t, http.StatusBadRequest, w.Code)
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Failed to meet expectations, got error: %v", err)
			}
		})
	}
}

func TestCreateHackathonRepoError1(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	// For repo.GetUserByUsername
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
	fmt.Println(user.ID)
	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, user.Username, user.FullName, user.HashedPassword, user.Email, user.LinkedIn, user.GitHub, user.WebLink, user.Organisation, user.IsAdmin)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("admin").
		WillReturnRows(userRows)

	// For repo.CreateHackathon
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `hackathons` (`created_at`,`updated_at`,`deleted_at`,`title`,`starting_time`,`ending_time`,`result_time`,`organisation_name`,`organiser_id`,`description`,`admin_approved`) VALUES (?,?,?,?,?,?,?,?,?,?,?)")).
		WithArgs().
		WillReturnError(&mySQL.MySQLError{Number: utils.DuplicateRecordErrorCode})
	mock.ExpectRollback()

	r.Use(func(c *gin.Context) {
		c.Set("authorization_payload", &token.Payload{
			Username:  "admin",
			IssuedAt:  time.Now(),
			ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
		})
	})

	jsonStr := `{
		"title" : "debangan1",
		"starting_time" : "2021-08-02T00:00:00+05:30",
		"ending_time" : "2021-08-02T01:00:00+05:30",
		"result_time" : "2021-08-02T02:00:00+05:30",
		"description" : "ABC",
		"organisation_name" : "ABCD"
	}`
	r.POST("/hackathons", CreateHackathon)
	req, _ := http.NewRequest("POST", "/hackathons", strings.NewReader(jsonStr))
	r.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusBadRequest, w.Code)
}

func TestCreateHackathonRepoError2(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	// For repo.GetUserByUsername
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
	fmt.Println(user.ID)
	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, user.Username, user.FullName, user.HashedPassword, user.Email, user.LinkedIn, user.GitHub, user.WebLink, user.Organisation, user.IsAdmin)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("admin").
		WillReturnRows(userRows)

	// For repo.CreateHackathon
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `hackathons` (`created_at`,`updated_at`,`deleted_at`,`title`,`starting_time`,`ending_time`,`result_time`,`organisation_name`,`organiser_id`,`description`,`admin_approved`) VALUES (?,?,?,?,?,?,?,?,?,?,?)")).
		WithArgs().
		WillReturnError(errors.New("Custom Error"))
	mock.ExpectRollback()

	r.Use(func(c *gin.Context) {
		c.Set("authorization_payload", &token.Payload{
			Username:  "admin",
			IssuedAt:  time.Now(),
			ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
		})
	})

	jsonStr := `{
		"title" : "debangan1",
		"starting_time" : "2021-08-02T00:00:00+05:30",
		"ending_time" : "2021-08-02T01:00:00+05:30",
		"result_time" : "2021-08-02T02:00:00+05:30",
		"description" : "ABC",
		"organisation_name" : "ABCD"
	}`
	r.POST("/hackathons", CreateHackathon)
	req, _ := http.NewRequest("POST", "/hackathons", strings.NewReader(jsonStr))
	r.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusBadRequest, w.Code)
}

func TestGetParticipants(t *testing.T) {
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
	}

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	// For repo.GetUserByUsername
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

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs(user.Username).
		WillReturnRows(userRows)

	// For repo.HackathonFromHackathonID
	var created_at = time.Now()

	hackathonMockRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_time", "ending_time", "res" +
		"ult_time", "organisation_name", "organisation_id", "description", "admin_approved"}).AddRow(1, created_at, created_at, created_at, "Test Hackathon", created_at, created_at, created_at, "Winathon", user.ID, "It is a test Hackathon", true)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs("1").WillReturnRows(hackathonMockRows)

	// For repo.GetSubmissions
	participantMockRows := sqlmock.NewRows([]string{"hackathon_id", "user_id", "demo_url", "code_url", "score"}).
		AddRow(1, 1, "abc", "abc", 10)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `participants` WHERE hackathon_id = ? ")).WithArgs("1").WillReturnRows(participantMockRows)

	// For repo.CountParticipants
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT count(*) FROM `participants` WHERE hackathon_id = ?")).WithArgs("1").WillReturnRows(sqlmock.NewRows(
		[]string{"count(*)"}).AddRow(3))

	// For repo.UserFromUserID
	var user_id = 1
	userRowsNew := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, "name", "fullname", "password", "email", "abc", "abc", "abc", "abc", false)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE ID = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).WithArgs(user_id).WillReturnRows(userRowsNew)

	GetParticipants(ctx)
	if ctx.Writer.Status() != http.StatusOK {
		t.Fatal("Wrong error code")
	}

	expected := `{"participants":[{"id":0,"username":"name","full_name":"fullname","email":"email","linked_in":"abc","github":"abc","web_link":"abc","organisation":"abc"}],"total_participants":3}`
	ctx.Writer.Flush()

	assert.Equal(t, w.Body.String(), expected)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestGetParticipantsRepoError1(t *testing.T) {
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
	}

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	// For repo.GetUserByUsername
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("admin").
		WillReturnError(errors.New("Custom Error"))

	GetParticipants(ctx)
	if ctx.Writer.Status() != http.StatusInternalServerError {
		t.Fatal("Wrong error code")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestGetParticipantsRepoError2(t *testing.T) {
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
	}

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	// For repo.GetUserByUsername
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

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs(user.Username).
		WillReturnRows(userRows)

	// For repo.HackathonFromHackathonID
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs("1").WillReturnError(errors.New("Custom Error"))

	GetParticipants(ctx)
	if ctx.Writer.Status() != http.StatusBadRequest {
		t.Fatal("Wrong error code")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestGetParticipantsRepoError3(t *testing.T) {
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
	}

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	// For repo.GetUserByUsername
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

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs(user.Username).
		WillReturnRows(userRows)

	// For repo.HackathonFromHackathonID
	var created_at = time.Now()

	hackathonMockRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_time", "ending_time", "res" +
		"ult_time", "organisation_name", "organisation_id", "description", "admin_approved"}).AddRow(1, created_at, created_at, created_at, "Test Hackathon", created_at, created_at, created_at, "Winathon", user.ID, "It is a test Hackathon", true)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs("1").WillReturnRows(hackathonMockRows)

	// For repo.GetSubmissions
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `participants` WHERE hackathon_id = ? ")).WithArgs("1").WillReturnError(errors.New("Custom Error"))

	GetParticipants(ctx)
	if ctx.Writer.Status() != http.StatusInternalServerError {
		t.Fatal("Wrong error code")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestGetParticipantsRepoError4(t *testing.T) {
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
	}

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	// For repo.GetUserByUsername
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

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs(user.Username).
		WillReturnRows(userRows)

	// For repo.HackathonFromHackathonID
	var created_at = time.Now()

	hackathonMockRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_time", "ending_time", "res" +
		"ult_time", "organisation_name", "organisation_id", "description", "admin_approved"}).AddRow(1, created_at, created_at, created_at, "Test Hackathon", created_at, created_at, created_at, "Winathon", user.ID, "It is a test Hackathon", true)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs("1").WillReturnRows(hackathonMockRows)

	// For repo.GetSubmissions
	participantMockRows := sqlmock.NewRows([]string{"hackathon_id", "user_id", "demo_url", "code_url", "score"}).
		AddRow(1, 1, "abc", "abc", 10)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `participants` WHERE hackathon_id = ? ")).WithArgs("1").WillReturnRows(participantMockRows)

	// For repo.CountParticipants
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT count(*) FROM `participants` WHERE hackathon_id = ?")).WithArgs("1").
		WillReturnError(errors.New("Custom Error"))

	GetParticipants(ctx)
	if ctx.Writer.Status() != http.StatusInternalServerError {
		t.Fatal("Wrong error code")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestGetParticipantsRepoError5(t *testing.T) {
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
	}

	driver, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	database.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: driver, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening the database connection", err)
	}

	// For repo.GetUserByUsername
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

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs(user.Username).
		WillReturnRows(userRows)

	// For repo.HackathonFromHackathonID
	var created_at = time.Now()

	hackathonMockRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_time", "ending_time", "res" +
		"ult_time", "organisation_name", "organisation_id", "description", "admin_approved"}).AddRow(1, created_at, created_at, created_at, "Test Hackathon", created_at, created_at, created_at, "Winathon", user.ID, "It is a test Hackathon", true)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs("1").WillReturnRows(hackathonMockRows)

	// For repo.GetSubmissions
	participantMockRows := sqlmock.NewRows([]string{"hackathon_id", "user_id", "demo_url", "code_url", "score"}).
		AddRow(1, 1, "abc", "abc", 10)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `participants` WHERE hackathon_id = ? ")).WithArgs("1").WillReturnRows(participantMockRows)

	// For repo.CountParticipants
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT count(*) FROM `participants` WHERE hackathon_id = ?")).WithArgs("1").WillReturnRows(sqlmock.NewRows(
		[]string{"count(*)"}).AddRow(3))

	// For repo.UserFromUserID
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE ID = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).WithArgs(1).
		WillReturnError(errors.New("Custom Error"))

	GetParticipants(ctx)
	if ctx.Writer.Status() != http.StatusInternalServerError {
		t.Fatal("Wrong error code")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestListOrganisedHackathons(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
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

	// For repo.GetUserByUsername
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

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("admin").
		WillReturnRows(userRows)

	// For repo.ListOrganisedHackathons
	date := time.Date(2021, 8, 4, 0, 0, 0, 0, time.UTC)

	hackathonMockRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_" +
		"time", "ending_time", "result_time", "organisation_name", "organiser_id", "description", "admin_" +
		"approved"}).AddRow(1, date, date, nil, "Test Hackathon", date, date, date, "Winathon", int(user.ID), "It is a test"+
		" Hackathon", true)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE (organiser_id = ?) AND `hackathons`.`deleted_at` IS NULL ORDER BY starting_time desc")).
		WithArgs(int(user.ID)).
		WillReturnRows(hackathonMockRow)

	// For repo.Countparticipants
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT count(*) FROM `participants` WHERE hackathon_id = ?")).WithArgs("1").WillReturnRows(sqlmock.NewRows(
		[]string{"count(*)"}).AddRow(1))

	ListOrganisedHackathons(ctx)
	if ctx.Writer.Status() != http.StatusOK {
		t.Fatal("Wrong error code")
	}

	expected := "{\"organised_hackathons\":[{\"id\":1,\"title\":\"\",\"starting_time\":\"2021-08-04T00:00:00Z\",\"ending_time\":\"2021-08-04T00:00:00Z\",\"result_time\":\"2021-08-04T00:00:00Z\",\"description\":\"It is a test Hackathon\",\"organisation_name\":\"Winathon\",\"admin_approved\":true,\"participants\":1}],\"status\":\"successful\"}"

	ctx.Writer.Flush()

	assert.Equal(t, w.Body.String(), expected)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestListOrganisedHackathonsRepoError1(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
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

	// For repo.GetUserByUsername
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("admin").
		WillReturnError(errors.New("Custom Error"))

	ListOrganisedHackathons(ctx)
	if ctx.Writer.Status() != http.StatusInternalServerError {
		t.Fatal("Wrong error code")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestListOrganisedHackathonsRepoError2(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
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

	// For repo.GetUserByUsername
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

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("admin").
		WillReturnRows(userRows)

	// For repo.ListOrganisedHackathons
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE (organiser_id = ?) AND `hackathons`.`deleted_at` IS NULL ORDER BY starting_time desc")).
		WithArgs(int(user.ID)).
		WillReturnError(errors.New("Custom Error"))

	ListOrganisedHackathons(ctx)
	if ctx.Writer.Status() != http.StatusInternalServerError {
		t.Fatal("Wrong error code")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestListOrganisedHackathonsRepoError3(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
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

	// For repo.GetUserByUsername
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

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("admin").
		WillReturnRows(userRows)

	// For repo.ListOrganisedHackathons
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE (organiser_id = ?) AND `hackathons`.`deleted_at` IS NULL ORDER BY starting_time desc")).
		WithArgs(int(user.ID)).
		WillReturnRows(sqlmock.NewRows(nil))

	ListOrganisedHackathons(ctx)
	if ctx.Writer.Status() != http.StatusOK {
		t.Fatal("Wrong error code")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestListOrganisedHackathonsRepoError4(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
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

	// For repo.GetUserByUsername
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

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("admin").
		WillReturnRows(userRows)

	// For repo.ListOrganisedHackathons
	date := time.Date(2021, 8, 4, 0, 0, 0, 0, time.UTC)

	hackathonMockRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_" +
		"time", "ending_time", "result_time", "organisation_name", "organiser_id", "description", "admin_" +
		"approved"}).AddRow(1, date, date, nil, "Test Hackathon", date, date, date, "Winathon", int(user.ID), "It is a test"+
		" Hackathon", true)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE (organiser_id = ?) AND `hackathons`.`deleted_at` IS NULL ORDER BY starting_time desc")).
		WithArgs(int(user.ID)).
		WillReturnRows(hackathonMockRow)

	// For repo.Countparticipants
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT count(*) FROM `participants` WHERE hackathon_id = ?")).WithArgs("1").
		WillReturnError(errors.New("Custom Error"))

	ListOrganisedHackathons(ctx)
	if ctx.Writer.Status() != http.StatusInternalServerError {
		t.Fatal("Wrong error code")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}
