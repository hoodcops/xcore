package db

import (
	"time"

	"github.com/jmoiron/sqlx"
)

// UserProfile models the profile information
// of mobile users
type UserProfile struct {
	ID        int          `db:"id" json:"id"`
	UserID    int          `db:"userId" json:"userId"`
	Title     string       `db:"title" json:"title"`
	Fullname  string       `db:"fullname" json:"fullname"`
	Street    string       `db:"street" json:"street"`
	City      string       `db:"city" json:"city"`
	PostCode  string       `db:"post_code" json:"postCode"`
	GeoLng    string       `db:"geo_lng" json:"geoLng"`
	GeoLat    string       `db:"geo_lat" json:"geoLat"`
	CreatedAt time.Time    `db:"created_at" json:"createdAt"`
	UpdatedAt NullableTime `db:"updated_at" json:"updatedAt"`
}

// UserProfilesRepo defines methods for interacting with user
// profile records in the database
type UserProfilesRepo struct {
	db *sqlx.DB
}

// NewUserProfilesRepo returns a new user profiles repo
func NewUserProfilesRepo(db *sqlx.DB) *UserProfilesRepo {
	return &UserProfilesRepo{
		db: db,
	}
}

// Create saves a new user profile into the database, update the value
// of ID with auto-generated value from database, and returns the user
// profile or error if the operation fails
func (repo *UserProfilesRepo) Create(profile *UserProfile) (*UserProfile, error) {
	query := "INSERT INTO mobile_user_profiles (user_id, title, fullname, street, city, post_code, geo_lng, geo_lat) VALUES(?, ?, ?, ?, ?, ?, ?, ?)"
	res, err := repo.db.Exec(
		query,
		profile.UserID,
		profile.Title,
		profile.Fullname,
		profile.Street,
		profile.City,
		profile.PostCode,
		profile.GeoLng,
		profile.GeoLat,
	)

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	profile.ID = int(id)
	return profile, nil
}
