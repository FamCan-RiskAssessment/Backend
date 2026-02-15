package enum

type Relative uint

const (
	Father Relative = iota + 1
	Mother
	Brother
	Sister
	HalfBrother
	HalfSister
	PaternalGrandFather
	PaternalGrandMother
	MaternalGrandFather
	MaternalGrandMother
	PaternalUncle
	PaternalAunt
	MaternalUncle
	MaternalAunt
	Daughter
	Son
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
	case HalfBrother:
		return "برادر ناتنی"
	case HalfSister:
		return "خواهر ناتنی"
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
	case Daughter:
		return "فرزند دختر"
	case Son:
		return "فرزند پسر"
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
		HalfBrother,
		HalfSister,
		PaternalGrandFather,
		PaternalGrandMother,
		MaternalGrandFather,
		MaternalGrandMother,
		PaternalUncle,
		PaternalAunt,
		MaternalUncle,
		MaternalAunt,
		Daughter,
		Son,
		DistantRelative,
	}
}
