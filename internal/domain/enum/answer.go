package enum

type Answer uint

const (
	AnswerYes Answer = iota + 1
	AnswerNo
	AnswerDontKnow
	AnswerUncertain
	AnswerNoInfo
	AnswerPossible
	AnswerAgo
	AnswerLongAgo
)

func (g Answer) String() string {
	switch g {
	case AnswerYes:
		return "بله"
	case AnswerNo:
		return "خیر"
	case AnswerDontKnow:
		return "نمی دانم"
	case AnswerUncertain:
		return "نامعین"
	case AnswerNoInfo:
		return "اطلاع ندارم"
	case AnswerPossible:
		return "احتمال دارد اما دقیق اطلاع ندارم"
	case AnswerAgo:
		return "سابقاً مصرف می کرده ام اما کمتر از ۱۵ سال است که ترک کرده ام"
	case AnswerLongAgo:
		return "سابقاً مصرف می کردم اما بیش از ۱۵ سال است که ترک کرده ام"
	}
	return "وارد نشده"
}

func GetAllAnswers() []Answer {
	return []Answer{
		AnswerYes,
		AnswerNo,
		AnswerDontKnow,
		AnswerUncertain,
		AnswerNoInfo,
		AnswerPossible,
		AnswerAgo,
		AnswerLongAgo,
	}
}
