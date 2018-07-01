package metrics

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/KyleBanks/dockerstats"
	. "github.com/mlabouardy/mon-put-instance-data/services"
	"github.com/shirou/gopsutil/docker"
)

// Docker metric entity
type Docker struct{}

// On older systems, the control groups might be mounted on /cgroup
func getCgroupMountPath() (string, error) {
	out, err := exec.Command("grep", "-m1", "cgroup", "/proc/mounts").Output()
	if err != nil {
		return "", errors.New("Cannot figure out where control groups are mounted")
	}
	res := strings.Fields(string(out))
	if strings.HasPrefix(res[1], "/cgroup") {
		return "/cgroup", nil
	}
	return "/sys/fs/cgroup", nil
}

//Collect CPU & Memory usage per Docker Container
func (d Docker) Collect(instanceID string, c CloudWatchService, namespace string) {
	containers, err := docker.GetDockerStat()
	if err != nil {
		log.Fatal(err)
	}

	containerStats, statsErr := dockerstats.Current()
	if statsErr != nil {
		log.Fatal(statsErr)
	}

	base, err := getCgroupMountPath()
	if err != nil {
		log.Fatal(err)
	}

	for _, container := range containers {
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
			Value: &container.ContainerID,
		})
		dimensionKey3 := "ContainerName"
		dimensions = append(dimensions, cloudwatch.Dimension{
			Name:  &dimensionKey3,
			Value: &container.Name,
		})
		dimensionKey4 := "DockerImage"
		dimensions = append(dimensions, cloudwatch.Dimension{
			Name:  &dimensionKey4,
			Value: &container.Image,
		})

		containerMemory, err := docker.CgroupMem(container.ContainerID, fmt.Sprintf("%s/memory/docker", base))
		if err != nil {
			log.Fatal(err)
		}

		containerMemoryData := constructMetricDatum("ContainerMemory", float64(containerMemory.MemUsageInBytes), cloudwatch.StandardUnitBytes, dimensions)
		c.Publish(containerMemoryData, namespace)

		var iCpu float64 
		fmt.Sscan(stats.CPU, &iCpu)
		containerCPUUsageData := constructMetricDatum("ContainerCPUUsage", iCpu, cloudwatch.StandardUnitPercent, dimensions)
		c.Publish(containerCPUUsageData, namespace)

		log.Printf("Docker - Container:%s Memory:%v CPU:%v%%\n", container.Name, containerMemory.MemMaxUsageInBytes, iCpu)
	}
}
