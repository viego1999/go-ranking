package controller

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go-ranking/cache"
	"go-ranking/models"
)

type PlayerController struct{}

func (p PlayerController) GetPlayers(ctx *gin.Context) {
	players, err := models.GetPlayers()
	if err != nil {
		ReturnError(ctx, 4004, "没有相关信息")
		return
	}
	ReturnSuccess(ctx, 0, "success", players, int64(len(players)))
}

func (p PlayerController) GetPlayersByAid(ctx *gin.Context) {
	aidStr := ctx.DefaultPostForm("aid", "0")
	aid, _ := strconv.Atoi(aidStr)

	rs, err := models.GetPlayersByAid(aid, "id asc")
	if err != nil {
		ReturnError(ctx, 4004, "没有相关信息")
		return
	}
	ReturnSuccess(ctx, 0, "success", rs, int64(len(rs)))
}

func (p PlayerController) GetRanking(ctx *gin.Context) {
	aidStr := ctx.DefaultPostForm("aid", "0")
	aid, _ := strconv.Atoi(aidStr)

	// 先查询 redis
	redisKey := "ranking:" + aidStr
	rs, err := cache.Rdb.ZRevRange(cache.Rctx, redisKey, 0, -1).Result()
	if err == nil && len(rs) > 0 {
		var players []models.Player
		for _, val := range rs {
			id, _ := strconv.Atoi(val)
			player, _ := models.GetPlayerById(id)
			if player.Id > 0 {
				players = append(players, player)
			}
		}
		ReturnSuccess(ctx, 0, "success", players, int64(len(players)))
		return
	}
	// 当 redis 中不存在，从数据库中进行查询
	rsDb, errDb := models.GetPlayersByAid(aid, "score desc")
	if errDb == nil {
		// 将记录插入 redis 中
		for _, value := range rsDb {
			cache.Rdb.ZAdd(cache.Rctx, redisKey, cache.Zscore(value.Id, value.Score)).Err()
		}
		// 设置过期时间
		cache.Rdb.Expire(cache.Rctx, redisKey, 24*time.Hour)
		ReturnSuccess(ctx, 0, "success", rsDb, int64(len(rsDb)))
		return
	}
	ReturnError(ctx, 4004, "没有相关信息")
}
