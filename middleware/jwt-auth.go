package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/myomyintko/strategy_robot/helper"
	"github.com/myomyintko/strategy_robot/service"
	"log"
	"net/http"
	"strings"
)

//AuthorizeJWT validates the token user given, return 401 if not valid
func AuthorizeJWT(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			response := helper.BuildErrorResponse("Failed to process request", "No token found", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		splitToken := strings.Split(authHeader, "Bearer ")
		authHeader = splitToken[1]
		token, err := jwtService.ValidateToken(authHeader)
		if err != nil {
			return
		}
		if !token.Valid {
			response := helper.BuildErrorResponse("Token is not valid", "Token error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		log.Println("Claim[user_id]: ", claims["user_id"])
	}
}
