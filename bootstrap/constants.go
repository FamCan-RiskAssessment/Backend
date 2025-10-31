package bootstrap

import "fmt"

type Constants struct {
	Context      Context
	RedisKey     RedisKey
	S3BucketPath BucketPath
	Field        Field
	BucketPath   BucketPath
	Tag          Tag
	SMSTemplates SMSTemplates
	JWTKeysPath  JWTKeysPath
}

type Context struct {
	Translator string
	ID         string
}

type RedisKey struct {
}

type BucketPath struct {
}

type Field struct {
	User             string
	Form             string
	OTP              string
	Role             string
	Page             string
	Permission       string
	Premm5Result     string
	BCRAResult       string
	MamoGraphyInfo   string
	FamilyCancerInfo string
}

type Tag struct {
	Expired      string
	Invalid      string
	NotFound     string
	AlreadyExist string
}

type SMSTemplates struct {
	OTP string
}

type JWTKeysPath struct {
	PrivateKey string
	PublicKey  string
}

func NewConstants() *Constants {
	return &Constants{
		Context: Context{
			Translator: "translator",
			ID:         "ID",
		},
		Field: Field{
			User:             "user",
			Form:             "form",
			OTP:              "otp",
			Role:             "role",
			Page:             "page",
			Permission:       "permission",
			Premm5Result:     "premm5_result",
			BCRAResult:       "bcra_result",
			MamoGraphyInfo:   "mamography_info",
			FamilyCancerInfo: "family_cancer_info",
		},
		Tag: Tag{
			Expired:      "expired",
			Invalid:      "invalid",
			NotFound:     "not found",
			AlreadyExist: "already exist",
		},
		SMSTemplates: SMSTemplates{
			OTP: "otp",
		},
		JWTKeysPath: JWTKeysPath{
			PrivateKey: "internal/infrastructure/jwt/privateKey.pem",
			PublicKey:  "internal/infrastructure/jwt/publicKey.pem",
		},
	}
}

func (r *RedisKey) GenerateOTPKey(value string) string {
	return fmt.Sprintf("otp:%s", value)
}

func (path *BucketPath) GetMamoGraphyPath(formID uint, fileName string) string {
	return fmt.Sprintf("Radiology/%d/MamoGraphy/%s", formID, fileName)
}

func (r *RedisKey) GenerateOperatorValidationOTPKey(operatorID uint, phone string) string {
	return fmt.Sprintf("operator:validation:otp:%d:%s", operatorID, phone)
}

func (r *RedisKey) GenerateOperatorValidationTokenKey(operatorID uint, userID uint) string {
	return fmt.Sprintf("operator:validation:token:%d:%d", operatorID, userID)
}
