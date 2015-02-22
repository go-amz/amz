//
// goamz - Go packages to interact with the Amazon Web Services.
//
//   https://wiki.ubuntu.com/goamz
//
package s3

import (
	"gopkg.in/amz.v2/aws"
)

var originalStrategy = attempts

func SetAttemptStrategy(s *aws.AttemptStrategy) {
	if s == nil {
		attempts = originalStrategy
	} else {
		attempts = *s
	}
}

func AttemptStrategy() aws.AttemptStrategy {
	return attempts
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
