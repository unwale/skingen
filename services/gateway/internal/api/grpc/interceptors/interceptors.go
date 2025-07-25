package interceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/unwale/skingen/pkg/constants"
	"github.com/unwale/skingen/pkg/contextutil"
)

func CorrelationIDInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		correlationID := contextutil.CorrelationIDFromContext(ctx)
		ctx = metadata.AppendToOutgoingContext(ctx, constants.CorrelationIDKey, correlationID)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
