package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB, rdb *redis.Client) *gin.Engine {
	r := gin.Default()

	// TODO: 在这里添加路由配置

	return r
}
