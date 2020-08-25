package api

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/noah-blockchain/autodeleg/internal/env"
	"github.com/noah-blockchain/autodeleg/internal/gate"
	"net/http"
)

// Run API
func Run(gateService *gate.NoahGate) {
	router := SetupRouter(gateService)
	aDelegLink := fmt.Sprintf("%s:%s", env.GetEnv(env.AdelegApiHostEnv, ""), env.GetEnv(env.AdelegApiPortEnv, ""))
	err := router.Run(aDelegLink)
	if err != nil {
		panic(err)
	}
}

//Setup router
func SetupRouter(gateService *gate.NoahGate) *gin.Engine {
	router := gin.Default()
	if !env.GetEnvAsBool(env.DebugModeEnv, true) {
		gin.SetMode(gin.ReleaseMode)
	}

	router.Use(cors.Default())             // CORS
	router.Use(gin.ErrorLogger())          // print all errors
	router.Use(gin.Recovery())             // returns 500 on any code panics
	router.Use(apiMiddleware(gateService)) // init global context

	router.GET("/", Index)

	v1 := router.Group("/api/v1")
	{
		v1.POST("/transactions", Delegate)
	}
	// Default handler 404
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": 404, "log": "Resource not found."}})
	})
	return router
}

//Add necessary services to global context
func apiMiddleware(gateService *gate.NoahGate) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("gate", gateService)
		c.Next()
	}
}
