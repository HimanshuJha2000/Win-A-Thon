package controllers

import (
	"database/sql"
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

func TestUpdateSubmission(t *testing.T) {
	gin.SetMode(gin.TestMode)
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
		Username:       "tim",
		FullName:       "abcdef ghijkl",
		HashedPassword: "fjaskfjaslfj",
		Email:          "abc.defghi@gmail.com",
		LinkedIn:       "",
		GitHub:         "",
		WebLink:        "",
		Organisation:   "Razorpay",
		IsAdmin:        true,
	}

	userRows := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "username", "full_name", "hashed_password", "email", "linked_in", "git_hub", "web_link", "organisation", "is_admin"}).AddRow(time.Now(), time.Now(), nil, user.Username, user.FullName, user.HashedPassword, user.Email, user.LinkedIn, user.GitHub, user.WebLink, user.Organisation, user.IsAdmin)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).WithArgs("tim").WillReturnRows(userRows)

	participantMockRows := sqlmock.NewRows([]string{"hackathon_id", "user_id", "demo_url", "code_url", "score"}).AddRow(1, 0, "", "", 0)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `participants` WHERE hackathon_id = ? AND user_id = ? ORDER BY `participants`.`hackathon_id` LIMIT 1")).WithArgs("1", 0).WillReturnRows(participantMockRows)

	var created_at = time.Now()
	date := time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC)
	hackathonMockRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "titles", "starting_time", "ending_time", "result_time", "organisation_name", "organiser_id", "description", "admin_approved"}).AddRow("1", created_at, created_at, sql.NullTime{Time: time.Time{}, Valid: false}, "Test Hackathon", created_at, date, date, "Winathon", "1", "It is a test Hackathon", true)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `hackathons` WHERE (id = ? AND admin_approved = ?)")).WithArgs("1", true).WillReturnRows(hackathonMockRows)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `participants` WHERE hackathon_id = ?")).WithArgs("1").WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))
	participantMockRows2 := sqlmock.NewRows([]string{"hackathon_id", "user_id", "demo_url", "code_url", "score"}).AddRow(1, 0, "", "", 0)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `participants` WHERE hackathon_id = ? AND user_id = ? AND `participants`.`hackathon_id` = ? ORDER BY `participants`.`hackathon_id` LIMIT 1")).WithArgs(1, 0, 1).WillReturnRows(participantMockRows2)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `participants` (`hackathon_id`,`user_id`,`demo_url`,`code_url`,`score`) VALUES (?,?,?,?,?)")).WithArgs(1, 0, "abc", "xyz", 0).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	jsonParams := `{"demo_url":"abc","code_url":"xyz"}`

	r.Use(func(c *gin.Context) {
		c.Set("authorization_payload", &token.Payload{
			Username:  "tim",
			IssuedAt:  time.Now(),
			ExpiredAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
		})
	})
	r.PATCH("/:hackathon_id/submit", UpdateSubmission)
	req, _ := http.NewRequest("PATCH", "/1/submit", strings.NewReader(jsonParams))
	r.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusOK, w.Code)
}
