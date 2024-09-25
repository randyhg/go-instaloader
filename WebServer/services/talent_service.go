package services

import (
	"fmt"
	"go-instaloader/WebServer/caches"
	"go-instaloader/config"
	"go-instaloader/models"
	"go-instaloader/utils/myDb"
	"go-instaloader/utils/rlog"
	"gorm.io/gorm/clause"
	"net/http"
	"net/url"
	"path/filepath"
	"runtime"
)

var TalentService = new(talentService)

type talentService struct{}

func (t *talentService) GetTalentList(tableName string, page, limit int) []*models.Talent {
	talents := caches.TalentCache.GetAllTalents(tableName)
	if len(talents) == 0 {
		return make([]*models.Talent, 0)
	}

	// pagination
	start := (page - 1) * limit
	end := start + limit
	return talents[start:end]
}

func (t *talentService) UpsertTalentData(talent *models.Talent) error {
	tableName := myDb.GetMonthTableName(models.Talent{})

	if err := myDb.GetDb().Table(tableName).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: talent.SheetId}},
		DoUpdates: clause.AssignmentColumns([]string{"status", "story_img_url", "story_img_path", "url"}),
	}).Create(&talent).Error; err != nil {
		rlog.Error(err)
		return err
	}
	caches.TalentCache.Invalidate(talent.Username)
	caches.TalentCache.InvalidateAllTalents(tableName)
	return nil
}

func (t *talentService) UpdateTalentData(talent *models.Talent) error {
	tableName := myDb.GetMonthTableName(models.Talent{})
	var existingTalent *models.Talent
	if err := myDb.GetDb().Table(tableName).First(&existingTalent, talent.Id).Error; err != nil {
		rlog.Error(err)
		return err
	}
	existingTalent.Username = talent.Username
	existingTalent.Url = talent.Url
	existingTalent.Status = talent.Status
	existingTalent.StoryImgUrl = talent.StoryImgUrl

	if err := myDb.GetDb().Table(tableName).Save(&existingTalent).Error; err != nil {
		rlog.Error(err)
		return err
	}
	caches.TalentCache.Invalidate(talent.Username)
	caches.TalentCache.InvalidateAllTalents(tableName)
	return nil
}

func (c *talentService) DeleteTalentData(talent *models.Talent) error {
	tableName := myDb.GetMonthTableName(models.Talent{})
	err := myDb.GetDb().Table(tableName).Where("username = ?", talent.Username).Delete(&models.Talent{}).Error
	if err != nil {
		rlog.Error(err)
		return err
	}
	caches.TalentCache.Invalidate(talent.Username)
	caches.TalentCache.InvalidateAllTalents(tableName)
	return nil
}

func ErrorHandler(err error) {
	_, file, line, _ := runtime.Caller(1)
	file = filepath.Base(file)

	// send error log to tg
	errorMsg := fmt.Sprintf("*error occurred on go-instaloader*\n\n`%s:%d:`\n%s", file, line, err.Error())
	fullUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.Instance.TeleBotToken)
	data := url.Values{}
	data.Set("chat_id", config.Instance.TeleGroupId)
	data.Set("text", errorMsg)
	data.Set("parse_mode", "Markdown")

	_, err = http.PostForm(fullUrl, data)
	if err != nil {
		rlog.Error(err)
		return
	}

	//rlog.Info("sent successfully!")
}
