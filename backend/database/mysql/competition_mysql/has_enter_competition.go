package competition_mysql

import (
	"GO1/global"
	"GO1/models/competition_model"
)

func HasEnterCompetition(competitionID, userID int64) (data []competition_model.CompetitionPlayer, err error) {
	err = global.DB.
		Select("id").
		Where("competition_id = ? AND user_id = ? AND is_cancel = ?", competitionID, userID, false).
		Find(&data).Error

	return
}
