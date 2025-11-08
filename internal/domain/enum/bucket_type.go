package enum

type BucketType uint

const (
	BucketTypeMamography BucketType = iota + 1
	BucketTypeCancer
)

func (bt BucketType) String() string {
	switch bt {
	case BucketTypeMamography:
		return "ماموگرافی"
	case BucketTypeCancer:
		return "سرطان"
	}
	return "unknown"
}

func GetAllBucketTypes() []BucketType {
	return []BucketType{
		BucketTypeMamography,
		BucketTypeCancer,
	}
}
