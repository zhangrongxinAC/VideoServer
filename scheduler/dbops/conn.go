package dbops

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbConn *sql.DB
	err    error
)

func init() {
	dbConn, err = sql.Open("mysql", "ubuntu:zrxAC1999&@/video_server?charset=utf8")
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("dbConn +%v\n", dbConn)
}
