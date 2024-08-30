package constants

import "time"

const AllTalentCacheExp = 300 * time.Second

const (
	AllTalentCacheKey = "cache/talent/all"
	TalentCacheKey    = "cache/talent"
)

const (
	// event name
	ErrorNotice = "errorNotice"
	TestEvent   = "test"
	FetchEvent  = "fetch"
	VerifyEvent = "verify"
)
