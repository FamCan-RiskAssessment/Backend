package bootstrap

import "fmt"

type Constants struct {
	Context      Context
	RedisKey     RedisKey
	Field        Field
	Tag          Tag
	SMSTemplates SMSTemplates
	JWTKeysPath  JWTKeysPath
}

type Context struct {
	Translator string
}

type RedisKey struct {
}

type Field struct {
	User       string
	OTP        string
	Role       string
	Permission string
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
		},
		Field: Field{
			User:       "user",
			OTP:        "otp",
			Role:       "role",
			Permission: "permission",
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
