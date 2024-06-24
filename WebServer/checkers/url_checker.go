package checkers

import (
	"go-instaloader/utils/rlog"
	"regexp"
)

const TheURL = "youtube.com"

func CheckUrl(theUrl, nodeUrl string) bool {
	rlog.Info("the URL:", theUrl)
	re := regexp.MustCompile(theUrl)
	if re.MatchString(nodeUrl) {
		return true
	}
	return false
}
