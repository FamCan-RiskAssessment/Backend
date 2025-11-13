package enum

type CancerType uint

const (
	CancerTypeBreast CancerType = iota + 1
	CancerTypeBreastBilateral
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
	CancerTypeEndometrial
	CancerTypeUterine
	CancerTypeBladder
	CancerTypeKidneyUrinary
	CancerTypeMetastatic
	CancerTypeOther
)

func (ct CancerType) String() string {
	switch ct {
	case CancerTypeBreast:
		return "سرطان پستان"
	case CancerTypeBreastBilateral:
		return "سرطان پستان دو بار (یا دوطرفه)"
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
		return "تومور مغزی"
	case CancerTypeLeukemia:
		return "لوسمی"
	case CancerTypeLymphoma:
		return "لنفوم"
	case CancerTypeEndometrial:
		return "آندومتر"
	case CancerTypeUterine:
		return "سرطان رحم"
	case CancerTypeBladder:
		return "سرطان مثانه"
	case CancerTypeKidneyUrinary:
		return "سرطان کلیه یا مجاری ادرار"
	case CancerTypeMetastatic:
		return "سرطان به صورت متاستاتیک بوده"
	case CancerTypeOther:
		return "سرطان های دیگر"
	}
	return "نامشخص"
}

func GetAllCancerTypes() []CancerType {
	return []CancerType{
		CancerTypeBreast,
		CancerTypeBreastBilateral,
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
		CancerTypeEndometrial,
		CancerTypeUterine,
		CancerTypeBladder,
		CancerTypeKidneyUrinary,
		CancerTypeMetastatic,
		CancerTypeOther,
	}
}
