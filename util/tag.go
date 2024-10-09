package util

import (
	"regexp"
	"strings"
)

var (
	specialCharRe *regexp.Regexp
	multiDashRe   *regexp.Regexp
	ltDashRe      *regexp.Regexp
)

func init() {
	specialCharRe = regexp.MustCompile("(&quot;)|([^a-zA-Z-]+)")
	multiDashRe = regexp.MustCompile("-+")
	ltDashRe = regexp.MustCompile("(^-+)|(-+$)")
}

func InputToTag(input string) string {
	if input == "" {
		return ""
	}

	splitFn := func(c rune) bool {
		return c == ','
	}
	tagArray := strings.FieldsFunc(strings.ToLower(input), splitFn)

	tag := tagArray[0]
	tag = strings.Replace(tag, " ", "-", -1)
	tag = specialCharRe.ReplaceAllString(tag, "")
	tag = multiDashRe.ReplaceAllString(tag, "-")
	tag = ltDashRe.ReplaceAllString(tag, "")

	return tag
}

// GetTags transforms tags from the user inputs into a standard format.
// dog walking -> dog-walking (one word)
func FormatTags(tags []string) []string {
	encountered := map[string]bool{}
	formatted := make([]string, 0, len(tags))

	for _, tag := range tags {
		tag = strings.ToLower(tag)
		tag = strings.Replace(tag, " ", "-", -1)
		tag = specialCharRe.ReplaceAllString(tag, "")
		tag = multiDashRe.ReplaceAllString(tag, "-")
		tag = ltDashRe.ReplaceAllString(tag, "")
		if len(tag) == 0 {
			continue
		}
		// remove duplicates
		if !encountered[tag] {
			formatted = append(formatted, tag)
			encountered[tag] = true
		}
	}

	return formatted
}

// ToSearchTags transforms tags from user inputs into searching tags.
// dog walking -> dog, walking (two words)
func ToSearchTags(words string) []string {
	splitFn := func(c rune) bool {
		return c == ',' || c == ' '
	}
	tags := strings.FieldsFunc(strings.ToLower(words), splitFn)
	return FormatTags(tags)
}
