package enum

type LifeStatus uint

const (
	Alive LifeStatus = iota + 1
	Dead
)

func (g LifeStatus) String() string {
	switch g {
	case Alive:
		return "در قید حیات"
	case Dead:
		return "فوت شده"
	}
	return "نامشخص"
}

func GetAllLifeStatuss() []LifeStatus {
	return []LifeStatus{
		Alive,
		Dead,
	}
}
