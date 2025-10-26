package enum

type HyperplasiaInBiopsyStatus uint

const (
	NoHyperplasiaInBiopsy HyperplasiaInBiopsyStatus = iota
	HasHyperplasiaInBiopsy
	UnkownHyperplasiaInBiopsy = iota + 97
)

func (HIB HyperplasiaInBiopsyStatus) String() string {
	switch HIB {
	case NoHyperplasiaInBiopsy:
		return "فاقد هایپرپلازی"
	case HasHyperplasiaInBiopsy:
		return "داری هایپرپلازی"
	case UnkownHyperplasiaInBiopsy:
		return "وضعیت هایپرپلازی نامشخص"
	}
	return "نامشخص"
}

func GetAllHyperplasiaInBiopsyStatuses() []HyperplasiaInBiopsyStatus {
	return []HyperplasiaInBiopsyStatus{
		NoHyperplasiaInBiopsy,
		HasHyperplasiaInBiopsy,
		UnkownHyperplasiaInBiopsy,
	}
}
