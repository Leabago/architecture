//go:generate oapi-codegen -generate types,server -package api openAPI.yaml > api/server.gen.go

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"open-api/api"

	"github.com/gin-gonic/gin"
	middleware "github.com/oapi-codegen/gin-middleware"
)

func NewGinPetServer(petStore *api.PetStore, port string) *http.Server {
	swagger, err := api.GetSwagger()
	if err != nil {
		log.Fatalf("Error loading swagger spec\n: %s", err)
	}

	for i := range len(swagger.Servers) {
		fmt.Println("swagger123", swagger.Servers[i].URL)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	route := gin.Default()

	// Serve raw OpenAPI spec
	route.GET("/openapi.json", func(c *gin.Context) {
		c.File("api/api.yaml") // or generate dynamically
	})
	// Route to serve single HTML file
	route.GET("/openapi", func(c *gin.Context) {
		c.File("doc/index.html")
	})

	route.GET("/hello", hello)

	route.Use(middleware.OapiRequestValidator(swagger))

	api.RegisterHandlers(route, petStore)

	s := &http.Server{
		Handler: route,
		Addr:    net.JoinHostPort("0.0.0.0", port),
	}

	return s
}

func main() {
	port := flag.String("port", "8080", "Port for test HTTP server")
	flag.Parse()
	// Create an instance of our handler which satisfies the generated interface
	petStore := api.NewPetStore()
	g := NewGinPetServer(petStore, *port)

	log.Fatal(g.ListenAndServe())

}

func hello(ctx *gin.Context) {
	ctx.String(http.StatusOK, "hello")
}
