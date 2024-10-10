package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func RequestLogger(log *slog.Logger) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		fn := func(w http.ResponseWriter, r *http.Request) {

			t1 := time.Now()
			defer func() {
				log.Info(
					"request accepted"+r.Method+r.Host+r.RequestURI+r.RemoteAddr,
					slog.Int64("took_ms", time.Since(t1).Abs().Milliseconds()),
				)
				fmt.Printf(
					"time %03d",
					time.Since(t1).Milliseconds(),
				)
			}()

			next.ServeHTTP(w, r)

		}

		return http.HandlerFunc(fn)
	}

}
