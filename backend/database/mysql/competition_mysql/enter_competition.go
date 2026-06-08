package competition_mysql

import (
	"GO1/global"
	"GO1/models/competition_model"

	"gorm.io/gorm"
)

func GetCompetitionAuthInfo(competitionID int64) (competition competition_model.Competition, exist bool, err error) {
	tx := global.DB.
		Select("id, visibility, password_hash").
		Where("id = ?", competitionID).
		Find(&competition)

	if tx.Error != nil {
		return competition, false, tx.Error
	}

	return competition, tx.RowsAffected > 0, nil
}

func EnterCompetition(enterPlayer *competition_model.CompetitionPlayer) (bool, error) {
	tx := global.DB.
		Model(&competition_model.Competition{}).
		Where("id = ?", enterPlayer.CompetitionID).
		Update("player_count", gorm.Expr("player_count + ?", 1))

	if tx.Error != nil {
		return false, tx.Error
	}

	if tx.RowsAffected == 0 {
		return false, nil
	}

	if err := global.DB.
		Model(&competition_model.CompetitionPlayer{}).
		Create(enterPlayer).Error; err != nil {
		return false, err
	}

	return true, nil
}
