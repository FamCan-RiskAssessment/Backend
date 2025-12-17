package enum

type FormType uint

const (
	Bahar FormType = iota + 1
	Navid
)

func (g FormType) String() string {
	switch g {
	case Bahar:
		return "فرم بهار"
	case Navid:
		return "فرم نوید"
	}
	return "فرم نامشخص"
}

func GetAllFormTypes() []FormType {
	return []FormType{
		Bahar,
		Navid,
	}
}
