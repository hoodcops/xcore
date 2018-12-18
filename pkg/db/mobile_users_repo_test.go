package db

import (
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestMobileUsersRepo_Create_ShouldPass(t *testing.T) {
	sql := `^INSERT INTO mobile_users \(msisdn\) VALUES\(\?\)$`
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening database connection", err)
	}
	defer db.Close()

	user := MobileUser{
		Msisdn: "+233200662782",
	}

	mock.ExpectExec(sql).
		WithArgs(
			user.Msisdn,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	dbMock := sqlx.NewDb(db, "sqlmock")

	repo := NewMobileUsersRepo(dbMock)
	savedUser, err := repo.Create(&user)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if savedUser == nil {
		t.Fatalf("expected user, got nil")
	}

	if savedUser.ID != 1 {
		t.Fatalf("expected id %d, got %d", 1, savedUser.ID)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMobileUsersRepo_Create_ShouldFail(t *testing.T) {
	sql := `^INSERT INTO mobile_users \(msisdn\) VALUES\(\?\)$`
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening database connection", err)
	}
	defer db.Close()

	user := MobileUser{
		Msisdn: "+233200662782",
	}

	mock.ExpectExec(sql).
		WithArgs(
			user.Msisdn,
		).
		WillReturnError(errors.New("some database error"))

	dbMock := sqlx.NewDb(db, "sqlmock")

	repo := NewMobileUsersRepo(dbMock)
	savedUser, err := repo.Create(&user)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if savedUser != nil {
		t.Fatalf("expected nil, got %v", savedUser)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
