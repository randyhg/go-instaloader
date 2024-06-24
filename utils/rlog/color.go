package rlog

import "fmt"

type FontStyle int8

const (
	ForegroundBlack FontStyle = iota + 30
	ForegroundRed
	ForegroundGreen
	ForegroundYellow
	ForegroundBlue
	ForegroundMagenta
	ForegroundCyan
	ForegroundWhite
	ForegroundDefault FontStyle = 39
)

const (
	ForegroundDarkGray FontStyle = iota + 90
	ForegroundLightRed
	ForegroundLightGreen
	ForegroundLightYellow
	ForegroundLightBlue
	ForegroundLightMagenta
	ForegroundLightCyan
	ForegroundLightWhite
)

const (
	Reset FontStyle = iota
	Bold
	Fuzzy
	Italic
	Underscore
	Blink
	FastBlink
	Reverse
	Concealed
	Strikethrough
)

const (
	BackgroundBlack FontStyle = iota + 40
	BackgroundRed
	BackgroundGreen
	BackgroundYellow
	BackgroundBlue
	BackgroundMagenta
	BackgroundCyan
	BackgroundWhite
	BackgroundDefault FontStyle = 49
)

const (
	BackgroundDarkGray FontStyle = iota + 100
	BackgroundLightRed
	BackgroundLightGreen
	BackgroundLightYellow
	BackgroundLightBlue
	BackgroundLightMagenta
	BackgroundLightCyan
	BackgroundLightWhite
)

const tpl = "\x1b[%dm"

var reset = Reset.String()

func (s FontStyle) String() string {
	return fmt.Sprintf(tpl, s)
}

func PrintWithColor(message interface{}, styles ...FontStyle) string {
	var res string
	for _, style := range styles {
		res += style.String()
	}

	return res + fmt.Sprint(message) + reset
}
