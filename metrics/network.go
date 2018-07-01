package metrics

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/shirou/gopsutil/net"
)

// Network metric entity
type Network struct{}

// Collect Network Traffic metrics
func (n Network) Collect(instanceID string, c *[]cloudwatch.MetricDatum) {
	networkMetrics, err := net.IOCounters(false)
	if err != nil {
		log.Fatal(err)
	}

	for _, iocounter := range networkMetrics {
		dimensions := make([]cloudwatch.Dimension, 0)
		dimensionKey1 := "InstanceId"
		dimensions = append(dimensions, cloudwatch.Dimension{
			Name:  &dimensionKey1,
			Value: &instanceID,
		})
		dimensionKey2 := "IOCounter"
		dimensions = append(dimensions, cloudwatch.Dimension{
			Name:  &dimensionKey2,
			Value: &iocounter.Name,
		})

		*c= append(*c, constructMetricDatum("BytesIn", float64(iocounter.BytesRecv), cloudwatch.StandardUnitBytes, dimensions))
		*c= append(*c, constructMetricDatum("BytesOut", float64(iocounter.BytesSent), cloudwatch.StandardUnitBytes, dimensions))
		
		*c= append(*c, constructMetricDatum("PacketsIn", float64(iocounter.PacketsRecv), cloudwatch.StandardUnitBytes, dimensions))
		*c= append(*c, constructMetricDatum("PacketsOut", float64(iocounter.PacketsSent), cloudwatch.StandardUnitBytes, dimensions))
		
		*c= append(*c, constructMetricDatum("ErrorsIn", float64(iocounter.Errin), cloudwatch.StandardUnitBytes, dimensions))
		*c= append(*c, constructMetricDatum("ErrorsOut", float64(iocounter.Errout), cloudwatch.StandardUnitBytes, dimensions))
		
		log.Printf("Network - %s Bytes In/Out: %v/%v Packets In/Out: %v/%v Errors In/Out: %v/%v\n",
			iocounter.Name, iocounter.BytesRecv, iocounter.BytesSent, iocounter.Errin,
			iocounter.Errout, iocounter.PacketsRecv, iocounter.PacketsSent)
	}
}
