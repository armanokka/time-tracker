package projects

import (
	"github.com/gin-gonic/gin"
)

type Handlers interface {
	GetByID() gin.HandlerFunc
	Create() gin.HandlerFunc
	Update() gin.HandlerFunc
	Delete() gin.HandlerFunc
	GetMembers() gin.HandlerFunc
	AddMember() gin.HandlerFunc
	RemoveMember() gin.HandlerFunc
	GetMemberProductivity() gin.HandlerFunc
}

type TaskHandlers interface {
	Get() gin.HandlerFunc
	Create() gin.HandlerFunc
	Update() gin.HandlerFunc
	Delete() gin.HandlerFunc

	Start() gin.HandlerFunc
	Stop() gin.HandlerFunc

	GetMembers() gin.HandlerFunc
	AddMember() gin.HandlerFunc
	DeleteMember() gin.HandlerFunc
}
