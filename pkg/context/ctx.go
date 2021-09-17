// 方便 context 中设置内容
// idea from https://github.com/grpc-ecosystem/go-grpc-middleware/tree/master/tags
package context

import "context"

type ctxMarker struct{}

var (
	ctxMarkerKey = &ctxMarker{}
)

func MustExtract(ctx context.Context) *MapTag {
	t, ok := ctx.Value(ctxMarkerKey).(*MapTag)
	if !ok {
		panic("tag must be set first")
	}
	return t
}

func NewTagContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxMarkerKey, NewMapTag())
}
