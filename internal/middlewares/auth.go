package middlewares

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
)

var jwtKey = []byte("your_secret_key")

func JWTAuthMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		tokenString := string(ctx.GetHeader("Authorization"))
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Authorization header 缺失"})
			ctx.Abort()
			return
		}

		// 判断签名方法是否是HMAC
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtKey, nil
		})

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "token解析失败"})
			ctx.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx.Set("uid", claims["uid"])
		} else {
			ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "无效的token claims"})
			ctx.Abort()
			return
		}

		ctx.Next(c)
	}
}
