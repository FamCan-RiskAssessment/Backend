package bootstrap

import (
	"fmt"

	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
)

type Constants struct {
	Context      Context
	RedisKey     RedisKey
	S3BucketPath BucketPath
	BucketPath   BucketPath
	Field        Field
	Tag          Tag
	SMSTemplates SMSTemplates
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
	Profile          string
	Form             string
	OTP              string
	Role             string
	Page             string
	Permission       string
	Record           string
	MamoGraphyInfo   string
	FamilyCancerInfo string
	Cancer           string
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

func NewConstants() *Constants {
	return &Constants{
		Context: Context{
			Translator: "translator",
			ID:         "ID",
		},
		Field: Field{
			User:             "user",
			Profile:          "profile",
			Form:             "form",
			OTP:              "otp",
			Role:             "role",
			Page:             "page",
			Permission:       "permission",
			Record:           "record",
			MamoGraphyInfo:   "mamography_info",
			FamilyCancerInfo: "family_cancer_info",
			Cancer:           "cancer",
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
	}
}

func (r *RedisKey) GenerateOTPKey(value string) string {
	return fmt.Sprintf("otp:%s", value)
}

func (path *BucketPath) GetMamoGraphyPath(formID uint, fileName string) string {
	return fmt.Sprintf("%d/%s", formID, fileName)
}

func (path *BucketPath) GetCancerPath(formID uint, cancerType enum.CancerType, fileName string) string {
	return fmt.Sprintf("%d/%s/%s", formID, cancerType.String(), fileName)
}

func (path *BucketPath) GetFamilyCancerPath(formID uint, cancerType enum.CancerType, fileName string) string {
	return fmt.Sprintf("familyCancer/%d/%s/%s", formID, cancerType.String(), fileName)
}

func (path *BucketPath) GetGeneticTestPath(formID uint, testType string, fileName string) string {
	return fmt.Sprintf("geneticTest/%d/%s/%s", formID, testType, fileName)
}

func (path *BucketPath) GetFatherGeneticTestPath(formID uint, fileName string) string {
	return fmt.Sprintf("fatherGeneticTest/%d/%s", formID, fileName)
}

func (path *BucketPath) GetMotherGeneticTestPath(formID uint, fileName string) string {
	return fmt.Sprintf("motherGeneticTest/%d/%s", formID, fileName)
}

func (r *RedisKey) GenerateOperatorValidationOTPKey(operatorID uint, phone string) string {
	return fmt.Sprintf("operator:validation:otp:%d:%s", operatorID, phone)
}

func (r *RedisKey) GenerateOperatorValidationTokenKey(operatorID uint, userID uint) string {
	return fmt.Sprintf("operator:validation:token:%d:%d", operatorID, userID)
}
