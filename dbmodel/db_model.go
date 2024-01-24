package dbmodel

import "time"

type Post struct {
	ID           uint64 `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Title        string `gorm:"not null"`
	Content      string `gorm:"not null"`
	Author       string `gorm:"not null; unique"`
	Hero         string `json:"Hero"`
	Published_At string `json:"PublishedAt"`
	Updated_At   string `json:"UpdateAt"`
}

type UserType int

const (
	UserType_Normal = 0
	UserType_Admin  = 1
)

type User struct {
	ID         uint64       `sql:"AUTO_INCREMENT" gorm:"primaryKey"`
	Email      string       `gorm:"index;unique"`
	Password   string       `gorm:""`
	Type       UserType     `gorm:"default:0"`
	AuthTokens []OAuthToken `gorm:"foreignKey:UserId"`
}

type OAuthToken struct {
	ID           uint64    `sql:"AUTO_INCREMENT" gorm:"primaryKey"`
	Version      string    `gorm:"default:2"`
	Provider     string    `gorm:"not null"`
	AccessToken  string    `gorm:"not null"`
	RefreshToken string    `gorm:"not null"`
	Expiry       time.Time `gorm:"not null"`
	LastRefresh  time.Time `gorm:"not null"`
	UserId       uint64    `gorm:"index"`
}

// Models defined here will be auto migrated into the database
// when the application starts.
var Models = []interface{}{
	&User{},
	&OAuthToken{},
	&Post{},
}
