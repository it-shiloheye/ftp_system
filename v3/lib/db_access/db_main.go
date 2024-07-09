package db

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"

	ftp_context "github.com/it-shiloheye/ftp_system/v3/lib/context"
	db_access "github.com/it-shiloheye/ftp_system/v3/lib/db_access/generated"
	"github.com/jackc/pgx/v5"
	// "github.com/jackc/pgx/v5/pgtype"
)

var DBPool *DBPoolStruct
var DB = &db_access.Queries{}

var db_url string

type DBPoolStruct struct {
	conns chan *pgx.Conn
	count atomic.Int32
}

func (dbp *DBPoolStruct) Len() int {
	return int(dbp.count.Load())
}

func (dbp *DBPoolStruct) GetConn() *pgx.Conn {
	return <-dbp.conns
}

func (dbp *DBPoolStruct) Return(pc *pgx.Conn) {
	dbp.conns <- pc
}

func (dbp *DBPoolStruct) PopulateConns(ctx ftp_context.Context, i int) {
	if dbp.conns == nil {
		dbp.conns = make(chan *pgx.Conn, i+1)
	} else {
		var tmp chan *pgx.Conn
		tmp, dbp.conns = dbp.conns, make(chan *pgx.Conn, int(dbp.count.Load())+i)
		dbp.conns = tmp
	}
	for ; i > 0; i -= 1 {
		db_conn, err1 := pgx.Connect(ctx, db_url)
		if err1 != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err1)
			os.Exit(1)
		}
		dbp.conns <- db_conn
		dbp.count.Add(1)
	}
}

func (dbp *DBPoolStruct) KillConns(ctx ftp_context.Context, i int) {
	for conn := range dbp.conns {
		conn.Close(ctx)
		i -= 1
		if i < 1 {
			break
		}
	}
}

func init() {
	db_url = os.Getenv("DATABASE_URL")
	if len(db_url) > 0 {
		return
	}

	db_fields_list := []string{
		"DATABASE_DBNAME", "DATABASE_HOST", "DATABASE_USER", "DATABASE_PASSWORD",
	}
	db_fields := map[string]string{}

	for _, db_field := range db_fields_list {
		db_field_let := os.Getenv(db_field)
		if len(db_field_let) < 1 {
			log.Fatalf(`Fatal: "%s" missing in .env`, db_field)
		}
		db_fields[db_field] = db_field_let

	}

	db_url = fmt.Sprintf("postgres://%s:%s@%s:5432/%s", db_fields["DATABASE_USER"], db_fields["DATABASE_PASSWORD"], db_fields["DATABASE_HOST"], db_fields["DATABASE_DBNAME"])
	if len(db_url) > 0 {
		log.Println("db_url: ", db_url)
	}

}

func ConnectToDB(ctx ftp_context.Context) {
	if DBPool == nil {

		DBPool = &DBPoolStruct{}
	}

	DBPool.PopulateConns(ctx, 10)

}
