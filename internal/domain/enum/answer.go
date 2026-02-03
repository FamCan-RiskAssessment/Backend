package enum

type Answer uint

const (
	AnswerYes Answer = iota + 1
	AnswerNo
	AnswerDontKnow
	LongAgo
)

func (g Answer) String() string {
	switch g {
	case AnswerYes:
		return "بله"
	case AnswerNo:
		return "خیر"
	case AnswerDontKnow:
		return "نمی دانم"
	case LongAgo:
		return "بیشتر از ۱۵ سال است که مصرف نکردم"
	}
	return "نامشخص"
}

func GetAllAnswers() []Answer {
	return []Answer{
		AnswerYes,
		AnswerNo,
		AnswerDontKnow,
		LongAgo,
	}
}
