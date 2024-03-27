package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"go-ranking/cache"
	"go-ranking/models"
	"go-ranking/utils"
)

type VoteController struct{}

func (v VoteController) AddVote(ctx *gin.Context) {
	// 获取用户Id（投票人Id）选手Id
	params := make(map[string]interface{})
	err := ctx.BindJSON(&params)
	if err != nil {
		ReturnError(ctx, 5000, "服务内部错误："+err.Error())
		return
	}
	userId := int(utils.GetOrDefault(params, "userId", 0).(float64))
	playerId := int(utils.GetOrDefault(params, "playerId", 0).(float64))

	if userId == 0 || playerId == 0 {
		ReturnError(ctx, 4001, "请输入正确信息")
		return
	}

	user, _ := models.GetUserById(userId)
	if user.Id == 0 {
		ReturnError(ctx, 4004, "投票用户不存在")
		return
	}
	player, _ := models.GetPlayerById(playerId)
	if player.Id == 0 {
		ReturnError(ctx, 4004, "参赛选手不存在")
		return
	}
	vote, _ := models.GetVoteInfo(userId, playerId)
	if vote.Id != 0 {
		ReturnError(ctx, 4001, "已投票")
		return
	}
	rs, err := models.AddVote(userId, playerId)
	if err == nil {
		// 更新参赛选手分数字段，自增1
		err = models.UpdatePlayerScore(playerId)
		if err != nil {
			ReturnError(ctx, 4002, "更新选手分数失败："+err.Error())
			return
		}
		// 同时更新redis
		redisKey := "ranking:" + strconv.Itoa(player.Aid)
		cache.Rdb.ZIncrBy(cache.Rctx, redisKey, 1, strconv.Itoa(playerId))
		ReturnSuccess(ctx, 0, "投票成功", rs, 1)
		return
	}
	ReturnError(ctx, 4001, "投票失败，请联系管理员")
}
