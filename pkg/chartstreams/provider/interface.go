package provider

import (
	"k8s.io/helm/pkg/repo"

	"github.com/otaviof/chart-streams/pkg/chartstreams/chart"
)

// ChartProvider group of methods to initialize a backend able to retrieve index and charts.
type ChartProvider interface {
	Initialize() error
	GetChart(name, version string) (*chart.Package, error)
	GetIndexFile() (*repo.IndexFile, error)
}
