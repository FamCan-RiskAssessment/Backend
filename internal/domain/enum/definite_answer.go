package enum

type Answer uint

const (
	Yes Answer = iota + 1
	No
	DontKnow
)

func (g Answer) String() string {
	switch g {
	case Yes:
		return "آری"
	case No:
		return "خیر"
	case DontKnow:
		return "نمی دانم"
	}
	return "نامشخص"
}

func GetAllAnswers() []Answer {
	return []Answer{
		Yes,
		No,
		DontKnow,
	}
}
