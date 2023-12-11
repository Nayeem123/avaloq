package api

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func Authentication(w http.ResponseWriter, r *http.Request) bool {
	// We can obtain the session token from the requests cookies, which come with every request
	//c, err := r.Cookie("token")
	//fmt.Println("/api/v1/execute api called")
	//fmt.Println("c=",c)

	tokenString := extractTokenFromHeader(r)
	if tokenString == "" {
		http.Error(w, "Token Missing!!! Pass Token with Valid Prefix", http.StatusUnauthorized)
		return false
	}
	//fmt.Println("tokenString=", tokenString)
	// if err != nil {
	// 	if err == http.ErrNoCookie {
	// 		// If the cookie is not set, return an unauthorized status

	// 		w.WriteHeader(http.StatusUnauthorized)
	// 		return
	// 	}
	// 	// For any other type of error, return a bad request status
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	// Get the JWT string from the cookie
	//tknStr := c.Value
	//fmt.Println("tknStr=",tknStr)
	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	// fmt.Println("tkn=", tkn)
	// fmt.Println("err", err)
	if err != nil {
		//if err == jwt.ErrSignatureInvalid {
		http.Error(w, "Token signature is invalid", http.StatusUnauthorized)
		//w.WriteHeader(http.StatusUnauthorized)
		return false
		//}
		// w.WriteHeader(http.StatusBadRequest)
		// return
	}
	if !tkn.Valid {
		http.Error(w, "Token Invalid", http.StatusUnauthorized)
		//w.WriteHeader(http.StatusUnauthorized)
		return false
	}
	// Finally, return the welcome message to the user, along with their
	// username given in the token
	fmt.Println("Token is valid")
	//w.Write([]byte(fmt.Sprintf("Welcome %s : You are Accessing authenticated API! \n", claims.Username)))
	return true
}

func extractTokenFromHeader(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if len(bearerToken) > 7 && bearerToken[:7] == "Bearer " {
		return bearerToken[7:]
	}

	return ""
}
