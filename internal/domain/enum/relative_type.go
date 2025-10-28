package enum

type Relative uint

const (
	Self Relative = iota + 1
	Father
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
)

func (g Relative) String() string {
	switch g {
	case Self:
		return "خود"
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
	}
	return "نامشخص"
}

func GetAllRelatives() []Relative {
	return []Relative{
		Self,
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
	}
}
