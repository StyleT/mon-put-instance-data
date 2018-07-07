package metrics

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/KyleBanks/dockerstats"
	"github.com/shirou/gopsutil/docker"
	"code.cloudfoundry.org/bytefmt"
)

// Docker metric entity
type Docker struct{}


//Collect CPU & Memory usage per Docker Container
func (d Docker) Collect(instanceID string, c *[]cloudwatch.MetricDatum) {
	containers, err := docker.GetDockerStat()
	if err != nil {
		log.Fatal(err)
	}

	containerStats, statsErr := dockerstats.Current()
	if statsErr != nil {
		log.Fatal(statsErr)
	}

	for ii, container := range containers {
		var stats dockerstats.Stats
		for _, v := range containerStats {
			if strings.HasPrefix(container.ContainerID, v.Container) {
				stats = v
				break
			}
		}

		dimensions := make([]cloudwatch.Dimension, 0)
		dimensionKey1 := "InstanceId"
		dimensions = append(dimensions, cloudwatch.Dimension{
			Name:  &dimensionKey1,
			Value: &instanceID,
		})
		dimensionKey2 := "ContainerId"
		dimensions = append(dimensions, cloudwatch.Dimension{
			Name:  &dimensionKey2,
			Value: &containers[ii].ContainerID,
		})
		dimensionKey3 := "ContainerName"
		dimensions = append(dimensions, cloudwatch.Dimension{
			Name:  &dimensionKey3,
			Value: &containers[ii].Name,
		})
		dimensionKey4 := "DockerImage"
		dimensions = append(dimensions, cloudwatch.Dimension{
			Name:  &dimensionKey4,
			Value: &containers[ii].Image,
		})

		var fRamUsage float64 
		fmt.Sscan(stats.Memory.Percent, &fRamUsage)
		*c= append(*c, constructMetricDatum("ContainerMemoryUsage", fRamUsage, cloudwatch.StandardUnitPercent, dimensions))

		usedMemoryStr := stats.Memory.Raw[:strings.Index(stats.Memory.Raw, " / ")]
		usedMemory, err := bytefmt.ToBytes(usedMemoryStr)
		if err != nil {
			log.Fatal(err)
		}
		*c= append(*c, constructMetricDatum("ContainerMemory", float64(usedMemory), cloudwatch.StandardUnitBytes, dimensions))

		var fCpu float64 
		fmt.Sscan(stats.CPU, &fCpu)
		*c= append(*c, constructMetricDatum("ContainerCPUUsage", fCpu, cloudwatch.StandardUnitPercent, dimensions))

		log.Printf("Docker - Container:%s Memory:%v CPU:%v%%\n", container.Name, usedMemoryStr, fCpu)
	}
}
