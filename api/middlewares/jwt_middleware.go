package middlewares

import (
	"context"
	"net/http"

	helpers "github.com/rahul108/order_management_system/api/helper/auth"
	utils "github.com/rahul108/order_management_system/api/utils/jwt"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := utils.ExtractBearerToken(r)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		username, err := utils.ValidateAccessToken(tokenString)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, helpers.CreateAuthorizationFailureMessage())
			return
		}

		// Add username to context for use in handlers
		ctx := context.WithValue(r.Context(), "username", username)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
