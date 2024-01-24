package dbmodel

type Post struct {
	ID           uint64 `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Title        string `gorm:"not null"`
	Content      string `gorm:"not null"`
	Author       string `gorm:"not null; unique"`
	Hero         string `json:"Hero"`
	Published_At string `json:"PublishedAt"`
	Updated_At   string `json:"UpdateAt"`
}

// Models defined here will be auto migrated into the database
// when the application starts.
var Models = []interface{}{
	&Post{},
}
