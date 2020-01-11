## README

### Answers

#### Task 1

> Load Data from CSV into Clickhouse with your own sample real or fake data(1 million rows+)

- 连接数据库

```sh
clickhouse-client --password *******
```

- 创建表

```SQL
CREATE TABLE IF NOT EXISTS sales
(
  region String,
  country String,
  item_type String,
  sales_channel FixedString(7),
  order_priority FixedString(1),
  order_id FixedString(9),
  units_sold UInt16,
  unit_price Float32,
  unit_cost Float32,
  total_revenue Float32,
  total_cost Float32,
  total_profit Float32
)
ENGINE = Memory
```

- 导入样本数据

[1 Million Sales Records](./1MillionSalesRecords.csv)

```sh
clickhouse-client --query="INSERT INTO sales FORMAT CSV" < 1MillionSalesRecords.csv --password ******
```

- 查询样本数据

```SQL
SELECT * FROM sales LIMIT 100
```

Output

```
#=> 100 rows in set. Elapsed: 0.005 sec. Processed 1.00 million rows, 99.30 MB (201.91 million rows/s., 20.05 GB/s.)
```

#### Task 2

> Connect to Clickhouse; Use Golang Clickhouse Driver

[./connect_clickhouse.go](./connect_clickhouse.go)

```go
package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ClickHouse/clickhouse-go"
)

func main() {
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

```

```sh
go run connect_clickhouse.go
```

#### Task 3

> Provide Golang APIs to query data(Restful or GraphQL)

- 启动服务

```sh
go run main.go
```

- 测试 API

浏览器打开 [localhost:8080/api](http://localhost:8080/api)

query

```graphql
query {
  ping
}
```

output

```json
{
  "data": {
    "ping": "pong"
  }
}
```
