package chartstreams

import (
	"github.com/gin-gonic/gin"
)

// Server represents the chart-streams server offering its API. The server puts together the routes,
// and bootstrap steps in order to respond as a valid Helm repository.
type Server struct {
	config *Config
}

// Start executes the boostrap steps in order to start listening on configured address. It can return
// errors from "listen" method.
func (s *Server) Start() error {
	return s.listen()
}

// listen on configured address, after adding the route handlers to the framework. It can return
// errors coming from Gin.
func (s *Server) listen() error {
	g := gin.Default()

	g.GET("/index.yaml", IndexHandler)
	g.GET("/chart/:name/*version", DirectLinkHandler)

	return g.Run(s.config.ListenAddr)
}

// NewServer instantiate a new server instance.
func NewServer(config *Config) *Server {
	return &Server{config: config}
}