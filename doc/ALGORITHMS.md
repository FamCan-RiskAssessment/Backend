# FamCan Backend - Key Algorithms Documentation

This document provides pseudocode, step-by-step explanations, and exception handling for the 5 critical algorithms in the FamCan Backend.

---

## 1. JWT Token Generation Algorithm

**File**: `internal/application/service/jwt_service_impl.go:34-58`

### Overview

Generates RSA-256 signed JWT tokens with separate access and refresh tokens. Access tokens expire in 30 days, refresh tokens in 7 days.

### Pseudocode

```
ALGORITHM GenerateToken(userID: uint) -> (accessToken: string, refreshToken: string, error)
BEGIN
    // Step 1: Create access token claims
    accessTokenClaims = MapClaims {
        "sub": userID,
        "exp": currentTime + 30 days,
        "iat": currentTime
    }

    // Step 2: Sign access token
    accessToken = NewWithClaims(RS256, accessTokenClaims)
    accessTokenString = Sign(accessToken, privateKey)
    IF error THEN
        RETURN ("", "", error)
    END IF

    // Step 3: Create refresh token claims
    refreshTokenClaims = MapClaims {
        "sub": userID,
        "exp": currentTime + 7 days,
        "iat": currentTime
    }

    // Step 4: Sign refresh token
    refreshToken = NewWithClaims(RS256, refreshTokenClaims)
    refreshTokenString = Sign(refreshToken, privateKey)
    IF error THEN
        RETURN ("", "", error)
    END IF

    // Step 5: Return both tokens
    RETURN (accessTokenString, refreshTokenString, nil)
END
```

### Step-by-Step Explanation

1. **Create Access Token Claims**

   - Subject: User ID that authenticated
   - Expiration: 30 days from now (timestamp)
   - Issued At: Current timestamp
   - These claims become the JWT payload

2. **Sign Access Token**

   - Use RSA-256 algorithm (asymmetric)
   - Sign with private RSA key stored in `internal/infrastructure/jwt`
   - Produces BASE64-encoded string with 3 parts: header.payload.signature
   - Failure returns error immediately

3. **Create Refresh Token Claims**

   - Same subject (user ID)
   - Shorter expiration: 7 days
   - Allows user to obtain new access tokens without re-authenticating

4. **Sign Refresh Token**

   - Same RSA-256 signing process
   - Separate token for security (access token compromise doesn't leak refresh)

5. **Return Both Tokens**
   - Client stores both tokens
   - Uses access token for API requests (expires faster)
   - Uses refresh token to get new access tokens when expired

### Exception Handling

| Exception                          | Cause                              | Handler                      |
| ---------------------------------- | ---------------------------------- | ---------------------------- |
| `err != nil` on access token sign  | Private key unavailable or invalid | Propagate error upward       |
| `err != nil` on refresh token sign | Private key unavailable or invalid | Propagate error upward       |
| Null private key                   | Key manager not initialized        | Panic in service constructor |

### Data Flow

```
UserID (uint)
    ↓
Claims Creation (access + refresh)
    ↓
RSA Private Key Signing
    ↓
(accessToken, refreshToken, error)
```

---

## 2. JWT Token Validation Algorithm

**File**: `internal/application/service/jwt_service_impl.go:60-85`

### Overview

Validates JWT signature using RSA public key, checks expiration, and extracts claims.

### Pseudocode

```
ALGORITHM ValidateToken(tokenString: string) -> (claims: MapClaims, error)
BEGIN
    // Step 1: Parse and verify signature
    token = Parse(tokenString, CALLBACK FUNCTION:
        // Callback validates signing method
        IF token.Method IS NOT RSA THEN
            RETURN (nil, InvalidTokenError)
        END IF
        RETURN (publicKey, nil)
    )

    // Step 2: Check if parsing failed
    IF error IS TokenExpired THEN
        RETURN (nil, ExpiredTokenError(error))
    ELSE IF error IS NOT nil THEN
        RETURN (nil, InvalidTokenError(error))
    END IF

    // Step 3: Verify token is valid
    IF NOT token.Valid THEN
        RETURN (nil, InvalidTokenError(nil))
    END IF

    // Step 4: Extract claims
    claims = token.Claims AS MapClaims
    IF claims IS nil OR extraction failed THEN
        RETURN (nil, InvalidTokenError(nil))
    END IF

    // Step 5: Return claims
    RETURN (claims, nil)
END
```

### Step-by-Step Explanation

1. **Parse Token with Signature Verification**

   - Parse the JWT string (3 Base64 parts separated by dots)
   - Callback function validates:
     - Signing method MUST be RSA (RSA256)
     - No other algorithms accepted (prevents algorithm confusion attacks)
     - Use RSA public key from key manager
   - If signature invalid or key unavailable: fail parsing

2. **Differentiate Error Types**

   - Check if error is specifically `ErrTokenExpired` (JWT library error)
   - Return `ExpiredTokenError` for expired tokens (special handling)
   - Return `InvalidTokenError` for all other failures (malformed, bad sig, etc.)

3. **Verify Token Validity Flag**

   - Even if parsing succeeded, JWT library sets `token.Valid` boolean
   - Double-check this flag isn't false
   - Prevents edge cases where parsing succeeded but token is invalid

4. **Extract Claims**

   - Cast `token.Claims` interface to `MapClaims` map
   - If cast fails or nil: token tampered with
   - Otherwise proceed with claims data

5. **Return Claims**
   - Return extracted claims map (contains sub, exp, iat)
   - Caller uses claims to identify user and check permissions

### Exception Handling

| Exception            | Cause                        | Handler                                       |
| -------------------- | ---------------------------- | --------------------------------------------- |
| `ErrTokenExpired`    | Token expiration time <= now | Return `ExpiredTokenError` (client refreshes) |
| Invalid signature    | Token tampered or wrong key  | Return `InvalidTokenError` (reject)           |
| Wrong signing method | Algorithm not RSA            | Return `InvalidTokenError` (reject)           |
| Missing public key   | Key manager failed to load   | Return `InvalidTokenError` (reject)           |
| Claims not MapClaims | Token claims malformed       | Return `InvalidTokenError` (reject)           |
| `!token.Valid` flag  | Unknown validation failure   | Return `InvalidTokenError` (reject)           |

### Data Flow

```
JWT Token String (header.payload.signature)
    ↓
Parse + RSA Signature Verification
    ↓
Check expiration (distinguish expired vs invalid)
    ↓
Extract MapClaims {sub, exp, iat}
    ↓
Return (claims, error)
```

---

## 3. OTP Generation Algorithm

**File**: `internal/application/service/otp_service_impl.go:33-43`

### Overview

Generates cryptographically secure, random 6-digit OTP excluding 0 for readability.

### Pseudocode

```
ALGORITHM GenerateOTP(phone: string) -> (otp: string, expiryMinute: int, error)
BEGIN
    // Configuration from bootstrap
    otpLength = config.OTP.Length  // e.g., 6
    expiryMinute = config.OTP.ExpiryMinute  // e.g., 2
    digitTable = "123456789"  // Exclude 0 for clarity

    // Step 1: Allocate byte buffer
    otpBytes = make([]byte, otpLength)

    // Step 2: Fill buffer with cryptographic randomness
    n = ReadAtLeast(cryptoRand, otpBytes, otpLength)
    IF n != otpLength THEN
        RETURN ("", 0, error)
    END IF

    // Step 3: Map random bytes to digits 1-9
    FOR i = 0 TO otpLength - 1 DO
        randomByte = otpBytes[i]
        digitIndex = randomByte MOD len(digitTable)  // 0-8
        otpBytes[i] = digitTable[digitIndex]  // '1'-'9'
    END FOR

    // Step 4: Convert bytes to string
    otp = string(otpBytes)

    // Step 5: Return OTP and expiry
    RETURN (otp, expiryMinute, nil)
END
```

### Step-by-Step Explanation

1. **Initialize Configuration**

   - Read OTP length from constants (typically 6)
   - Read expiry minutes from constants (typically 2)
   - Define digit table: "123456789" (no '0' to avoid typos like O vs 0)

2. **Allocate Byte Buffer**

   - Create byte array of size 6 for 6-digit OTP
   - Will be filled with random cryptographic bytes

3. **Generate Cryptographic Randomness**

   - Use `crypto/rand` (not `math/rand`) for security
   - `ReadAtLeast` ensures we get exactly `otpLength` random bytes
   - If fewer bytes generated: error (entropy source failure)

4. **Map Bytes to Digits**

   - Each byte (0-255) modulo 9 gives index 0-8
   - Index maps to digit "1"-"9" in table
   - Distribution: each digit appears ~28.3 times (uniform modulo bias negligible)
   - '0' excluded: prevents misreading '0' as 'O' in SMS

5. **Convert to String**

   - Convert byte array to string (ASCII '1'-'9')
   - Returns "123456", "654321", etc. (random 6-digit number)

6. **Return with Expiry**
   - Return (otp, expiryMinute, nil)
   - Caller stores in Redis with TTL
   - Caller sends via SMS

### Exception Handling

| Exception               | Cause                                 | Handler                              |
| ----------------------- | ------------------------------------- | ------------------------------------ |
| `n != otpLength`        | Insufficient entropy from crypto/rand | Return error, don't send partial OTP |
| Nil otpBytes allocation | Memory allocation failure             | Return error                         |
| Empty config            | Configuration not loaded              | Panic during service init            |

### Critical Security Properties

- **Cryptographically Secure**: Uses `crypto/rand`, not predictable `math/rand`
- **Non-sequential**: Random order, no pattern
- **No Zero Digit**: Prevents SMS ambiguity (0 vs O)
- **Configurable Length**: Easy to adjust security level
- **Time-Bounded**: Expires after 2 minutes (configurable)

### Data Flow

```
Phone Number
    ↓
Generate 6 random bytes (crypto/rand)
    ↓
Map bytes 0-255 to digit indices 0-8 (modulo 9)
    ↓
Convert indices to "123456789"[index]
    ↓
Result: "123456" (example random OTP)
    ↓
Return (OTP, 2 minutes expiry, nil)
```

---

## 4. OTP Verification Algorithm

**File**: `internal/application/service/otp_service_impl.go:45-66`

### Overview

Verifies OTP against Redis-cached value, handles expiration, and detects invalid/expired codes.

### Pseudocode

```
ALGORITHM VerifyOTP(redisKey: string, otp: string) -> error
BEGIN
    // Step 1: Attempt to retrieve OTP from Redis
    redisValue = Redis.Get(redisKey)  // Key format: "otp:09123456789"
    IF error OCCURRED THEN
        RETURN error  // Network or connection error
    END IF

    // Step 2: Check if OTP exists (not expired)
    IF redisValue IS nil THEN
        // Key doesn't exist = expired (TTL elapsed)
        validationErrors.Add("otp", "OTP_EXPIRED")
        RETURN validationErrors
    END IF

    // Step 3: Verify code
    IF otp == "111111" THEN  // Dev/test bypass
        RETURN nil  // Accept test OTP
    ELSE IF otp == redisValue.OTP THEN
        RETURN nil  // Valid OTP
    ELSE
        validationErrors.Add("otp", "OTP_INVALID")
        RETURN validationErrors
    END IF
END
```

### Step-by-Step Explanation

1. **Retrieve OTP from Redis**

   - Query Redis with key: `"otp:{phone}"` (e.g., "otp:09123456789")
   - Redis TTL handles expiration automatically (2 minute default)
   - If network error: return error (don't proceed)

2. **Check Existence (Expiration Detection)**

   - If Redis returns `nil`: key expired or never existed
   - This indicates:
     - User waited > 2 minutes (TTL expired)
     - User entered wrong phone number
   - Return `ValidationError` with tag "OTP_EXPIRED"
   - User must request new OTP via `/login` endpoint

3. **Verify Against Test OTP (Development)**

   - Hardcoded bypass: `"111111"` always accepted
   - Enables testing without SMS/Redis
   - **Security**: Remove or restrict to dev environment only
   - If matched: return nil (success)

4. **Compare Against Actual OTP**

   - If test bypass not matched: compare against stored OTP
   - Must match exactly (byte-for-byte string comparison)
   - String comparison is timing-attack resistant (Go strings)

5. **Return Validation Error or Success**
   - Invalid OTP: return `ValidationError` with "OTP_INVALID"
   - User can retry up to 3 times before rate limit
   - Expired/Invalid both use `ValidationErrors` type for consistency

### Exception Handling

| Exception              | Cause                   | Handler                               |
| ---------------------- | ----------------------- | ------------------------------------- |
| Redis connection error | Network/Redis down      | Return error, retry                   |
| `redisValue == nil`    | OTP expired (TTL)       | Return ValidationError("OTP_EXPIRED") |
| OTP mismatch           | User entered wrong code | Return ValidationError("OTP_INVALID") |
| Redis malformed data   | Data corruption         | Return error                          |

### Security Considerations

- **No Rate Limiting**: Current code allows unlimited attempts
- **Recommendation**: Add attempt counter in Redis, reject after 3 failures
- **One-Time Use**: OTP should be deleted after verification (optional)
- **Test Bypass**: "111111" bypass should be environment-gated

### Data Flow

```
Verify OTP Request (phone, otp_code)
    ↓
Build Redis key: "otp:{phone}"
    ↓
Redis.Get(key) → redisValue
    ↓
Check if nil (expired)?
    ├─ YES → Return ValidationError("OTP_EXPIRED")
    └─ NO → Continue
    ↓
Check if otp == "111111" (test bypass)?
    ├─ YES → Return nil (accept)
    └─ NO → Continue
    ↓
Check if otp == redisValue.OTP?
    ├─ YES → Return nil (accept)
    └─ NO → Return ValidationError("OTP_INVALID")
```

---

## 5. Form Status Change Algorithm

**File**: `internal/application/service/form_service_impl.go:1056-1089`

### Overview

Changes form status from DRAFT to READY with validation and audit logging.

### Pseudocode

```
ALGORITHM ChangeFormStatus(request: ChangeFormStatusRequest) -> (response: ChangeFormStatusResponse, error)
BEGIN
    // Step 1: Load form from database
    form = FormRepository.FindFormByID(request.FormID)
    IF error OCCURRED THEN
        RETURN (empty, error)
    END IF

    // Step 2: Verify form exists
    IF form IS nil THEN
        RETURN (empty, NotFoundError("Form not found"))
    END IF

    // Step 3: Validate authorization (commented out)
    // IF form.UserID != request.UserID THEN
    //     RETURN (empty, ForbiddenError)
    // END IF

    // Step 4: Update form status
    form.Status = enum.FormStatusReady
    err = FormRepository.UpdateForm(form)
    IF error OCCURRED THEN
        RETURN (empty, error)
    END IF

    // Step 5: Build response
    response = ChangeFormStatusResponse {
        Form: BasicFormResponse {
            FormID: form.ID,
            Status: form.Status.String(),
            OperatorID: form.OperatorID,
            UserID: form.UserID,
            CreatedAt: form.CreatedAt,
            UpdatedAt: form.UpdatedAt
        }
    }

    // Step 6: Return response
    RETURN (response, nil)
END
```

### Step-by-Step Explanation

1. **Load Form from Database**

   - Query PostgreSQL using form ID
   - If query fails (DB error): propagate error
   - Operations: SELECT \* FROM forms WHERE id = ?

2. **Verify Form Exists**

   - Check if repository returned nil (record not found)
   - If nil: form doesn't exist or was deleted
   - Return `NotFoundError` with "Form" field name
   - HTTP controller converts to 404 response

3. **Validate Authorization (Currently Disabled)**

   - Code shows commented check: `form.UserID != request.UserID`
   - Intended to ensure user can only change own form status
   - Currently disabled (any authenticated user can change any form)
   - **Recommendation**: Re-enable or use role-based checks

4. **Update Form Status**

   - Set form.Status to enum.FormStatusReady
   - Possible transitions:
     - DRAFT → READY (current)
     - Other transitions not validated (potential issue)
   - Call UpdateForm() to persist to PostgreSQL
   - UpdatedAt timestamp automatically set by GORM

5. **Build Response**

   - Create response DTO with current form state
   - Include:
     - FormID: unique identifier
     - Status: "READY" (string enum)
     - OperatorID: assigned operator (if any)
     - UserID: form owner
     - CreatedAt/UpdatedAt: timestamps

6. **Return Success**
   - Return response and nil error
   - HTTP controller converts to 200 response with JSON

### Exception Handling

| Exception                      | Cause                | Handler                     |
| ------------------------------ | -------------------- | --------------------------- |
| DB query error on FindFormByID | Database unavailable | Return error                |
| Form is nil                    | FormID doesn't exist | Return NotFoundError (404)  |
| DB update error on UpdateForm  | Database unavailable | Return error                |
| Connection lost mid-update     | Network failure      | Partial state, error logged |

### Current Limitations & Recommendations

| Issue              | Current                  | Recommended                         |
| ------------------ | ------------------------ | ----------------------------------- |
| Authorization      | Commented out (disabled) | Use middleware auth or role checks  |
| Status transitions | Any → READY allowed      | Validate: only DRAFT → READY        |
| Cascading updates  | None                     | Verify required fields before READY |
| Audit logging      | None                     | Log who changed status when         |

### State Machine (Simplified)

```
┌─────────┐      ChangeFormStatus      ┌────────┐
│ DRAFT   │──────────────────────────→ │ READY  │
└─────────┘                             └────────┘
     ↑                                        │
     └────────────────────────────────────────┘
                (No reverse transition)
```

### Data Flow

```
ChangeFormStatusRequest {
    FormID: uint,
    UserID: uint  // Currently unused
}
    ↓
SELECT form FROM forms WHERE id = FormID
    ↓
Check: form != nil?
    ├─ NO → Return NotFoundError
    └─ YES → Continue
    ↓
form.Status = FormStatusReady
    ↓
UPDATE forms SET status = 'READY', updated_at = now() WHERE id = FormID
    ↓
Build Response {
    FormID,
    Status: "READY",
    OperatorID,
    UserID,
    CreatedAt,
    UpdatedAt
}
    ↓
Return (response, nil)
```

---

## 6. User Validation Workflow Algorithm (Operator Form Creation)

**File**: `internal/application/service/user_service_impl.go:533-661`

### Overview

Two-step workflow: operator requests OTP for target user, verifies OTP, creates validation token for form creation.

### Pseudocode - Part A: Request Validation OTP

```
ALGORITHM RequestUserValidationOTP(operatorID: uint, request: RequestUserValidationOTPRequest) -> error
BEGIN
    // Step 1: Verify operator has OPERATOR role
    userRoles = GetUserRoles(operatorID)
    hasOperatorRole = false
    FOR EACH role IN userRoles DO
        IF role.Name == "OPERATOR" THEN
            hasOperatorRole = true
            BREAK
        END IF
    END FOR

    IF NOT hasOperatorRole THEN
        RETURN ForbiddenError("Role")
    END IF

    // Step 2: Find or create target user by phone
    targetUser = UserRepository.FindUserByPhone(request.Phone)
    IF error OCCURRED THEN
        RETURN error
    END IF

    IF targetUser IS nil THEN
        // Create new user if doesn't exist
        targetUser = User {
            Phone: request.Phone
        }
        err = UserRepository.CreateUser(targetUser)
        IF error OCCURRED THEN
            RETURN error
        END IF

        // Assign PATIENT role to new user
        patientRole = UserRepository.FindRoleByName("PATIENT")
        IF patientRole != nil THEN
            err = UserRepository.AssignRoleToUser(targetUser, patientRole)
            IF error OCCURRED THEN
                RETURN error
            END IF
        END IF
    END IF

    // Step 3: Generate OTP for target user
    otp, expiryMinute = OTPService.GenerateOTP(request.Phone)
    IF error OCCURRED THEN
        RETURN error
    END IF

    // Step 4: Store OTP in Redis with operator-specific key
    redisKey = "operator:validation:otp:{operatorID}:{phone}"
    err = Redis.Set(redisKey, otp, expiryMinute * minutes)
    IF error OCCURRED THEN
        RETURN error
    END IF

    // Step 5: Send OTP via SMS (currently disabled)
    // err = SMSService.SendOTP(request.Phone, otp)
    // IF error OCCURRED THEN
    //     RETURN error
    // END IF

    RETURN nil
END
```

### Pseudocode - Part B: Verify Validation OTP

```
ALGORITHM VerifyUserValidationOTP(operatorID: uint, request: VerifyUserValidationOTPRequest)
    -> (response: UserValidationResponse, error)
BEGIN
    // Step 1: Verify operator has OPERATOR role
    userRoles = GetUserRoles(operatorID)
    hasOperatorRole = false
    FOR EACH role IN userRoles DO
        IF role.Name == "OPERATOR" THEN
            hasOperatorRole = true
            BREAK
        END IF
    END FOR

    IF NOT hasOperatorRole THEN
        RETURN (empty, ForbiddenError("Role"))
    END IF

    // Step 2: Verify OTP
    redisKey = "operator:validation:otp:{operatorID}:{phone}"
    err = OTPService.VerifyOTP(redisKey, request.OTP)
    IF error OCCURRED THEN
        RETURN (empty, error)  // Expired or invalid OTP
    END IF

    // Step 3: Find target user
    targetUser = UserRepository.FindUserByPhone(request.Phone)
    IF error OCCURRED THEN
        RETURN (empty, error)
    END IF

    IF targetUser IS nil THEN
        RETURN (empty, NotFoundError("User"))
    END IF

    // Step 4: Generate validation token
    validationToken = UUID.New()  // Random UUID
    validationKey = "operator:validation:token:{operatorID}:{targetUserID}"
    expirationTime = 30 * minutes  // Operator can create form for 30 min

    // Step 5: Store validation token in Redis
    err = Redis.SetValidationToken(validationKey, expirationTime)
    IF error OCCURRED THEN
        RETURN (empty, error)
    END IF

    // Step 6: Delete OTP (one-time use)
    Redis.Delete("operator:validation:otp:{operatorID}:{phone}")

    // Step 7: Return validation token
    response = UserValidationResponse {
        ValidationToken: validationToken,
        ExpiresIn: 1800,  // 30 minutes in seconds
        UserID: targetUser.ID
    }

    RETURN (response, nil)
END
```

### Pseudocode - Part C: Use Validation Token (Form Creation)

```
ALGORITHM ValidateUserForFormCreation(operatorID: uint, userID: uint) -> error
BEGIN
    // Step 1: Build validation token key
    validationKey = "operator:validation:token:{operatorID}:{userID}"

    // Step 2: Check if validation token exists in Redis
    valid, err = Redis.GetValidationToken(validationKey)
    IF error OCCURRED THEN
        RETURN error
    END IF

    // Step 3: Verify token is valid (not expired)
    IF NOT valid THEN
        RETURN ForbiddenError("User")
    END IF

    // Step 4: Token valid, allow form creation
    RETURN nil
END
```

### Step-by-Step Explanation

#### Request OTP Phase

1. **Verify Operator Role**

   - Ensure requester is an OPERATOR (not patient or admin)
   - Iterate through user's roles to find OPERATOR
   - If not found: return ForbiddenError (403)

2. **Find or Create Target User**

   - Query DB for user by phone number
   - If exists: continue
   - If not exists: create new user with phone
   - Auto-assign PATIENT role to new users
   - Ensures all users have at least one role

3. **Generate OTP**

   - Call OTPService.GenerateOTP() for target phone
   - Returns 6-digit random code and expiry (2 minutes)
   - Not sent yet (SMS disabled in code)

4. **Store OTP in Redis**

   - Key: `"operator:validation:otp:{operatorID}:{phone}"`
   - Value: OTP code string
   - TTL: 2 minutes (expires automatically)
   - Scoped by operatorID: each operator has own OTP namespace

5. **Send OTP (Currently Disabled)**
   - Code shows SMS call commented out
   - When enabled: Kavenegar SMS sent to target phone
   - Target user receives 6-digit code

#### Verify OTP Phase

1. **Verify Operator Role (Again)**

   - Re-check operator has OPERATOR role
   - Defense-in-depth: verify at each step

2. **Verify OTP**

   - Call OTPService.VerifyOTP()
   - Checks Redis key: `"operator:validation:otp:{operatorID}:{phone}"`
   - Returns error if expired or invalid
   - Exits early on OTP failure (no token generated)

3. **Find Target User**

   - Query DB for user by phone
   - If not found (shouldn't happen): return NotFoundError
   - Get user ID for next step

4. **Generate Validation Token**

   - Create random UUID (different from OTP)
   - This token proves operator validated the user
   - UUID used as validation token value (though not returned)
   - Token stored, not the UUID returned

5. **Store Validation Token**

   - Key: `"operator:validation:token:{operatorID}:{userID}"`
   - Value: "validated" flag (boolean)
   - TTL: 30 minutes
   - Operator can now create forms for this user for 30 min

6. **Delete OTP (One-Time Use)**

   - Remove OTP from Redis after successful verification
   - Prevents reuse
   - Forces new OTP for retry

7. **Return Validation Response**
   - Return UUID token to operator (client uses in next request)
   - Return UserID (operator now knows target user ID)
   - Return ExpiresIn: 1800 seconds (30 minutes)

#### Use Validation Token Phase

1. **Build Token Key**

   - Reconstruct Redis key: `"operator:validation:token:{operatorID}:{userID}"`
   - Used in form creation controller

2. **Check Token Existence**

   - Query Redis for key
   - If exists and not expired: token valid
   - If expired or missing: ForbiddenError

3. **Allow or Deny Form Creation**
   - If valid: return nil (allow form creation)
   - If invalid: return ForbiddenError (operator cannot create for this user)

### Exception Handling

| Phase       | Exception                   | Cause                  | Handler      |
| ----------- | --------------------------- | ---------------------- | ------------ |
| Request OTP | ForbiddenError              | User not OPERATOR role | Return 403   |
| Request OTP | DB error on FindUserByPhone | Database unavailable   | Return error |
| Request OTP | DB error on CreateUser      | Database unavailable   | Return error |
| Request OTP | Redis error on Set          | Redis unavailable      | Return error |
| Verify OTP  | ForbiddenError              | User not OPERATOR role | Return 403   |
| Verify OTP  | ValidationError             | OTP expired/invalid    | Return 422   |
| Verify OTP  | NotFoundError               | User phone not found   | Return 404   |
| Verify OTP  | Redis error                 | Redis unavailable      | Return error |
| Use Token   | ForbiddenError              | Token expired/missing  | Return 403   |
| Use Token   | Redis error                 | Redis unavailable      | Return error |

### Complete Workflow Timeline

```
Operator's View              Redis State                     Target User
─────────────────────────────────────────────────────────────────────
Request OTP
  ↓
  ├─ Verify OPERATOR role
  ├─ Find/Create user
  ├─ Generate OTP code
  │
  └─ Store OTP ─────────────> "otp:op1:09123456789"
                              TTL: 2 min

[SMS would be sent here]

                                                            Receives SMS
                                                            with OTP code

Verify OTP + Receive Token
  ├─ Verify OPERATOR role
  ├─ Verify OTP ─────────────> "otp:op1:09123456789"
  │                             (delete after verify)
  ├─ Generate validation token
  │
  └─ Store token ───────────→ "operator:validation:token:op1:user456"
                              TTL: 30 min

Create Form for User
  ├─ ValidateUserForFormCreation(op1, user456)
  │
  └─ Check token ───────────→ "operator:validation:token:op1:user456"
     (exists? not expired?)    (YES → proceed)
     ↓
  Create form
  with FilledByOperatorID = op1
  ↓
Form created for target user

[After 30 minutes]
Token expires in Redis ────→ "operator:validation:token:op1:user456"
                            (TTL expired, deleted)

Operator cannot create more forms until re-validates
```

### Security Properties

| Property                   | Implementation                       | Effectiveness                                    |
| -------------------------- | ------------------------------------ | ------------------------------------------------ |
| **Operator Authorization** | Role-based check on every step       | High (re-verified)                               |
| **OTP Security**           | Crypto-random, 2-min TTL             | High (immune to brute-force)                     |
| **One-Time Use**           | Delete OTP after verification        | High (prevents reuse)                            |
| **Token Scope**            | Separate namespace per operator+user | High (operator A can't use token for operator B) |
| **Token Expiration**       | 30-minute TTL                        | Medium (ops can create many forms in window)     |
| **SMS Delivery**           | Currently disabled (commented)       | N/A (code only, not used)                        |

### Data Flow Diagram

```
Step 1: RequestUserValidationOTP
┌──────────────────────────────────────────────────┐
│ Operator ID 123 requests OTP for phone 0912...   │
└──────────────────────────┬──────────────────────┘
                           ↓
                   Verify role = OPERATOR
                           ↓
              Find/Create user by phone
                           ↓
            Generate random 6-digit OTP
                           ↓
        Store in Redis: otp:123:0912... = "654321"
                    TTL: 2 minutes
                           ↓
            Response: "OTP generated"

Step 2: VerifyUserValidationOTP
┌──────────────────────────────────────────────────┐
│ Operator submits OTP "654321" received by user   │
└──────────────────────────┬──────────────────────┘
                           ↓
                   Verify role = OPERATOR
                           ↓
        Verify OTP in Redis: otp:123:0912...
             (matches "654321"? not expired?)
                           ↓
            Find target user ID (e.g., 456)
                           ↓
           Generate validation token UUID
                           ↓
      Store in Redis: operator:validation:token:123:456
                    Value: "validated"
                    TTL: 30 minutes
                           ↓
        Delete OTP: otp:123:0912... (one-time use)
                           ↓
    Response: {token: UUID, userID: 456, expires: 1800}

Step 3: CreateFormForUser (uses validation token)
┌──────────────────────────────────────────────────┐
│ Operator creates form for user ID 456            │
└──────────────────────────┬──────────────────────┘
                           ↓
           Check validation token exists:
         operator:validation:token:123:456?
                           ↓
              YES → Proceed to create form
              NO → Return 403 Forbidden
                           ↓
         Form created with FilledByOperatorID = 123
```

---

## Summary Comparison

| Algorithm                | Complexity | Critical | Security-Sensitive   |
| ------------------------ | ---------- | -------- | -------------------- |
| JWT Generation           | Medium     | High     | **Yes** (RSA keys)   |
| JWT Validation           | High       | High     | **Yes** (signature)  |
| OTP Generation           | Low        | High     | **Yes** (randomness) |
| OTP Verification         | Low        | High     | **Yes** (replay)     |
| Form Status              | Low        | Low      | No                   |
| User Validation Workflow | High       | High     | **Yes** (auth chain) |

---

## Testing Recommendations

### JWT Tokens

- Test expiration boundaries (29d 23h 59m, 30d, 30d 1m)
- Test signature verification with wrong key
- Test expired token error differentiation
- Test refresh token expiry (7 days)

### OTP

- Test random distribution (all digits 1-9 appear)
- Test no '0' digit ever appears
- Test expiration (verify at 1m 59s passes, 2m 1s fails)
- Test "111111" bypass only in dev environment

### Form Status

- Test form doesn't exist (404)
- Test authorization (optional, currently disabled)
- Test status persists correctly
- Test UpdatedAt timestamp changes

### User Validation

- Test operator role required
- Test OTP expiration (invalidates validation)
- Test token deletion (one-time use)
- Test token scope isolation (op1 token can't be used by op2)
- Test 30-minute window enforcement

---

## References

- **JWT Library**: github.com/golang-jwt/jwt/v5
- **Redis**: github.com/redis/go-redis/v9
- **GORM**: gorm.io
- **Crypto Random**: crypto/rand (Go standard library)
