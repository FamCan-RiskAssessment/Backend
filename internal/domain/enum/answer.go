package enum

type Answer uint

const (
	AnswerYes Answer = iota + 1
	AnswerNo
	AnswerDontKnow
)

func (g Answer) String() string {
	switch g {
	case AnswerYes:
		return "بله"
	case AnswerNo:
		return "خیر"
	case AnswerDontKnow:
		return "نمی دانم"
	}
	return "نامشخص"
}

func GetAllAnswer() []Answer {
	return []Answer{
		AnswerYes,
		AnswerNo,
		AnswerDontKnow,
	}
}
