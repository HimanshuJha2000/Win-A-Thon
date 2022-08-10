package main

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"win-a-thon/database"
	"win-a-thon/routes"
)

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func Test(t *testing.T) {
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

	// Grab our router
	router, err := routes.Setup()
	if err != nil {
		t.Fatal()
	}

	// Perform a GET request with that handler.
	w := performRequest(router, "GET", "/hackathons")
	assert.Equal(t, http.StatusOK, w.Code)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
	assert.Nil(t, err)

	expected := "{\"hackathons\":[{\"id\":1,\"title\":\"\",\"starting_time\":\"2021-08-02T00:00:00Z\",\"ending_time\":\"2021-08-02T00:00:00Z\",\"participants\":1}],\"status\":\"successful\"}"
	assert.Equal(t, w.Body.String(), expected)
}
