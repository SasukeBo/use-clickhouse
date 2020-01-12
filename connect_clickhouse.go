package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ClickHouse/clickhouse-go"
)

func connectClickhouse() {
	var host = "192.168.9.39"
	var port = "9000"
	var password = "Wb922149@...S"
	var database = "default"

	var configStr = fmt.Sprintf(
		"tcp://%s:%s?password=%s&database=%s",
		host,
		port,
		password,
		database,
	)

	connect, err := sql.Open("clickhouse", configStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println(err)
		}
		return
	}

	rows, err := connect.Query("SELECT region, country, order_id, total_cost FROM sales LIMIT 10")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			region    string
			country   string
			orderID   string
			totalCost float32
		)
		if err := rows.Scan(&region, &country, &orderID, &totalCost); err != nil {
			log.Fatal(err)
		}
		log.Printf("region: %s, country: %s, order id: %s, total cost: %v", region, country, orderID, totalCost)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
