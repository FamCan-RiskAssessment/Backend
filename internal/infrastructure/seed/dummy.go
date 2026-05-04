package seed

import (
	"fmt"
	"time"

	"github.com/FamCan-RiskAssessment/Backend/internal/application/usecase"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
	repository "github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/postgres"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
)

type DummySeeder struct {
	db             database.Database
	userRepository repository.UserRepository
	formRepository repository.FormRepository
	passwordHasher usecase.PasswordHasher
}

func NewDummySeeder(
	db database.Database,
	userRepository repository.UserRepository,
	formRepository repository.FormRepository,
	passwordHasher usecase.PasswordHasher,
) *DummySeeder {
	return &DummySeeder{
		db:             db,
		userRepository: userRepository,
		formRepository: formRepository,
		passwordHasher: passwordHasher,
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

	hashedDummyPassword, err := d.passwordHasher.HashPassword("123456")
	if err != nil {
		panic(err)
	}

	users := []struct {
		phone string
		role  string
	}{
		{"09123456771", "مراجعه کننده"},
		{"09123456772", "مراجعه کننده"},
		{"09123456773", "مراجعه کننده"},
		{"09123456774", "مراجعه کننده"},
		{"09123456775", "مراجعه کننده"},
		{"09123456776", "مراجعه کننده"},
		{"09123456777", "مراجعه کننده"},
		{"09123456778", "مراجعه کننده"},
		{"09123456779", "مراجعه کننده"},
		{"09111325792", "مراجعه کننده"},

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
			Password: hashedDummyPassword,
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
	statuses := enum.GetAllFormStatuses()
	menopausalStatuses := enum.GetAllMenopausalStatuses()
	lifeStatuses := enum.GetAllLifeStatus()
	hrtTypes := []string{"استروژن", "پروژسترون", "ترکیبی", "سایر"}
	insuranceStatuses := []string{"تأمین اجتماعی", "خدمات درمانی", "نیروهای مسلح", "خصوصی", "ندارد"}
	occupationalExposures := []string{
		"آزبست", "بنزن", "کروم", "نیکل", "آرسنیک", "رادون", "ذرات معلق", "دود سیگار",
		"مواد شیمیایی", "اشعه", "هیچکدام", "نامشخص",
	}
	relations := []string{"پدر", "مادر", "برادر", "خواهر", "عمو", "عمه", "دایی", "خاله"}
	lungDiseaseTypes := []string{"آسم", "برونشیت مزمن", "فیبروز ریوی", "COPD", "سایر"}
	smokingTypes := []string{"سیگار", "سیگار برگ", "پیپ", "قلیان", "چپق", "سیگار الکترونیکی"}
	cancerTypes := enum.GetAllCancerTypes()
	pastSmokingStatuses := []enum.Answer{enum.AnswerAgo, enum.AnswerNo, enum.AnswerLongAgo}
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
	degrees := []string{"ابتدایی یا کمتر", "دیپلم", "فوق دیپلم/مدارک فنی حرفه ای بالای دیپلم", "لیسانس", "فوق لیسانس", "دکتری حرفه ای یا تخصصی"}
	lungDiseases := []string{"مراجعه کنندهی انسداد ریوی مزمن (COPD)", "برونشیت مزمن", "آمفیزم", ""}
	genders := enum.GetAllGenders()

	roles, err := d.userRepository.FindAllRoles(d.db)
	if err != nil {
		panic(err)
	}

	var patientRole *entity.Role
	for _, role := range roles {
		if role.Name == "مراجعه کننده" {
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
			contact := d.createDummyContact(form.ID, j, names, addresses, provinces, cities, countries, degrees)
			if err := d.formRepository.CreateContact(d.db, contact); err != nil {
				panic(err)
			}
			lungCancer := d.createDummyLungCancer(form.ID, j, insuranceStatuses, occupationalExposures, lungDiseaseTypes, smokingTypes, pastSmokingStatuses, secondhandSmokeLocations, cancerTypes, relations, lungDiseases)
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
func (d *DummySeeder) createDummyBasicInfo(formID uint, formIndex int, months []string, genders []enum.Gender) *entity.BasicInfo {
	gender := genders[formIndex%len(genders)]

	BirthDate := time.Date(1970+(formIndex%40), time.Month(formIndex%len(months)+1), formIndex%28+1, 0, 0, 0, 0, time.UTC)
	// BirthDate.AddDate(1970 + (formIndex % 40), , (formIndex % 28) + 1)

	basicInfo := &entity.BasicInfo{
		FormID:               formID,
		Gender:               gender,
		BirthDate:            BirthDate,
		IsAtba:               formIndex%2 == 0,
		SocialSecurityNumber: fmt.Sprintf("%03d-%06d-%03d", formIndex%1000, formIndex%1000000, formIndex%1000),
		Height:               float64(150 + (formIndex % 50)), // 150-199 cm
		Weight:               float64(50 + (formIndex % 80)),  // 50-129 kg
	}
	return basicInfo
}
func (d *DummySeeder) createDummyGeneralHealth(formID uint, formIndex int) *entity.GeneralHealthInfo {
	var smokingNow = enum.AnswerNo
	if formIndex%3 == 0 {
		smokingNow = enum.AnswerYes
	}
	var drinksAlcohol = enum.AnswerNo
	if formIndex%4 == 0 {
		drinksAlcohol = enum.AnswerYes
	}

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
		SmokingNow:            &smokingNow,
		YearSmoke:             uintPtr(uint(15 + (formIndex % 20))),
		LeaveSmokingAge:       uintPtr(uint(20 + (formIndex % 30))),
		CountSmokingDaily:     stringPtr(fmt.Sprintf("%d", (formIndex%20)+1)),
		CountGheliandaily:     stringPtr(fmt.Sprintf("%d", (formIndex%10)+1)),
		CountSmokingDailyPast: stringPtr(fmt.Sprintf("%d", (formIndex%15)+1)),
		CountGheliandailyPast: stringPtr(fmt.Sprintf("%d", (formIndex%8)+1)),
	}
	return generalHealth
}
func (d *DummySeeder) createDummyMamography(formID uint, formIndex int, menopausalStatuses []enum.MenopausalStatus, hrtTypes []string) *entity.MamoGraphyInfo {
	hasChildren := enum.AnswerNo
	hrtType := hrtTypes[formIndex%len(hrtTypes)]
	menopausalStatus := menopausalStatuses[formIndex%len(menopausalStatuses)]
	SonCount := uint((formIndex % 3) + 1)
	DaughterCount := uint((formIndex % 3) + 1)
	var drinksAlcohol = enum.AnswerNo
	if formIndex%4 == 0 {
		drinksAlcohol = enum.AnswerYes
	}
	var lastFiveYearsHRTUse = enum.AnswerNo
	if formIndex%3 == 0 {
		lastFiveYearsHRTUse = enum.AnswerYes
	}
	if formIndex%3 == 0 {
		hasChildren = enum.AnswerYes
	}

	mamoGraphyInfo := &entity.MamoGraphyInfo{
		FormID:                       formID,
		GhaedeAge:                    uint(12 + (formIndex % 10)),
		HasChildren:                  &hasChildren,
		NumberOfChildren:             uintPtr(uint(SonCount + DaughterCount)),
		SonCount:                     &SonCount,
		DaughterCount:                &DaughterCount,
		AgeOfFirstBirth:              uintPtr(uint(18 + (formIndex % 20))),
		MenopausalStatus:             menopausalStatus,
		MenopauseAge:                 stringPtr(fmt.Sprintf("%d", 45+(formIndex%15))),
		HRT:                          &drinksAlcohol,
		HRTUseLength:                 uintPtr(uint((formIndex % 10) + 1)),
		LastFiveYearsHRTUse:          &lastFiveYearsHRTUse,
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
func (d *DummySeeder) createDummyCancer(formID uint, formIndex int, cancerTypes []enum.CancerType) *entity.CancerInfo {
	cancerInfo := &entity.CancerInfo{
		FormID:     formID,
		CancerAge:  uint(30 + (formIndex % 40)),
		CancerType: cancerTypes[formIndex%len(cancerTypes)],
	}
	return cancerInfo
}
func (d *DummySeeder) createDummyFamilyCancer(formID uint, formIndex int, cancerTypes []enum.CancerType, lifeStatuses []enum.LifeStatus, relations []string) *entity.FamilyCancerInfo {
	allRelatives := enum.GetAllRelatives()
	relative := allRelatives[formIndex%len(allRelatives)]

	var relativeRelation *string = nil
	if relative == enum.DistantRelative && len(relations) > 0 {
		relativeRelation = stringPtr(relations[formIndex%len(relations)])
	}

	name := stringPtr(fmt.Sprintf("%s %d", relative.String(), formIndex+1))

	lifeStatus := lifeStatuses[formIndex%len(lifeStatuses)]
	cancerAge := uint(25 + (formIndex % 30))
	cancerType := cancerTypes[formIndex%len(cancerTypes)]

	familyCancerInfo := &entity.FamilyCancerInfo{
		FormID:           formID,
		Relative:         relative,
		RelativeRelation: relativeRelation,
		LifeStatus:       &lifeStatus,
		Name:             name,
		CancerAge:        cancerAge,
		CancerType:       cancerType,
	}

	return familyCancerInfo
}
func (d *DummySeeder) createDummyContact(formID uint, formIndex int, names []string, addresses []string, provinces []string, cities []string, countries []string, degrees []string) *entity.ContactInfo {
	name := names[formIndex%len(names)]
	address := addresses[formIndex%len(addresses)]
	var hasTestGen = enum.AnswerNo
	if formIndex%5 == 0 {
		hasTestGen = enum.AnswerYes
	}
	var hasFmTestGen = enum.AnswerNo
	if formIndex%6 == 0 {
		hasFmTestGen = enum.AnswerYes
	}
	province := provinces[formIndex%len(provinces)]
	city := cities[formIndex%len(cities)]
	country := countries[formIndex%len(countries)]
	degree := degrees[formIndex%len(degrees)]

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
		Education:    degree,
	}
	return contactInfo
}
func (d *DummySeeder) createDummyLungCancer(formID uint, formIndex int, insuranceStatuses []string, occupationalExposures []string, lungDiseaseTypes []string, smokingTypes []string, pastSmokingStatuses []enum.Answer, secondhandSmokeLocations []string, cancerTypes []enum.CancerType, relations []string, lungDiseases []string) *entity.LungCancerInfo {
	var drinksAlcohol = enum.AnswerNo
	if formIndex%4 == 0 {
		drinksAlcohol = enum.AnswerYes
	}
	var SupplementaryInsuranceStatus = enum.AnswerNo
	if formIndex%2 == 0 {
		SupplementaryInsuranceStatus = enum.AnswerYes
	}
	var hasOtherCancerHistory = enum.AnswerNo
	if formIndex%8 == 0 {
		hasOtherCancerHistory = enum.AnswerYes
	}
	var hasLungCancerHistory = enum.AnswerNo
	if formIndex%8 == 0 {
		hasLungCancerHistory = enum.AnswerYes
	}
	var secondhandSmoke = enum.AnswerNo
	if formIndex%3 == 0 {
		secondhandSmoke = enum.AnswerYes
	}
	var currentSmoking = enum.AnswerNo
	if formIndex%4 == 0 {
		currentSmoking = enum.AnswerYes
	}
	var hasOtherCancerFamily = enum.AnswerNo
	if formIndex%10 == 0 {
		hasOtherCancerFamily = enum.AnswerYes
	}
	var hasLungCancerFamily = enum.AnswerNo
	if formIndex%9 == 0 {
		hasLungCancerFamily = enum.AnswerYes
	}
	insuranceStatus := insuranceStatuses[formIndex%len(insuranceStatuses)]
	occupationalExposure := occupationalExposures[formIndex%len(occupationalExposures)]
	lungDiseaseType := lungDiseaseTypes[formIndex%len(lungDiseaseTypes)]
	smokingType := smokingTypes[formIndex%len(smokingTypes)]
	pastSmokingStatus := pastSmokingStatuses[formIndex%len(pastSmokingStatuses)]
	secondhandSmokeLocation := secondhandSmokeLocations[formIndex%len(secondhandSmokeLocations)]
	lungDiseaseHistory := lungDiseases[formIndex%len(lungDiseases)]

	lungCancerInfo := &entity.LungCancerInfo{
		FormID:                       formID,
		InsuranceStatus:              &insuranceStatus,
		SupplementaryInsuranceStatus: &SupplementaryInsuranceStatus,
		SupplementaryInsurances:      stringPtr(fmt.Sprintf("بیمه %d", formIndex+1)),
		ChronicLungDisease:           &drinksAlcohol,
		ChronicLungDiseaseType:       &lungDiseaseType,
		LungCancerHistory:            &hasLungCancerHistory,
		OtherCancerHistory:           &hasOtherCancerHistory,
		OtherCancerType:              &cancerTypes[formIndex%len(cancerTypes)],
		LungCancerFamily:             &hasLungCancerFamily,
		LungCancerFamilyRelation:     stringPtr(relations[formIndex%len(relations)]),
		OtherCancerFamily:            &hasOtherCancerFamily,
		OtherCancerFamilyType:        &cancerTypes[formIndex%len(cancerTypes)],
		OtherCancerFamilyRelation:    stringPtr(relations[formIndex%len(relations)]),
		OccupationalExposure:         &occupationalExposure,
		CurrentSmoking:               &currentSmoking,
		SmokingStartAgeCurrent:       uintPtr(uint(15 + (formIndex % 20))),
		SmokingTypesCurrent:          &smokingType,
		CigarettesPerDayCurrent:      uintPtr(uint((formIndex % 40) + 1)),
		CigarPerDayCurrent:           uintPtr(uint((formIndex % 10) + 1)),
		ECigPerDayCurrent:            uintPtr(uint((formIndex % 20) + 1)),
		PipePerDayCurrent:            uintPtr(uint((formIndex % 15) + 1)),
		ChapoghPerDayCurrent:         uintPtr(uint((formIndex % 8) + 1)),
		SmokedOpiumPerDayCurrent:     uintPtr(uint((formIndex % 5) + 1)),
		ChewedOpiumPerDayCurrent:     uintPtr(uint((formIndex % 3) + 1)),
		HookahPerWeekCurrent:         uintPtr(uint((formIndex % 7) + 1)),
		PastSmoking:                  &pastSmokingStatus,
		LeaveSmoke:                   uintPtr(uint((formIndex % 20) + 1)),
		SmokingStartAgePast:          uintPtr(uint(15 + (formIndex % 20))),
		SmokingTypesPast:             &smokingType,
		CigarettesPerDayPast:         uintPtr(uint((formIndex % 30) + 1)),
		CigarPerDayPast:              uintPtr(uint((formIndex % 8) + 1)),
		ECigPerDayPast:               uintPtr(uint((formIndex % 15) + 1)),
		PipePerDayPast:               uintPtr(uint((formIndex % 12) + 1)),
		ChapoghPerDayPast:            uintPtr(uint((formIndex % 6) + 1)),
		SmokedOpiumPerDayPast:        uintPtr(uint((formIndex % 4) + 1)),
		ChewedOpiumPerDayPast:        uintPtr(uint((formIndex % 2) + 1)),
		HookahPerWeekPast:            uintPtr(uint((formIndex % 5) + 1)),
		SecondhandSmoke:              &secondhandSmoke,
		SecondhandSmokeLocation:      &secondhandSmokeLocation,
		LungDiseaseHistory:           stringPtr(lungDiseaseHistory),
	}
	return lungCancerInfo
}

func stringPtr(s string) *string {
	return &s
}

func uintPtr(u uint) *uint {
	return &u
}
