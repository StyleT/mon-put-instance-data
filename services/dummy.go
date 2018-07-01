package services

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
)

// CloudWatchService entity
type DummyService struct {}

// Publish save metrics to cloudwatch using AWS CloudWatch API
func (c DummyService) Publish(metricData []cloudwatch.MetricDatum, namespace string) {
	for _, row := range metricData {
		log.Printf("Metric data '%v' unit '%v' value '%f'", *row.MetricName, row.Unit, *row.Value)
	}
}

func (c DummyService) GetContainer() []cloudwatch.MetricDatum {
	var s []cloudwatch.MetricDatum
	return s
}