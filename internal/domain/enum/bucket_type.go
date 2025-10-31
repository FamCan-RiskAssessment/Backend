package enum

type BucketType uint

const (
	BucketTypeMamography BucketType = iota + 1
)

func (bt BucketType) String() string {
	switch bt {
	case BucketTypeMamography:
		return "ماموگرافی"
	}
	return "unknown"
}

func GetAllBucketTypes() []BucketType {
	return []BucketType{
		BucketTypeMamography,
	}
}
