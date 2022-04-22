package response

import (
	"encoding/json"
	"eurus-backend/foundation"
	"eurus-backend/foundation/api/request"

	"net/http"
)

type ErrorMiddleware struct {
}

func (me *ErrorMiddleware) ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if ctx.Value("error") != nil {
			err := ctx.Value("error").(foundation.IServerError)
			if err == nil {
				err = foundation.NewErrorWithMessage(foundation.InternalServerError, foundation.InternalServerError.String())
			}

			var buffer []byte = make([]byte, 1024)
			length, _ := r.Body.Read(buffer)
			reqObj := request.RequestBase{}
			json.Unmarshal(buffer[:length], &reqObj)
			res := CreateErrorResponse(&reqObj, err.GetReturnCode(), err.Error())
			resByte, _ := json.Marshal(res)
			w.Write(resByte)

		} else {
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

func (me *ErrorMiddleware) JsonErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if ctx.Value("error") != nil {
			err := ctx.Value("error").(foundation.IServerError)
			if err == nil {
				err = foundation.NewErrorWithMessage(foundation.InternalServerError, foundation.InternalServerError.String())
			}

			var buffer []byte = make([]byte, 1024)
			length, _ := r.Body.Read(buffer)
			reqObj := request.RequestBase{}
			_ = json.Unmarshal(buffer[:length], &reqObj)
			res := CreateErrorResponse(&reqObj, err.GetReturnCode(), err.Error())
			resByte, _ := json.Marshal(res)

			var statusCode int
			switch err.GetReturnCode() {
			case foundation.LoginTokenInvalid, foundation.LoginTokenExired, foundation.UnauthorizedAccess:
				statusCode = http.StatusUnauthorized
			case foundation.InvalidSession, foundation.InvalidSignature, foundation.BadRequest:
				statusCode = http.StatusBadRequest
			case foundation.MethodNotFound, foundation.RecordNotFound, foundation.UserNotFound:
				statusCode = http.StatusNotFound
			default:
				statusCode = http.StatusInternalServerError
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(statusCode)
			_, _ = w.Write(resByte)
		} else {
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
