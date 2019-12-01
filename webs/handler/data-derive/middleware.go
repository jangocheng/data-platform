package data_derive

import (
	"context"
	"net/http"
	"platform/micro-common/utils"
	"time"
)

func (d *DataDeriveWeb) addStartTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "startTime", time.Now())
		next.ServeHTTP(w, r.WithContext(ctx))
		r.Header.Set("Content-Type", "application/json")
	})
}

func (d *DataDeriveWeb) addSwiftNumberMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx = context.Background()

		platformString := r.FormValue("_platform")
		reqPlatformParams := make(map[string]string)
		reqPlatformParams["_platform"] = platformString

		reqParamsStr := r.FormValue(dataApiParamName)

		ctx = context.WithValue(r.Context(), "swiftNumber", utils.GenSwiftNumber(reqParamsStr))

		ctx = context.WithValue(ctx, "parentSwiftNumber", utils.GetParentSwiftNumber(reqPlatformParams))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
