package services

import (
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
)

type PubliserService interface {
	Publish(metricData []cloudwatch.MetricDatum, namespace string)
	GetContainer() []cloudwatch.MetricDatum
}