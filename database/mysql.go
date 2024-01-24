package database

import (
	"fmt"
	"go-graphql-api/dbmodel"
	"go-graphql-api/util"
	"go-graphql-api/util/logger"
	"strconv"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// a variable to store database connection
var (
	_db_name              = "Blog_Posts"
	_db_instance *gorm.DB = nil
	_db_init_mtx sync.Mutex
)

// Var for error handling
var err error

func get_connection_url() (string, error) {
	user := util.EnvOrDefault("MYSQL_USER", "")
	pass := util.EnvOrDefault("MYSQL_PASS", "")
	host := util.EnvOrDefault("MYSQL_HOST", "")

	if len(user) == 0 || len(pass) == 0 || len(host) == 0 {
		return "", fmt.Errorf(fmt.Sprintf(`mysql connection variables not set:
		user: len=%d, psss: len=%d, host: len=%d
		`, len(user), len(pass), len(host)))
	}

	opts := "charset=utf8&parseTime=True&loc=Local"
	return fmt.Sprintf("%s:%s@tcp(%s)/?%s", user, pass, host, opts), nil
}

// Get and return the datase instance.
// The database is initialized the first time this is called.
func GetDbInstance() (*gorm.DB, error) {
	_db_init_mtx.Lock()
	defer _db_init_mtx.Unlock()

	if _db_instance == nil {
		return init_database()
	} else {
		return _db_instance, nil
	}
}

func init_database() (*gorm.DB, error) {
	err := connect_db()
	if err != nil {
		return nil, err
	}
	create_db()
	migrate_db()
	return _db_instance, nil
}

func connect_db() error {
	connection_string, err := get_connection_url()
	if err != nil {
		return err
	}
	_db_instance, err = gorm.Open("mysql", connection_string)
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

	_db_instance.DB().SetMaxIdleConns(max_idle_conn)
	_db_instance.DB().SetMaxOpenConns(max_open_conn)

	// log all database operations performed by this connection
	_db_instance.LogMode(true)
	return nil
}

func create_db() {
	_db_instance.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", _db_name))
	_db_instance.Exec(fmt.Sprintf("USE %s", _db_name))
}

func migrate_db() {
	logger.Info("Migrating %d model(s)", len(dbmodel.Models))
	for _, model := range dbmodel.Models {
		logger.Info("> Migrating model: %#v", model)
		_db_instance.AutoMigrate(model)
	}

	logger.Info("Database migration completed....")
}
