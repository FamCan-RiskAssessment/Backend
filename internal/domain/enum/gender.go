package enum

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
