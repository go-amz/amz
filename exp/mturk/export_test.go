package mturk

import (
	"gopkg.in/amz.v1/aws"
)

func Sign(auth aws.Auth, service, method, timestamp string, params map[string]string) {
	sign(auth, service, method, timestamp, params)
}
