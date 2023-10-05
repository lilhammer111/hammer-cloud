package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//username := c.Request.FormValue("username")
		token := c.GetHeader("Authorization")
		//err := r.ParseForm()
		//if err != nil {
		//	http.Error(w, "wrong parameter", http.StatusBadRequest)
		//	return
		//}
		//username := r.Form.Get("username")
		if !IsTokenValid(token) {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{"message": "unauthorized"})
			return
		}
		c.Next()
	}
}

func IsTokenValid(token string) bool {
	// TODO to judge if token expires, and validate token
	if len(token) != 40 {
		return false
	}
	return true
}
