package router

import (
	"github.com/gin-contrib/sessions"
	sessions_redis "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"go-ranking/config"
	"go-ranking/controller"
	"go-ranking/pkg/logger"
)

func Router() *gin.Engine {
	r := gin.Default()

	r.Use(gin.LoggerWithConfig(logger.LoggerToFile()))
	r.Use(logger.Recover)
	store, _ := sessions_redis.NewStore(10, "tcp", config.RedisAddress, "", []byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	user := r.Group("/user")
	{
		user.POST("/register", controller.UserController{}.Register)

		user.POST("/login", controller.UserController{}.Login)

		user.GET("/info/:id", controller.UserController{}.GetUserInfo)

		user.GET("/list", controller.UserController{}.GetUserList)
	}

	player := r.Group("/player")
	{
		player.GET("/list", controller.PlayerController{}.GetPlayers)

		player.POST("/list/aid", controller.PlayerController{}.GetPlayersByAid)

		player.POST("/ranking", controller.PlayerController{}.GetRanking)
	}

	vote := r.Group("/vote")
	{
		vote.POST("/add", controller.VoteController{}.AddVote)
	}

	return r
}
