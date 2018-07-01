package metrics

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/shirou/gopsutil/disk"
)

// Disk metric entity
type Disk struct{}

// Collect Disk used & free space
func (d Disk) Collect(instanceID string, c *[]cloudwatch.MetricDatum) {
	diskMetrics, err := disk.Usage("/")
	if err != nil {
		log.Fatal(err)
	}

	dimensionKey := "InstanceId"
	dimensions := []cloudwatch.Dimension{
		cloudwatch.Dimension{
			Name:  &dimensionKey,
			Value: &instanceID,
		},
	}

	*c= append(*c, constructMetricDatum("DiskUtilization", diskMetrics.UsedPercent, cloudwatch.StandardUnitPercent, dimensions))

	*c= append(*c, constructMetricDatum("DiskUsed", float64(diskMetrics.Used), cloudwatch.StandardUnitBytes, dimensions))

	*c= append(*c, constructMetricDatum("DiskFree", float64(diskMetrics.Free), cloudwatch.StandardUnitBytes, dimensions))

	log.Printf("Disk - Utilization:%v%% Used:%v Free:%v\n", diskMetrics.UsedPercent, diskMetrics.Used, diskMetrics.Free)
}
