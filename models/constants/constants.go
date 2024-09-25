package constants

import "time"

const AllTalentCacheExp = 300 * time.Second

const (
	AllTalentCacheKey = "cache/talent/all"
	TalentCacheKey    = "cache/talent"
)

type ConfigColumn int

const (
	UsernameCol ConfigColumn = iota
	StatusCol
	RemarkCol
	TalentUrlCol
	TalentCountCol
)
