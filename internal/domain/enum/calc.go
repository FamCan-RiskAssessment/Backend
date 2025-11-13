package enum

type Calc uint

const (
	CalcPremm5 Calc = iota + 1
	CalcBCRA
	CalcGail
	CalcPLCO
)

func (c Calc) String() string {
	switch c {
	case CalcPremm5:
		return "premm5"
	case CalcBCRA:
		return "bcra"
	case CalcGail:
		return "gail"
	case CalcPLCO:
		return "plco"
	}
	return "unknown"
}

func GetAllCalcs() []Calc {
	return []Calc{
		CalcPremm5,
		CalcBCRA,
		CalcGail,
		CalcPLCO,
	}
}
