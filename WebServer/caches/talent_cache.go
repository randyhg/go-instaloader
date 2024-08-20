package caches

import (
	"errors"
	"go-instaloader/WebServer/caches/redis_cache"
	"go-instaloader/models"
	"go-instaloader/models/constants"
	"go-instaloader/utils/myDb"
	"go-instaloader/utils/rlog"
	"gorm.io/gorm"
	"time"
)

const TalentCacheTimeOut = 3600 * 1

type talentCache struct {
	cache *redis_cache.RedisCache
}

var TalentCache = newTalentCache()

func newTalentCache() *talentCache {
	c := &talentCache{}
	c.cache = redis_cache.NewRedisCacheWithOption(
		constants.TalentCacheKey,
		c.loader,
		redis_cache.WithExpireAfterWrite(TalentCacheTimeOut*time.Second),
	)

	return c
}

func (c *talentCache) loader(key redis_cache.Key) (value redis_cache.Value, err error) {
	username := key2String(key)
	ret := &models.Talent{}
	tableName := myDb.GetMonthTableName(models.Talent{})

	if err = myDb.GetDb().Table(tableName).First(&ret, "username = ?", username).Error; err != nil {
		rlog.Error(err)
		return nil, err
	}

	return ret, nil
}

func (c *talentCache) Get(username string) *models.Talent {
	out := &models.Talent{}
	val, err := c.cache.GetString(username, out)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			rlog.Error(err, username)
		}
		return nil
	}
	if val == nil {
		return nil
	}

	user := val.(*models.Talent)
	if user != nil {
		user.StoryImgUrl = user.GetStoryUrls()
	}
	return user
}

func (c *talentCache) GetArray(usernames []string) []*models.Talent {
	var talentList []*models.Talent
	for _, username := range usernames {
		out := &models.Talent{}
		ptr, _ := c.cache.GetString(username, out)
		if ptr != nil {
			user := ptr.(*models.Talent)
			talentList = append(talentList, user)
		}
	}

	return talentList
}

func (c *talentCache) Invalidate(username string) {
	c.cache.InvalidateString(username)
}
