package enum

type BucketType uint

const (
	BucketTypeMamography BucketType = iota + 1
	BucketTypeCancer
	BucketTypeGeneticTest
	BucketTypeFatherGeneticTest
	BucketTypeMotherGeneticTest
)

func (bt BucketType) String() string {
	switch bt {
	case BucketTypeMamography:
		return "ماموگرافی"
	case BucketTypeCancer:
		return "سرطان"
	case BucketTypeGeneticTest:
		return "تست ژنتیک"
	case BucketTypeFatherGeneticTest:
		return "تست ژنتیک بابا"
	case BucketTypeMotherGeneticTest:
		return "تست ژنتیک مادر"
	}
	return "unknown"
}

func GetAllBucketTypes() []BucketType {
	return []BucketType{
		BucketTypeMamography,
		BucketTypeCancer,
		BucketTypeGeneticTest,
		BucketTypeFatherGeneticTest,
		BucketTypeMotherGeneticTest,
	}
}
