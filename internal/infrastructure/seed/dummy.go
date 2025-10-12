package seed

import (
	"fmt"

	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
	repository "github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/postgres"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
)

type DummySeeder struct {
	db             database.Database
	userRepository repository.UserRepository
	formRepository repository.FormRepository
}

func NewDummySeeder(
	db database.Database,
	userRepository repository.UserRepository,
	formRepository repository.FormRepository,
) *DummySeeder {
	return &DummySeeder{
		db:             db,
		userRepository: userRepository,
		formRepository: formRepository,
	}
}

func (d *DummySeeder) SeedDummy() {
	d.seedUsers()
	d.seedForms()
}

func (d *DummySeeder) seedUsers() {
	roles, err := d.userRepository.FindAllRoles(d.db)
	if err != nil {
		panic(err)
	}

	roleMap := make(map[string]*entity.Role)
	for _, role := range roles {
		roleMap[role.Name] = role
	}

	users := []struct {
		phone string
		role  string
	}{
		{"09123456771", "بیمار"},
		{"09123456772", "بیمار"},
		{"09123456773", "بیمار"},
		{"09123456774", "بیمار"},
		{"09123456775", "بیمار"},
		{"09123456776", "بیمار"},
		{"09123456777", "بیمار"},
		{"09123456778", "بیمار"},
		{"09123456779", "بیمار"},
		{"09111325792", "بیمار"},

		{"09123456780", "اپراتور"},
		{"09123456781", "اپراتور"},
		{"09123456782", "اپراتور"},
		{"09123456783", "اپراتور"},
		{"09123456784", "اپراتور"},
		{"09123456785", "اپراتور"},
		{"09111325793", "اپراتور"},

		{"09123456786", "مدیر"},
		{"09123456787", "مدیر"},
		{"09123456788", "مدیر"},
		{"09123456789", "مدیر"},
	}

	for _, userData := range users {
		existingUser, err := d.userRepository.FindUserByPhone(d.db, userData.phone)
		if err != nil {
			panic(err)
		}
		if existingUser != nil {
			continue
		}

		user := &entity.User{
			Phone:    userData.phone,
			Password: "123456",
		}

		if err := d.userRepository.CreateUser(d.db, user); err != nil {
			panic(err)
		}

		if role, exists := roleMap[userData.role]; exists {
			if err := d.userRepository.AssignRoleToUser(d.db, user, role); err != nil {
				panic(err)
			}
		}
	}

}

func (d *DummySeeder) seedForms() {
	statuses := []enum.FormStatus{enum.FormStatusPending, enum.FormStatusApproved, enum.FormStatusRejected, enum.FormStatusInComplete, enum.FormStatusReady}
	menopausalStatuses := []string{"قبل از یائسگی", "یائسه", "نامشخص"}
	lifeStatuses := []string{"زنده", "فوت شده", "نامشخص"}
	hrtTypes := []string{"استروژن", "پروژسترون", "ترکیبی", "سایر"}
	insuranceStatuses := []string{"تأمین اجتماعی", "خدمات درمانی", "نیروهای مسلح", "خصوصی", "ندارد"}
	occupationalExposures := []string{
		"آزبست", "بنزن", "کروم", "نیکل", "آرسنیک", "رادون", "ذرات معلق", "دود سیگار",
		"مواد شیمیایی", "اشعه", "هیچکدام", "نامشخص",
	}
	relations := []string{"پدر", "مادر", "برادر", "خواهر", "عمو", "عمه", "دایی", "خاله"}
	lungDiseaseTypes := []string{"آسم", "برونشیت مزمن", "فیبروز ریوی", "COPD", "سایر"}
	smokingTypes := []string{"سیگار", "سیگار برگ", "پیپ", "قلیان", "چپق", "سیگار الکترونیکی"}
	cancerTypes := []string{"سرطان سینه", "سرطان ریه", "سرطان کولون", "سرطان پروستات", "سرطان تخمدان"}
	pastSmokingStatuses := []string{"ترک کرده", "هرگز", "نامشخص"}
	secondhandSmokeLocations := []string{"خانه", "محل کار", "مکان عمومی", "هیچکدام"}
	names := []string{
		"علی احمدی", "فاطمه محمدی", "حسن رضایی", "زهرا کریمی", "محمد حسینی",
		"مریم صادقی", "احمد نوری", "نرگس احمدی", "رضا محمدی", "سارا رضایی",
	}
	addresses := []string{
		"تهران، خیابان ولیعصر، پلاک 123", "اصفهان، خیابان چهارباغ، پلاک 456",
		"شیراز، خیابان زند، پلاک 789", "مشهد، خیابان امام رضا، پلاک 321",
		"تبریز، خیابان آزادی، پلاک 654", "کرج، خیابان فردوسی، پلاک 987",
	}
	provinces := []string{"تهران", "اصفهان", "فارس", "خراسان رضوی", "آذربایجان شرقی", "البرز"}
	cities := []string{"تهران", "اصفهان", "شیراز", "مشهد", "تبریز", "کرج"}
	countries := []string{"ایران", "ترکیه", "آلمان", "کانادا", "آمریکا"}
	months := []string{"فروردین", "اردیبهشت", "خرداد", "تیر", "مرداد", "شهریور",
		"مهر", "آبان", "آذر", "دی", "بهمن", "اسفند"}
	genders := []string{"مرد", "زن"}

	roles, err := d.userRepository.FindAllRoles(d.db)
	if err != nil {
		panic(err)
	}

	var patientRole *entity.Role
	for _, role := range roles {
		if role.Name == "بیمار" {
			patientRole = role
			break
		}
	}

	users, err := d.userRepository.FindUsersByRoleID(d.db, patientRole.ID)
	if err != nil {
		panic(err)
	}

	for i, user := range users {
		existingForms, err := d.formRepository.FindFormsByUserID(d.db, user.ID, nil)
		if err != nil {
			panic(err)
		}
		if len(existingForms) > 0 {
			continue
		}

		numForms := (i % 3) + 1
		for j := 0; j < numForms; j++ {
			form := d.createDummyForm(user.ID, j, statuses)
			if err := d.formRepository.CreateForm(d.db, form); err != nil {
				panic(err)
			}
			basicInfo := d.createDummyBasicInfo(form.ID, j, months, genders)
			if err := d.formRepository.CreateBasicInfo(d.db, basicInfo); err != nil {
				panic(err)
			}
			generalHealth := d.createDummyGeneralHealth(form.ID, j)
			if err := d.formRepository.CreateGeneralHealth(d.db, generalHealth); err != nil {
				panic(err)
			}
			mamography := d.createDummyMamography(form.ID, j, menopausalStatuses, hrtTypes)
			if err := d.formRepository.CreateMamography(d.db, mamography); err != nil {
				panic(err)
			}
			cancer := d.createDummyCancer(form.ID, j, cancerTypes)
			if err := d.formRepository.CreateCancer(d.db, cancer); err != nil {
				panic(err)
			}
			familyCancer := d.createDummyFamilyCancer(form.ID, j, cancerTypes, lifeStatuses, relations)
			if err := d.formRepository.CreateFamilyCancer(d.db, familyCancer); err != nil {
				panic(err)
			}
			contact := d.createDummyContact(form.ID, j, names, addresses, provinces, cities, countries)
			if err := d.formRepository.CreateContact(d.db, contact); err != nil {
				panic(err)
			}
			lungCancer := d.createDummyLungCancer(form.ID, j, insuranceStatuses, occupationalExposures, lungDiseaseTypes, smokingTypes, pastSmokingStatuses, secondhandSmokeLocations, cancerTypes, relations)
			if err := d.formRepository.CreateLungCancer(d.db, lungCancer); err != nil {
				panic(err)
			}
		}
	}
}

func (d *DummySeeder) createDummyForm(userID uint, formIndex int, statuses []enum.FormStatus) *entity.Form {
	status := statuses[formIndex%len(statuses)]

	form := &entity.Form{
		Status: status,
		UserID: userID,
	}
	return form
}
func (d *DummySeeder) createDummyBasicInfo(formID uint, formIndex int, months []string, genders []string) *entity.BasicInfo {
	birthMonth := months[formIndex%len(months)]
	gender := genders[formIndex%len(genders)]

	basicInfo := &entity.BasicInfo{
		FormID:               formID,
		Gender:               gender,
		BirthYear:            uint(1970 + (formIndex % 40)), // 1970-2009
		BirthMonth:           birthMonth,
		BirthDay:             uint((formIndex % 28) + 1), // 1-28
		IsAtba:               formIndex%2 == 0,
		SocialSecurityNumber: fmt.Sprintf("%03d-%06d-%03d", formIndex%1000, formIndex%1000000, formIndex%1000),
		Height:               float64(150 + (formIndex % 50)), // 150-199 cm
		Weight:               float64(50 + (formIndex % 80)),  // 50-129 kg
	}
	return basicInfo
}
func (d *DummySeeder) createDummyGeneralHealth(formID uint, formIndex int) *entity.GeneralHealthInfo {
	smokingNow := formIndex%3 == 0
	drinksAlcohol := formIndex%4 == 0

	generalHealth := &entity.GeneralHealthInfo{
		FormID:                  formID,
		DrinksAlcohol:           &drinksAlcohol,
		CupsPerWeek:             stringPtr(fmt.Sprintf("%d", formIndex%20)),
		LastMonthSabzijatMeal:   fmt.Sprintf("%d", (formIndex%30)+1),
		LastMonthSabzijatWeight: fmt.Sprintf("%.1f", float64(100+(formIndex%500))/10.0),

		MediumActivityMonthInYear: uint((formIndex % 12) + 1),
		MediumActivityHourInWeek:  fmt.Sprintf("%d", (formIndex%20)+1),
		HardActivityMonthInYear:   uint((formIndex % 12) + 1),
		HardActivityHourInWeek:    fmt.Sprintf("%d", (formIndex%15)+1),

		SmokeAtLeast100:       &smokingNow,
		SmokingAge:            uintPtr(uint(15 + (formIndex % 20))),
		SmokingNow:            smokingNow,
		LeaveSmokingAge:       uintPtr(uint(20 + (formIndex % 30))),
		CountSmokingDaily:     stringPtr(fmt.Sprintf("%d", (formIndex%20)+1)),
		CountGheliandaily:     stringPtr(fmt.Sprintf("%d", (formIndex%10)+1)),
		CountSmokingDailyPast: stringPtr(fmt.Sprintf("%d", (formIndex%15)+1)),
		CountGheliandailyPast: stringPtr(fmt.Sprintf("%d", (formIndex%8)+1)),
	}
	return generalHealth
}
func (d *DummySeeder) createDummyMamography(formID uint, formIndex int, menopausalStatuses []string, hrtTypes []string) *entity.MamoGraphyInfo {
	drinksAlcohol := formIndex%4 == 0
	hasChildren := formIndex%3 == 0
	hrtType := hrtTypes[formIndex%len(hrtTypes)]
	menopausalStatus := menopausalStatuses[formIndex%len(menopausalStatuses)]

	mamoGraphyInfo := &entity.MamoGraphyInfo{
		FormID:                       formID,
		GhaedeAge:                    uint(12 + (formIndex % 10)),
		HasChildren:                  hasChildren,
		NumberOfChildren:             uintPtr(uint((formIndex % 5) + 1)),
		AgeOfFirstBirth:              uintPtr(uint(18 + (formIndex % 20))),
		MenopausalStatus:             menopausalStatus,
		MenopauseAge:                 stringPtr(fmt.Sprintf("%d", 45+(formIndex%15))),
		HRT:                          &drinksAlcohol,
		HRTUseLength:                 uintPtr(uint((formIndex % 10) + 1)),
		LastFiveYearsHRTUse:          formIndex%3 == 0,
		CurrentHRTUse:                &drinksAlcohol,
		IntendedHRTUse:               uintPtr(uint((formIndex % 5) + 1)),
		HRTType:                      &hrtType,
		Oral:                         &drinksAlcohol,
		OralDuration:                 stringPtr(fmt.Sprintf("%d", (formIndex%10)+1)),
		OralTwoLastYears:             &drinksAlcohol,
		MamoGraphy:                   &drinksAlcohol,
		Falop:                        &drinksAlcohol,
		Andometrioz:                  &drinksAlcohol,
		LeavePestan:                  formIndex%4 == 0,
		LeaveTokhmdan:                formIndex%5 == 0,
		LaDeColon:                    &drinksAlcohol,
		LaDePol:                      &drinksAlcohol,
		AspLaMo:                      &drinksAlcohol,
		NsaiDLaMo:                    &drinksAlcohol,
		LastFiveYearBloodTestInStool: &drinksAlcohol,
	}
	return mamoGraphyInfo
}
func (d *DummySeeder) createDummyCancer(formID uint, formIndex int, cancerTypes []string) *entity.CancerInfo {
	hasCancer := formIndex%4 == 0

	cancerInfo := &entity.CancerInfo{
		FormID:     formID,
		Cancer:     hasCancer,
		CancerType: stringPtr(cancerTypes[formIndex%len(cancerTypes)]),
		CancerAge:  uintPtr(uint(30 + (formIndex % 40))),
	}
	return cancerInfo
}
func (d *DummySeeder) createDummyFamilyCancer(formID uint, formIndex int, cancerTypes []string, lifeStatuses []string, relations []string) *entity.FamilyCancerInfo {

	hasChildCancer := formIndex%6 == 0
	hasMotherCancer := formIndex%5 == 0
	hasFatherCancer := formIndex%7 == 0
	hasSiblingCancer := formIndex%8 == 0
	hasAmeAmoCancer := formIndex%9 == 0
	hasKhaleDaeiCancer := formIndex%10 == 0
	hasOtherRelativeCancer := formIndex%11 == 0

	familyCancerInfo := &entity.FamilyCancerInfo{
		FormID:          formID,
		ChildCancer:     hasChildCancer,
		ChildName:       stringPtr(fmt.Sprintf("فرزند %d", formIndex+1)),
		ChildCancerType: stringPtr(cancerTypes[formIndex%len(cancerTypes)]),
		ChildCancerAge:  uintPtr(uint(5 + (formIndex % 20))),
		ChildLifeStatus: stringPtr(lifeStatuses[formIndex%len(lifeStatuses)]),

		MotherCancer:     hasMotherCancer,
		MotherName:       stringPtr(fmt.Sprintf("مادر %d", formIndex+1)),
		MotherLifeStatus: stringPtr(lifeStatuses[formIndex%len(lifeStatuses)]),
		MotherCancerType: stringPtr(cancerTypes[formIndex%len(cancerTypes)]),
		MotherCancerAge:  uintPtr(uint(40 + (formIndex % 40))),

		FatherCancer:     hasFatherCancer,
		FatherName:       stringPtr(fmt.Sprintf("پدر %d", formIndex+1)),
		FatherLifeStatus: stringPtr(lifeStatuses[formIndex%len(lifeStatuses)]),
		FatherCancerType: stringPtr(cancerTypes[formIndex%len(cancerTypes)]),
		FatherCancerAge:  uintPtr(uint(45 + (formIndex % 35))),

		SiblingCancer:     hasSiblingCancer,
		SiblingName:       stringPtr(fmt.Sprintf("خواهر/برادر %d", formIndex+1)),
		SiblingLifeStatus: stringPtr(lifeStatuses[formIndex%len(lifeStatuses)]),
		SiblingCancerType: stringPtr(cancerTypes[formIndex%len(cancerTypes)]),
		SiblingCancerAge:  uintPtr(uint(25 + (formIndex % 30))),

		AmeAmoCancer:     hasAmeAmoCancer,
		AmeAmoName:       stringPtr(fmt.Sprintf("عمو/عمه %d", formIndex+1)),
		AmeAmoLifeStatus: stringPtr(lifeStatuses[formIndex%len(lifeStatuses)]),
		AmeAmoCancerType: stringPtr(cancerTypes[formIndex%len(cancerTypes)]),
		AmeAmoCancerAge:  uintPtr(uint(35 + (formIndex % 35))),

		KhaleDaeiCancer:     hasKhaleDaeiCancer,
		KhaleDaeiName:       stringPtr(fmt.Sprintf("خاله/دایی %d", formIndex+1)),
		KhaleDaeiLifeStatus: stringPtr(lifeStatuses[formIndex%len(lifeStatuses)]),
		KhaleDaeiCancerType: stringPtr(cancerTypes[formIndex%len(cancerTypes)]),
		KhaleDaeiCancerAge:  uintPtr(uint(30 + (formIndex % 30))),

		OtherRelativeCancer:     &hasOtherRelativeCancer,
		OtherRelativeName:       stringPtr(fmt.Sprintf("فامیل %d", formIndex+1)),
		OtherRelativeRelation:   stringPtr(relations[formIndex%len(relations)]),
		OtherRelativeLifeStatus: stringPtr(lifeStatuses[formIndex%len(lifeStatuses)]),
		OtherRelativeCancerType: stringPtr(cancerTypes[formIndex%len(cancerTypes)]),
		OtherRelativeCancerAge:  uintPtr(uint(25 + (formIndex % 35))),
	}
	return familyCancerInfo
}
func (d *DummySeeder) createDummyContact(formID uint, formIndex int, names []string, addresses []string, provinces []string, cities []string, countries []string) *entity.ContactInfo {
	name := names[formIndex%len(names)]
	address := addresses[formIndex%len(addresses)]
	hasTestGen := formIndex%5 == 0
	hasFmTestGen := formIndex%6 == 0
	province := provinces[formIndex%len(provinces)]
	city := cities[formIndex%len(cities)]
	country := countries[formIndex%len(countries)]

	contactInfo := &entity.ContactInfo{
		FormID:       formID,
		Name:         name,
		TestGen:      &hasTestGen,
		FmTestGen:    &hasFmTestGen,
		CallExpert:   true,
		BirthCountry: &country,
		Province:     &province,
		City:         &city,
		Country:      &country,
		Address:      address,
		PostalCode:   fmt.Sprintf("%05d", 10000+(formIndex%90000)),
	}
	return contactInfo
}
func (d *DummySeeder) createDummyLungCancer(formID uint, formIndex int, insuranceStatuses []string, occupationalExposures []string, lungDiseaseTypes []string, smokingTypes []string, pastSmokingStatuses []string, secondhandSmokeLocations []string, cancerTypes []string, relations []string) *entity.LungCancerInfo {
	drinksAlcohol := formIndex%4 == 0
	insuranceStatus := insuranceStatuses[formIndex%len(insuranceStatuses)]
	hasHypertension := formIndex%4 == 0
	hasHeartDisease := formIndex%5 == 0
	hasDiabetes := formIndex%6 == 0
	hasLungCancerHistory := formIndex%7 == 0
	hasOtherCancerHistory := formIndex%8 == 0
	hasLungCancerFamily := formIndex%9 == 0
	hasOtherCancerFamily := formIndex%10 == 0
	currentSmoking := formIndex%4 == 0
	secondhandSmoke := formIndex%3 == 0
	occupationalExposure := occupationalExposures[formIndex%len(occupationalExposures)]
	lungDiseaseType := lungDiseaseTypes[formIndex%len(lungDiseaseTypes)]
	smokingType := smokingTypes[formIndex%len(smokingTypes)]
	pastSmokingStatus := pastSmokingStatuses[formIndex%len(pastSmokingStatuses)]
	secondhandSmokeLocation := secondhandSmokeLocations[formIndex%len(secondhandSmokeLocations)]

	lungCancerInfo := &entity.LungCancerInfo{
		FormID:                    formID,
		InsuranceStatus:           &insuranceStatus,
		SupplementaryInsurances:   stringPtr(fmt.Sprintf("بیمه %d", formIndex+1)),
		Hypertension:              hasHypertension,
		HypertensionTreatment:     &hasHypertension,
		HeartDisease:              hasHeartDisease,
		HeartDiseaseTreatment:     &hasHeartDisease,
		Diabetes:                  hasDiabetes,
		DiabetesTreatment:         &hasDiabetes,
		ChronicLungDisease:        &drinksAlcohol,
		ChronicLungDiseaseType:    &lungDiseaseType,
		LungCancerHistory:         hasLungCancerHistory,
		OtherCancerHistory:        hasOtherCancerHistory,
		OtherCancerType:           stringPtr(cancerTypes[formIndex%len(cancerTypes)]),
		LungCancerFamily:          &hasLungCancerFamily,
		LungCancerFamilyRelation:  stringPtr(relations[formIndex%len(relations)]),
		OtherCancerFamily:         &hasOtherCancerFamily,
		OtherCancerFamilyType:     stringPtr(cancerTypes[formIndex%len(cancerTypes)]),
		OtherCancerFamilyRelation: stringPtr(relations[formIndex%len(relations)]),
		OccupationalExposure:      &occupationalExposure,
		CurrentSmoking:            currentSmoking,
		SmokingStartAgeCurrent:    uintPtr(uint(15 + (formIndex % 20))),
		SmokingTypesCurrent:       &smokingType,
		CigarettesPerDayCurrent:   uintPtr(uint((formIndex % 40) + 1)),
		CigarPerDayCurrent:        uintPtr(uint((formIndex % 10) + 1)),
		ECigPerDayCurrent:         uintPtr(uint((formIndex % 20) + 1)),
		PipePerDayCurrent:         uintPtr(uint((formIndex % 15) + 1)),
		ChapoghPerDayCurrent:      uintPtr(uint((formIndex % 8) + 1)),
		SmokedOpiumPerDayCurrent:  uintPtr(uint((formIndex % 5) + 1)),
		ChewedOpiumPerDayCurrent:  uintPtr(uint((formIndex % 3) + 1)),
		HookahPerWeekCurrent:      uintPtr(uint((formIndex % 7) + 1)),
		PastSmoking:               &pastSmokingStatus,
		SmokingStartAgePast:       uintPtr(uint(15 + (formIndex % 20))),
		SmokingTypesPast:          &smokingType,
		CigarettesPerDayPast:      uintPtr(uint((formIndex % 30) + 1)),
		CigarPerDayPast:           uintPtr(uint((formIndex % 8) + 1)),
		ECigPerDayPast:            uintPtr(uint((formIndex % 15) + 1)),
		PipePerDayPast:            uintPtr(uint((formIndex % 12) + 1)),
		ChapoghPerDayPast:         uintPtr(uint((formIndex % 6) + 1)),
		SmokedOpiumPerDayPast:     uintPtr(uint((formIndex % 4) + 1)),
		ChewedOpiumPerDayPast:     uintPtr(uint((formIndex % 2) + 1)),
		HookahPerWeekPast:         uintPtr(uint((formIndex % 5) + 1)),
		SecondhandSmoke:           secondhandSmoke,
		SecondhandSmokeLocation:   &secondhandSmokeLocation,
	}
	return lungCancerInfo
}

func stringPtr(s string) *string {
	return &s
}

func uintPtr(u uint) *uint {
	return &u
}
