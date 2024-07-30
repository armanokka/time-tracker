package auth

import (
	"github.com/gin-gonic/gin"
)

type Handlers interface {
	Register() gin.HandlerFunc
	Login() gin.HandlerFunc
	//Logout() gin.HandlerFunc // TODO
	Update() gin.HandlerFunc
	Delete() gin.HandlerFunc
	GetUserByID() gin.HandlerFunc
	SearchUsers() gin.HandlerFunc
}
