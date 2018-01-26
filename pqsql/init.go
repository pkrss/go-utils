package pqsql

import (
	"database/sql"
	"os"
	"regexp"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func init() {

	initSql()
}

func initSql() {
	dbuser := os.Getenv("MY_DB_PGSQL_USER")
	dbpsw := os.Getenv("MY_DB_PGSQL_PASSWORD")

	// jdbc:postgresql://rm-bp1waaqgdq6929vswo.pg.rds.aliyuncs.com:3432/sx98?autoReconnect=true&useUnicode=true&characterEncoding=utf-8
	dbUrl := os.Getenv("MY_DB_PGSQL_URL")

	reg := regexp.MustCompile(`://([^:]+):(\d+)/([^\?]+)`) // (`://[^\:]+:\d+/[^\?]+`)
	ss := reg.FindStringSubmatch(dbUrl)

	dbHost := ss[1]
	dbPort := ss[2]
	dbName := ss[3]

	// dbName = "go"

	var conninfo string = "user=" + dbuser + " password=" + dbpsw + " dbname=" + dbName + " host=" + dbHost + " port=" + dbPort + " sslmode=disable"
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		panic(err)
	}

	Db = db
}
