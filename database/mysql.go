package database

import (
	"fmt"
	"go-graphql-api/dbmodel"
	"go-graphql-api/util"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// a variable to store database connection
var (
	_db_name            = "Blog_Posts"
	DBInstance *gorm.DB = nil
)

// Var for error handling
var err error

func get_connection_url() (string, error) {
	user := util.EnvOrDefault("MYSQL_USER", "")
	pass := util.EnvOrDefault("MYSQL_PASS", "")
	host := util.EnvOrDefault("MYSQL_ADDR", "")

	if len(user) == 0 || len(pass) == 0 || len(host) == 0 {
		return "", fmt.Errorf(fmt.Sprintf(`mysql connection variables not set:
		user: len=%d, psss: len=%d, host: len=%d
		`, len(user), len(pass), len(host)))
	}

	opts := "charset=utf8&parseTime=True&loc=Local"
	return fmt.Sprintf("%s:%s@tcp(%s)/?%s", user, pass, host, opts), nil
}

// connecting to the db
func ConnectDB() error {
	connection_string, err := get_connection_url()
	if err != nil {
		return err
	}
	DBInstance, err = gorm.Open("mysql", connection_string)
	if err != nil {
		return err
	}

	max_idle_conn_str := util.EnvOrDefault("MYSQL_IDLE_CONNECTIONS", "10")
	max_idle_conn, err := strconv.Atoi(max_idle_conn_str)
	if err != nil {
		return err
	}
	max_open_conn_str := util.EnvOrDefault("MYSQL_OPEN_CONNECTIONS", "100")
	max_open_conn, err := strconv.Atoi(max_open_conn_str)
	if err != nil {
		return err
	}

	DBInstance.DB().SetMaxIdleConns(max_idle_conn)
	DBInstance.DB().SetMaxOpenConns(max_open_conn)

	// log all database operations performed by this connection
	DBInstance.LogMode(true)
	return nil
}

func CreateDB() {
	DBInstance.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", _db_name))
	DBInstance.Exec(fmt.Sprintf("USE %s", _db_name))
}

func MigrateDB() {
	fmt.Printf("Migrating %d model(s)\n", len(dbmodel.Models))
	for _, model := range dbmodel.Models {
		fmt.Printf("> Migrating model: %#v\n", model)
		DBInstance.AutoMigrate(model)
	}

	fmt.Println("Database migration completed....")
}
