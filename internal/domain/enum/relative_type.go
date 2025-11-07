package enum

type Relative uint

const (
	Father Relative = iota + 1
	Mother
	Brother
	Sister
	PaternalGrandFather
	PaternalGrandMother
	MaternalGrandFather
	MaternalGrandMother
	PaternalUncle
	PaternalAunt
	MaternalUncle
	MaternalAunt
	Child
	DistantRelative
)

func (g Relative) String() string {
	switch g {
	case Father:
		return "پدر"
	case Mother:
		return "مادر"
	case Brother:
		return "برادر"
	case Sister:
		return "خواهر"
	case PaternalGrandFather:
		return "پدربزرگ پدری"
	case PaternalGrandMother:
		return "مادربزرگ پدری"
	case MaternalGrandFather:
		return "پدربزرگ مادری"
	case MaternalGrandMother:
		return "مادربزرگ مادری"
	case PaternalUncle:
		return "عمو"
	case PaternalAunt:
		return "عمه"
	case MaternalUncle:
		return "دایی"
	case MaternalAunt:
		return "خاله"
	case Child:
		return "فرزند"
	case DistantRelative:
		return "فامیل دور"
	}
	return "نامشخص"
}

func GetAllRelatives() []Relative {
	return []Relative{
		Father,
		Mother,
		Brother,
		Sister,
		PaternalGrandFather,
		PaternalGrandMother,
		MaternalGrandFather,
		MaternalGrandMother,
		PaternalUncle,
		PaternalAunt,
		MaternalUncle,
		MaternalAunt,
		Child,
		DistantRelative,
	}
}
