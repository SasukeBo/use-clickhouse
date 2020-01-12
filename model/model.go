package model

import (
	"database/sql"
	"fmt"
	"github.com/SasukeBo/use-clickhouse/database"
	"reflect"
	"strings"
)

// Sale model
type Sale struct {
	Region        string  `column:"region"`
	Country       string  `column:"country"`
	ItemType      string  `column:"item_type"`
	SalesChannel  string  `column:"sales_channel"`
	OrderPriority string  `column:"order_priority"`
	OrderID       string  `column:"order_id"`
	UnitsSold     uint16  `column:"units_sold"`
	UnitPrice     float32 `column:"unit_price"`
	UnitCost      float32 `column:"unit_cost"`
	TotalRevenue  float32 `column:"total_revenue"`
	TotalCost     float32 `column:"total_cost"`
	TotalProfit   float32 `column:"total_profit"`
}

// SimpleQuery _
func SimpleQuery(filters []interface{}, fields []interface{}, limit, offset int) (interface{}, error) {
	selectField := parseSelectField(fields)
	whereField := parseWhereField(filters)

	sql := fmt.Sprintf(
		"SELECT %s FROM sales %sLIMIT %d OFFSET %d",
		selectField,
		whereField,
		limit,
		offset,
	)

	rows, err := database.DB.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sales := scanValue(rows, fields)

	sql = fmt.Sprintf("SELECT count() FROM sales %s", whereField)

	rows, err = database.DB.Query(sql)
	if err != nil {
		return nil, err
	}

	total := 0

	for rows.Next() {
		if err := rows.Scan(&total); err != nil {
			continue
		}
	}

	return struct {
		Total int
		Sales []Sale
	}{total, sales}, nil
}

func scanValue(rows *sql.Rows, fields []interface{}) []Sale {
	sales := []Sale{}
	sale := Sale{}
	rts := reflect.TypeOf(sale)

	names := []string{}

	for _, field := range fields {
		_, ok := rts.FieldByName(field.(string))
		if ok {
			names = append(names, field.(string))
		}
	}

	for rows.Next() {
		s := Sale{}
		rvs := reflect.ValueOf(&s)

		values := make([]interface{}, len(names))
		valueps := make([]interface{}, len(names))

		for i := range values {
			valueps[i] = &values[i]
		}

		if err := rows.Scan(valueps...); err != nil {
			continue
		}

		for i, v := range values {
			value := reflect.ValueOf(v)
			rvs.Elem().FieldByName(names[i]).Set(value)
		}

		sales = append(sales, s)
	}

	return sales
}

func parseWhereField(filters []interface{}) string {
	s := Sale{}
	rts := reflect.TypeOf(s)

	conditions := []string{}

	for _, item := range filters {
		filter := item.(map[string]interface{})
		fieldName, ok := rts.FieldByName(filter["field"].(string))

		if ok {
			conditions = append(conditions, fmt.Sprintf("position(lcase(%s), lcase('%v')) > 0", fieldName.Tag.Get("column"), filter["value"]))
		}
	}

	if len(conditions) == 0 {
		return ""
	}

	return fmt.Sprintf("WHERE %s ", strings.Join(conditions, " and "))
}

func parseSelectField(fields []interface{}) string {
	s := Sale{}
	rts := reflect.TypeOf(s)

	columns := []string{}

	for _, field := range fields {
		fieldName, ok := rts.FieldByName(field.(string))
		if ok {
			columns = append(columns, fieldName.Tag.Get("column"))
		}
	}

	if len(columns) == 0 {
		return "*"
	}

	return strings.Join(columns, ", ")
}
