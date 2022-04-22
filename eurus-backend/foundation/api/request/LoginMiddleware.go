package request

import (
	"context"
	"net/http"
	"strings"

	"eurus-backend/foundation"
	"eurus-backend/foundation/auth_base"
)

type LoginMiddleware struct {
	AuthClient               auth_base.IAuth
	SuccessLoginTokenHandler func(token string, loginToken auth_base.ILoginToken) bool
}

//Depreciated
// func (me *LoginMiddleware) VerifyLoginToken(next http.Handler) http.Handler {
// 	return me.newHttpHandler(next, network.VerifyModeUser|network.VerifyModeService)
// }

func (me *LoginMiddleware) VerifyUserLoginToken(next http.Handler) http.Handler {
	return me.newHttpHandler(next, auth_base.VerifyModeUser)
}

func (me *LoginMiddleware) VerifyServiceLoginToken(next http.Handler) http.Handler {
	return me.newHttpHandler(next, auth_base.VerifyModeService)
}

func (me *LoginMiddleware) newHttpHandler(next http.Handler, verifyMode auth_base.VerifyMode) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		loginToken, serverErr := parseAuthorizationHeader(r)

		if serverErr != nil {
			var ctx context.Context
			ctx = context.WithValue(r.Context(), "error", serverErr)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			_, verifiedLoginToken, err := me.AuthClient.VerifyLoginToken(loginToken, verifyMode)
			var ctx context.Context
			if err != nil {
				ctx = context.WithValue(r.Context(), "error", err)
			} else {
				var isSuccess bool = true
				if me.SuccessLoginTokenHandler != nil {
					isSuccess = me.SuccessLoginTokenHandler(loginToken, verifiedLoginToken)
				}
				if isSuccess {
					ctx = context.WithValue(r.Context(), "loginToken", verifiedLoginToken)
				}
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

func parseAuthorizationHeader(r *http.Request) (string, *foundation.ServerError) {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")

	if len(splitToken) != 2 {
		return "", foundation.NewError(foundation.LoginTokenInvalid)
	}
	return splitToken[1], nil

}
