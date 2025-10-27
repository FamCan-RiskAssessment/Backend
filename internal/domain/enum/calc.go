package enum

type Calc uint

const (
	CalcPremm5 Calc = iota + 1
	CalcBCRA
	CalcGBR
)

func (c Calc) String() string {
	switch c {
	case CalcPremm5:
		return "premm5"
	case CalcBCRA:
		return "bcra"
	case CalcGBR:
		return "gbr"
	}
	return "unknown"
}

func GetAllCalcs() []Calc {
	return []Calc{
		CalcPremm5,
		CalcBCRA,
		CalcGBR,
	}
}
