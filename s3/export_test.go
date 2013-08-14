package s3

import (
	"launchpad.net/goamz/aws"
)

var originalStrategy = eventualConsistency

func SetAttemptStrategy(s *aws.AttemptStrategy) {
	if s == nil {
		eventualConsistency = originalStrategy
	} else {
		eventualConsistency = *s
	}
}

func Sign(auth aws.Auth, method, path string, params, headers map[string][]string) {
	sign(auth, method, path, params, headers)
}

func SetListPartsMax(n int) {
	listPartsMax = n
}

func SetListMultiMax(n int) {
	listMultiMax = n
}
