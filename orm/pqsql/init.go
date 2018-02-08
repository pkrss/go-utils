package pqsql

import (
	"log"
	"regexp"

	"github.com/pkrss/go-utils/profile"

	"github.com/go-pg/pg"
)

var Db *pg.DB

// func init() {

// 	initSql()
// }

type ConnectOption struct {
	User     string
	Password string
	Database string
	SrvHost  string
	SrvPort  string
}

func ConnectOptionFromProfile() *ConnectOption {
	dbuser := profile.ProfileReadString("MY_DB_PGSQL_USER")
	dbpsw := profile.ProfileReadString("MY_DB_PGSQL_PASSWORD")

	// jdbc:postgresql://domain:port/dbname?autoReconnect=true&useUnicode=true&characterEncoding=utf-8
	dbUrl := profile.ProfileReadString("MY_DB_PGSQL_URL")

	reg := regexp.MustCompile(`://([^:]+):(\d+)/([^\?]+)`) // (`://[^\:]+:\d+/[^\?]+`)
	ss := reg.FindStringSubmatch(dbUrl)

	dbHost := ss[1]
	dbPort := ss[2]
	dbName := ss[3]

	return &ConnectOption{User: dbuser,
		Password: dbpsw,
		Database: dbName,
		SrvHost:  dbHost,
		SrvPort:  dbPort,
	}
}

func CreatePgSql(opts ...*ConnectOption) *pg.DB {
	var opt *ConnectOption
	if len(opts) > 0 {
		opt = opts[0]
	}
	if opt == nil {
		opt = ConnectOptionFromProfile()
	}

	log.Printf("pgsql connect: %s:%s/%s\n", opt.SrvHost, opt.SrvPort, opt.Database)

	db := pg.Connect(&pg.Options{
		User:     opt.User,
		Password: opt.Password,
		Database: opt.Database,
		Addr:     opt.SrvHost + ":" + opt.SrvPort,
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
