# FamCan Backend - Module Documentation

## Table of Contents
1. [Authentication & User Management](#authentication--user-management)
2. [Form Management](#form-management)
3. [Risk Calculation](#risk-calculation)
4. [Action Logging](#action-logging)
5. [Enum & Reference Data](#enum--reference-data)
6. [Data Models](#data-models)

---

## Authentication & User Management

### Module Overview
Handles user authentication via phone-based OTP verification, role-based access control (RBAC), and permission management for the system.

**File References:**
- [user_service.go](../internal/application/usecase/user_service.go) - Service interface
- [user_service_impl.go](../internal/application/service/user_service_impl.go) - Implementation
- [user_general.go](../internal/presentation/controller/user/general.go) - User endpoints
- [user_admin.go](../internal/presentation/controller/user/admin.go) - Admin endpoints

### Submodule: Phone Login with OTP

#### Inputs
**Endpoint:** `POST /auth/login`

```json
{
  "phone": "string (format: +98XXXXXXXXXX or 09XXXXXXXXX)",
  "deviceId": "string (optional, for session tracking)"
}
```

**Validation Rules:**
- Phone number is required
- Phone must be in valid Persian/Iranian format
- Phone format: `+989XXXXXXXXX` or `09XXXXXXXXX`

#### Outputs
**Success (200):**
```json
{
  "success": true,
  "message": "OTP sent successfully",
  "data": {
    "expiresAt": "2025-11-17T14:30:00Z",
    "resendAfter": 30
  }
}
```

**Error (400) - Invalid Phone:**
```json
{
  "success": false,
  "message": "Invalid phone number format",
  "error": "VALIDATION_ERROR",
  "details": {
    "field": "phone",
    "issue": "Phone must be in format +989XXXXXXXXX or 09XXXXXXXXX"
  }
}
```

#### Logic
1. Validate phone number format (Persian/Iranian phone)
2. Generate a 6-digit OTP code
3. Store OTP in Redis with 10-minute expiration
4. Send OTP via SMS using Asanak provider
5. Return success response with OTP expiration time
6. Implement rate limiting (max 3 requests per phone per 5 minutes)

---

### Submodule: OTP Verification & Token Generation

#### Inputs
**Endpoint:** `POST /auth/verify-otp`

```json
{
  "phone": "string (format: +98XXXXXXXXXX or 09XXXXXXXXX)",
  "otp": "string (6 digits)",
  "deviceId": "string (optional)"
}
```

**Validation Rules:**
- Phone and OTP are required
- OTP must be exactly 6 digits
- Phone format must match login request

#### Outputs
**Success (200):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "tokenType": "Bearer",
    "expiresIn": 86400,
    "user": {
      "id": "user-uuid",
      "phone": "09XXXXXXXXX",
      "roles": [
        {
          "id": "role-uuid",
          "name": "USER"
        }
      ]
    }
  }
}
```

**Error (401) - Invalid OTP:**
```json
{
  "success": false,
  "message": "Invalid OTP",
  "error": "AUTH_ERROR",
  "details": {
    "attemptRemaining": 2
  }
}
```

**Error (401) - OTP Expired:**
```json
{
  "success": false,
  "message": "OTP has expired",
  "error": "AUTH_ERROR"
}
```

#### Logic
1. Validate phone and OTP format
2. Retrieve OTP from Redis cache
3. Check if OTP has expired (10-minute window)
4. Verify OTP matches stored value
5. Find or create user with phone number
6. Load user roles and permissions
7. Generate JWT token (24-hour expiration)
8. Clear OTP from Redis after successful verification
9. Return JWT token and user information

---

### Submodule: Admin Password Login

#### Inputs
**Endpoint:** `POST /auth/admin/login`

```json
{
  "phone": "string",
  "password": "string (minimum 8 characters)"
}
```

**Validation Rules:**
- Both phone and password are required
- Password minimum 8 characters
- Only admin users can use this endpoint

#### Outputs
**Success (200):**
```json
{
  "success": true,
  "message": "Admin login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "tokenType": "Bearer",
    "expiresIn": 86400,
    "user": {
      "id": "user-uuid",
      "phone": "09XXXXXXXXX",
      "roles": [
        {
          "id": "role-uuid",
          "name": "ADMIN"
        }
      ]
    }
  }
}
```

**Error (401) - Invalid Credentials:**
```json
{
  "success": false,
  "message": "Invalid phone or password",
  "error": "AUTH_ERROR"
}
```

#### Logic
1. Validate phone and password format
2. Find user by phone in database
3. Verify user has admin role
4. Compare provided password with stored hash (bcrypt)
5. Increment failed login attempts on failure
6. Lock account after 5 failed attempts
7. Generate JWT token on success
8. Return token and user information

---

### Submodule: Role Management

#### Inputs
**Endpoint:** `POST /admin/role` (Create)

```json
{
  "name": "string (3-50 characters, unique)",
  "description": "string (optional, max 500 characters)",
  "permissionIds": ["uuid", "uuid"]
}
```

**Endpoint:** `PUT /admin/role/:roleID` (Update)

```json
{
  "name": "string (optional)",
  "description": "string (optional)",
  "permissionIds": ["uuid"]
}
```

#### Outputs
**Success (201) - Create:**
```json
{
  "success": true,
  "message": "Role created successfully",
  "data": {
    "id": "role-uuid",
    "name": "OPERATOR",
    "description": "Form operator with limited permissions",
    "permissions": [
      {
        "id": "perm-uuid",
        "type": "CategoryFormManagement"
      }
    ]
  }
}
```

**Success (200) - Update:**
```json
{
  "success": true,
  "message": "Role updated successfully",
  "data": {
    "id": "role-uuid",
    "name": "OPERATOR",
    "permissions": [...]
  }
}
```

**Error (409) - Duplicate Name:**
```json
{
  "success": false,
  "message": "Role name already exists",
  "error": "CONFLICT_ERROR"
}
```

**Error (400) - Validation Failed:**
```json
{
  "success": false,
  "message": "Validation failed",
  "error": "VALIDATION_ERROR",
  "details": [
    {
      "field": "name",
      "issue": "Name must be between 3 and 50 characters"
    }
  ]
}
```

#### Logic
1. Validate role name uniqueness (case-insensitive)
2. Validate all permission IDs exist in database
3. Create or update role in database
4. Associate permissions with role
5. Return created/updated role with permissions

---

### Submodule: Permission Management

#### Inputs
**Endpoint:** `GET /admin/permission` (List all)

Query Parameters:
```
page: integer (default: 1)
limit: integer (default: 10, max: 100)
```

**Endpoint:** `GET /admin/permission/:permissionID/roles` (Get roles with permission)

#### Outputs
**Success (200) - List Permissions:**
```json
{
  "success": true,
  "data": [
    {
      "id": "perm-uuid",
      "type": "CategoryFormManagement",
      "description": "Can manage all forms"
    },
    {
      "id": "perm-uuid-2",
      "type": "PermissionCreateFormForUser",
      "description": "Can create forms for other users"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 15
  }
}
```

**Success (200) - Get Permission Roles:**
```json
{
  "success": true,
  "data": {
    "permissionId": "perm-uuid",
    "permissionType": "CategoryFormManagement",
    "roles": [
      {
        "id": "role-uuid",
        "name": "ADMIN",
        "userCount": 5
      },
      {
        "id": "role-uuid-2",
        "name": "OPERATOR",
        "userCount": 12
      }
    ]
  }
}
```

#### Logic
1. Query all permissions or specific permission's roles
2. Apply pagination filters if listing
3. Load associated roles and user counts
4. Return formatted permission/role data

---

### Submodule: User Role Assignment

#### Inputs
**Endpoint:** `PUT /admin/user/:userID/role` (Update user roles)

```json
{
  "roleIds": ["role-uuid-1", "role-uuid-2"]
}
```

#### Outputs
**Success (200):**
```json
{
  "success": true,
  "message": "User roles updated successfully",
  "data": {
    "userId": "user-uuid",
    "phone": "09XXXXXXXXX",
    "roles": [
      {
        "id": "role-uuid",
        "name": "OPERATOR"
      }
    ]
  }
}
```

**Error (404) - User Not Found:**
```json
{
  "success": false,
  "message": "User not found",
  "error": "NOT_FOUND_ERROR"
}
```

#### Logic
1. Find user by ID
2. Validate all role IDs exist
3. Remove all current user roles
4. Assign new roles to user
5. Load and return updated user with roles

---

### Submodule: User Password Management

#### Inputs
**Endpoint:** `PUT /admin/user/password` (Set user password - Admin only)

```json
{
  "userId": "uuid",
  "password": "string (minimum 8 characters)"
}
```

**Validation Rules:**
- Password minimum 8 characters
- Should contain mix of uppercase, lowercase, numbers
- Cannot reuse last 3 passwords

#### Outputs
**Success (200):**
```json
{
  "success": true,
  "message": "Password updated successfully",
  "data": {
    "userId": "user-uuid",
    "phone": "09XXXXXXXXX"
  }
}
```

**Error (400) - Weak Password:**
```json
{
  "success": false,
  "message": "Password does not meet requirements",
  "error": "VALIDATION_ERROR",
  "details": {
    "requirements": [
      "At least 8 characters",
      "Must contain uppercase letter",
      "Must contain lowercase letter",
      "Must contain number"
    ]
  }
}
```

#### Logic
1. Validate password strength
2. Check password history (prevent reuse)
3. Hash password using bcrypt (cost: 10)
4. Update user password in database
5. Clear any existing sessions/tokens
6. Return success message

---

### Submodule: User Listing & Filtering

#### Inputs
**Endpoint:** `GET /admin/user` (List users with filters)

Query Parameters:
```
page: integer (default: 1)
limit: integer (default: 10, max: 100)
roleId: uuid (optional - filter by role)
phone: string (optional - search by phone)
status: string (optional - active/inactive)
```

#### Outputs
**Success (200):**
```json
{
  "success": true,
  "data": [
    {
      "id": "user-uuid",
      "phone": "09XXXXXXXXX",
      "roles": [
        {
          "id": "role-uuid",
          "name": "OPERATOR"
        }
      ],
      "createdAt": "2025-10-15T10:30:00Z",
      "lastLogin": "2025-11-17T09:20:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 45
  }
}
```

#### Logic
1. Apply filters (role, phone search, status)
2. Query users with pagination
3. Load user roles and last login info
4. Return paginated user list

---

## Form Management

### Module Overview
Handles creation, modification, and management of medical forms with multiple subsections for health information collection.

**File References:**
- [form_service.go](../internal/application/usecase/form_service.go) - Service interface
- [form_service_impl.go](../internal/application/service/form_service_impl.go) - Implementation
- [form_customer.go](../internal/presentation/controller/form/customer.go) - Customer endpoints
- [form_admin.go](../internal/presentation/controller/form/admin.go) - Admin endpoints
- [form_general.go](../internal/presentation/controller/form/general.go) - General endpoints
- [form.go](../internal/domain/entity/form.go) - Form entity

### Submodule: Basic Form Creation

#### Inputs
**Endpoint:** `POST /form/basic` (Customer) or `POST /admin/operator/form` (Admin on behalf of user)

Customer Request:
```json
{
  "userId": "uuid (auto-filled from auth token)"
}
```

Admin Operator Request:
```json
{
  "userId": "uuid (target user)",
  "operatorId": "uuid (auto-filled from auth token)",
  "notes": "string (optional)"
}
```

#### Outputs
**Success (201):**
```json
{
  "success": true,
  "message": "Form created successfully",
  "data": {
    "id": "form-uuid",
    "userId": "user-uuid",
    "status": "Draft",
    "createdAt": "2025-11-17T10:00:00Z",
    "createdBy": "user-uuid",
    "lastModifiedAt": "2025-11-17T10:00:00Z",
    "completionPercentage": 0
  }
}
```

**Error (409) - Draft Form Already Exists:**
```json
{
  "success": false,
  "message": "User already has a draft form in progress",
  "error": "CONFLICT_ERROR",
  "details": {
    "existingFormId": "form-uuid"
  }
}
```

#### Logic
1. Validate user exists and is not suspended
2. Check if user already has a draft form
3. Create new form with status "Draft"
4. Initialize empty subsection placeholders
5. Set creation timestamp and creator
6. Return created form details

---

### Submodule: Basic Information Form Section

#### Inputs
**Endpoint:** `PUT /form/:formID/basic` (Customer) or `PATCH /admin/form/:formID/basic` (Admin)

```json
{
  "gender": "enum (Male | Female | Other)",
  "birthDate": "string (format: YYYY-MM-DD, must be 18+)",
  "socialSecurityNumber": "string (optional, 10 digits)",
  "height": "number (cm, 100-250)",
  "weight": "number (kg, 20-300)"
}
```

**Validation Rules:**
- Gender is required
- Birth date is required and age must be 18+
- Social Security Number must be 10 digits if provided
- Height: 100-250 cm
- Weight: 20-300 kg

#### Outputs
**Success (200):**
```json
{
  "success": true,
  "message": "Basic information updated successfully",
  "data": {
    "id": "basic-info-uuid",
    "formId": "form-uuid",
    "gender": "Female",
    "birthDate": "1990-05-15",
    "age": 35,
    "socialSecurityNumber": "1234567890",
    "height": 168,
    "weight": 65,
    "bmi": 23.0,
    "updatedAt": "2025-11-17T10:15:00Z"
  }
}
```

**Error (400) - Age Validation:**
```json
{
  "success": false,
  "message": "Validation failed",
  "error": "VALIDATION_ERROR",
  "details": [
    {
      "field": "birthDate",
      "issue": "User must be at least 18 years old"
    }
  ]
}
```

**Error (404) - Form Not Found:**
```json
{
  "success": false,
  "message": "Form not found",
  "error": "NOT_FOUND_ERROR"
}
```

#### Logic
1. Find form by ID and verify user has access
2. Validate all input fields
3. Calculate age from birth date
4. Calculate BMI from height and weight
5. Create or update BasicInfo record
6. Update form lastModifiedAt timestamp
7. Recalculate form completion percentage
8. Return updated basic info with calculated fields

---

### Submodule: General Health Information

#### Inputs
**Endpoint:** `PUT /form/:formID/generalhealth`

```json
{
  "drinksAlcohol": {
    "drinks": "enum (Never | Former | Current)",
    "frequency": "string (optional, e.g., 'weekly', 'daily')",
    "unitsPerWeek": "number (optional, 0-50)"
  },
  "smokingStatus": {
    "status": "enum (Never | Former | Current)",
    "yearsSmoked": "number (optional, 0-80)",
    "quitYear": "number (optional, if Former)",
    "packYears": "number (optional, calculated as: (cigarettes/day * years) / 20)"
  },
  "diet": {
    "vegetableServingsPerDay": "number (0-10)",
    "fruitServingsPerDay": "number (0-10)",
    "wholeGrainServingsPerDay": "number (0-10)",
    "processedMeatServingsPerWeek": "number (0-20)"
  },
  "exercise": {
    "minutesPerWeek": "number (0-500)",
    "intensity": "enum (Light | Moderate | Vigorous)",
    "frequency": "enum (Daily | EveryOtherDay | TwiceAWeek | Once | Never)"
  }
}
```

#### Outputs
**Success (200):**
```json
{
  "success": true,
  "message": "General health information updated successfully",
  "data": {
    "id": "health-info-uuid",
    "formId": "form-uuid",
    "drinksAlcohol": {
      "drinks": "Current",
      "frequency": "weekly",
      "unitsPerWeek": 4
    },
    "smokingStatus": {
      "status": "Former",
      "yearsSmoked": 10,
      "quitYear": 2015,
      "packYears": 15
    },
    "diet": {
      "vegetableServingsPerDay": 3,
      "fruitServingsPerDay": 2,
      "wholeGrainServingsPerDay": 2,
      "processedMeatServingsPerWeek": 2
    },
    "exercise": {
      "minutesPerWeek": 150,
      "intensity": "Moderate",
      "frequency": "EveryOtherDay"
    },
    "updatedAt": "2025-11-17T10:20:00Z"
  }
}
```

#### Logic
1. Validate all lifestyle values are within acceptable ranges
2. Calculate pack-years for smoking history
3. Validate alcohol units per week
4. Validate exercise frequency and intensity consistency
5. Create or update GeneralHealthInfo record
6. Update form modification timestamp
7. Recalculate form completion percentage

---

### Submodule: Mamography Information

#### Inputs
**Endpoint:** `PUT /form/:formID/mamography`

```json
{
  "numberOfChildren": "number (0-20, optional)",
  "ageAtFirstBirth": "number (optional, 16-50)",
  "menopauseStatus": "enum (PreMenopausal | Perimenopausal | PostMenopausal | Unknown)",
  "ageAtMenopause": "number (optional if PostMenopausal, 35-65)",
  "onHormoneReplacementTherapy": "boolean",
  "hormoneReplacementTherapyDuration": "number (optional, years, 0-40)",
  "numberOfBiopsies": "number (0-10, optional)",
  "hadHyperplasiaInBiopsy": "boolean",
  "hyperplasiaType": "enum (Usual | Atypical | None | Unknown)"
}
```

#### Outputs
**Success (200):**
```json
{
  "success": true,
  "message": "Mamography information updated successfully",
  "data": {
    "id": "mamography-uuid",
    "formId": "form-uuid",
    "numberOfChildren": 2,
    "ageAtFirstBirth": 28,
    "menopauseStatus": "PostMenopausal",
    "ageAtMenopause": 52,
    "onHormoneReplacementTherapy": true,
    "hormoneReplacementTherapyDuration": 5,
    "numberOfBiopsies": 1,
    "hadHyperplasiaInBiopsy": false,
    "hyperplasiaType": "None",
    "riskFactorsCount": 3,
    "updatedAt": "2025-11-17T10:25:00Z"
  }
}
```

#### Logic
1. Validate age at menopause only if PostMenopausal
2. Validate age at first birth is less than current age
3. Validate HRT duration only if on HRT
4. Validate hyperplasia type is set if biopsies exist
5. Count risk factors (children, HRT, biopsies, hyperplasia)
6. Create or update MamoGraphyInfo record
7. Update form modification timestamp
8. Return updated mamography info with risk factor count

---

### Submodule: Personal Cancer History

#### Inputs
**Endpoint:** `POST /form/:formID/cancer` (Add)

```json
{
  "cancerType": "enum (Breast | Lung | Ovarian | Uterine | Pancreatic | Prostate | Colorectal | Other)",
  "cancerAge": "number (18-120, optional)",
  "picturePath": "string (optional, S3 object key)",
  "notes": "string (optional, max 500 chars)"
}
```

**Endpoint:** `PUT /form/:formID/cancer/:cancerID` (Update)
Same structure as POST

**Endpoint:** `DELETE /form/:formID/cancer/:cancerID` (Delete)
No body required

#### Outputs
**Success (201) - Create:**
```json
{
  "success": true,
  "message": "Cancer record created successfully",
  "data": {
    "id": "cancer-uuid",
    "formId": "form-uuid",
    "cancerType": "Breast",
    "cancerAge": 45,
    "picturePath": "forms/user-uuid/cancer/image-123.jpg",
    "notes": "Diagnosed in 2015, treated with chemotherapy",
    "createdAt": "2025-11-17T10:30:00Z"
  }
}
```

**Success (200) - Update:**
```json
{
  "success": true,
  "message": "Cancer record updated successfully",
  "data": {
    "id": "cancer-uuid",
    "formId": "form-uuid",
    "cancerType": "Breast",
    "cancerAge": 46,
    "updatedAt": "2025-11-17T10:35:00Z"
  }
}
```

**Success (204) - Delete:** No content

**Error (400) - Age Validation:**
```json
{
  "success": false,
  "message": "Validation failed",
  "error": "VALIDATION_ERROR",
  "details": [
    {
      "field": "cancerAge",
      "issue": "Cancer age must be greater than birth date and less than current age"
    }
  ]
}
```

#### Logic
1. Validate cancer age is reasonable (18-120 and less than current age)
2. Validate cancer type is from enum
3. If picture provided, validate it exists in S3
4. Create or update cancer record
5. Update form lastModifiedAt
6. Return cancer record with metadata

---

### Submodule: Family Cancer History

#### Inputs
**Endpoint:** `POST /form/:formID/familycancer` (Add)

```json
{
  "relativeType": "enum (Mother | Father | Sister | Brother | Daughter | Son | Grandmother | Grandfather | Aunt | Uncle | Cousin)",
  "cancerType": "enum (Breast | Lung | Ovarian | Uterine | Pancreatic | Prostate | Colorectal | Other)",
  "relativeAge": "number (optional, at cancer diagnosis, 18-120)",
  "relativeDeceased": "boolean",
  "deceasedAge": "number (optional, if deceased, 18-120)",
  "notes": "string (optional, max 500 chars)"
}
```

**Endpoint:** `PUT /form/:formID/familycancer/:familyCancerID` (Update)

**Endpoint:** `DELETE /form/:formID/familycancer/:familyCancerID` (Delete)

#### Outputs
**Success (201) - Create:**
```json
{
  "success": true,
  "message": "Family cancer record created successfully",
  "data": {
    "id": "family-cancer-uuid",
    "formId": "form-uuid",
    "relativeType": "Mother",
    "cancerType": "Breast",
    "relativeAge": 52,
    "relativeDeceased": false,
    "deceasedAge": null,
    "notes": "Diagnosed 10 years ago, currently in remission",
    "createdAt": "2025-11-17T10:40:00Z"
  }
}
```

**Error (400) - Age Consistency:**
```json
{
  "success": false,
  "message": "Validation failed",
  "error": "VALIDATION_ERROR",
  "details": [
    {
      "field": "deceasedAge",
      "issue": "Deceased age must be greater than cancer age"
    }
  ]
}
```

#### Logic
1. Validate relative type is from enum
2. Validate cancer type is from enum
3. If deceased, validate deceased age > cancer age
4. Create or update family cancer record
5. Update form lastModifiedAt
6. Return family cancer record

---

### Submodule: Contact Information

#### Inputs
**Endpoint:** `PUT /form/:formID/contact`

```json
{
  "fullName": "string (required, 3-100 chars)",
  "email": "string (optional, valid email format)",
  "address": "string (optional, max 500 chars)",
  "city": "string (optional)",
  "province": "string (optional)",
  "country": "string (optional)",
  "birthCountry": "string (optional)",
  "zipCode": "string (optional, format depends on country)"
}
```

#### Outputs
**Success (200):**
```json
{
  "success": true,
  "message": "Contact information updated successfully",
  "data": {
    "id": "contact-uuid",
    "formId": "form-uuid",
    "fullName": "فاطمه احمدی",
    "email": "fatima.ahmadi@email.com",
    "address": "خیابان شریف، پلاک 123",
    "city": "تهران",
    "province": "تهران",
    "country": "ایران",
    "birthCountry": "ایران",
    "zipCode": "1234567890",
    "updatedAt": "2025-11-17T10:45:00Z"
  }
}
```

#### Logic
1. Validate full name length and format
2. Validate email format if provided
3. Normalize address and location data
4. Create or update ContactInfo record
5. Update form lastModifiedAt
6. Return updated contact information

---

### Submodule: Lung Cancer Information

#### Inputs
**Endpoint:** `PUT /form/:formID/lungcancer`

```json
{
  "smokingHistory": {
    "status": "enum (Never | Former | Current)",
    "yearsSmoked": "number (optional, 0-80)",
    "quitYear": "number (optional, if Former)",
    "packYears": "number (optional, cigarettes/day * years / 20)"
  },
  "occupationalExposure": {
    "exposed": "boolean",
    "exposureTypes": ["string (asbestos | silica | radon | diesel | other)"],
    "yearsExposed": "number (optional, 0-60)"
  },
  "insuranceStatus": "enum (Insured | Uninsured | PartiallyInsured)",
  "comorbidities": {
    "hasCOPD": "boolean",
    "hasAsthma": "boolean",
    "hasPulmonaryFibrosis": "boolean",
    "previousLungDisease": "boolean"
  }
}
```

#### Outputs
**Success (200):**
```json
{
  "success": true,
  "message": "Lung cancer information updated successfully",
  "data": {
    "id": "lung-cancer-uuid",
    "formId": "form-uuid",
    "smokingHistory": {
      "status": "Former",
      "yearsSmoked": 20,
      "quitYear": 2010,
      "packYears": 40
    },
    "occupationalExposure": {
      "exposed": true,
      "exposureTypes": ["asbestos"],
      "yearsExposed": 15
    },
    "insuranceStatus": "Insured",
    "comorbidities": {
      "hasCOPD": false,
      "hasAsthma": true,
      "hasPulmonaryFibrosis": false,
      "previousLungDisease": true
    },
    "riskLevel": "High",
    "updatedAt": "2025-11-17T10:50:00Z"
  }
}
```

#### Logic
1. Validate smoking status and history
2. Validate occupational exposure types
3. Validate exposure duration if exposed
4. Assess risk level based on smoking + occupational history
5. Create or update LungCancerInfo record
6. Update form lastModifiedAt
7. Return lung cancer info with risk assessment

---

### Submodule: Form Status Management

#### Inputs
**Endpoint:** `PUT /form/:formID/status` (Customer)

```json
{
  "status": "enum (Draft | Submitted)"
}
```

**Endpoint:** `PUT /admin/form/:formID/accept` or `/reject` (Admin)

```json
{
  "notes": "string (optional, max 500 chars)"
}
```

#### Outputs
**Success (200):**
```json
{
  "success": true,
  "message": "Form status changed to Submitted",
  "data": {
    "id": "form-uuid",
    "userId": "user-uuid",
    "status": "Submitted",
    "completionPercentage": 85,
    "submittedAt": "2025-11-17T11:00:00Z",
    "lastModifiedAt": "2025-11-17T11:00:00Z"
  }
}
```

**Success (200) - Admin Accept:**
```json
{
  "success": true,
  "message": "Form accepted successfully",
  "data": {
    "id": "form-uuid",
    "userId": "user-uuid",
    "status": "Approved",
    "approvedAt": "2025-11-17T11:05:00Z",
    "approvedBy": "admin-user-uuid",
    "approvalNotes": "All information verified"
  }
}
```

**Error (400) - Incomplete Form:**
```json
{
  "success": false,
  "message": "Cannot submit incomplete form",
  "error": "VALIDATION_ERROR",
  "details": {
    "completionPercentage": 65,
    "missingFields": ["mamography", "contact"]
  }
}
```

#### Logic
1. Validate status transition is allowed
2. If submitting: verify form completion (min 70%)
3. If admin approving: verify all sections are filled
4. Update form status and timestamps
5. Create action log entry
6. If rejected: notify user with rejection reason
7. Return updated form

---

### Submodule: Form Listing & Pagination

#### Inputs
**Endpoint:** `GET /form` (Customer list own forms)

Query Parameters:
```
page: integer (default: 1)
limit: integer (default: 10, max: 100)
status: string (optional, Draft | Submitted | Approved | Rejected)
```

**Endpoint:** `GET /admin/form` (Admin list all forms)

Query Parameters:
```
page: integer (default: 1)
limit: integer (default: 10, max: 100)
userId: uuid (optional - filter by user)
status: string (optional)
operatorId: uuid (optional - assigned operator)
dateFrom: string (optional, YYYY-MM-DD)
dateTo: string (optional, YYYY-MM-DD)
```

#### Outputs
**Success (200):**
```json
{
  "success": true,
  "data": [
    {
      "id": "form-uuid",
      "userId": "user-uuid",
      "userPhone": "09XXXXXXXXX",
      "status": "Approved",
      "completionPercentage": 95,
      "createdAt": "2025-11-15T10:00:00Z",
      "submittedAt": "2025-11-16T14:30:00Z",
      "approvedAt": "2025-11-17T08:15:00Z",
      "operatorId": "operator-uuid",
      "operatorName": "علی محمدی"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 45,
    "pages": 5
  }
}
```

#### Logic
1. Apply filters based on user role and permissions
2. Query forms with pagination
3. Load related user and operator information
4. Calculate completion percentage for each form
5. Sort by most recent first
6. Return paginated list

---

## Risk Calculation

### Module Overview
Integrates with external medical risk assessment models to calculate cancer risk based on patient medical history.

**File References:**
- [calc_service.go](../internal/application/usecase/calc_service.go) - Service interface
- [calc_service_impl.go](../internal/application/service/calc_service_impl.go) - Implementation
- [calc_admin.go](../internal/presentation/controller/calc/admin.go) - Admin endpoints
- [calc_general.go](../internal/presentation/controller/calc/general.go) - General endpoints

### Submodule: PREMM5 Risk Calculation (Lynch Syndrome)

#### Inputs
**Endpoint:** `POST /admin/calc/model`

```json
{
  "formId": "uuid",
  "modelType": "enum (PREMM5 | BCRA | GAIL | PLCO)",
  "recalculate": "boolean (optional, default: false)"
}
```

Or Query-based: `GET /admin/calc/premm5/:formID`

**Required Form Data for PREMM5:**
- Age
- Gender (Female required)
- Personal cancer history (colorectal, endometrial, ovarian)
- Family history (colorectal, endometrial, ovarian cancers)
- Microsatellite instability (MSI) status if available

#### Outputs
**Success (200):**
```json
{
  "success": true,
  "message": "PREMM5 calculation completed",
  "data": {
    "formId": "form-uuid",
    "modelType": "PREMM5",
    "calculationDate": "2025-11-17T11:30:00Z",
    "results": {
      "mlh1MutationRisk": 0.015,
      "msh2MutationRisk": 0.022,
      "msh6MutationRisk": 0.008,
      "pms2MutationRisk": 0.005,
      "anyLynchSyndromeRisk": 0.050,
      "riskCategory": "Moderate"
    },
    "interpretation": {
      "recommendation": "Consider genetic testing and counseling",
      "followUpActions": [
        "Refer to genetic counselor",
        "Consider colonoscopy screening",
        "Annual surveillance recommended"
      ]
    }
  }
}
```

**Error (400) - Insufficient Data:**
```json
{
  "success": false,
  "message": "Cannot calculate PREMM5 risk - insufficient form data",
  "error": "VALIDATION_ERROR",
  "details": {
    "missingData": [
      "Personal cancer history",
      "Complete family cancer history"
    ],
    "requiredFields": [
      "age",
      "gender (must be Female)",
      "cancer information"
    ]
  }
}
```

#### Logic
1. Validate form has required data (age, gender, cancer history)
2. Validate gender is Female (Lynch syndrome risk for women)
3. Extract cancer information from form
4. Prepare payload for external PREMM5 model
5. Send form data to external service
6. Parse and validate response
7. Store calculation results in database
8. Calculate risk interpretation and recommendations
9. Return results with clinical interpretation

---

### Submodule: BCRA Risk Calculation (BRCA Mutations)

#### Inputs
**Endpoint:** `GET /admin/calc/bcra/:formID`

**Required Form Data for BCRA:**
- Age
- Gender (Female required)
- Menopausal status
- Personal breast/ovarian cancer history
- Family history (breast, ovarian cancers)
- Age at menarche (if available)
- Age at first live birth
- Number of breast biopsies
- Hyperplasia status in biopsies

#### Outputs
**Success (200):**
```json
{
  "success": true,
  "data": {
    "formId": "form-uuid",
    "modelType": "BCRA",
    "calculationDate": "2025-11-17T11:35:00Z",
    "results": {
      "brca1MutationRisk": 0.008,
      "brca2MutationRisk": 0.012,
      "anyBrcaMutationRisk": 0.020,
      "lifetimeBreastCancerRisk": 0.18,
      "lifetimeOvarianCancerRisk": 0.04,
      "riskCategory": "Moderate"
    },
    "relativeRisk": {
      "breastCancerVsPopulation": 2.5,
      "ovarianCancerVsPopulation": 3.2
    },
    "interpretation": {
      "recommendation": "BRCA testing recommended",
      "followUpActions": [
        "Genetic counseling",
        "BRCA testing (blood test)",
        "Enhanced breast imaging (MRI) if positive",
        "Increased surveillance"
      ]
    }
  }
}
```

#### Logic
1. Validate form has required data for BRCA calculation
2. Verify gender is Female
3. Extract breast cancer risk factors
4. Send form data to external BCRA model
5. Receive and validate mutation risk probabilities
6. Calculate lifetime cancer risks
7. Calculate relative risks vs. population
8. Generate clinical recommendations
9. Store results and return to user

---

### Submodule: GAIL Risk Calculation (Breast Cancer)

#### Inputs
**Endpoint:** `GET /admin/calc/gail/:formID`

**Required Form Data for GAIL:**
- Age
- Gender (Female required)
- Age at menarche
- Age at first live birth
- Number of breast biopsies
- Hyperplasia status
- Family history of breast cancer
- Race/ethnicity (implicit in model)

#### Outputs
**Success (200):**
```json
{
  "success": true,
  "data": {
    "formId": "form-uuid",
    "modelType": "GAIL",
    "calculationDate": "2025-11-17T11:40:00Z",
    "results": {
      "fiveYearBreastCancerRisk": 0.032,
      "lifetimeBreastCancerRisk": 0.142,
      "riskCategory": "Average",
      "riskScore": 1.8
    },
    "absoluteRisk": {
      "fiveYearAbsoluteRisk": "3.2%",
      "lifetimeAbsoluteRisk": "14.2%"
    },
    "relativeRisk": {
      "comparedToAverageWoman": 1.1
    },
    "interpretation": {
      "riskLevel": "Average Risk",
      "recommendation": "Standard screening guidelines apply",
      "followUpActions": [
        "Annual mammography starting at age 40",
        "Clinical breast examination annually",
        "Self-examination monthly"
      ]
    }
  }
}
```

#### Logic
1. Validate form has all required GAIL inputs
2. Verify gender is Female
3. Extract breast cancer risk factors
4. Send data to GAIL model
5. Receive 5-year and lifetime absolute risks
6. Calculate relative risk vs. average population
7. Classify risk category (Low, Average, Moderate, High)
8. Generate screening recommendations based on risk
9. Store results and return

---

### Submodule: PLCO Risk Calculation (Lung Cancer)

#### Inputs
**Endpoint:** `GET /admin/calc/plco/:formID`

**Required Form Data for PLCO:**
- Age (50-75 years)
- Gender
- Smoking status and pack-years
- Education level
- Occupational exposure (asbestos)
- Family history of lung cancer
- Personal history of cancer
- COPD status

#### Outputs
**Success (200):**
```json
{
  "success": true,
  "data": {
    "formId": "form-uuid",
    "modelType": "PLCO",
    "calculationDate": "2025-11-17T11:45:00Z",
    "results": {
      "threeYearLungCancerRisk": 0.018,
      "sixYearLungCancerRisk": 0.045,
      "riskCategory": "Elevated",
      "screeningEligible": true
    },
    "riskBreakdown": {
      "smokingContribution": 0.038,
      "ageContribution": 0.005,
      "occupationalContribution": 0.002
    },
    "interpretation": {
      "recommendation": "Eligible for low-dose CT screening",
      "followUpActions": [
        "Low-dose CT scan screening annually",
        "Smoking cessation counseling",
        "Occupational exposure reduction",
        "Annual follow-up assessment"
      ]
    }
  }
}
```

**Error (400) - Age Out of Range:**
```json
{
  "success": false,
  "message": "PLCO model applies to ages 50-75 only",
  "error": "VALIDATION_ERROR",
  "details": {
    "userAge": 48,
    "validRange": "50-75"
  }
}
```

#### Logic
1. Validate age is within 50-75 range
2. Extract lung cancer risk factors
3. Calculate pack-years from smoking history
4. Validate occupational exposure data
5. Send data to PLCO model
6. Receive 3-year and 6-year risk projections
7. Determine CT screening eligibility
8. Break down risk contributions
9. Generate personalized recommendations
10. Store results and return

---

### Submodule: Get All Calculation Models

#### Inputs
**Endpoint:** `GET /admin/calc/all-models`

No parameters required.

#### Outputs
**Success (200):**
```json
{
  "success": true,
  "data": [
    {
      "modelType": "PREMM5",
      "name": "Lynch Syndrome Risk",
      "description": "Predicts probability of Lynch syndrome mutations",
      "applicableTo": "Women with personal or family history of colorectal/endometrial cancer",
      "requiredData": [
        "age",
        "gender (Female)",
        "cancer history",
        "family cancer history"
      ],
      "outputMetrics": [
        "MLH1 mutation risk",
        "MSH2 mutation risk",
        "MSH6 mutation risk",
        "PMS2 mutation risk"
      ]
    },
    {
      "modelType": "BCRA",
      "name": "BRCA Mutation Risk",
      "description": "Calculates probability of BRCA1/2 mutations",
      "applicableTo": "Women with breast, ovarian, or pancreatic cancer",
      "requiredData": [
        "age",
        "menopausal status",
        "cancer history",
        "family cancer history"
      ],
      "outputMetrics": [
        "BRCA1 risk",
        "BRCA2 risk",
        "Breast cancer lifetime risk",
        "Ovarian cancer lifetime risk"
      ]
    },
    {
      "modelType": "GAIL",
      "name": "Breast Cancer Risk (GAIL)",
      "description": "Estimates 5-year and lifetime breast cancer risk",
      "applicableTo": "Women aged 35+",
      "requiredData": [
        "age",
        "age at menarche",
        "age at first birth",
        "biopsy history",
        "family history"
      ],
      "outputMetrics": [
        "5-year risk",
        "Lifetime risk",
        "Relative risk"
      ]
    },
    {
      "modelType": "PLCO",
      "name": "Lung Cancer Risk (PLCO)",
      "description": "Predicts 3-year and 6-year lung cancer risk",
      "applicableTo": "Current/former smokers aged 50-75",
      "requiredData": [
        "age (50-75)",
        "smoking status",
        "pack-years",
        "occupational exposure"
      ],
      "outputMetrics": [
        "3-year risk",
        "6-year risk",
        "CT screening eligibility"
      ]
    }
  ]
}
```

#### Logic
1. Return metadata for all available calculation models
2. Include descriptions and requirements for each
3. Help users understand which models apply to their situation

---

## Action Logging

### Module Overview
Comprehensive audit trail system that tracks all significant user actions for compliance and monitoring.

**File References:**
- [action_log_service.go](../internal/application/usecase/action_log_service.go) - Service interface
- [action_log_impl.go](../internal/application/service/action_log_impl.go) - Implementation
- [action_log_admin.go](../internal/presentation/controller/action_log/admin.go) - Admin endpoints

### Submodule: Log Retrieval & Filtering

#### Inputs
**Endpoint:** `GET /admin/log` (Get all action logs)

Query Parameters:
```
page: integer (default: 1)
limit: integer (default: 20, max: 100)
actorId: uuid (optional - filter by who performed action)
targetId: uuid (optional - filter by who was affected)
actionType: string (optional - filter by action type)
resource: string (optional - filter by resource type)
dateFrom: string (optional, YYYY-MM-DD)
dateTo: string (optional, YYYY-MM-DD)
```

**Endpoint:** `GET /admin/log/types` (Get all action types)

#### Outputs
**Success (200) - List Logs:**
```json
{
  "success": true,
  "data": [
    {
      "id": "log-uuid",
      "timestamp": "2025-11-17T10:30:00Z",
      "actor": {
        "id": "user-uuid",
        "phone": "09XXXXXXXXX",
        "role": "ADMIN"
      },
      "target": {
        "id": "user-uuid-2",
        "phone": "09XXXXXXXXX",
        "role": "USER"
      },
      "actionType": "UserRoleAssignment",
      "resource": "User",
      "resourceId": "user-uuid-2",
      "details": {
        "assignedRoles": ["OPERATOR"],
        "previousRoles": ["USER"]
      },
      "ipAddress": "192.168.1.100",
      "userAgent": "Mozilla/5.0...",
      "severity": "INFO"
    },
    {
      "id": "log-uuid-2",
      "timestamp": "2025-11-17T11:15:00Z",
      "actor": {
        "id": "user-uuid",
        "phone": "09XXXXXXXXX",
        "role": "ADMIN"
      },
      "target": null,
      "actionType": "FormAccepted",
      "resource": "Form",
      "resourceId": "form-uuid",
      "details": {
        "formStatus": "Approved",
        "notes": "All information verified"
      },
      "ipAddress": "192.168.1.100",
      "severity": "INFO"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 1250,
    "pages": 63
  }
}
```

**Success (200) - List Action Types:**
```json
{
  "success": true,
  "data": [
    {
      "type": "UserCreated",
      "category": "UserManagement",
      "description": "A new user account was created",
      "severity": "INFO"
    },
    {
      "type": "UserRoleAssignment",
      "category": "UserManagement",
      "description": "User roles were changed",
      "severity": "INFO"
    },
    {
      "type": "PasswordChanged",
      "category": "Security",
      "description": "User password was changed",
      "severity": "WARNING"
    },
    {
      "type": "FormCreated",
      "category": "FormManagement",
      "description": "A new medical form was created",
      "severity": "INFO"
    },
    {
      "type": "FormAccepted",
      "category": "FormManagement",
      "description": "A form was approved by admin",
      "severity": "INFO"
    },
    {
      "type": "FormRejected",
      "category": "FormManagement",
      "description": "A form was rejected by admin",
      "severity": "WARNING"
    },
    {
      "type": "CancerHistoryAdded",
      "category": "FormData",
      "description": "Cancer history record was added to form",
      "severity": "INFO"
    },
    {
      "type": "PermissionDenied",
      "category": "Security",
      "description": "User attempted action without permission",
      "severity": "WARNING"
    },
    {
      "type": "FailedLogin",
      "category": "Security",
      "description": "Failed login attempt",
      "severity": "WARNING"
    },
    {
      "type": "DataExported",
      "category": "Compliance",
      "description": "User data was exported",
      "severity": "WARNING"
    }
  ]
}
```

#### Logic
1. Query action logs with filters applied
2. Load actor and target user information
3. Apply date range filtering
4. Apply action type and resource filtering
5. Paginate results (most recent first)
6. Enrich with metadata (severity, category)
7. Return formatted logs with pagination

---

### Submodule: Automatic Log Creation (Internal)

#### Tracked Events

**User Management Events:**
- User login (success/failure)
- User account creation
- Password change
- Role assignment/removal
- Permission changes

**Form Management Events:**
- Form creation
- Form status changes (Draft → Submitted → Approved/Rejected)
- Form data modifications
- Form deletion

**Data Management Events:**
- Cancer history addition/modification/deletion
- Family cancer history addition/modification/deletion
- Form section updates (Basic Info, Health Info, etc.)

**Security Events:**
- Failed login attempts
- Permission denied attempts
- Suspicious activities
- OTP verification

**Compliance Events:**
- Data exports
- Bulk operations
- Sensitive data access

#### Log Structure

```go
type ActionLog struct {
  ID          string                 // UUID
  Timestamp   time.Time              // When action occurred
  ActorID     string                 // User who performed action
  TargetID    *string                // User affected (if applicable)
  ActionType  ActionType             // Type of action
  Resource    string                 // Resource type (User, Form, etc.)
  ResourceID  string                 // Specific resource ID
  Details     map[string]interface{} // Action-specific details
  IPAddress   string                 // Source IP
  UserAgent   string                 // Browser info
  Severity    LogSeverity            // INFO, WARNING, ERROR
}
```

#### Logic for Log Creation
1. Capture action initiation
2. Extract actor ID from JWT token
3. Determine target user if applicable
4. Collect action-specific details
5. Record IP address and user agent
6. Assign severity level
7. Store in PostgreSQL database
8. Ensure logs are immutable (no updates/deletes)

---

## Enum & Reference Data

### Module Overview
Enumeration types used throughout the system for consistent data validation and categorization.

**File Location:** [enum/](../internal/domain/enum/)

### Enumerations

#### Gender Enum
```
Values: Male, Female, Other
Used in: BasicInfo form
API Response Example:
{
  "type": "Gender",
  "values": ["Male", "Female", "Other"]
}
```

#### Cancer Type Enum
```
Values:
- Breast
- Lung
- Ovarian
- Uterine
- Pancreatic
- Prostate
- Colorectal
- Skin
- Thyroid
- Lymphoma
- Other

Used in: Personal Cancer History, Family Cancer History
```

#### Form Status Enum
```
Values:
- Draft (user still filling)
- Submitted (user submitted, awaiting approval)
- Approved (admin reviewed and approved)
- Rejected (admin reviewed and rejected)
- Archived (completed forms older than 1 year)

State Transitions:
Draft → Submitted (user action)
Submitted → Approved or Rejected (admin action)
Approved/Rejected → Archived (system, after 12 months)
```

#### Menopausal Status Enum
```
Values:
- PreMenopausal (still menstruating)
- Perimenopausal (transitional period, irregular periods)
- PostMenopausal (no menstruation for 12+ months)
- Unknown (status not determined)

Used in: MamoGraphyInfo form
```

#### Life Status Enum
```
Values:
- Alive
- Deceased

Used in: Family Cancer History
```

#### Relative Type Enum
```
Values:
- Mother
- Father
- Sister
- Brother
- Daughter
- Son
- Grandmother
- Grandfather
- Aunt
- Uncle
- Cousin
- Other

Used in: Family Cancer History
```

#### Hyperplasia Type Enum
```
Values:
- None (no hyperplasia found)
- Usual (usual hyperplasia without atypia)
- Atypical (atypical hyperplasia, higher risk)
- Unknown (status not determined)

Used in: MamoGraphyInfo biopsy results
```

#### Permission Type Enum
```
Values:
- CategoryFormManagement (manage all forms)
- PermissionCreateFormForUser (create forms for others)
- CategoryLogManagement (access audit logs)
- CategoryUserManagement (manage users/roles)
- CategoryReportGeneration (generate reports)
- CategoryDataExport (export user data)
- CategorySystemConfiguration (system settings)

Used in: RBAC permission assignment
```

#### Role Type Enum
```
Values:
- ADMIN (full system access)
- OPERATOR (can fill forms for users)
- USER (regular user, can only fill own forms)
- GUEST (read-only access)

Used in: User role assignment
```

#### Action Type Enum
```
Values (Sample):
- UserCreated
- UserLoginSuccess
- UserLoginFailed
- PasswordChanged
- UserRoleAssignment
- RoleCreated
- FormCreated
- FormSubmitted
- FormApproved
- FormRejected
- CancerHistoryAdded
- FormDataModified
- PermissionDenied
- DataExported

Used in: Action Logging
```

#### Bucket Type Enum
```
Values:
- UserProfileImages
- FormDocuments
- CancerHistoryImages
- SystemDocuments

Used in: S3 storage organization
```

#### Calculation Model Type Enum
```
Values:
- PREMM5 (Lynch Syndrome risk)
- BCRA (BRCA mutation risk)
- GAIL (Breast cancer risk)
- PLCO (Lung cancer risk)

Used in: Risk calculation endpoints
```

---

## Data Models

### Core Entities

#### User Entity
```
id: string (UUID)
phone: string (unique, required)
password: string (optional, hashed, for admin users)
roles: Role[] (many-to-many)
createdAt: datetime
updatedAt: datetime
lastLoginAt: datetime (nullable)
status: enum (Active | Suspended | Deleted)
```

**Relationships:**
- Has many Forms (1:M)
- Has many ActionLogs as actor (1:M)
- Has many ActionLogs as target (1:M)
- Has many Roles (M:M)

---

#### Form Entity
```
id: string (UUID)
userId: string (FK to User)
operatorId: string (FK to User, nullable - who filled form)
filledByOperatorId: string (FK to User, nullable - who can fill form on behalf of user)
status: FormStatus
completionPercentage: integer (0-100)
createdAt: datetime
createdBy: string (FK to User)
submittedAt: datetime (nullable)
approvedAt: datetime (nullable)
approvedBy: string (FK to User, nullable)
rejectedAt: datetime (nullable)
rejectedBy: string (FK to User, nullable)
rejectionReason: string (nullable)

Nested Objects (1:1 relationships):
- basicInfo: BasicInfo
- generalHealth: GeneralHealthInfo
- mamography: MamoGraphyInfo
- contact: ContactInfo
- lungCancer: LungCancerInfo

Child Collections (1:M relationships):
- cancers: CancerInfo[]
- familyCancers: FamilyCancerInfo[]
```

**Relationships:**
- Belongs to User (M:1)
- Has one BasicInfo (1:1)
- Has one GeneralHealthInfo (1:1)
- Has one MamoGraphyInfo (1:1)
- Has one ContactInfo (1:1)
- Has one LungCancerInfo (1:1)
- Has many CancerInfo (1:M)
- Has many FamilyCancerInfo (1:M)
- Has many CalculationResult (1:M)

---

#### BasicInfo Entity
```
id: string (UUID)
formId: string (FK to Form, unique)
gender: enum (Male | Female | Other)
birthDate: date
age: integer (calculated, 18+)
socialSecurityNumber: string (optional, encrypted)
height: number (cm)
weight: number (kg)
bmi: number (calculated)
updatedAt: datetime

Calculated Fields:
- age: today - birthDate
- bmi: weight / ((height/100) ^ 2)
```

---

#### GeneralHealthInfo Entity
```
id: string (UUID)
formId: string (FK to Form, unique)

Nested: AlcoholInfo
  drinks: enum (Never | Former | Current)
  frequency: string
  unitsPerWeek: number

Nested: SmokingInfo
  status: enum (Never | Former | Current)
  yearsSmoked: number
  quitYear: number (if Former)
  packYears: number (calculated)

Nested: DietInfo
  vegetableServingsPerDay: number
  fruitServingsPerDay: number
  wholeGrainServingsPerDay: number
  processedMeatServingsPerWeek: number

Nested: ExerciseInfo
  minutesPerWeek: number
  intensity: enum (Light | Moderate | Vigorous)
  frequency: enum (Daily | EveryOtherDay | TwiceAWeek | Once | Never)

updatedAt: datetime
```

---

#### MamoGraphyInfo Entity
```
id: string (UUID)
formId: string (FK to Form, unique)

numberOfChildren: number
ageAtFirstBirth: number (optional)
menopauseStatus: enum (PreMenopausal | Perimenopausal | PostMenopausal | Unknown)
ageAtMenopause: number (optional)
onHormoneReplacementTherapy: boolean
hormoneReplacementTherapyDuration: number (years, optional)

numberOfBiopsies: number
hadHyperplasiaInBiopsy: boolean
hyperplasiaType: enum (Usual | Atypical | None | Unknown)

riskFactorsCount: integer (calculated)
updatedAt: datetime
```

---

#### CancerInfo Entity
```
id: string (UUID)
formId: string (FK to Form)
cancerType: enum (Breast | Lung | Ovarian | ... | Other)
cancerAge: number (18-120)
picturePath: string (optional, S3 key)
notes: string (optional, max 500)
createdAt: datetime
updatedAt: datetime
```

---

#### FamilyCancerInfo Entity
```
id: string (UUID)
formId: string (FK to Form)
relativeType: enum (Mother | Father | Sister | ... | Other)
cancerType: enum (Breast | Lung | Ovarian | ... | Other)
relativeAge: number (optional, age at cancer diagnosis)
relativeDeceased: boolean
deceasedAge: number (optional, age at death)
notes: string (optional)
createdAt: datetime
updatedAt: datetime
```

---

#### ContactInfo Entity
```
id: string (UUID)
formId: string (FK to Form, unique)

fullName: string (required, 3-100 chars)
email: string (optional, valid email)
address: string (optional, max 500)
city: string (optional)
province: string (optional)
country: string (optional)
birthCountry: string (optional)
zipCode: string (optional)

updatedAt: datetime
```

---

#### LungCancerInfo Entity
```
id: string (UUID)
formId: string (FK to Form, unique)

Nested: SmokingHistory
  status: enum (Never | Former | Current)
  yearsSmoked: number
  quitYear: number (if Former)
  packYears: number (calculated)

Nested: OccupationalExposure
  exposed: boolean
  exposureTypes: string[] (asbestos | silica | radon | diesel | other)
  yearsExposed: number

insuranceStatus: enum (Insured | Uninsured | PartiallyInsured)

Nested: Comorbidities
  hasCOPD: boolean
  hasAsthma: boolean
  hasPulmonaryFibrosis: boolean
  previousLungDisease: boolean

riskLevel: string (calculated: Low | Moderate | High | VeryHigh)
updatedAt: datetime
```

---

#### Calculation Results Entities

**Premm5Result:**
```
id: string (UUID)
formId: string (FK to Form)
mlh1MutationRisk: number (0-1, probability)
msh2MutationRisk: number (0-1)
msh6MutationRisk: number (0-1)
pms2MutationRisk: number (0-1)
anyLynchSyndromeRisk: number (0-1)
riskCategory: enum (Low | Moderate | High)
calculatedAt: datetime
```

**BCRAResult:**
```
id: string (UUID)
formId: string (FK to Form)
brca1MutationRisk: number (0-1)
brca2MutationRisk: number (0-1)
anyBrcaMutationRisk: number (0-1)
lifetimeBreastCancerRisk: number (0-1)
lifetimeOvarianCancerRisk: number (0-1)
relativeBreastCancerRisk: number
relativeOvarianCancerRisk: number
riskCategory: enum (Low | Moderate | High)
calculatedAt: datetime
```

**GailResult:**
```
id: string (UUID)
formId: string (FK to Form)
fiveYearRisk: number (0-1)
lifetimeRisk: number (0-1)
relativeRisk: number
riskScore: number
riskCategory: enum (Low | Average | Moderate | High)
calculatedAt: datetime
```

**PLCOResult:**
```
id: string (UUID)
formId: string (FK to Form)
threeYearRisk: number (0-1)
sixYearRisk: number (0-1)
riskCategory: enum (Low | Moderate | High | VeryHigh)
screeningEligible: boolean
riskBreakdown: object {
  smokingContribution: number,
  ageContribution: number,
  occupationalContribution: number
}
calculatedAt: datetime
```

---

#### Role Entity
```
id: string (UUID)
name: string (unique, 3-50 chars)
description: string (optional, max 500)
permissions: Permission[] (many-to-many)
createdAt: datetime
updatedAt: datetime
users: User[] (many-to-many, reference back)
```

---

#### Permission Entity
```
id: string (UUID)
type: enum (CategoryFormManagement | PermissionCreateFormForUser | ... )
description: string
category: string (UserManagement | FormManagement | Security | etc.)
createdAt: datetime
roles: Role[] (many-to-many, reference back)
```

---

#### ActionLog Entity
```
id: string (UUID)
timestamp: datetime (immutable)
actorId: string (FK to User, who performed action)
targetId: string (FK to User, nullable - who was affected)
actionType: enum (UserCreated | FormAccepted | etc.)
resource: string (User | Form | etc.)
resourceId: string (ID of affected resource)
details: object (action-specific details)
ipAddress: string
userAgent: string
severity: enum (INFO | WARNING | ERROR | CRITICAL)
createdAt: datetime

Properties:
- Immutable (no updates/deletes after creation)
- Indexed on: timestamp, actorId, targetId, actionType, resource
```

---

This comprehensive module documentation provides clear input/output specifications, business logic flows, and data model definitions for all major modules in the FamCan Backend system.
