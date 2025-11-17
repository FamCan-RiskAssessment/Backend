# FamCan Backend - Context Flow Diagram (CFD)

## System Overview
**FamCan Backend** is a healthcare risk assessment platform built with Go/Gin that manages patient cancer risk forms and integrates with multiple risk calculation models.

---

## 1. System Context Diagram

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           External Systems & Actors                          │
└─────────────────────────────────────────────────────────────────────────────┘

┌──────────────┐              ┌──────────────┐              ┌──────────────┐
│   Patients   │              │  Operators   │              │    Admins    │
│  (Mobile/Web)│              │ (Healthcare) │              │ (Management) │
└──────┬───────┘              └──────┬───────┘              └──────┬───────┘
       │                             │                             │
       │ OTP Login                   │ Password/OTP                │ Password
       │ Submit Forms                │ Create/Validate Forms       │ Manage Users/Roles
       │ View Results                │ Accept/Reject Forms         │ View Logs
       │                             │ Assign Cases                │ Trigger Calculations
       │                             │                             │
       └─────────────────────────────┼─────────────────────────────┘
                                     │
                                     │ HTTPS/REST API
                                     ▼
         ┌───────────────────────────────────────────────────────────┐
         │                                                           │
         │                  FamCan Backend (Go/Gin)                  │
         │                     Port: 8080                            │
         │                                                           │
         │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐      │
         │  │Presentation │  │ Application │  │Infrastructure│      │
         │  │   Layer     │──│    Layer    │──│    Layer     │      │
         │  └─────────────┘  └─────────────┘  └─────────────┘      │
         │                                                           │
         └───────────────────────────┬───────────────────────────────┘
                                     │
         ┌───────────────────────────┼───────────────────────────────┐
         │                           │                               │
         ▼                           ▼                               ▼
┌─────────────────┐       ┌─────────────────┐           ┌─────────────────┐
│  PostgreSQL 17  │       │    Redis 8      │           │   S3 Storage    │
│  (Main DB)      │       │   (Cache/OTP)   │           │ (File Upload)   │
│                 │       │                 │           │                 │
│ • Users/Roles   │       │ • OTP Tokens    │           │ • Mamography    │
│ • Forms Data    │       │ • Validation    │           │ • Cancer Docs   │
│ • Results       │       │ • Sessions      │           │                 │
│ • Audit Logs    │       │                 │           │                 │
└─────────────────┘       └─────────────────┘           └─────────────────┘

         │                           │                               │
         └───────────────────────────┴───────────────────────────────┘
                                     │
         ┌───────────────────────────┼───────────────────────────────┐
         │                           │                               │
         ▼                           ▼                               ▼
┌─────────────────┐       ┌─────────────────┐           ┌─────────────────┐
│  Kavenegar SMS  │       │   PREMM5 API    │           │    BCRA API     │
│   (OTP Send)    │       │ (Lynch Risk)    │           │ (BRCA Mutation) │
└─────────────────┘       └─────────────────┘           └─────────────────┘

         ▼                           ▼
┌─────────────────┐       ┌─────────────────┐
│    GAIL API     │       │    PLCO API     │
│(Breast Cancer)  │       │ (Lung Cancer)   │
└─────────────────┘       └─────────────────┘
```

---

## 2. Data Flow Diagram - Authentication

```
┌─────────┐                                                    ┌──────────┐
│  User   │                                                    │  Redis   │
└────┬────┘                                                    └────┬─────┘
     │                                                              │
     │ 1. POST /auth/login                                         │
     │    { phone: "09123456789" }                                │
     ├────────────────────────────────────────┐                   │
     │                                        │                   │
     │                          ┌─────────────▼──────────────┐    │
     │                          │  GeneralUserController     │    │
     │                          │    .Login()                │    │
     │                          └─────────────┬──────────────┘    │
     │                                        │                   │
     │                          ┌─────────────▼──────────────┐    │
     │                          │     UserService            │    │
     │                          │  .Login(phone)             │    │
     │                          └─────────────┬──────────────┘    │
     │                                        │                   │
     │                                        │ 2. Find/Create    │
     │                          ┌─────────────▼──────────────┐    │
     │                          │   UserRepository (GORM)    │    │
     │                          │   .FindByPhone()           │    │
     │                          └─────────────┬──────────────┘    │
     │                                        │                   │
     │                          ┌─────────────▼──────────────┐    │
     │                          │     OTPService             │    │
     │                          │  .GenerateOTP()            │    │
     │                          └─────────────┬──────────────┘    │
     │                                        │                   │
     │                                        │ 3. Store OTP      │
     │                                        ├──────────────────►│
     │                                        │   Key: otp:{phone}│
     │                                        │   TTL: 2 minutes  │
     │                          ┌─────────────▼──────────────┐    │
     │                          │  AsanakSMSService          │    │
     │                          │   .SendOTP(phone, code)    │    │
     │                          └─────────────┬──────────────┘    │
     │                                        │                   │
     │◄───────────────────────────────────────┤                   │
     │  { message: "OTP sent" }               │                   │
     │                                                             │
     │                                                             │
     │ 4. POST /auth/verify-otp                                   │
     │    { phone: "09123456789", otp: "123456" }                │
     ├────────────────────────────────────────┐                   │
     │                                        │                   │
     │                          ┌─────────────▼──────────────┐    │
     │                          │  GeneralUserController     │    │
     │                          │    .VerifyOTP()            │    │
     │                          └─────────────┬──────────────┘    │
     │                                        │                   │
     │                          ┌─────────────▼──────────────┐    │
     │                          │     OTPService             │    │
     │                          │  .VerifyOTP()              │    │
     │                          └─────────────┬──────────────┘    │
     │                                        │                   │
     │                                        │ 5. Verify OTP     │
     │                                        ├──────────────────►│
     │                                        │   Get otp:{phone} │
     │                                        │◄──────────────────┤
     │                          ┌─────────────▼──────────────┐    │
     │                          │     JWTService             │    │
     │                          │  .GenerateToken(user)      │    │
     │                          │  (RSA Private Key Sign)    │    │
     │                          └─────────────┬──────────────┘    │
     │◄───────────────────────────────────────┤                   │
     │  {                                                         │
     │    accessToken: "eyJhbGc...",                              │
     │    refreshToken: "eyJhbGc..."                              │
     │  }                                                         │
     │                                                             │
     └─────────────────────────────────────────────────────────────┘
```

---

## 3. Data Flow Diagram - Form Submission (Patient)

```
┌─────────┐                                          ┌──────────────┐
│ Patient │                                          │  PostgreSQL  │
└────┬────┘                                          └──────┬───────┘
     │                                                      │
     │ 1. POST /form/basic                                 │
     │    Authorization: Bearer {JWT}                      │
     │    { gender: "MALE", birthDate: "1990-01-01" }     │
     ├────────────────────────────────────┐               │
     │                                    │               │
     │                  ┌─────────────────▼────────────┐  │
     │                  │  AuthMiddleware              │  │
     │                  │  - Validate JWT              │  │
     │                  │  - Extract UserID            │  │
     │                  └─────────────────┬────────────┘  │
     │                                    │               │
     │                  ┌─────────────────▼────────────┐  │
     │                  │ CustomerFormController       │  │
     │                  │   .CreateForm()              │  │
     │                  └─────────────────┬────────────┘  │
     │                                    │               │
     │                  ┌─────────────────▼────────────┐  │
     │                  │   FormService                │  │
     │                  │   .CreateBasicInfoForm()     │  │
     │                  └─────────────────┬────────────┘  │
     │                                    │               │
     │                                    │ 2. Insert     │
     │                  ┌─────────────────▼────────────┐  │
     │                  │   FormRepository             │  │
     │                  │   .CreateForm()              │──►
     │                  └─────────────────┬────────────┘  │
     │                                    │               │
     │                  ┌─────────────────▼────────────┐  │
     │                  │   ActionLogService           │  │
     │                  │   .LogAction()               │  │
     │                  └─────────────────┬────────────┘  │
     │◄───────────────────────────────────┤               │
     │  { formID: "uuid-123", status: "DRAFT" }          │
     │                                                    │
     │                                                    │
     │ 3. PUT /form/{formID}/generalhealth               │
     │    { hasAlcohol: true, smokingStatus: "NEVER" }   │
     ├────────────────────────────────────┐              │
     │                                    │              │
     │                  ┌─────────────────▼────────────┐ │
     │                  │ CustomerFormController       │ │
     │                  │   .UpsertGeneralHealthInfo() │ │
     │                  └─────────────────┬────────────┘ │
     │                                    │              │
     │                  ┌─────────────────▼────────────┐ │
     │                  │   FormService                │ │
     │                  │   .UpsertGeneralHealthInfo() │ │
     │                  └─────────────────┬────────────┘ │
     │                                    │              │
     │                                    │ 4. Upsert    │
     │                  ┌─────────────────▼────────────┐ │
     │                  │   FormRepository             │ │
     │                  │   .UpsertGeneralHealthInfo() ├─►
     │                  └──────────────────────────────┘ │
     │◄───────────────────────────────────┤              │
     │  { message: "Updated successfully" }              │
     │                                                    │
     │                                                    │
     │ 5. PUT /form/{formID}/status                      │
     │    { status: "SUBMITTED" }                        │
     ├────────────────────────────────────┐              │
     │                                    │              │
     │                  ┌─────────────────▼────────────┐ │
     │                  │   FormService                │ │
     │                  │   .ChangeFormStatus()        │ │
     │                  └─────────────────┬────────────┘ │
     │                                    │              │
     │                                    │ 6. Update    │
     │                                    ├─────────────►│
     │◄───────────────────────────────────┤              │
     │  { message: "Form submitted" }                    │
     │                                                    │
     └────────────────────────────────────────────────────┘
```

---

## 4. Data Flow Diagram - Operator Workflow

```
┌──────────┐                              ┌─────────┐         ┌──────────────┐
│ Operator │                              │  Redis  │         │  PostgreSQL  │
└────┬─────┘                              └────┬────┘         └──────┬───────┘
     │                                         │                     │
     │ 1. POST /operator/validate-user/request │                     │
     │    { phone: "09123456789" }             │                     │
     ├────────────────────────────────┐        │                     │
     │                                │        │                     │
     │              ┌─────────────────▼──────────────┐               │
     │              │ AdminFormController            │               │
     │              │  .RequestUserValidationOTP()   │               │
     │              └─────────────────┬──────────────┘               │
     │                                │                              │
     │              ┌─────────────────▼──────────────┐               │
     │              │   UserService                  │               │
     │              │  .RequestUserValidationOTP()   │               │
     │              └─────────────────┬──────────────┘               │
     │                                │                              │
     │                                │ 2. Store OTP                 │
     │                                ├────────────►                 │
     │                                │   operator:validation:otp:   │
     │                                │   {operatorID}:{phone}       │
     │◄───────────────────────────────┤                              │
     │  { message: "OTP sent" }       │                              │
     │                                                                │
     │ 3. POST /operator/validate-user/verify                        │
     │    { phone: "09123456789", otp: "123456" }                   │
     ├────────────────────────────────┐                              │
     │                                │                              │
     │              ┌─────────────────▼──────────────┐               │
     │              │   UserService                  │               │
     │              │  .VerifyUserValidationOTP()    │               │
     │              └─────────────────┬──────────────┘               │
     │                                │                              │
     │                                │ 4. Verify & Create Token     │
     │                                ├────────────►                 │
     │                                │   operator:validation:token: │
     │                                │   {operatorID}:{userID}      │
     │◄───────────────────────────────┤                              │
     │  { userID: "uuid-456", token: "validation-token" }           │
     │                                                                │
     │                                                                │
     │ 5. POST /operator/form                                        │
     │    { userID: "uuid-456", basicInfo: {...} }                  │
     ├────────────────────────────────┐                              │
     │                                │                              │
     │              ┌─────────────────▼──────────────┐               │
     │              │ AdminFormController            │               │
     │              │  .CreateFormForUser()          │               │
     │              └─────────────────┬──────────────┘               │
     │                                │                              │
     │              ┌─────────────────▼──────────────┐               │
     │              │   FormService                  │               │
     │              │  .CreateBasicInfoForm()        │               │
     │              │  (Set FilledByOperatorID)      │               │
     │              └─────────────────┬──────────────┘               │
     │                                │                              │
     │                                │ 6. Insert with OperatorID    │
     │                                ├────────────────────────────►│
     │◄───────────────────────────────┤                              │
     │  { formID: "uuid-789" }                                       │
     │                                                                │
     │                                                                │
     │ 7. PUT /operator/form/{formID}/accept                         │
     ├────────────────────────────────┐                              │
     │                                │                              │
     │              ┌─────────────────▼──────────────┐               │
     │              │   FormService                  │               │
     │              │  .AcceptForm()                 │               │
     │              └─────────────────┬──────────────┘               │
     │                                │                              │
     │                                │ 8. Update Status: ACCEPTED   │
     │                                ├────────────────────────────►│
     │◄───────────────────────────────┤                              │
     │  { message: "Form accepted" }                                 │
     │                                                                │
     └────────────────────────────────────────────────────────────────┘
```

---

## 5. Data Flow Diagram - Risk Calculation

```
┌────────┐                                      ┌──────────────┐   ┌────────────┐
│ Admin  │                                      │  PostgreSQL  │   │ External   │
└───┬────┘                                      └──────┬───────┘   │ Calc APIs  │
    │                                                  │           └─────┬──────┘
    │ 1. POST /calc/model                             │                 │
    │    { formID: "uuid-123", model: "PREMM5" }     │                 │
    ├──────────────────────────────┐                 │                 │
    │                              │                 │                 │
    │            ┌─────────────────▼──────────────┐  │                 │
    │            │ AdminCalcController            │  │                 │
    │            │  .SendFormToCalc()             │  │                 │
    │            └─────────────────┬──────────────┘  │                 │
    │                              │                 │                 │
    │            ┌─────────────────▼──────────────┐  │                 │
    │            │   CalcService                  │  │                 │
    │            │  .SendFormToCalc()             │  │                 │
    │            └─────────────────┬──────────────┘  │                 │
    │                              │                 │                 │
    │                              │ 2. Get Form     │                 │
    │                              ├────────────────►│                 │
    │                              │◄────────────────┤                 │
    │                              │  Form + Cancers │                 │
    │                              │  + Family Data  │                 │
    │                              │                 │                 │
    │                              │ 3. POST to API                    │
    │                              ├──────────────────────────────────►│
    │                              │   PREMM5_API_URL                  │
    │                              │   { form data mapped }            │
    │                              │◄──────────────────────────────────┤
    │                              │   { MLH1: 0.15, MSH2: 0.12, ... } │
    │                              │                                   │
    │                              │ 4. Store Result                   │
    │                              ├────────────────►│                 │
    │                              │  Premm5Result   │                 │
    │◄─────────────────────────────┤                 │                 │
    │  { message: "Calculation complete" }           │                 │
    │                                                 │                 │
    │                                                 │                 │
    │ 5. GET /calc/premm5/{formID}                   │                 │
    ├──────────────────────────────┐                 │                 │
    │                              │                 │                 │
    │            ┌─────────────────▼──────────────┐  │                 │
    │            │   CalcService                  │  │                 │
    │            │  .GetPremm5Results()           │  │                 │
    │            └─────────────────┬──────────────┘  │                 │
    │                              │                 │                 │
    │                              │ 6. Query Result │                 │
    │                              ├────────────────►│                 │
    │◄─────────────────────────────┤◄────────────────┤                 │
    │  {                           │  Premm5Result   │                 │
    │    MLH1: 0.15,                                 │                 │
    │    MSH2: 0.12,                                 │                 │
    │    MSH6: 0.08,                                 │                 │
    │    PMS2: 0.05,                                 │                 │
    │    EPCAM: 0.02                                 │                 │
    │  }                                             │                 │
    │                                                 │                 │
    └─────────────────────────────────────────────────┴─────────────────┘

Note: Similar flows exist for BCRA, GAIL, and PLCO models
```

---

## 6. Component Interaction Flow - File Upload

```
┌────────┐                                ┌──────────┐         ┌──────────────┐
│ Client │                                │ Backend  │         │  S3 Storage  │
└───┬────┘                                └────┬─────┘         └──────┬───────┘
    │                                          │                      │
    │ 1. PUT /form/{formID}/mamography        │                      │
    │    Content-Type: multipart/form-data    │                      │
    │    { file: mamography.jpg, ... }        │                      │
    ├─────────────────────────────────────────►                      │
    │                                          │                      │
    │                        ┌─────────────────▼──────────────┐      │
    │                        │ CustomerFormController         │      │
    │                        │  .UpsertMamoGraphyInfo()       │      │
    │                        └─────────────────┬──────────────┘      │
    │                                          │                      │
    │                        ┌─────────────────▼──────────────┐      │
    │                        │   FormService                  │      │
    │                        │  .UpsertMamoGraphyInfo()       │      │
    │                        └─────────────────┬──────────────┘      │
    │                                          │                      │
    │                                          │ 2. Upload File       │
    │                        ┌─────────────────▼──────────────┐      │
    │                        │   S3Storage                    │      │
    │                        │  .UploadFile(file, bucket)     │      │
    │                        └─────────────────┬──────────────┘      │
    │                                          │                      │
    │                                          │ 3. PUT Object        │
    │                                          ├─────────────────────►│
    │                                          │   Bucket: mamography │
    │                                          │◄─────────────────────┤
    │                                          │   URL: s3://bucket/  │
    │                                          │        uuid.jpg      │
    │                        ┌─────────────────▼──────────────┐      │
    │                        │   FormRepository               │      │
    │                        │  .UpsertMamoGraphyInfo()       │      │
    │                        │  (Save S3 path to DB)          │      │
    │                        └─────────────────┬──────────────┘      │
    │◄─────────────────────────────────────────┤                      │
    │  { picturePath: "s3://bucket/uuid.jpg" }                       │
    │                                                                 │
    │                                                                 │
    │ 4. GET /form/{formID}/mamography                               │
    ├─────────────────────────────────────────►                      │
    │                                          │                      │
    │                        ┌─────────────────▼──────────────┐      │
    │                        │   FormService                  │      │
    │                        │  .GetMamoGraphyInfo()          │      │
    │                        └─────────────────┬──────────────┘      │
    │                                          │                      │
    │                                          │ 5. Generate          │
    │                        ┌─────────────────▼──────────────┐      │
    │                        │   S3Storage                    │      │
    │                        │  .GetPresignedURL(path)        │      │
    │                        └─────────────────┬──────────────┘      │
    │                                          │                      │
    │                                          │ 6. Presigned URL     │
    │                                          ├─────────────────────►│
    │                                          │◄─────────────────────┤
    │                                          │   Signed URL         │
    │◄─────────────────────────────────────────┤   (15-min expiry)    │
    │  {                                                              │
    │    picturePath: "https://s3.../uuid.jpg?signature=..."         │
    │  }                                                              │
    │                                                                 │
    └─────────────────────────────────────────────────────────────────┘
```

---

## 7. Authorization Flow - Role-Based Access Control

```
┌────────┐                                          ┌──────────────┐
│ Client │                                          │  PostgreSQL  │
└───┬────┘                                          └──────┬───────┘
    │                                                      │
    │ Request with JWT                                    │
    │ GET /admin/user                                     │
    │ Authorization: Bearer eyJhbGc...                    │
    ├──────────────────────────────┐                     │
    │                              │                     │
    │            ┌─────────────────▼──────────────┐      │
    │            │   AuthMiddleware               │      │
    │            │  (routes/admin.go applied)     │      │
    │            └─────────────────┬──────────────┘      │
    │                              │                     │
    │                              │ 1. Validate JWT     │
    │            ┌─────────────────▼──────────────┐      │
    │            │   JWTService                   │      │
    │            │  .ValidateToken()              │      │
    │            │  (RSA Public Key Verify)       │      │
    │            └─────────────────┬──────────────┘      │
    │                              │                     │
    │                              │ 2. Extract UserID   │
    │                              │                     │
    │            ┌─────────────────▼──────────────┐      │
    │            │   UserRepository               │      │
    │            │  .GetUserByID()                │      │
    │            └─────────────────┬──────────────┘      │
    │                              │                     │
    │                              │ 3. Load User        │
    │                              ├────────────────────►│
    │                              │◄────────────────────┤
    │                              │  User + Roles +     │
    │                              │  Permissions        │
    │                              │                     │
    │            ┌─────────────────▼──────────────┐      │
    │            │   AuthMiddleware               │      │
    │            │  .CheckPermission()            │      │
    │            │  Required: USER_MANAGEMENT_READ│      │
    │            └─────────────────┬──────────────┘      │
    │                              │                     │
    │                              │ 4. Has Permission?  │
    │                              │                     │
    │                    ┌─────────┴─────────┐           │
    │                    │                   │           │
    │                   YES                 NO           │
    │                    │                   │           │
    │                    ▼                   ▼           │
    │         ┌──────────────────┐  ┌────────────────┐  │
    │         │ Continue to      │  │ Return 403     │  │
    │         │ Controller       │  │ Forbidden      │  │
    │         └────────┬─────────┘  └────────┬───────┘  │
    │                  │                     │           │
    │                  ▼                     │           │
    │    ┌──────────────────────┐           │           │
    │    │ AdminUserController  │           │           │
    │    │  .GetUsers()         │           │           │
    │    └──────────┬───────────┘           │           │
    │               │                       │           │
    │◄──────────────┴───────────────────────┘           │
    │  Response (200 or 403)                            │
    │                                                    │
    └────────────────────────────────────────────────────┘
```

---

## 8. Layer Architecture Flow

```
┌───────────────────────────────────────────────────────────────────────┐
│                         PRESENTATION LAYER                            │
│                     (/internal/presentation)                          │
├───────────────────────────────────────────────────────────────────────┤
│                                                                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────────┐   │
│  │  Middleware  │  │  Controllers │  │       Routes             │   │
│  │              │  │              │  │                          │   │
│  │ • CORS       │  │ • User       │  │ • /auth (General)        │   │
│  │ • Auth       │  │ • Form       │  │ • /form (Customer)       │   │
│  │ • Recovery   │  │ • Calc       │  │ • /admin (Protected)     │   │
│  │ • i18n       │  │ • ActionLog  │  │                          │   │
│  └──────┬───────┘  └──────┬───────┘  └──────────────────────────┘   │
│         │                 │                                          │
└─────────┼─────────────────┼──────────────────────────────────────────┘
          │                 │
          │ HTTP Request    │ DTO Objects
          │                 │
          ▼                 ▼
┌───────────────────────────────────────────────────────────────────────┐
│                        APPLICATION LAYER                              │
│                     (/internal/application)                           │
├───────────────────────────────────────────────────────────────────────┤
│                                                                       │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │                        Services                              │   │
│  │                                                              │   │
│  │  • UserService       - Authentication, RBAC                 │   │
│  │  • FormService       - Form CRUD, validation                │   │
│  │  • CalcService       - Risk calculations                    │   │
│  │  • JWTService        - Token generation                     │   │
│  │  • OTPService        - OTP generation/verification          │   │
│  │  • ActionLogService  - Audit logging                        │   │
│  │                                                              │   │
│  └──────────────────────────┬───────────────────────────────────┘   │
│                             │                                        │
│  ┌──────────────────────────▼───────────────────────────────────┐   │
│  │                         DTOs                                 │   │
│  │                                                              │   │
│  │  • Request DTOs  - Input validation                         │   │
│  │  • Response DTOs - Output formatting                        │   │
│  │                                                              │   │
│  └──────────────────────────────────────────────────────────────┘   │
│                             │                                        │
└─────────────────────────────┼────────────────────────────────────────┘
                              │ Domain Entities
                              │ Repository Interfaces
                              ▼
┌───────────────────────────────────────────────────────────────────────┐
│                           DOMAIN LAYER                                │
│                        (/internal/domain)                             │
├───────────────────────────────────────────────────────────────────────┤
│                                                                       │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐   │
│  │    Entities      │  │  Repository      │  │    Exceptions    │   │
│  │                  │  │  Interfaces      │  │                  │   │
│  │ • User           │  │                  │  │ • AuthError      │   │
│  │ • Role           │  │ • UserRepo       │  │ • NotFoundError  │   │
│  │ • Permission     │  │ • FormRepo       │  │ • ValidationErr  │   │
│  │ • Form           │  │ • LogRepo        │  │ • ForbiddenErr   │   │
│  │ • BasicInfo      │  │ • CacheRepo      │  │                  │   │
│  │ • CancerInfo     │  │                  │  │                  │   │
│  │ • CalcResults    │  │                  │  │                  │   │
│  └──────────────────┘  └──────────────────┘  └──────────────────┘   │
│                                                                       │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐   │
│  │     Enums        │  │ Communication    │  │    Storage       │   │
│  │                  │  │ Interfaces       │  │  Interfaces      │   │
│  │ • Gender         │  │                  │  │                  │   │
│  │ • CancerType     │  │ • SmsService     │  │ • S3Storage      │   │
│  │ • FormStatus     │  │                  │  │                  │   │
│  │ • LifeStatus     │  │                  │  │                  │   │
│  └──────────────────┘  └──────────────────┘  └──────────────────┘   │
│                             │                        │               │
└─────────────────────────────┼────────────────────────┼───────────────┘
                              │ Interface              │ Interface
                              │ Implementation         │ Implementation
                              ▼                        ▼
┌───────────────────────────────────────────────────────────────────────┐
│                      INFRASTRUCTURE LAYER                             │
│                    (/internal/infrastructure)                         │
├───────────────────────────────────────────────────────────────────────┤
│                                                                       │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │                    Database Layer                            │   │
│  │                                                              │   │
│  │  • PostgresDatabase  - GORM wrapper (Singleton)             │   │
│  │  • RedisDatabase     - Redis client wrapper                 │   │
│  │                                                              │   │
│  └──────────────────────────┬───────────────────────────────────┘   │
│                             │                                        │
│  ┌──────────────────────────▼───────────────────────────────────┐   │
│  │                    Repositories                              │   │
│  │                                                              │   │
│  │  • UserRepository (PostgreSQL)                              │   │
│  │  • FormRepository (PostgreSQL)                              │   │
│  │  • ActionLogRepository (PostgreSQL)                         │   │
│  │  • UserCacheRepository (Redis)                              │   │
│  │                                                              │   │
│  └──────────────────────────────────────────────────────────────┘   │
│                                                                       │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐   │
│  │  Communication   │  │    Storage       │  │      JWT         │   │
│  │                  │  │                  │  │                  │   │
│  │ • AsanakSMS      │  │ • S3Storage      │  │ • JWTKeyManager  │   │
│  │   (Kavenegar)    │  │   (AWS SDK)      │  │   (RSA Keys)     │   │
│  └──────────────────┘  └──────────────────┘  └──────────────────┘   │
│                                                                       │
│  ┌──────────────────┐  ┌──────────────────┐                         │
│  │  Localization    │  │      Seed        │                         │
│  │                  │  │                  │                         │
│  │ • i18nService    │  │ • RoleSeeder     │                         │
│  │   (Multi-lang)   │  │ • DummySeeder    │                         │
│  └──────────────────┘  └──────────────────┘                         │
│         │                      │                                     │
└─────────┼──────────────────────┼─────────────────────────────────────┘
          │                      │
          ▼                      ▼
┌─────────────────┐    ┌─────────────────┐
│   PostgreSQL    │    │      Redis      │
│   (Port 5432)   │    │   (Port 6379)   │
└─────────────────┘    └─────────────────┘
```

---

## 9. Dependency Injection Flow (Wire)

```
┌───────────────────────────────────────────────────────────────────┐
│                      Application Bootstrap                        │
│                         (main.go)                                 │
└───────────────────────────────────┬───────────────────────────────┘
                                    │
                                    │ InitializeApplication()
                                    ▼
┌───────────────────────────────────────────────────────────────────┐
│                       Wire Container                              │
│                     (wire/wire_gen.go)                            │
├───────────────────────────────────────────────────────────────────┤
│                                                                   │
│  1. Bootstrap Config                                             │
│     └─> env.NewEnv() → bootstrap.Config                          │
│                                                                   │
│  2. Database Providers                                           │
│     ├─> PostgresDatabase.GetInstance()                           │
│     └─> RedisDatabase.GetInstance()                              │
│                                                                   │
│  3. Repository Providers                                         │
│     ├─> NewUserRepository(db) → IUserRepository                  │
│     ├─> NewFormRepository(db) → IFormRepository                  │
│     ├─> NewActionLogRepository(db) → IActionLogRepository        │
│     └─> NewUserCacheRepository(rdb) → IUserCacheRepository       │
│                                                                   │
│  4. Adapter Providers (Infrastructure)                           │
│     ├─> NewJWTKeyManager(config) → IJWTKeyManager                │
│     ├─> NewAsanakSMSService(config) → ISmsService                │
│     ├─> NewS3Storage(config) → IS3Storage                        │
│     └─> NewTranslationService() → ITranslationService            │
│                                                                   │
│  5. Service Providers (Application Layer)                        │
│     ├─> NewUserService(userRepo, cacheRepo, smsService, ...)     │
│     ├─> NewFormService(formRepo, s3Storage, logService, ...)     │
│     ├─> NewCalcService(formRepo, config)                         │
│     ├─> NewJWTService(keyManager, config)                        │
│     ├─> NewOTPService(cacheRepo, config)                         │
│     └─> NewActionLogService(logRepo)                             │
│                                                                   │
│  6. Middleware Providers                                         │
│     ├─> NewAuthMiddleware(jwtService, userService, i18n)         │
│     ├─> NewCORSMiddleware()                                      │
│     ├─> NewRecoveryMiddleware(i18n)                              │
│     └─> NewLocalizationMiddleware(i18n)                          │
│                                                                   │
│  7. Controller Providers                                         │
│     ├─> NewGeneralUserController(userService, jwtService, ...)   │
│     ├─> NewAdminUserController(userService, ...)                 │
│     ├─> NewCustomerFormController(formService, ...)              │
│     ├─> NewAdminFormController(formService, ...)                 │
│     ├─> NewAdminCalcController(calcService, ...)                 │
│     └─> NewActionLogController(logService, ...)                  │
│                                                                   │
│  8. Route Setup                                                  │
│     └─> NewRoute(controllers, middleware, router)                │
│                                                                   │
│  9. Seed Providers                                               │
│     ├─> NewRoleSeeder(userService)                               │
│     └─> NewDummySeeder(userService, formService)                 │
│                                                                   │
│  10. Application                                                 │
│     └─> NewApplication(route, seeder)                            │
│                                                                   │
└───────────────────────────────────┬───────────────────────────────┘
                                    │
                                    │ Run()
                                    ▼
                          ┌──────────────────┐
                          │   Gin Server     │
                          │   Port: 8080     │
                          └──────────────────┘
```

---

## 10. Error Handling & Recovery Flow

```
┌────────┐
│ Client │
└───┬────┘
    │
    │ HTTP Request
    │ POST /form/basic { invalid_data }
    │
    ▼
┌─────────────────────────────────────────────────────┐
│              Middleware Chain                       │
├─────────────────────────────────────────────────────┤
│                                                     │
│  1. CORS Middleware                                │
│     └─> Set CORS headers                           │
│                                                     │
│  2. Localization Middleware                        │
│     └─> Detect language (Accept-Language header)   │
│     └─> Set context: c.Set("language", "fa")       │
│                                                     │
│  3. Recovery Middleware                            │
│     └─> defer recover() {                          │
│           if r := recover(); r != nil {            │
│             // Handle panic                        │
│           }                                        │
│         }                                          │
│                                                     │
│  4. Auth Middleware (if protected route)           │
│     └─> Validate JWT                               │
│     └─> Check permissions                          │
│                                                     │
└─────────────────────┬───────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────┐
│              Controller Handler                     │
│         (CustomerFormController.CreateForm)         │
├─────────────────────────────────────────────────────┤
│                                                     │
│  1. Bind Request                                   │
│     └─> c.ShouldBindJSON(&dto)                     │
│     └─> If error → panic(BindingError)             │
│                                                     │
│  2. Validate DTO                                   │
│     └─> validator.Struct(dto)                      │
│     └─> If error → panic(ValidationError)          │
│                                                     │
│  3. Business Logic                                 │
│     └─> formService.CreateBasicInfoForm(dto)       │
│         └─> Check duplicates                       │
│         └─> If exists → panic(ConflictError)       │
│         └─> Database insert                        │
│         └─> If error → panic(InternalError)        │
│                                                     │
└─────────────────────┬───────────────────────────────┘
                      │
                      │ Exception thrown (panic)
                      ▼
┌─────────────────────────────────────────────────────┐
│          Recovery Middleware (defer)                │
├─────────────────────────────────────────────────────┤
│                                                     │
│  1. Catch Panic                                    │
│     └─> r := recover()                             │
│                                                     │
│  2. Type Switch on Exception                       │
│     switch err := r.(type) {                       │
│       case *ValidationError:                       │
│         └─> statusCode = 422                       │
│       case *AuthError:                             │
│         └─> statusCode = 401                       │
│       case *ForbiddenError:                        │
│         └─> statusCode = 403                       │
│       case *NotFoundError:                         │
│         └─> statusCode = 404                       │
│       case *ConflictError:                         │
│         └─> statusCode = 409                       │
│       default:                                     │
│         └─> statusCode = 500                       │
│     }                                              │
│                                                     │
│  3. Translate Error Message                       │
│     └─> lang := c.GetString("language")            │
│     └─> message := i18n.Translate(err.Key, lang)   │
│                                                     │
│  4. Build Response                                 │
│     └─> response := {                              │
│           "status_code": statusCode,               │
│           "message": message,                      │
│           "data": err.FieldErrors                  │
│         }                                          │
│                                                     │
│  5. Return JSON                                    │
│     └─> c.JSON(statusCode, response)               │
│                                                     │
└─────────────────────┬───────────────────────────────┘
                      │
                      ▼
                  ┌────────┐
                  │ Client │
                  │        │
                  │ {      │
                  │   "status_code": 422,             │
                  │   "message": "اطلاعات نامعتبر",   │
                  │   "data": {                       │
                  │     "birthDate": {                │
                  │       "required": "تاریخ تولد الزامی است" │
                  │     }                             │
                  │   }                               │
                  │ }                                 │
                  └────────────────────────────────────┘
```

---

## 11. Cache Strategy Flow (Redis)

```
┌─────────────────────────────────────────────────────────────────┐
│                     OTP Cache Pattern                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Key Format: "otp:{phone}"                                     │
│  Value: "123456"                                               │
│  TTL: 2 minutes (configurable)                                 │
│                                                                 │
│  Write Path:                                                   │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │  OTPService.GenerateOTP(phone)                          │ │
│  │    └─> Generate random 6-digit code                     │ │
│  │    └─> Redis.Set("otp:09123456789", "123456", 2min)    │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                 │
│  Read Path:                                                    │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │  OTPService.VerifyOTP(phone, code)                      │ │
│  │    └─> stored := Redis.Get("otp:09123456789")          │ │
│  │    └─> if stored == code { valid = true }              │ │
│  │    └─> Redis.Del("otp:09123456789")  // One-time use   │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│               Operator Validation Cache Pattern                 │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  OTP Key: "operator:validation:otp:{operatorID}:{phone}"       │
│  Value: "123456"                                               │
│  TTL: 2 minutes                                                │
│                                                                 │
│  Token Key: "operator:validation:token:{operatorID}:{userID}"  │
│  Value: "validated"                                            │
│  TTL: 5 minutes                                                │
│                                                                 │
│  Flow:                                                         │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │  1. Request OTP                                         │ │
│  │     └─> Redis.Set(                                      │ │
│  │           "operator:validation:otp:op123:09123456789",  │ │
│  │           "654321",                                     │ │
│  │           2min                                          │ │
│  │         )                                               │ │
│  │                                                         │ │
│  │  2. Verify OTP                                          │ │
│  │     └─> Validate code from Redis                       │ │
│  │     └─> Redis.Set(                                      │ │
│  │           "operator:validation:token:op123:user456",    │ │
│  │           "validated",                                  │ │
│  │           5min                                          │ │
│  │         )                                               │ │
│  │                                                         │ │
│  │  3. Create Form (must have valid token)                │ │
│  │     └─> exists := Redis.Exists(                         │ │
│  │           "operator:validation:token:op123:user456"     │ │
│  │         )                                               │ │
│  │     └─> if !exists { return Forbidden }                │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                  User Cache Pattern (Optional)                  │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Key Format: "user:{userID}"                                   │
│  Value: JSON serialized User + Roles + Permissions             │
│  TTL: 15 minutes                                               │
│                                                                 │
│  Read-Through Cache:                                           │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │  UserService.GetUserByID(userID)                        │ │
│  │    ├─> cached := Redis.Get("user:uuid-123")             │ │
│  │    │   if cached != nil {                               │ │
│  │    │     return deserialize(cached)  // Cache HIT       │ │
│  │    │   }                                                 │ │
│  │    │                                                     │ │
│  │    └─> user := DB.Query(userID)  // Cache MISS          │ │
│  │        └─> Redis.Set("user:uuid-123", serialize(user))  │ │
│  │        └─> return user                                  │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                 │
│  Invalidation:                                                 │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │  On user/role/permission change:                        │ │
│  │    └─> Redis.Del("user:uuid-123")                       │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 12. Transaction Flow Pattern

```
┌─────────────────────────────────────────────────────────────────┐
│              Form Creation with Transaction                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  FormService.CreateBasicInfoForm(dto)                          │
│                                                                 │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │  1. Begin Transaction                                   │ │
│  │     db.Begin()                                          │ │
│  │     └─> tx := GORM.Transaction                          │ │
│  │                                                         │ │
│  │  2. Create Form                                         │ │
│  │     tx.Create(&Form{                                    │ │
│  │       UserID: userID,                                   │ │
│  │       Status: DRAFT                                     │ │
│  │     })                                                  │ │
│  │     └─> If error: Rollback & Return                     │ │
│  │                                                         │ │
│  │  3. Create BasicInfo                                    │ │
│  │     tx.Create(&BasicInfo{                               │ │
│  │       FormID: form.ID,                                  │ │
│  │       Gender: dto.Gender,                               │ │
│  │       BirthDate: dto.BirthDate                          │ │
│  │     })                                                  │ │
│  │     └─> If error: Rollback & Return                     │ │
│  │                                                         │ │
│  │  4. Create ActionLog                                    │ │
│  │     tx.Create(&ActionLog{                               │ │
│  │       FormID: form.ID,                                  │ │
│  │       UserID: userID,                                   │ │
│  │       ActionType: CREATE_FORM                           │ │
│  │     })                                                  │ │
│  │     └─> If error: Rollback & Return                     │ │
│  │                                                         │ │
│  │  5. Commit Transaction                                  │ │
│  │     tx.Commit()                                         │ │
│  │     └─> All changes atomic                              │ │
│  │                                                         │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                 │
│  Benefits:                                                     │
│  • Atomicity: All-or-nothing execution                        │
│  • Consistency: Database remains in valid state               │
│  • Isolation: Concurrent transactions don't interfere         │
│  • Durability: Committed changes are permanent                │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 13. Key Design Patterns

| Pattern | Location | Purpose |
|---------|----------|---------|
| **Clean Architecture** | Entire codebase | Separation of concerns, testability |
| **Repository Pattern** | `infrastructure/repository/` | Abstract data access |
| **Dependency Injection** | `wire/` | Manage dependencies, testability |
| **Singleton** | `database/PostgresDatabase` | Single DB connection pool |
| **Middleware Chain** | `presentation/middleware/` | Request processing pipeline |
| **DTO Pattern** | `application/dto/` | Data transfer, validation |
| **Service Layer** | `application/service/` | Business logic encapsulation |
| **Exception Handling** | `domain/exception/` | Centralized error management |
| **Strategy Pattern** | SMS, Storage interfaces | Swappable implementations |
| **Factory Pattern** | Wire providers | Object creation |

---

## 14. Security Measures

```
┌─────────────────────────────────────────────────────────────────┐
│                     Security Layers                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  1. Authentication                                             │
│     ├─> Phone-based OTP (2-factor)                             │
│     ├─> JWT with RSA-256 signing                               │
│     ├─> Token expiry (configurable)                            │
│     └─> Refresh token rotation                                 │
│                                                                 │
│  2. Authorization                                              │
│     ├─> Role-Based Access Control (RBAC)                       │
│     ├─> Granular permissions                                   │
│     ├─> Permission categories (9 types)                        │
│     └─> Middleware enforcement                                 │
│                                                                 │
│  3. Data Protection                                            │
│     ├─> Password hashing (assumed bcrypt/argon2)               │
│     ├─> OTP one-time use + expiry                              │
│     ├─> Validation tokens with TTL                             │
│     └─> HTTPS enforcement                                      │
│                                                                 │
│  4. Input Validation                                           │
│     ├─> DTO validation (go-validator)                          │
│     ├─> Request binding checks                                 │
│     ├─> Type safety (strong typing)                            │
│     └─> SQL injection prevention (GORM parameterized queries)  │
│                                                                 │
│  5. Rate Limiting                                              │
│     ├─> OTP max attempts (default: 3)                          │
│     └─> RateLimitError exception type                          │
│                                                                 │
│  6. Audit Trail                                                │
│     ├─> ActionLog for all operations                           │
│     ├─> User tracking                                          │
│     └─> Timestamp logging                                      │
│                                                                 │
│  7. File Upload Security                                       │
│     ├─> S3 presigned URLs (time-limited)                       │
│     ├─> Bucket isolation (mamography, cancer)                  │
│     └─> File type validation (assumed)                         │
│                                                                 │
│  8. CORS Protection                                            │
│     ├─> CORS middleware                                        │
│     └─> Allowed origins configuration                          │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 15. Scalability Considerations

```
┌─────────────────────────────────────────────────────────────────┐
│                  Horizontal Scaling Strategy                    │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Load Balancer (e.g., Nginx)                                   │
│         │                                                       │
│         ├─────────┬─────────┬─────────┐                        │
│         ▼         ▼         ▼         ▼                        │
│    Backend 1  Backend 2  Backend 3  Backend N                  │
│         │         │         │         │                        │
│         └─────────┴─────────┴─────────┘                        │
│                     │                                           │
│         ┌───────────┼───────────┐                              │
│         │           │           │                              │
│         ▼           ▼           ▼                              │
│   PostgreSQL    Redis       S3 Storage                         │
│   (Master +     (Cluster)   (Distributed)                      │
│    Replicas)                                                   │
│                                                                 │
│  Stateless Design:                                             │
│  • No server-side sessions (JWT only)                          │
│  • All state in Redis/PostgreSQL                               │
│  • Each backend instance is identical                          │
│                                                                 │
│  Database Scaling:                                             │
│  • Read replicas for queries                                   │
│  • Master for writes                                           │
│  • Connection pooling (GORM)                                   │
│                                                                 │
│  Cache Scaling:                                                │
│  • Redis Cluster for high availability                         │
│  • Consistent hashing for key distribution                     │
│                                                                 │
│  File Storage Scaling:                                         │
│  • S3 auto-scales                                              │
│  • CDN for static assets (optional)                            │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## Summary

This Context Flow Diagram provides a comprehensive view of the FamCan Backend system architecture, including:

1. **System Context** - External actors and services
2. **Authentication Flows** - OTP and JWT-based login
3. **Form Management** - Patient and operator workflows
4. **Risk Calculation** - Integration with external APIs
5. **File Upload** - S3 storage integration
6. **Authorization** - RBAC implementation
7. **Layer Architecture** - Clean architecture separation
8. **Dependency Injection** - Wire-based DI
9. **Error Handling** - Centralized exception recovery
10. **Caching Strategy** - Redis patterns
11. **Transaction Management** - ACID compliance
12. **Design Patterns** - Architectural patterns used
13. **Security** - Multi-layered security approach
14. **Scalability** - Horizontal scaling strategy

The system is designed for maintainability, scalability, and security in a healthcare context.