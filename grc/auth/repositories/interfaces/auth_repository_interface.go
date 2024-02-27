package interfaces

import "github.com/gin-gonic/gin"

type AuthRepositoryInterface interface {
	Authenticate(c *gin.Context, email string, password string) (map[string]string, error)
}
