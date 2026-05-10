package grpc

import (
	"fmt"

	"google.golang.org/grpc/codes"
)

var errUnauthenticated = fmt.Errorf("%s", codes.Unauthenticated)
