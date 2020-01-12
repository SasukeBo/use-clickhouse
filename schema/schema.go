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
	gString = graphql.String
	gInt    = graphql.Int
	gNInt   = graphql.NewNonNull(graphql.Int)
	gFloat  = graphql.Float
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
		"ping":   ping,
		"simple": simpleQuery,
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
		"fields":  arg(graphql.NewList(gString), nil, "query fields"),
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
		"field": &graphql.InputObjectFieldConfig{Type: gString},
		"value": &graphql.InputObjectFieldConfig{Type: gString},
	},
})

var saleList = graphql.NewObject(graphql.ObjectConfig{
	Name: "SaleList",
	Fields: graphql.Fields{
		"sales": &graphql.Field{Type: graphql.NewList(sale), Description: "list of sale"},
		"total": &graphql.Field{Type: gInt, Description: "count of total record"},
	},
})

// sale object
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
