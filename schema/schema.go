package schema

import (
	"github.com/graphql-go/graphql"
	"log"
)

// Schema _
var Schema graphql.Schema

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
		"ping": ping,
	},
})

var ping = &graphql.Field{
	Type: graphql.String,
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		return "pong", nil
	},
}
