package middleware

import (
	"errors"
	"net/http"
	"scheduler/internal/app/models"
	"scheduler/utils"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func RoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")
		authHeader := strings.Split(bearerToken, ":")
		if len(authHeader) != 2 {
			handleJWTValidationError(c)
			return
		}
		token := authHeader[1]
		if token == "" {
			utils.Fail(c, http.StatusUnauthorized, utils.InvalidTokenError)
			return
		}

		accountRole, userID, err := GetUserIDAndAccountRoleFromToken(token)
		if err != nil {
			utils.Fail(c, http.StatusUnauthorized, utils.InvalidTokenError)
			return
		}

		c.Set("currentUserID", userID)

		for _, role := range requiredRoles {
			if accountRole == role {
				c.Next()
				return
			}
		}

		utils.Fail(c, http.StatusUnauthorized, "Access Denied")
		c.Abort()
	}
}

// VerifyToken : to Verify JWT Token
func VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")
		authHeader := strings.Split(bearerToken, ":")
		if len(authHeader) != 2 {
			handleJWTValidationError(c)
			return
		}
		jwtToken := authHeader[1]
		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(utils.GetEnvKey("JWTKEY")), nil
		})
		if err != nil {
			handleJWTValidationError(c)
			return
		}
		if !token.Valid {
			handleJWTValidationError(c)
			return
		}
		setUserContext(c, claims)
		c.Next()
	}
}

func handleJWTValidationError(c *gin.Context) {
	utils.Fail(c, http.StatusBadRequest, utils.InvalidTokenError)
	c.Abort()
}

func setUserContext(c *gin.Context, claims *models.Claims) {
	c.Set("id", strconv.FormatUint(uint64(claims.ID), 10))
	c.Set("email", claims.Email)
}

func GetUserIDAndAccountRoleFromToken(tokenString string) (string, uint, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.GetEnvKey("JWTKEY")), nil
	})

	if err != nil {
		return "", 0, err
	}

	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		return claims.Role, claims.ID, nil
	}

	return "", 0, errors.New("invalid token or claims")
}
