package services

import (
	"go-instaloader/WebServer/caches"
	"go-instaloader/models"
	"go-instaloader/utils/myDb"
	"go-instaloader/utils/rlog"
	"gorm.io/gorm/clause"
)

var TalentService = new(talentService)

type talentService struct{}

func (t *talentService) UpsertTalentData(talent *models.Talent) error {
	tableName := myDb.GetMonthTableName(models.Talent{})

	if err := myDb.GetDb().Table(tableName).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: talent.Uuid}},
		DoUpdates: clause.AssignmentColumns([]string{"status", "story_img_url"}),
	}).Create(&talent).Error; err != nil {
		rlog.Error(err)
		return err
	}
	caches.TalentCache.Invalidate(talent.Username)
	return nil
}

func (t *talentService) UpdateTalentData(talent *models.Talent) error {
	tableName := myDb.GetMonthTableName(models.Talent{})
	var existingTalent *models.Talent
	if err := myDb.GetDb().Table(tableName).First(&existingTalent, talent.Id).Error; err != nil {
		return err
	}
	existingTalent.Username = talent.Username
	existingTalent.Url = talent.Url
	existingTalent.Status = talent.Status
	existingTalent.StoryImgUrl = talent.StoryImgUrl

	err := myDb.GetDb().Table(tableName).Save(&existingTalent).Error
	return err
}
