package bootstrap

import "fmt"

type Constants struct {
	RedisKey     RedisKey
	Field        Field
	Tag          Tag
	SMSTemplates SMSTemplates
	JWTKeysPath  JWTKeysPath
}

type RedisKey struct {
}

type Field struct {
	User string
	OTP  string
}

type Tag struct {
	Expired  string
	Invalid  string
	NotFound string
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
		Field: Field{
			User: "user",
			OTP:  "otp",
		},
		Tag: Tag{
			Expired:  "expired",
			Invalid:  "invalid",
			NotFound: "not found",
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
