package models

import (
	"go-ranking/dao"
	"gorm.io/gorm"
)

type Player struct {
	Id          int    `json:"id"`
	Aid         int    `json:"aid"`
	Ref         string `json:"ref"`
	Nickname    string `json:"nickname"`
	Declaration string `json:"declaration"`
	Avatar      string `json:"avatar"`
	Score       int    `json:"score"`
}

func (Player) TableName() string {
	return "player"
}

func GetPlayers() ([]Player, error) {
	var players []Player
	err := dao.Db.Find(&players).Error
	return players, err
}

func GetPlayerById(id int) (Player, error) {
	var player Player
	err := dao.Db.Where("id = ?", id).Find(&player).Error
	return player, err
}

func GetPlayersByAid(aid int, sort string) ([]Player, error) {
	var players []Player
	err := dao.Db.Where("aid = ?", aid).Order(sort).Find(&players).Error
	return players, err
}

func UpdatePlayerScore(id int) error {
	var player Player
	err := dao.Db.Model(&player).Where("id=?", id).UpdateColumn("score", gorm.Expr("score+?", 1)).Error
	return err
}