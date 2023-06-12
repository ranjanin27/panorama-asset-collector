package prometheus

import (
	"testing"
)

func TestSetupMetricCollector(t *testing.T) {
	metricCollectionHandler := NewMetricCollectionHandler()
	metricCollectionHandler.SetupMetricCollector()
}
