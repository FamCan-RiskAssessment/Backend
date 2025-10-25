package enum

type CancerType uint

const (
	CancerTypeBreast CancerType = iota + 1
	CancerTypeLung
	CancerTypeColon
	CancerTypeProstate
	CancerTypeCervical
	CancerTypeOvarian
	CancerTypeLiver
	CancerTypeStomach
	CancerTypePancreatic
	CancerTypeBrain
	CancerTypeLeukemia
	CancerTypeLymphoma
	CancerTypeOther
)

func (ct CancerType) String() string {
	switch ct {
	case CancerTypeBreast:
		return "سرطان پستان"
	case CancerTypeLung:
		return "سرطان ریه"
	case CancerTypeColon:
		return "سرطان روده بزرگ"
	case CancerTypeProstate:
		return "سرطان پروستات"
	case CancerTypeCervical:
		return "سرطان دهانه رحم"
	case CancerTypeOvarian:
		return "سرطان تخمدان"
	case CancerTypeLiver:
		return "سرطان کبد"
	case CancerTypeStomach:
		return "سرطان معده"
	case CancerTypePancreatic:
		return "سرطان پانکراس"
	case CancerTypeBrain:
		return "سرطان مغز"
	case CancerTypeLeukemia:
		return "لوسمی"
	case CancerTypeLymphoma:
		return "لنفوم"
	case CancerTypeOther:
		return "سایر"
	}
	return "نامشخص"
}

func GetAllCancerTypes() []CancerType {
	return []CancerType{
		CancerTypeBreast,
		CancerTypeLung,
		CancerTypeColon,
		CancerTypeProstate,
		CancerTypeCervical,
		CancerTypeOvarian,
		CancerTypeLiver,
		CancerTypeStomach,
		CancerTypePancreatic,
		CancerTypeBrain,
		CancerTypeLeukemia,
		CancerTypeLymphoma,
		CancerTypeOther,
	}
}
