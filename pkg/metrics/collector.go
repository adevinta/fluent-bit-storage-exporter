package metrics

import (
	"log"

	"github.com/adevinta/fluent-bit-storage-exporter/pkg/client"
	"github.com/dustin/go-humanize"
	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	fbStorageChunks            *prometheus.Desc
	fbStorageChunksMem         *prometheus.Desc
	fbStorageChunksFs          *prometheus.Desc
	fbStorageChunksFsUp        *prometheus.Desc
	fbStorageChunksFsDown      *prometheus.Desc
	fbStorageInputOverlimit    *prometheus.Desc
	fbStorageInputMemBytes     *prometheus.Desc
	fbStorageInputLimitBytes   *prometheus.Desc
	fbStorageInputChunks       *prometheus.Desc
	fbStorageInputChunksFsDown *prometheus.Desc
	fbStorageInputChunksBusy   *prometheus.Desc
	fbStorageInputBusyBytes    *prometheus.Desc
	fluentBitCalls             client.FluentBitMethods
}

const (
	AuditLabel     string = "audit"
	ContainerLabel string = "containers"
	SystemdLabel   string = "systemd"
)

func NewCollector(fluentBitCalls client.FluentBitMethods) *Collector {
	return &Collector{

		fluentBitCalls: fluentBitCalls,
		fbStorageChunks: prometheus.NewDesc("fluentbit_storage_chunks",
			"Amount of currently used chunks", nil, nil,
		),
		fbStorageChunksMem: prometheus.NewDesc("fluentbit_storage_chunks_mem",
			"Amount of chunks currently in memory", nil, nil,
		),
		fbStorageChunksFs: prometheus.NewDesc("fluentbit_storage_chunks_fs",
			"Amount of chunks currently in filesystem", nil, nil,
		),
		fbStorageChunksFsUp: prometheus.NewDesc("fluentbit_storage_chunks_fs_up",
			"Amount of chunks currently up", nil, nil,
		),
		fbStorageChunksFsDown: prometheus.NewDesc("fluentbit_storage_chunks_fs_down",
			"Amount of chunks currently down", nil, nil,
		),
		fbStorageInputOverlimit: prometheus.NewDesc("fluentbit_storage_input_overlimit",
			"Memory buffer limit reached for input", []string{"name"}, nil,
		),
		fbStorageInputMemBytes: prometheus.NewDesc("fluentbit_storage_input_mem_bytes",
			"Currently used memory buffer for input in bytes", []string{"name"}, nil,
		),
		fbStorageInputLimitBytes: prometheus.NewDesc("fluentbit_storage_input_limit_bytes",
			"Memory buffer limit for input in bytes", []string{"name"}, nil,
		),
		fbStorageInputChunks: prometheus.NewDesc("fluentbit_storage_input_chunks",
			"Amount of chunks currently used for input", []string{"name"}, nil,
		),
		fbStorageInputChunksFsDown: prometheus.NewDesc("fluentbit_storage_input_chunks_fs_down",
			"Amount of chunks for input currently down", []string{"name"}, nil,
		),
		fbStorageInputChunksBusy: prometheus.NewDesc("fluentbit_storage_input_chunks_busy",
			"Amount of chunks for input currently busy", []string{"name"}, nil,
		),
		fbStorageInputBusyBytes: prometheus.NewDesc("fluentbit_storage_input_busy_bytes",
			"Size of chunks for input currently busy in bytes", []string{"name"}, nil,
		),
	}
}

func (collector *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.fbStorageChunks
	ch <- collector.fbStorageChunksMem
	ch <- collector.fbStorageChunksFs
	ch <- collector.fbStorageChunksFsUp
	ch <- collector.fbStorageChunksFsDown
	ch <- collector.fbStorageInputOverlimit
	ch <- collector.fbStorageInputMemBytes
	ch <- collector.fbStorageInputLimitBytes
	ch <- collector.fbStorageInputChunks
	ch <- collector.fbStorageInputChunksFsDown
	ch <- collector.fbStorageInputChunksBusy
	ch <- collector.fbStorageInputBusyBytes
}

func (collector *Collector) Collect(ch chan<- prometheus.Metric) {
	metrics, err := collector.fluentBitCalls.GetMetricData()
	if err != nil {
		log.Fatal(err)
	}
	ch <- prometheus.MustNewConstMetric(collector.fbStorageChunks, prometheus.GaugeValue, float64(metrics.StorageLayer.Chunks.TotalChunks))
	ch <- prometheus.MustNewConstMetric(collector.fbStorageChunksMem, prometheus.GaugeValue, float64(metrics.StorageLayer.Chunks.MemChunks))
	ch <- prometheus.MustNewConstMetric(collector.fbStorageChunksFs, prometheus.GaugeValue, float64(metrics.StorageLayer.Chunks.FsChunks))
	ch <- prometheus.MustNewConstMetric(collector.fbStorageChunksFsUp, prometheus.GaugeValue, float64(metrics.StorageLayer.Chunks.FsChunksUp))
	ch <- prometheus.MustNewConstMetric(collector.fbStorageChunksFsDown, prometheus.GaugeValue, float64(metrics.StorageLayer.Chunks.FsChunksDown))

	//Check inputChunks is empty
	if (metrics.InputChunks != client.InputChunks{}) {
		for _, labelName := range []string{AuditLabel, ContainerLabel, SystemdLabel} {
			var cdata client.GenericInputStruct
			switch labelName {
			case AuditLabel:
				cdata = metrics.InputChunks.Audit
			case ContainerLabel:
				cdata = metrics.InputChunks.Containers
			case SystemdLabel:
				cdata = metrics.InputChunks.Systemd
			}
			memsize, err := convertStringToInt(cdata.Status.MemSize)
			if err != nil {
				log.Fatal(err)
			}
			memlimit, err := convertStringToInt(cdata.Status.MemLimit)
			if err != nil {
				log.Fatal(err)
			}
			busyBytes, err := convertStringToInt(cdata.Chunks.BusySize)
			if err != nil {
				log.Fatal(err)
			}
			ch <- prometheus.MustNewConstMetric(collector.fbStorageInputOverlimit, prometheus.GaugeValue, float64(convertBoolToInt(cdata.Status.Overlimit)), labelName)
			ch <- prometheus.MustNewConstMetric(collector.fbStorageInputMemBytes, prometheus.GaugeValue, float64(*memsize), labelName)
			ch <- prometheus.MustNewConstMetric(collector.fbStorageInputLimitBytes, prometheus.GaugeValue, float64(*memlimit), labelName)
			ch <- prometheus.MustNewConstMetric(collector.fbStorageInputChunks, prometheus.GaugeValue, float64(cdata.Chunks.Total), labelName)
			ch <- prometheus.MustNewConstMetric(collector.fbStorageInputChunksFsDown, prometheus.GaugeValue, float64(cdata.Chunks.Down), labelName)
			ch <- prometheus.MustNewConstMetric(collector.fbStorageInputChunksBusy, prometheus.GaugeValue, float64(cdata.Chunks.Busy), labelName)
			ch <- prometheus.MustNewConstMetric(collector.fbStorageInputBusyBytes, prometheus.GaugeValue, float64(*busyBytes), labelName)

		}
	}
}

func convertBoolToInt(b bool) uint32 {
	if b {
		return 1
	}
	return 0
}

func convertStringToInt(s string) (*uint64, error) {
	result, err := humanize.ParseBytes(s)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
