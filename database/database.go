package database

import (
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	"log"
)

// DB a connection to clickhouse
var DB *sql.DB

func init() {
	var configStr = fmt.Sprintf(
		"tcp://%s:%s?password=%s&database=%s",
		"192.168.9.39",
		"9000",
		"Wb922149@...S",
		"default",
	)

	var err error
	DB, err = sql.Open("clickhouse", configStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := DB.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println(err)
		}
		return
	}

	fmt.Printf("connect to clickhouse successful!")
}
