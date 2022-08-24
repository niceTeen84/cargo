package component

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

var Db *sqlx.DB

func init() {

	db, err := sqlx.Open("mysql", "root:123456qq@tcp(127.0.0.1:13306)/test")
	if err != nil {
		log.Fatal("init db failed", err)
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(100)
	db.SetConnMaxLifetime(30 * time.Second)
	err = db.Ping()
	if err != nil {
		log.Fatal("ping db failed", err)
	}
	Db = db
}

func SingleQuery() {
	var ret string
	row := Db.QueryRow("select version()")
	row.Scan(&ret)
	fmt.Println(ret)
}
