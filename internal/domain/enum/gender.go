package enum

import (
	"strconv"
	"strings"
)

type Gender uint

const (
	GenderMale Gender = iota + 1
	GenderFemale
	GenderOther
)

func (g Gender) String() string {
	switch g {
	case GenderMale:
		return "مرد"
	case GenderFemale:
		return "زن"
	case GenderOther:
		return "سایر"
	}
	return "نامشخص"
}

func GetAllGenders() []Gender {
	return []Gender{
		GenderMale,
		GenderFemale,
		GenderOther,
	}
}

// ParseGenderFilter accepts numeric IDs (1–3), English names (male/female/other), or Persian labels.
func ParseGenderFilter(value string) (Gender, bool) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	switch normalized {
	case "1", "male":
		return GenderMale, true
	case "2", "female":
		return GenderFemale, true
	case "3", "other":
		return GenderOther, true
	}
	switch strings.TrimSpace(value) {
	case GenderMale.String():
		return GenderMale, true
	case GenderFemale.String():
		return GenderFemale, true
	case GenderOther.String():
		return GenderOther, true
	}
	if id, err := strconv.ParseUint(normalized, 10, 32); err == nil {
		gender := Gender(id)
		if gender >= GenderMale && gender <= GenderOther {
			return gender, true
		}
	}
	return 0, false
}
