package utils

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"transaction_api/model"
	"transaction_api/service"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddlewareGin(src *service.UserService) gin.HandlerFunc {

	return func(c *gin.Context) {
		var payload struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusUnauthorized, model.Response{
				ResponseCode: "0002",
				Message:      "[Failed] Invalid request body",
			})
			c.Abort()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		if err := sonic.Unmarshal(bodyBytes, &payload); err != nil {
			c.JSON(http.StatusUnauthorized, model.Response{
				ResponseCode: "0002",
				Message:      "[Failed] Invalid username or password",
			})
			c.Abort()
			return
		}

		log.Println("REQ PATH:", c.Request.Method, c.Request.URL.Path)

		authHeader := c.GetHeader("Authorization")

		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			c.JSON(http.StatusUnauthorized, model.Response{
				ResponseCode: "0002",
				Message:      "Unauthorized",
			})
			c.Abort()
			return
		}

		user, err := src.GetUserDetail(c, payload.Username)
		if err != nil || user == nil {
			c.JSON(http.StatusUnauthorized, model.Response{
				ResponseCode: "0002",
				Message:      "[Failed] Username is invalid",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, prefix)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// pastikan algoritma signing sesuai (mencegah alg confusion attack)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, model.Response{
				ResponseCode: "0002",
				Message:      "[Failed] Token is invalid or expire",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, model.Response{
				ResponseCode: "0002",
				Message:      "[Failed] Token is invalid",
			})
			c.Abort()
			return
		}

		if tokenString != *user.Token {
			c.JSON(http.StatusUnauthorized, model.Response{
				ResponseCode: "0002",
				Message:      "Invalid token",
			})
			c.Abort()
			return
		}

		// Lanjut ke handler berikutnya
		c.Set("user", claims)
		c.Set("currentUser", user)

		c.Next()

	}
}
