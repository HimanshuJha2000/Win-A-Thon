package controllers

import (
	"fmt"
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

func TestGetLeaderboard(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

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

	participantMockRows := sqlmock.NewRows([]string{"hackathon_id", "user_id", "demo_url", "code_url", "score"}).
		AddRow(1, 1, "abc", "abc", 10)

	// For repo.GetLeaderboard
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `participants` WHERE hackathon_id = ? ORDER BY score desc")).WithArgs().WillReturnRows(participantMockRows)

	// For repo.UserFromUserID
	var created_at = time.Now()

	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, "name", "fullname", "password", "email", "abc", "abc", "abc", "abc", false)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE ID = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).WithArgs().WillReturnRows(userRows)

	// For repo.HackathonFromHackathonID
	hackathonMockRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_time", "ending_time", "res" +
		"ult_time", "organisation_name", "organisation_id", "description", "admin_approved"}).AddRow(1, created_at, created_at, created_at, "Test Hackathon", created_at, created_at, created_at, "Winathon", 1, "It is a test Hackathon", true)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs().WillReturnRows(hackathonMockRows)

	GetLeaderboard(ctx)
	if ctx.Writer.Status() != http.StatusOK {
		t.Fatal("Wrong error code")
	}

	expected := "{\"leaderboard\":[{\"user_name\":\"name\",\"full_name\":\"fullname\",\"score\":10}]}"

	ctx.Writer.Flush()

	assert.Equal(t, w.Body.String(), expected)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestGetSubmissions(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Params = []gin.Param{
		{
			Key:   "hackathon_id",
			Value: "1",
		},
	}

	ctx.Keys = make(map[string]interface{})
	ctx.Keys["authorization_payload"] = &token.Payload{
		Username:  "name",
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
		Username:       "name",
		FullName:       "abcdef ghijkl",
		HashedPassword: "fjaskfjaslfj",
		Email:          "abc.defghi@gmail.com",
		LinkedIn:       "",
		GitHub:         "",
		WebLink:        "",
		Organisation:   "Razorpay",
		IsAdmin:        true,
	}

	created_at := time.Now()

	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, "name", user.FullName, "password", "email", "abc", "abc", "abc", "abc", false)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs().
		WillReturnRows(userRows)

	// For repo.HackathonFromHackathonID
	hackathonMockRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_time", "ending_time", "res" +
		"ult_time", "organisation_name", "organisation_id", "description", "admin_approved"}).AddRow(1, created_at, created_at, created_at, "Test Hackathon", created_at, created_at, created_at, "Winathon", 1, "It is a test Hackathon", true)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs().WillReturnRows(hackathonMockRows)

	// For repo.GetSubmission
	participantMockRows := sqlmock.NewRows([]string{"hackathon_id", "user_id", "demo_url", "code_url", "score"}).
		AddRow(1, 1, "abc", "abc", 10)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `participants` WHERE hackathon_id = ? ")).WithArgs().WillReturnRows(participantMockRows)

	// For repo.UserFromUserID
	var user_id = 1

	userRows = sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, "name", "fullname", "password", "email", "abc", "abc", "abc", "abc", false)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE ID = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).WithArgs(user_id).WillReturnRows(userRows)

	GetSubmissions(ctx)
	if ctx.Writer.Status() != http.StatusOK {
		t.Fatal("Wrong error code")
	}

	expected := `{"submissions":[{"hackathon_id":1,"username":"name","code_url":"abc","demo_url":"abc"}],"total_submissions":1}`
	ctx.Writer.Flush()

	assert.Equal(t, w.Body.String(), expected)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestGetSubmissionsUserNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Params = []gin.Param{
		{
			Key:   "hackathon_id",
			Value: "1",
		},
	}

	ctx.Keys = make(map[string]interface{})
	ctx.Keys["authorization_payload"] = &token.Payload{
		Username:  "name",
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
		WithArgs().
		WillReturnRows(sqlmock.NewRows(nil))

	GetSubmissions(ctx)
	if ctx.Writer.Status() != http.StatusBadRequest {
		t.Fatal("Wrong error code")
	}

	expected := `{"message":"User doesn't exist"}`
	ctx.Writer.Flush()

	assert.Equal(t, w.Body.String(), expected)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestGetSubmissionsHackathonNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Params = []gin.Param{
		{
			Key:   "hackathon_id",
			Value: "1",
		},
	}

	ctx.Keys = make(map[string]interface{})
	ctx.Keys["authorization_payload"] = &token.Payload{
		Username:  "name",
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
		Username:       "name",
		FullName:       "abcdef ghijkl",
		HashedPassword: "fjaskfjaslfj",
		Email:          "abc.defghi@gmail.com",
		LinkedIn:       "",
		GitHub:         "",
		WebLink:        "",
		Organisation:   "Razorpay",
		IsAdmin:        true,
	}

	//created_at := time.Now()

	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, "name", user.FullName, "password", "email", "abc", "abc", "abc", "abc", false)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs().
		WillReturnRows(userRows)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs().WillReturnRows(sqlmock.NewRows(nil))

	GetSubmissions(ctx)
	if ctx.Writer.Status() != http.StatusNotFound {
		t.Fatal("Wrong error code")
	}

	expected := `{"message":"This hackathon id doesnt exist"}`
	ctx.Writer.Flush()

	//fmt.Println(w.Body.String())
	assert.Equal(t, w.Body.String(), expected)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestGetSubmissionsUnauthorized(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Params = []gin.Param{
		{
			Key:   "hackathon_id",
			Value: "1",
		},
	}

	ctx.Keys = make(map[string]interface{})
	ctx.Keys["authorization_payload"] = &token.Payload{
		Username:  "notorganiser",
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
		Username:       "name",
		FullName:       "abcdef ghijkl",
		HashedPassword: "fjaskfjaslfj",
		Email:          "abc.defghi@gmail.com",
		LinkedIn:       "",
		GitHub:         "",
		WebLink:        "",
		Organisation:   "Razorpay",
		IsAdmin:        true,
	}

	created_at := time.Now()

	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, "name", user.FullName, "password", "email", "abc", "abc", "abc", "abc", false)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs().
		WillReturnRows(userRows)

	// For repo.HackathonFromHackathonID
	hackathonMockRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_time", "ending_time", "res" +
		"ult_time", "organisation_name", "organiser_id", "description", "admin_approved"}).AddRow(1, created_at, created_at, created_at, "Test Hackathon", created_at, created_at, created_at, "Winathon", 10, "It is a test Hackathon", true)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs().WillReturnRows(hackathonMockRows)

	GetSubmissions(ctx)
	if ctx.Writer.Status() != http.StatusUnauthorized {
		t.Fatal("Wrong error code")
	}

	expected := `{"message":"Not authorised to view submissions"}`
	ctx.Writer.Flush()

	//fmt.Println(w.Body.String())
	assert.Equal(t, w.Body.String(), expected)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestGetSubmissionsTimeConstraint(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Params = []gin.Param{
		{
			Key:   "hackathon_id",
			Value: "1",
		},
	}

	ctx.Keys = make(map[string]interface{})
	ctx.Keys["authorization_payload"] = &token.Payload{
		Username:  "notorganiser",
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
		Username:       "name",
		FullName:       "abcdef ghijkl",
		HashedPassword: "fjaskfjaslfj",
		Email:          "abc.defghi@gmail.com",
		LinkedIn:       "",
		GitHub:         "",
		WebLink:        "",
		Organisation:   "Razorpay",
		IsAdmin:        true,
	}

	created_at := time.Now()

	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, "name", user.FullName, "password", "email", "abc", "abc", "abc", "abc", false)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs().
		WillReturnRows(userRows)

	// For repo.HackathonFromHackathonID
	end_time := time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC)

	hackathonMockRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_time", "ending_time", "res" +
		"ult_time", "organisation_name", "organiser_id", "description", "admin_approved"}).AddRow(1, created_at, created_at, created_at, "Test Hackathon", created_at, end_time, end_time, "Winathon", 0, "It is a test Hackathon", true)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs().WillReturnRows(hackathonMockRows)

	GetSubmissions(ctx)
	if ctx.Writer.Status() != http.StatusUnauthorized {
		t.Fatal("Wrong error code")
	}

	expected := `{"message":"Can't view submissions right now. Try again once the hackathon ends!!"}`

	ctx.Writer.Flush()

	//fmt.Println(w.Body.String())
	assert.Equal(t, w.Body.String(), expected)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestGetSubmissionsUser(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Params = []gin.Param{
		{
			Key:   "hackathon_id",
			Value: "1",
		},
	}

	ctx.Keys = make(map[string]interface{})
	ctx.Keys["authorization_payload"] = &token.Payload{
		Username:  "name",
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
		Username:       "name",
		FullName:       "abcdef ghijkl",
		HashedPassword: "fjaskfjaslfj",
		Email:          "abc.defghi@gmail.com",
		LinkedIn:       "",
		GitHub:         "",
		WebLink:        "",
		Organisation:   "Razorpay",
		IsAdmin:        true,
	}

	created_at := time.Now()

	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, "name", user.FullName, "password", "email", "abc", "abc", "abc", "abc", false)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs().
		WillReturnRows(userRows)

	// For repo.HackathonFromHackathonID
	hackathonMockRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_time", "ending_time", "res" +
		"ult_time", "organisation_name", "organisation_id", "description", "admin_approved"}).AddRow(1, created_at, created_at, created_at, "Test Hackathon", created_at, created_at, created_at, "Winathon", 1, "It is a test Hackathon", true)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs().WillReturnRows(hackathonMockRows)

	//For repo.GetSubmission
	participantMockRows := sqlmock.NewRows([]string{"hackathon_id", "user_id", "demo_url", "code_url", "score"}).
		AddRow(1, 2, "abc", "abc", 10)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `participants` WHERE hackathon_id = ? ")).WithArgs("1").WillReturnRows(participantMockRows)

	// For repo.UserFromUserID
	var user_id = 2

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE ID = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).WithArgs(user_id).WillReturnRows(sqlmock.NewRows(nil))

	GetSubmissions(ctx)
	if ctx.Writer.Status() != http.StatusInternalServerError {
		t.Fatal("Wrong error code")
	}

	ctx.Writer.Flush()

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestGetSubmissionOfParticipant(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Params = []gin.Param{
		{
			Key:   "hackathon_id",
			Value: "1",
		},
		{
			Key:   "username",
			Value: "name",
		},
	}
	ctx.Keys = make(map[string]interface{})
	ctx.Keys["authorization_payload"] = &token.Payload{
		Username:  "name",
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
		Username:       "name",
		FullName:       "abcdef ghijkl",
		HashedPassword: "fjaskfjaslfj",
		Email:          "abc.defghi@gmail.com",
		LinkedIn:       "",
		GitHub:         "",
		WebLink:        "",
		Organisation:   "Razorpay",
		IsAdmin:        true,
	}

	created_at := time.Now()

	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, "name", user.FullName, "password", "email", "abc", "abc", "abc", "abc", false)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs().
		WillReturnRows(userRows)

	// For repo.HackathonFromHackathonID
	hackathonMockRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_time", "ending_time", "res" +
		"ult_time", "organisation_name", "organisation_id", "description", "admin_approved"}).AddRow(1, created_at, created_at, created_at, "Test Hackathon", created_at, created_at, created_at, "Winathon", 1, "It is a test Hackathon", true)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs().WillReturnRows(hackathonMockRows)

	//user from user name
	userRows = sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, "name", user.FullName, "password", "email", "abc", "abc", "abc", "abc", false)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs().
		WillReturnRows(userRows)

	//participant from user id
	participantMockRows := sqlmock.NewRows([]string{"hackathon_id", "user_id", "demo_url", "code_url", "score"}).
		AddRow(1, 1, "abc", "abc", 10)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `participants` WHERE hackathon_id = ? AND user_id = ?")).WithArgs().WillReturnRows(participantMockRows)

	GetSubmissionOfParticipant(ctx)

	if ctx.Writer.Status() != http.StatusOK {
		t.Fatal("Wrong error code")
	}

	expected := `{"code_url":"abc","demo_url":"abc","hackathon_id":1,"user name":"name"}`
	ctx.Writer.Flush()

	assert.Equal(t, w.Body.String(), expected)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestGetSubmissionOfParticipantUserNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Params = []gin.Param{
		{
			Key:   "hackathon_id",
			Value: "1",
		},
		{
			Key:   "username",
			Value: "name",
		},
	}
	ctx.Keys = make(map[string]interface{})
	ctx.Keys["authorization_payload"] = &token.Payload{
		Username:  "name",
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
		WithArgs("name").
		WillReturnRows(sqlmock.NewRows(nil))

	GetSubmissionOfParticipant(ctx)

	if ctx.Writer.Status() != http.StatusBadRequest {
		t.Fatal("Wrong error code")
	}

	expected := `{"message":"User doesn't exist"}`
	ctx.Writer.Flush()

	assert.Equal(t, w.Body.String(), expected)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestGetSubmissionOfParticipantHackathonNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Params = []gin.Param{
		{
			Key:   "hackathon_id",
			Value: "1",
		},
		{
			Key:   "username",
			Value: "name",
		},
	}
	ctx.Keys = make(map[string]interface{})
	ctx.Keys["authorization_payload"] = &token.Payload{
		Username:  "name",
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

	//For repo.GetUserByUsername
	user := &models.User{
		Username:       "name",
		FullName:       "abcdef ghijkl",
		HashedPassword: "fjaskfjaslfj",
		Email:          "abc.defghi@gmail.com",
		LinkedIn:       "",
		GitHub:         "",
		WebLink:        "",
		Organisation:   "Razorpay",
		IsAdmin:        true,
	}

	//created_at := time.Now()

	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, "name", user.FullName, "password", "email", "abc", "abc", "abc", "abc", false)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("name").
		WillReturnRows(userRows)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs().WillReturnRows(sqlmock.NewRows(nil))

	GetSubmissionOfParticipant(ctx)

	if ctx.Writer.Status() != http.StatusNotFound {
		t.Fatal("Wrong error code")
	}

	expected := `{"message":"This hackathon id doesnt exist"}`
	ctx.Writer.Flush()

	//fmt.Println(w.Body.String())
	assert.Equal(t, w.Body.String(), expected)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestGetSubmissionOfParticipantUnauthorized(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Params = []gin.Param{
		{
			Key:   "hackathon_id",
			Value: "1",
		},
		{
			Key:   "username",
			Value: "name",
		},
	}
	ctx.Keys = make(map[string]interface{})
	ctx.Keys["authorization_payload"] = &token.Payload{
		Username:  "name",
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

	//For repo.GetUserByUsername
	user := &models.User{
		Username:       "name",
		FullName:       "abcdef ghijkl",
		HashedPassword: "fjaskfjaslfj",
		Email:          "abc.defghi@gmail.com",
		LinkedIn:       "",
		GitHub:         "",
		WebLink:        "",
		Organisation:   "Razorpay",
		IsAdmin:        true,
	}

	created_at := time.Now()

	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, "name", user.FullName, "password", "email", "abc", "abc", "abc", "abc", false)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("name").
		WillReturnRows(userRows)

	//For repo.HackathonFromHackathonID
	hackathonMockRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_time", "ending_time", "res" +
		"ult_time", "organisation_name", "organiser_id", "description", "admin_approved"}).AddRow(1, created_at, created_at, created_at, "Test Hackathon", created_at, created_at, created_at, "Winathon", 10, "It is a test Hackathon", true)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs().WillReturnRows(hackathonMockRows)

	GetSubmissionOfParticipant(ctx)

	if ctx.Writer.Status() != http.StatusUnauthorized {
		t.Fatal("Wrong error code")
	}

	expected := `{"message":"Not authorised to view submission!"}`
	ctx.Writer.Flush()

	assert.Equal(t, w.Body.String(), expected)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestGetSubmissionOfParticipantTimeConstraint(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Params = []gin.Param{
		{
			Key:   "hackathon_id",
			Value: "1",
		},
		{
			Key:   "username",
			Value: "name",
		},
	}
	ctx.Keys = make(map[string]interface{})
	ctx.Keys["authorization_payload"] = &token.Payload{
		Username:  "name",
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

	//For repo.GetUserByUsername
	user := &models.User{
		Username:       "name",
		FullName:       "abcdef ghijkl",
		HashedPassword: "fjaskfjaslfj",
		Email:          "abc.defghi@gmail.com",
		LinkedIn:       "",
		GitHub:         "",
		WebLink:        "",
		Organisation:   "Razorpay",
		IsAdmin:        true,
	}

	created_at := time.Now()

	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, "name", user.FullName, "password", "email", "abc", "abc", "abc", "abc", false)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("name").
		WillReturnRows(userRows)

	//For repo.HackathonFromHackathonID
	end_time := time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC)

	hackathonMockRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_time", "ending_time", "res" +
		"ult_time", "organisation_name", "organisation_id", "description", "admin_approved"}).AddRow(1, created_at, created_at, created_at, "Test Hackathon", created_at, end_time, end_time, "Winathon", 1, "It is a test Hackathon", true)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs().WillReturnRows(hackathonMockRows)

	GetSubmissionOfParticipant(ctx)

	if ctx.Writer.Status() != http.StatusUnauthorized {
		t.Fatal("Wrong error code")
	}

	expected := `{"message":"Can't view the submission right now. Try again once the hackathon ends!!"}`
	ctx.Writer.Flush()

	assert.Equal(t, w.Body.String(), expected)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestGetSubmissionOfParticipantUserNotFound2(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Params = []gin.Param{
		{
			Key:   "hackathon_id",
			Value: "1",
		},
		{
			Key:   "username",
			Value: "notfound",
		},
	}
	ctx.Keys = make(map[string]interface{})
	ctx.Keys["authorization_payload"] = &token.Payload{
		Username:  "name",
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

	//For repo.GetUserByUsername
	user := &models.User{
		Username:       "name",
		FullName:       "abcdef ghijkl",
		HashedPassword: "fjaskfjaslfj",
		Email:          "abc.defghi@gmail.com",
		LinkedIn:       "",
		GitHub:         "",
		WebLink:        "",
		Organisation:   "Razorpay",
		IsAdmin:        true,
	}

	created_at := time.Now()

	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, "name", user.FullName, "password", "email", "abc", "abc", "abc", "abc", false)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("name").
		WillReturnRows(userRows)

	//For repo.HackathonFromHackathonID
	hackathonMockRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_time", "ending_time", "res" +
		"ult_time", "organisation_name", "organisation_id", "description", "admin_approved"}).AddRow(1, created_at, created_at, created_at, "Test Hackathon", created_at, created_at, created_at, "Winathon", 1, "It is a test Hackathon", true)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs().WillReturnRows(hackathonMockRows)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("notfound").
		WillReturnRows(sqlmock.NewRows(nil))

	GetSubmissionOfParticipant(ctx)

	if ctx.Writer.Status() != http.StatusBadRequest {
		t.Fatal("Wrong error code")
	}

	expected := `{"message":"Username doesn't exist"}`
	ctx.Writer.Flush()

	assert.Equal(t, w.Body.String(), expected)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestPrizeDistribution(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Params = []gin.Param{
		{
			Key:   "hackathon_id",
			Value: "1",
		},
	}
	ctx.Keys = make(map[string]interface{})
	ctx.Keys["authorization_payload"] = &token.Payload{
		Username:  "name",
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
		Username:       "name",
		FullName:       "abcdef ghijkl",
		HashedPassword: "fjaskfjaslfj",
		Email:          "abc.defghi@gmail.com",
		LinkedIn:       "",
		GitHub:         "",
		WebLink:        "",
		Organisation:   "Razorpay",
		IsAdmin:        true,
	}

	created_at := time.Now()

	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, "name", user.FullName, "password", "email", "abc", "abc", "abc", "abc", false)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs().
		WillReturnRows(userRows)

	//hackathon from hackathon id
	hackathonMockRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_time", "ending_time", "res" +
		"ult_time", "organisation_name", "organisation_id", "description", "admin_approved"}).AddRow("1", created_at, created_at, created_at, "Test Hackathon", created_at, created_at, created_at, "Winathon", "1", "It is a test Hackathon", "1")

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs().WillReturnRows(hackathonMockRows)

	// For repo.GetWinners
	participantMockRows := sqlmock.NewRows([]string{"hackathon_id", "user_id", "demo_url", "code_url", "score"}).
		AddRow(1, 1, "abc", "abc", 10)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `participants` WHERE hackathon_id = ?  ORDER BY score desc LIMIT 3")).WithArgs("1").WillReturnRows(participantMockRows)

	// For repo.UserFromUserID
	var user_id = 1

	userRows = sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, "name", "fullname", "password", "team.winathon@gmail.com", "abc", "abc", "abc", "abc", false)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE ID = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).WithArgs(user_id).WillReturnRows(userRows)

	PrizeDistribution(ctx)
	if ctx.Writer.Status() != http.StatusInternalServerError {
		t.Fatal("Wrong error code")
	}

	expected := `{"message":"server unavailable"}`

	ctx.Writer.Flush()

	assert.Equal(t, w.Body.String(), expected)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestPrizeDistributionBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Params = []gin.Param{
		{
			Key:   "hackathon_id",
			Value: "2",
		},
	}
	ctx.Keys = make(map[string]interface{})
	ctx.Keys["authorization_payload"] = &token.Payload{
		Username:  "name",
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
		Username:       "name",
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
		AddRow(time.Now(), time.Now(), nil, "name", user.FullName, "password", "email", "abc", "abc", "abc", "abc", false)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs().
		WillReturnRows(userRows)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs("2").WillReturnRows(sqlmock.NewRows(nil))

	PrizeDistribution(ctx)
	if ctx.Writer.Status() != http.StatusNotFound {
		t.Fatal("Wrong error code")
	}

	expected := `{"message":"This hackathon id doesnt exist"}`

	ctx.Writer.Flush()
	//fmt.Println(w.Body.String())
	assert.Equal(t, w.Body.String(), expected)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestPrizeDistributionUserNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Params = []gin.Param{
		{
			Key:   "hackathon_id",
			Value: "1",
		},
	}
	ctx.Keys = make(map[string]interface{})
	ctx.Keys["authorization_payload"] = &token.Payload{
		Username:  "name",
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
		WithArgs("name").
		WillReturnRows(sqlmock.NewRows(nil))

	PrizeDistribution(ctx)
	if ctx.Writer.Status() != http.StatusBadRequest {
		t.Fatal("Wrong error code")
	}

	expected := `{"message":"User doesn't exist"}`

	ctx.Writer.Flush()
	fmt.Println(w.Body.String())
	assert.Equal(t, w.Body.String(), expected)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}

}

func TestPrizeDistributionUnauthorised(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Params = []gin.Param{
		{
			Key:   "hackathon_id",
			Value: "1",
		},
	}
	ctx.Keys = make(map[string]interface{})
	ctx.Keys["authorization_payload"] = &token.Payload{
		Username:  "name",
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

	//For repo.GetUserByUsername
	user := &models.User{
		Username:       "name",
		FullName:       "abcdef ghijkl",
		HashedPassword: "fjaskfjaslfj",
		Email:          "abc.defghi@gmail.com",
		LinkedIn:       "",
		GitHub:         "",
		WebLink:        "",
		Organisation:   "Razorpay",
		IsAdmin:        true,
	}

	created_at := time.Now()

	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, "name", user.FullName, "password", "email", "abc", "abc", "abc", "abc", false)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs("name").
		WillReturnRows(userRows)

	//hackathon from hackathon id
	hackathonMockRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_time", "ending_time", "res" +
		"ult_time", "organisation_name", "organiser_id", "description", "admin_approved"}).AddRow("1", created_at, created_at, created_at, "Test Hackathon", created_at, created_at, created_at, "Winathon", 10, "It is a test Hackathon", "1")

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs("1").WillReturnRows(hackathonMockRows)

	PrizeDistribution(ctx)
	if ctx.Writer.Status() != http.StatusUnauthorized {
		t.Fatal("Wrong error code")
	}

	expected := `{"message":"Not authorised! You are not the organiser/judge."}`

	ctx.Writer.Flush()
	assert.Equal(t, w.Body.String(), expected)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestPrizeDistributionTimeConstraint(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Params = []gin.Param{
		{
			Key:   "hackathon_id",
			Value: "1",
		},
	}
	ctx.Keys = make(map[string]interface{})
	ctx.Keys["authorization_payload"] = &token.Payload{
		Username:  "name",
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
		Username:       "name",
		FullName:       "abcdef ghijkl",
		HashedPassword: "fjaskfjaslfj",
		Email:          "abc.defghi@gmail.com",
		LinkedIn:       "",
		GitHub:         "",
		WebLink:        "",
		Organisation:   "Razorpay",
		IsAdmin:        true,
	}

	created_at := time.Now()

	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).
		AddRow(time.Now(), time.Now(), nil, "name", user.FullName, "password", "email", "abc", "abc", "abc", "abc", false)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
		WithArgs().
		WillReturnRows(userRows)

	//hackathon from hackathon id
	end_time := time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC)

	hackathonMockRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_time", "ending_time", "res" +
		"ult_time", "organisation_name", "organisation_id", "description", "admin_approved"}).AddRow("1", created_at, created_at, created_at, "Test Hackathon", created_at, end_time, end_time, "Winathon", "1", "It is a test Hackathon", "1")

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `hackathons` WHERE ID = ? AND `hackathons`.`deleted_at` IS NULL ORDER BY `hackathons`.`id` LIMIT 1")).WithArgs().WillReturnRows(hackathonMockRows)

	PrizeDistribution(ctx)
	if ctx.Writer.Status() != http.StatusUnauthorized {
		t.Fatal("Wrong error code")
	}

	expected := `{"message":"Results aren't declared yet!"}`

	ctx.Writer.Flush()

	assert.Equal(t, w.Body.String(), expected)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}
