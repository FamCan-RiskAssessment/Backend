package enum

type LifeStatus uint

const (
	Alive LifeStatus = iota + 1
	Deceased
)

func (g LifeStatus) String() string {
	switch g {
	case Alive:
		return "در قید حیات"
	case Deceased:
		return "فوت شده"
	}
	return "نامشخص"
}

func GetAllLifeStatus() []LifeStatus {
	return []LifeStatus{
		Alive,
		Deceased,
	}
}
