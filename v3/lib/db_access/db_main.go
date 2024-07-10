package db

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"

	ftp_context "github.com/it-shiloheye/ftp_system/v3/lib/context"
	db_access "github.com/it-shiloheye/ftp_system/v3/lib/db_access/generated"
	"github.com/it-shiloheye/ftp_system/v3/lib/logging"
	"github.com/it-shiloheye/ftp_system/v3/lib/logging/log_item"
	"github.com/jackc/pgx/v5"
	// "github.com/jackc/pgx/v5/pgtype"
)

var DBPool = NewDBPool()
var DB = &db_access.Queries{}

var Logger = logging.Logger

var db_url string

type DBPoolStruct struct {
	conns chan *pgx.Conn
	count atomic.Int32

	retry_chan chan int
}

func (dbp *DBPoolStruct) Len() int {
	return int(dbp.count.Load())
}

// panics if connection is closed
// checks if connection is nil
// makes 3 reconnect attempts in 3mins
func (dbp *DBPoolStruct) GetConn(ctx ftp_context.Context) *pgx.Conn {
	loc := log_item.Loc(`func (dbp *DBPoolStruct) GetConn(ctx ftp_context.Context) *pgx.Conn`)
	var err1 error
	pc := <-dbp.conns
	if pc == nil || pc.IsClosed() {
		dbp.count.Add(-1)
		Logger.LogErr(loc, &log_item.LogItem{
			Message: "connection is closed",
		})
		pc, err1 = dbp.retry_connections(ctx, 3)
		if err1 != nil {
			dbp.SignalConnLost()
			panic(Logger.LogErr(loc, err1))
		}
		return pc
	}
	return pc
}

// returns a connection to the connection pool
// if connection is nill or closed it retries to connect 3x
func (dbp *DBPoolStruct) Return(pc *pgx.Conn, ctx ftp_context.Context) {
	loc := log_item.Loc(`func (dbp *DBPoolStruct) Return(pc *pgx.Conn, ctx ftp_context.Context)`)
	var err1 error
	if pc == nil || pc.IsClosed() {
		dbp.count.Add(-1)
		Logger.LogErr(loc, &log_item.LogItem{
			Message: "connection is closed",
		})
		pc, err1 = dbp.retry_connections(ctx, 3)
		if err1 != nil {
			Logger.LogErr(loc, err1)
			dbp.SignalConnLost()
			return
		}

	}
	dbp.conns <- pc
}

// attempts to form connection to database;
// defaults to an infinite loop of re
func (dbp *DBPoolStruct) retry_connections(ctx ftp_context.Context, tries ...int) (pc *pgx.Conn, err1 error) {
	loc := log_item.Loc(`func (dbp *DBPoolStruct) retry_connections(ctx ftp_context.Context)(pc *pgx.Conn)`)
	l := -1

	if len(tries) > 0 {
		l = tries[0]
	}

	for ; l >= 0; l -= 1 {
		if a, b := ctx.NearDeadline(time.Millisecond); a && b {
			return
		} else {
			pc, err1 = dbp.create_connection(ctx)
			if err1 != nil {
				Logger.LogErr(loc, err1)
				select {
				case <-time.After(time.Minute):

				case <-ctx.Done():
					return
				}

			} else {
				break
			}
			Logger.Logf(loc, "attempting to reconnect to db")
		}
	}

	return
}

func (dbp *DBPoolStruct) create_connection(ctx ftp_context.Context) (*pgx.Conn, error) {
	db_conn, err1 := pgx.Connect(ctx, db_url)
	if err1 != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err1)
		return nil, err1
	}

	dbp.count.Add(1)
	return db_conn, nil
}

func (dbp *DBPoolStruct) PopulateConns(ctx ftp_context.Context, i int) {
	c := i
	if dbp.conns == nil {
		dbp.conns = make(chan *pgx.Conn, i+1)
	} else {
		var tmp chan *pgx.Conn
		tmp, dbp.conns = dbp.conns, make(chan *pgx.Conn, int(dbp.count.Load())+i)
		close(tmp)
		for c := range tmp {
			dbp.conns <- c
		}
	}
	for ; i > 0; i -= 1 {
		db_conn, err1 := dbp.retry_connections(ctx)
		if err1 != nil {
			<-time.After(time.Second)
			i += 1
			continue

		}
		dbp.conns <- db_conn
	}
	log.Println("added ", c, "connections to db")
}

func (dbp *DBPoolStruct) KillConns(ctx ftp_context.Context, i int) {
	close(dbp.conns)
	close(dbp.retry_chan)
	for conn := range dbp.conns {
		conn.Close(ctx)
		i -= 1
		if i < 1 {
			break
		}
	}
}

const (
	env_db_name     string = "DATABASE_DBNAME"
	env_db_host     string = "DATABASE_HOST"
	env_db_user     string = "DATABASE_USER"
	env_db_password string = "DATABASE_PASSWORD"
)

func init() {
	db_url = os.Getenv("DATABASE_URL")
	if len(db_url) > 0 {
		return
	}

	db_fields_list := []string{
		env_db_name, env_db_host, env_db_user, env_db_password,
	}
	db_fields := map[string]string{}

	for _, db_field := range db_fields_list {
		db_field_let := os.Getenv(db_field)
		if len(db_field_let) < 1 {
			log.Fatalf(`Fatal: "%s" missing in .env`, db_field)
		}
		db_fields[db_field] = db_field_let

	}

	db_url = fmt.Sprintf("postgres://%s:%s@%s:5432/%s", db_fields[env_db_user], db_fields[env_db_password], db_fields[env_db_host], db_fields[env_db_name])
	if len(db_url) > 0 {
		log.Println("db_url: ", db_url)
	}

}

func ConnectToDB(ctx ftp_context.Context) {
	log.Println("connecting to db")

	DBPool.PopulateConns(ctx, 10)
	log.Println("successfully connected to db")
}

func NewDBPool() (dbp *DBPoolStruct) {
	dbp = &DBPoolStruct{
		retry_chan: make(chan int),
	}

	go func() {
		for range dbp.retry_chan {
			pc, err1 := dbp.retry_connections(ftp_context.CreateNewContext())
			if err1 != nil {
				return
			}

			dbp.conns <- pc
		}

	}()

	return
}
func (dbp *DBPoolStruct) SignalConnLost() {

	dbp.retry_chan <- 1
}
