package db

import (
	"time"

	"github.com/jmoiron/sqlx"
)

// UserContact models phone contacts that are uploaded by
// users to be contacted in case of emergency
type UserContact struct {
	ID        int       `db:"id" json:"id"`
	UserID    int       `db:"user_id" json:"userId"`
	Msisdn    string    `db:"msisdn" json:"msisdn"`
	Fullname  string    `db:"fullname" json:"fullname"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

// UserContactsRepo provides methods for interacting with user
// contacts in the database
type UserContactsRepo struct {
	db *sqlx.DB
}

// NewUserContactsRepo returns an instance of UserContactsRepo
func NewUserContactsRepo(db *sqlx.DB) *UserContactsRepo {
	return &UserContactsRepo{
		db: db,
	}
}

// CreateContacts inserts contacts into the database
func (repo *UserContactsRepo) CreateContacts(contacts []*UserContact) ([]*UserContact, error) {
	return contacts, nil
}

// GetAll returns all mobile user contacts in the database
func (repo *UserContactsRepo) GetAll() ([]*UserContact, error) {
	query := "SELECT * FROM mobile_user_contacts"
	var contacts []*UserContact

	err := repo.db.Select(&contacts, query)
	if err != nil {
		return nil, err
	}

	return contacts, nil
}
