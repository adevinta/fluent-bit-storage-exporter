package metrics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/adevinta/fluent-bit-storage-exporter/pkg/client"
	dto "github.com/prometheus/client_model/go"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type FakeUtil struct {
	path    string
	counter int
}

func (c *FakeUtil) GetMetricData() (*client.Response, error) {
	jsonFile, _ := os.Open(c.path)
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result client.Response
	json.Unmarshal([]byte(byteValue), &result)

	c.counter++
	return &result, nil
}

func TestFullCollect(t *testing.T) {
	fakeUtil := FakeUtil{path: "testdata/full.json"}
	customColector := NewCollector(&fakeUtil)
	ch := make(chan prometheus.Metric)
	go func() {
		customColector.Collect(ch)
		close(ch)
	}()
	reportedMetrics := 0
	for element := range ch {
		reportedMetrics++
		metric := dto.Metric{}
		require.NoError(t, element.Write(&metric))
		switch element.Desc() {
		case customColector.fbStorageChunks:
			require.NotNil(t, metric.Gauge)
			assert.EqualValues(t, float64(2), metric.Gauge.GetValue())
		case customColector.fbStorageChunksMem:
			require.NotNil(t, metric.Gauge)
			assert.EqualValues(t, float64(20), metric.Gauge.GetValue())
		case customColector.fbStorageChunksFs:
			require.NotNil(t, metric.Gauge)
			assert.EqualValues(t, float64(0), metric.Gauge.GetValue())
		case customColector.fbStorageChunksFsUp:
			require.NotNil(t, metric.Gauge)
			assert.EqualValues(t, float64(0), metric.Gauge.GetValue())
		case customColector.fbStorageChunksFsDown:
			require.NotNil(t, metric.Gauge)
			assert.EqualValues(t, float64(0), metric.Gauge.GetValue())
		case customColector.fbStorageInputOverlimit:
			assertMetricWithMultipleLabels(t, metric, map[string]float64{
				"containers": float64(1),
				"systemd":    float64(0),
				"audit":      float64(0),
			})
		case customColector.fbStorageInputMemBytes:
			assertMetricWithMultipleLabels(t, metric, map[string]float64{
				"containers": float64(47900),
				"systemd":    float64(1024),
				"audit":      float64(0),
			})
		case customColector.fbStorageInputLimitBytes:
			assertMetricWithMultipleLabels(t, metric, map[string]float64{
				"containers": 6.1e+07,
				"systemd":    float64(0),
				"audit":      3.34e+07,
			})
		case customColector.fbStorageInputChunks:
			assertMetricWithMultipleLabels(t, metric, map[string]float64{
				"containers": float64(2),
				"systemd":    float64(0),
				"audit":      float64(0),
			})
		case customColector.fbStorageInputChunksFsDown:
			assertMetricWithMultipleLabels(t, metric, map[string]float64{
				"containers": float64(0),
				"systemd":    float64(0),
				"audit":      float64(0),
			})
		case customColector.fbStorageInputChunksBusy:
			assertMetricWithMultipleLabels(t, metric, map[string]float64{
				"containers": float64(2),
				"systemd":    float64(0),
				"audit":      float64(0),
			})
		case customColector.fbStorageInputBusyBytes:
			assertMetricWithMultipleLabels(t, metric, map[string]float64{
				"containers": float64(47900),
				"systemd":    float64(0),
				"audit":      float64(0),
			})
		default:
			t.Errorf("unsupported metric desc: %v", element.Desc())
		}
	}
	assert.Equal(t, 1, fakeUtil.counter)
	assert.Equal(t, 26, reportedMetrics)
}

func TestCollectWithOnlyStorage(t *testing.T) {
	fakeUtil := FakeUtil{path: "testdata/onlystorage.json"}
	customColector := NewCollector(&fakeUtil)
	ch := make(chan prometheus.Metric)
	go func() {
		customColector.Collect(ch)
		close(ch)
	}()
	reportedMetrics := 0
	for element := range ch {
		reportedMetrics++
		metric := dto.Metric{}
		require.NoError(t, element.Write(&metric))
		switch element.Desc() {
		case customColector.fbStorageChunks:
			require.NotNil(t, metric.Gauge)
			assert.EqualValues(t, float64(2), metric.Gauge.GetValue())
		case customColector.fbStorageChunksMem:
			require.NotNil(t, metric.Gauge)
			assert.EqualValues(t, float64(20), metric.Gauge.GetValue())
		case customColector.fbStorageChunksFs:
			require.NotNil(t, metric.Gauge)
			assert.EqualValues(t, float64(0), metric.Gauge.GetValue())
		case customColector.fbStorageChunksFsUp:
			require.NotNil(t, metric.Gauge)
			assert.EqualValues(t, float64(0), metric.Gauge.GetValue())
		case customColector.fbStorageChunksFsDown:
			require.NotNil(t, metric.Gauge)
			assert.EqualValues(t, float64(0), metric.Gauge.GetValue())
		case customColector.fbStorageInputOverlimit:
			require.NotNil(t, metric.Gauge)
			assert.EqualValues(t, float64(1), metric.Gauge.GetValue())
		case customColector.fbStorageInputMemBytes:
			require.NotNil(t, metric.Gauge)
			assert.EqualValues(t, float64(47900), metric.Gauge.GetValue())
		case customColector.fbStorageInputLimitBytes:
			require.NotNil(t, metric.Gauge)
			assert.EqualValues(t, float64(6.1e+07), metric.Gauge.GetValue())
		case customColector.fbStorageInputChunks:
			require.NotNil(t, metric.Gauge)
			assert.EqualValues(t, float64(2), metric.Gauge.GetValue())
		case customColector.fbStorageInputChunksFsDown:
			require.NotNil(t, metric.Gauge)
			assert.EqualValues(t, float64(0), metric.Gauge.GetValue())
		case customColector.fbStorageInputChunksBusy:
			require.NotNil(t, metric.Gauge)
			assert.EqualValues(t, float64(2), metric.Gauge.GetValue())
		case customColector.fbStorageInputBusyBytes:
			require.NotNil(t, metric.Gauge)
			assert.EqualValues(t, float64(47900), metric.Gauge.GetValue())

		default:
			t.Errorf("unsupported metric desc: %v", element.Desc())
		}
	}
	assert.Equal(t, 1, fakeUtil.counter)
	assert.Equal(t, 5, reportedMetrics)
}

func assertMetricWithMultipleLabels(t *testing.T, metric dto.Metric, options map[string]float64) {
	labels := metric.GetLabel()
	for _, v := range labels {
		if v.GetName() == "name" {
			switch v.GetValue() {
			case "containers":
				assert.EqualValues(t, float64(options["containers"]), metric.Gauge.GetValue(), fmt.Sprintf("Failed checking metric label %s", v.GetValue()))
			case "systemd":
				assert.EqualValues(t, float64(options["systemd"]), metric.Gauge.GetValue(), fmt.Sprintf("Failed checking metric label %s", v.GetValue()))
			case "audit":
				assert.EqualValues(t, float64(options["audit"]), metric.Gauge.GetValue(), fmt.Sprintf("Failed checking metric label %s", v.GetValue()))
			default:
				t.Errorf("unsupported label metric : %s", v.GetValue())

			}
		}
	}

}
