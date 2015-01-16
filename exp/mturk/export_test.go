package mturk

import (
	"gopkg.in/amz.v3-unstable/aws"
)

func Sign(auth aws.Auth, service, method, timestamp string, params map[string]string) {
	sign(auth, service, method, timestamp, params)
}
