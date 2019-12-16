package chartstreams

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/otaviof/chart-streams/pkg/chartstreams/config"
	"github.com/otaviof/chart-streams/pkg/chartstreams/provider"
)

// ChartStreamServer represents the chartstreams server offering its API. The server puts together
// the routes, and bootstrap steps in order to respond as a valid Helm repository.
type ChartStreamServer struct {
	config        *config.Config
	chartProvider provider.ChartProvider
}

// Start executes the boostrap steps in order to start listening on configured address. It can return
// errors from "listen" method.
func (s *ChartStreamServer) Start() error {
	if err := s.chartProvider.Initialize(); err != nil {
		return err
	}

	return s.listen()
}

// IndexHandler endpoint handler to render a index.yaml file.
func (s *ChartStreamServer) IndexHandler(c *gin.Context) {
	index, err := s.chartProvider.GetIndexFile()
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.YAML(200, index)
}

func (s *ChartStreamServer) runHelmInstall(c *gin.Context) {

	chartName := c.PostForm("chart")
	namespace := c.PostForm("namespace")
	bearerToken := c.GetHeader("Authorization")
	app := "/usr/local/bin/helm"

	arg1 := "install"
	arg0 := "--generate-name"
	arg2 := chartName //"https://technosophos.github.io/tscharts/mink-0.1.0.tgz"
	arg3 := "--namespace=" + namespace
	arg4 := "--token=" + bearerToken

	fmt.Println(app, arg1, arg0, arg2, arg3, arg4)

	cmd := exec.Command(app, arg1, arg0, arg2, arg3, arg4)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "XDG_CACHE_HOME=/tmp")

	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(stdout))
}

// DirectLinkHandler endpoint handler to directly load a chart tarball payload.
func (s *ChartStreamServer) DirectLinkHandler(c *gin.Context) {
	name := c.Param("name")
	version := c.Param("version")
	version = strings.TrimPrefix(version, "/")

	p, err := s.chartProvider.GetChart(name, version)
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Data(http.StatusOK, "application/gzip", p.Bytes())
}

// listen on configured address, after adding the route handlers to the framework. It can return
// errors coming from Gin.
func (s *ChartStreamServer) listen() error {
	g := gin.New()

	g.Use(ginrus.Ginrus(log.StandardLogger(), time.RFC3339, true))

	g.GET("/index.yaml", s.IndexHandler)
	g.GET("/chart/:name/*version", s.DirectLinkHandler)
	g.POST("helm/install", s.runHelmInstall)

	return g.Run(s.config.ListenAddr)
}

// NewServer instantiate a new server instance.
func NewServer(config *config.Config) *ChartStreamServer {
	p := provider.NewGitChartProvider(config)
	return &ChartStreamServer{
		config:        config,
		chartProvider: p,
	}
}
