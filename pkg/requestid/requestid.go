package requestid

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const (
	Header = "X-Request-ID"
	LogKey = "requestId"
)

func Empty() string {
	return uuid.Nil.String()
}

func Generator() string {
	rid, err := uuid.NewRandom()
	if err != nil {
		return Empty()
	}

	return rid.String()
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rid := r.Header.Get(Header)
			if rid == "" {
				rid = Generator()
			}

			w.Header().Set(Header, rid)

			ctx := Inject(r.Context(), rid)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type ctxKey struct{}

func Inject(ctx context.Context, rid string) context.Context {
	return context.WithValue(ctx, ctxKey{}, rid)
}

func Extract(ctx context.Context) string {
	rid, ok := ctx.Value(ctxKey{}).(string)
	if !ok {
		return Empty()
	}

	return rid
}
