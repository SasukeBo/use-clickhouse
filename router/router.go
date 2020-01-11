package router

import (
	"github.com/SasukeBo/use-clickhouse/controller"
	"github.com/gin-gonic/gin"
)

// Run start router engine
func Run() error {
	r := gin.Default()

	r.GET("/api", controller.GraphQLHander())
	r.POST("/api", controller.GraphQLHander())

	err := r.Run(":4000")
	return err
}
