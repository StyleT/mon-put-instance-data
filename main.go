package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	. "github.com/mlabouardy/mon-put-instance-data/metrics"
	. "github.com/mlabouardy/mon-put-instance-data/services"
	"github.com/urfave/cli"
)

// GetInstanceID return EC2 instance id
func GetInstanceID() (string, error) {
	value := os.Getenv("AWS_INSTANCE_ID")
	if len(value) > 0 {
		return value, nil
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://169.254.169.254/latest/meta-data/instance-id", nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Collect metrics about enabled metric
func Collect(metrics []Metric, c PubliserService, namespace string, instanceId string) {
	metricsContainer := c.GetContainer()

	for _, metric := range metrics {
		metric.Collect(instanceId, &metricsContainer)
	}

	c.Publish(metricsContainer, namespace)
}

func main() {
	app := cli.NewApp()
	app.Name = "mon-put-instance-data"
	app.Usage = "Publish Custom Metrics to CloudWatch"
	app.Version = "1.0.0"
	app.Author = "Mohamed Labouardy"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "memory",
			Usage: "Collect memory metrics",
		},
		cli.BoolFlag{
			Name:  "swap",
			Usage: "Collect swap metrics",
		},
		cli.BoolFlag{
			Name:  "disk",
			Usage: "Collect disk metrics",
		},
		cli.BoolFlag{
			Name:  "network",
			Usage: "Collect network metrics",
		},
		cli.BoolFlag{
			Name:  "docker",
			Usage: "Collect containers metrics",
		},
		cli.StringFlag{
			Name:  "region",
			Usage: "AWS region (default: instance region or us-east-1)",
			Value: "default",
		},
		cli.IntFlag{
			Name:  "interval",
			Usage: "Time interval",
			Value: 5,
		},
		cli.BoolFlag{
			Name:  "once",
			Usage: "Run once (i.e. not on an interval)",
		},
		cli.StringFlag{
			Name:  "namespace",
			Usage: "Namespace for the metric data",
			Value: "CustomMetrics",
		},
		cli.BoolFlag{
			Name:  "dummy",
			Usage: "Disables interaction with AWS and only prints data to stdout",
		},
	}
	app.Action = func(c *cli.Context) error {
		enabledMetrics := make([]string, 0)
		metrics := make([]Metric, 0)
		if c.Bool("memory") {
			metrics = append(metrics, Memory{})
			enabledMetrics = append(enabledMetrics, "memory")
		}
		if c.Bool("swap") {
			metrics = append(metrics, Swap{})
			enabledMetrics = append(enabledMetrics, "swap")
		}
		if c.Bool("disk") {
			metrics = append(metrics, Disk{})
			enabledMetrics = append(enabledMetrics, "disk")
		}
		if c.Bool("network") {
			metrics = append(metrics, Network{})
			enabledMetrics = append(enabledMetrics, "network")
		}
		if c.Bool("docker") {
			metrics = append(metrics, Docker{})
			enabledMetrics = append(enabledMetrics, "docker")
		}

		cfg, err := external.LoadDefaultAWSConfig()
		if err != nil {
			panic("Unable to load SDK config")
		}

		if c.String("region") == "default" && !c.Bool("dummy") {
			metadataSvc := ec2metadata.New(cfg)
			if !metadataSvc.Available() {
				log.Printf("Metadata service cannot be reached. Default us-east-1 region will be used")
				cfg.Region = "us-east-1"
			} else {
				region, err := metadataSvc.Region()
				if err != nil {
					panic("Unable to fetch instance region")
				}

				log.Printf("Region is: %s", region)
				cfg.Region = region
			}
		} else {
			cfg.Region = c.String("region")
		}

		var service PubliserService
		var instanceId string
		if !c.Bool("dummy") {
			service = CloudWatchService{
				Config: cfg,
			}
			instanceId, err = GetInstanceID()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			service = DummyService{}
			instanceId = "i-1234567890abcdef0"
		}
		fmt.Printf("Instance ID: %s\n", instanceId)
		

		interval := c.Int("interval")

		fmt.Printf("Features enabled: %s\n", strings.Join(enabledMetrics, ", "))

		var collect = func() {
			Collect(metrics, service, c.String("namespace"), instanceId)
		}

		if c.Bool("once") {
			collect()
		} else {
			fmt.Printf("Interval: %d minutes\n", interval)

			duration := time.Duration(interval) * time.Minute
			for range time.Tick(duration) {
				collect()
			}
		}

		return nil
	}
	app.Run(os.Args)
}
