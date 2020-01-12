package schema

import (
	// "fmt"
	"github.com/graphql-go/graphql"
	"log"

	"github.com/SasukeBo/use-clickhouse/model"
)

// Schema _
var Schema graphql.Schema

var (
	gString  = graphql.String
	gNString = graphql.NewNonNull(graphql.String)
	gInt     = graphql.Int
	gNInt    = graphql.NewNonNull(graphql.Int)
	gFloat   = graphql.Float
)

// create graphql argument config
// gt type
// dv defautValue
// opts[0] description
func arg(gt graphql.Input, dv interface{}, opts ...string) *graphql.ArgumentConfig {
	des := ""
	if len(opts) > 0 {
		des = opts[0]
	}

	return &graphql.ArgumentConfig{
		Type:         gt,
		Description:  des,
		DefaultValue: dv,
	}
}

func init() {
	var err error
	Schema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query: RootQueries,
	})
	if err != nil {
		log.Fatal(err)
	}
}

// RootQueries _
var RootQueries = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"ping":       ping,
		"simple":     simpleQuery,
		"aggregated": aggregatedQuery,
	},
})

var ping = &graphql.Field{
	Type: gString,
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		return "pong", nil
	},
}

var simpleQuery = &graphql.Field{
	Type: saleList,
	Args: graphql.FieldConfigArgument{
		"limit":   arg(gNInt, nil, "max return rows"),
		"offset":  arg(gInt, 0, "query offset"),
		"fields":  arg(graphql.NewList(allowedQueryField), nil, "query fields"),
		"filters": arg(graphql.NewList(simpleQueryFilter), nil, "query filters"),
	},
	Description: "query simple data by filters, and select fields",
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		fields, ok := params.Args["fields"].([]interface{})
		if !ok {
			fields = []interface{}{}
		}
		filters, ok := params.Args["filters"].([]interface{})
		if !ok {
			filters = []interface{}{}
		}
		limit := params.Args["limit"].(int)
		offset := params.Args["offset"].(int)
		return model.SimpleQuery(filters, fields, limit, offset)
	},
}

var simpleQueryFilter = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "SimpleQueryFilter",
	Fields: graphql.InputObjectConfigFieldMap{
		"field": &graphql.InputObjectFieldConfig{Type: allowedFilterField},
		"value": &graphql.InputObjectFieldConfig{Type: gString},
	},
})

var allowedQueryField = graphql.NewEnum(graphql.EnumConfig{
	Name: "AllowedQueryField",
	Values: graphql.EnumValueConfigMap{
		"Region":        &graphql.EnumValueConfig{Value: "Region"},
		"Country":       &graphql.EnumValueConfig{Value: "Country"},
		"ItemType":      &graphql.EnumValueConfig{Value: "ItemType"},
		"SalesChannel":  &graphql.EnumValueConfig{Value: "SalesChannel"},
		"OrderPriority": &graphql.EnumValueConfig{Value: "OrderPriority"},
		"OrderID":       &graphql.EnumValueConfig{Value: "OrderID"},
		"UnitsSold":     &graphql.EnumValueConfig{Value: "UnitsSold"},
		"UnitPrice":     &graphql.EnumValueConfig{Value: "UnitPrice"},
		"UnitCost":      &graphql.EnumValueConfig{Value: "UnitCost"},
		"TotalRevenue":  &graphql.EnumValueConfig{Value: "TotalRevenue"},
		"TotalCost":     &graphql.EnumValueConfig{Value: "TotalCost"},
		"TotalProfit":   &graphql.EnumValueConfig{Value: "TotalProfit"},
	},
})

var allowedFilterField = graphql.NewEnum(graphql.EnumConfig{
	Name: "AllowedFilterField",
	Values: graphql.EnumValueConfigMap{
		"Region":   &graphql.EnumValueConfig{Value: "Region"},
		"Country":  &graphql.EnumValueConfig{Value: "Country"},
		"ItemType": &graphql.EnumValueConfig{Value: "ItemType"},
		"OrderID":  &graphql.EnumValueConfig{Value: "OrderID"},
	},
})

var saleList = graphql.NewObject(graphql.ObjectConfig{
	Name: "SaleList",
	Fields: graphql.Fields{
		"sales": &graphql.Field{Type: graphql.NewList(sale), Description: "list of sale"},
		"total": &graphql.Field{Type: gInt, Description: "count of total record"},
	},
})

var sale = graphql.NewObject(graphql.ObjectConfig{
	Name: "Sale",
	Fields: graphql.Fields{
		"region":        &graphql.Field{Type: graphql.String},
		"country":       &graphql.Field{Type: graphql.String},
		"itemType":      &graphql.Field{Type: graphql.String},
		"salesChannel":  &graphql.Field{Type: graphql.String},
		"orderPriority": &graphql.Field{Type: graphql.String},
		"orderId":       &graphql.Field{Type: graphql.String},
		"unitsSold":     &graphql.Field{Type: graphql.Int},
		"unitPrice":     &graphql.Field{Type: graphql.Float},
		"unitCost":      &graphql.Field{Type: graphql.Float},
		"totalRevenue":  &graphql.Field{Type: graphql.Float},
		"totalCost":     &graphql.Field{Type: graphql.Float},
		"totalProfit":   &graphql.Field{Type: graphql.Float},
	},
})

var aggregatedQuery = &graphql.Field{
	Type: graphql.NewList(aggregatedData),
	Args: graphql.FieldConfigArgument{
		"groupBy": arg(graphql.NewNonNull(allowedGroupByField), nil, "group by"),
		"field":   arg(graphql.NewNonNull(allowedCalculateField), nil, "column to calculate, TotalCost | UnitsSold | UnitPrice | UnitCost | TotalRevenue | TotalCost | TotalProfit"),
		"filters": arg(graphql.NewList(simpleQueryFilter), nil, "where filters"),
		"havings": arg(graphql.NewList(aggregatedQueryHaving), nil, "having filters"),
	},
	Description: "query aggregated data of field by filters, sum, avg and count",
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		groupBy := params.Args["groupBy"].(string)
		field := params.Args["field"].(string)
		filters, ok := params.Args["filters"].([]interface{})
		if !ok {
			filters = []interface{}{}
		}
		havings, ok := params.Args["havings"].([]interface{})
		if !ok {
			filters = []interface{}{}
		}

		return model.AggregatedQuery(filters, havings, groupBy, field)
	},
}

var aggregatedQueryHaving = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "AggregatedQueryHaving",
	Fields: graphql.InputObjectConfigFieldMap{
		"field": &graphql.InputObjectFieldConfig{Type: allowedHavingField, Description: "can be Sum | Avg | Count"},
		"value": &graphql.InputObjectFieldConfig{Type: gFloat},
		"op":    &graphql.InputObjectFieldConfig{Type: allowedHavingOp, Description: "can be gt | lt | eq"},
	},
})

var allowedGroupByField = graphql.NewEnum(graphql.EnumConfig{
	Name: "AllowedGroupByField",
	Values: graphql.EnumValueConfigMap{
		"Region":        &graphql.EnumValueConfig{Value: "region"},
		"Country":       &graphql.EnumValueConfig{Value: "country"},
		"ItemType":      &graphql.EnumValueConfig{Value: "item_type"},
		"SalesChannel":  &graphql.EnumValueConfig{Value: "sales_channel"},
		"OrderPriority": &graphql.EnumValueConfig{Value: "order_priority"},
	},
})

var allowedCalculateField = graphql.NewEnum(graphql.EnumConfig{
	Name: "AllowedCalculateField",
	Values: graphql.EnumValueConfigMap{
		"UnitsSold":    &graphql.EnumValueConfig{Value: "units_sold"},
		"UnitPrice":    &graphql.EnumValueConfig{Value: "unit_price"},
		"UnitCost":     &graphql.EnumValueConfig{Value: "unit_cost"},
		"TotalRevenue": &graphql.EnumValueConfig{Value: "total_revenue"},
		"TotalCost":    &graphql.EnumValueConfig{Value: "total_cost"},
		"TotalProfit":  &graphql.EnumValueConfig{Value: "total_profit"},
	},
})

var allowedHavingOp = graphql.NewEnum(graphql.EnumConfig{
	Name: "AllowedHavingOp",
	Values: graphql.EnumValueConfigMap{
		"lt": &graphql.EnumValueConfig{Value: "<"},
		"gt": &graphql.EnumValueConfig{Value: ">"},
		"eq": &graphql.EnumValueConfig{Value: "="},
	},
})

var allowedHavingField = graphql.NewEnum(graphql.EnumConfig{
	Name: "AllowedHavingField",
	Values: graphql.EnumValueConfigMap{
		"Sum":   &graphql.EnumValueConfig{Value: "sum"},
		"Avg":   &graphql.EnumValueConfig{Value: "avg"},
		"Count": &graphql.EnumValueConfig{Value: "count"},
	},
})

var aggregatedData = graphql.NewObject(graphql.ObjectConfig{
	Name: "AggregatedData",
	Fields: graphql.Fields{
		"name":  &graphql.Field{Type: gString},
		"sum":   &graphql.Field{Type: gFloat},
		"avg":   &graphql.Field{Type: gFloat},
		"count": &graphql.Field{Type: gInt},
	},
})
