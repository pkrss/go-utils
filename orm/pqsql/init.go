package pqsql

import (
	"regexp"

	"github.com/pkrss/go-utils/profile"

	"github.com/go-pg/pg"
)

var Db *pg.DB

// func init() {

// 	initSql()
// }

func CreatePgSql() *pg.DB {

	dbuser := profile.ProfileReadString("MY_DB_PGSQL_USER")
	dbpsw := profile.ProfileReadString("MY_DB_PGSQL_PASSWORD")

	// jdbc:postgresql://domain:port/dbname?autoReconnect=true&useUnicode=true&characterEncoding=utf-8
	dbUrl := profile.ProfileReadString("MY_DB_PGSQL_URL")

	reg := regexp.MustCompile(`://([^:]+):(\d+)/([^\?]+)`) // (`://[^\:]+:\d+/[^\?]+`)
	ss := reg.FindStringSubmatch(dbUrl)

	dbHost := ss[1]
	dbPort := ss[2]
	dbName := ss[3]

	// dbName = "go"

	db := pg.Connect(&pg.Options{
		User:     dbuser,
		Password: dbpsw,
		Database: dbName,
		Addr:     dbHost + ":" + dbPort,
	})

	Db = db

	// var conninfo string = "user=" + dbuser + " password=" + dbpsw + " dbname=" + dbName + " host=" + dbHost + " port=" + dbPort + " sslmode=disable"
	// db, err := sql.Open("postgres", conninfo)
	// if err != nil {
	// 	panic(err)
	// }

	// Db = db

	return db
}
