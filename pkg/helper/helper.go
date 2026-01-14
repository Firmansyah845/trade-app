package helper

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
)

func EpochStringToFormattedTime(epochString string, layout string) (string, error) {

	epochSeconds, err := strconv.ParseInt(epochString, 10, 64)
	if err != nil {
		return "", err
	}

	epochTime := time.Unix(epochSeconds, 0)

	formattedTime := epochTime.Format(layout)

	return formattedTime, nil
}

func FilterBadWords(message string, badWords []string) string {
	for _, word := range badWords {
		re := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(word) + `\b`)
		//replace := string(word[0]) + strings.Repeat("*", len(word)-2) + string(word[len(word)-1])
		replace := strings.Repeat("*", len(word))

		message = re.ReplaceAllString(message, replace)
	}
	return message
}

func FilterScriptTag(message string) string {
	message = strings.ReplaceAll(message, "\n", "")
	message = strings.ReplaceAll(message, "\t", "")
	message = bluemonday.NewPolicy().AllowElementsContent("script", "iframe").AllowUnsafe(true).Sanitize(message)
	return message
}

func InArrayInt(input int, allowed []int) bool {
	for _, value := range allowed {
		if input == value {
			return true
		}
	}
	return false
}
