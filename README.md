[![Build Status](https://travis-ci.org/StyleT/mon-put-instance-data.svg?branch=master)](https://travis-ci.org/StyleT/mon-put-instance-data)
[![Docker Stars](https://img.shields.io/docker/pulls/mlabouardy/mon-put-instance-data.svg)](https://hub.docker.com/r/stylet/mon-put-instance-data/) 
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE) [![Docker Stars](https://img.shields.io/github/issues/mlabouardy/mon-put-instance-data.svg)](https://github.com/stylet/mon-put-instance-data/issues)  

## How to use

* Setup an IAM Policy:

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "1",
            "Effect": "Allow",
            "Action": "cloudwatch:PutMetricData",
            "Resource": "*"
        }
    ]
}
```

* Start metrics collector:

```
mon-put-instance-data --memory --swap --disk --network --docker --duration 1
```

## Docker

Use the official Docker image:

```
docker run --rm -it --privileged -v /var/run/docker.sock:/var/run/docker.sock:ro -v /sys:/sys:ro stylet/mon-put-instance-data --interval 1 --memory --swap --disk --network --docker
```

## Metrics

* Memory
    * Memory Utilization (%)
    * Memory Used (Mb)
    * Memory Available (Mb)
* Swap
    * Swap Utilization (%)
    * Swap Used (Mb)
* Disk
    * Disk Space Utilization (%)
    * Disk Space Used (Gb)
    * Disk Space Available (Gb)
* Network
    * Bytes In/Out
    * Packets In/Out
    * Errors In/Out
* Docker
    * Memory Utilization per Container
    * CPU usage per Container

## Supported AMI

* Amazon Linux
* Amazon Linux 2
* Ubuntu 16.04
* CoreOS

## Tutorial

* [Publish Custom Metrics to AWS CloudWatch](http://www.blog.labouardy.com/publish-custom-metrics-aws-cloudwatch/)

## Development

LDE:
```
$ docker run -it -v $PWD:/go/src/github.com/mlabouardy/mon-put-instance-data golang
$ curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
$ dep ensure
$ go build
```

Compile binaries:
```
$ go get github.com/mitchellh/gox
$ gox -osarch="linux/amd64 darwin/amd64"
```
