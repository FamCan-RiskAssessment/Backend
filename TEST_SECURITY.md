# Security Features Testing Guide

## Prerequisites
1. PostgreSQL running on `localhost:5433`
2. Redis running on `localhost:6380`
3. Database reset (fresh start)

---

## 1️⃣ Start the Server

```bash
# From Backend directory
go run main.go
```

Expected output:
```
Server running on port 8080
```

---

## 2️⃣ Test Password Hashing

### A. Create a User with Password

```bash
# Set password for a user (assumes user with ID 1 exists)
curl -X POST http://localhost:8080/user/set-password \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "password": "MySecurePassword123!"
  }'
```

### B. Verify Password is Hashed in Database

```bash
# Connect to PostgreSQL
psql -h localhost -p 5433 -U famcan0541 -d famcan

# Check the password field
SELECT id, phone, LEFT(password, 20) as password_preview
FROM users
WHERE id = 1;
```

**Expected Result:**
```
 id |    phone    | password_preview
----+-------------+------------------
  1 | 09111325794 | $2a$12$AbCdEf123...
```
✅ Password should start with `$2a$12$` (bcrypt hash, cost 12)

### C. Test Login with Hashed Password

```bash
# Login with correct password
curl -X POST http://localhost:8080/auth/login-password \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "09111325794",
    "password": "MySecurePassword123!"
  }'
```

**Expected Result:**
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "permissions": [...],
  "roles": [...]
}
```

### D. Test Login with Wrong Password

```bash
# Login with incorrect password
curl -X POST http://localhost:8080/auth/login-password \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "09111325794",
    "password": "WrongPassword"
  }'
```

**Expected Result:** HTTP 401 or 403 error

---

## 3️⃣ Test Rate Limiting

### A. Send 101 Requests from Same IP

```bash
# Bash script to test rate limiting
for i in {1..101}; do
  echo "Request $i:"
  curl -s -o /dev/null -w "HTTP %{http_code}\n" \
    http://localhost:8080/enums
  sleep 0.1
done
```

**Expected Result:**
```
Request 1: HTTP 200
Request 2: HTTP 200
...
Request 100: HTTP 200
Request 101: HTTP 429  ← Rate limit exceeded!
```

### B. Verify Rate Limit Error Message

```bash
# Request when rate limited
curl -X GET http://localhost:8080/enums
```

**Expected Result:**
```json
{
  "error": "بیشتر از حد مجاز درخواست ثبت کرده اید.",
  "message": "Too many requests from this IP"
}
```

### C. Check Redis for Rate Limit Keys

```bash
# Connect to Redis
redis-cli -h localhost -p 6380

# Check rate limit keys
KEYS rate_limit:*

# Check counter for your IP (replace with your IP)
GET rate_limit:127.0.0.1
```

**Expected Result:**
```
"100"  (or current count)
```

### D. Wait 60 Seconds and Retry

```bash
# Wait for rate limit window to expire
sleep 61

# Try again
curl -X GET http://localhost:8080/enums
```

**Expected Result:** HTTP 200 (rate limit reset)

---

## 4️⃣ Test OTP Backdoor Removed

### A. Request OTP

```bash
# Request OTP for phone number
curl -X POST http://localhost:8080/auth/send-otp \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "09123456789"
  }'
```

**Expected Result:**
```json
{
  "message": "OTP sent successfully"
}
```

### B. Try Backdoor OTP (Should Fail)

```bash
# Try the old backdoor OTP "111111"
curl -X POST http://localhost:8080/auth/verify-otp \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "09123456789",
    "otp": "111111"
  }'
```

**Expected Result:** HTTP 400/401 - Invalid OTP error
```json
{
  "error": "کد تایید نامعتبر است"
}
```

✅ Backdoor "111111" should NOT work anymore!

### C. Check Redis for Real OTP

```bash
# Connect to Redis
redis-cli -h localhost -p 6380

# Find the OTP key
KEYS otp:*

# Get the actual OTP (for testing purposes)
GET otp:09123456789
```

Use the real OTP to verify it works.

---

## 5️⃣ Test SSN Encryption

### A. Create Form with SSN

```bash
# Create basic info form (assumes user ID 1 exists and is authenticated)
# Replace YOUR_ACCESS_TOKEN with token from login

curl -X POST http://localhost:8080/form/basic-info \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "user_id": 1,
    "gender": 1,
    "birth_date": "1990-01-01",
    "is_atba": false,
    "social_security_number": "1234567890",
    "height": 175.5,
    "weight": 70.0
  }'
```

**Expected Result:**
```json
{
  "form_id": 1,
  "user_id": 1,
  ...
}
```

### B. Verify SSN is Encrypted in Database

```bash
# Connect to PostgreSQL
psql -h localhost -p 5433 -U famcan0541 -d famcan

# Check the encrypted SSN
SELECT id, form_id, LEFT(social_security_number, 50) as ssn_encrypted
FROM basic_infos;
```

**Expected Result:**
```
 id | form_id |              ssn_encrypted
----+---------+------------------------------------------
  1 |       1 | ymt2wHO9X6q3p+13kTvbS0KEb9bWTOjS...
```

✅ SSN should be base64-encoded ciphertext, NOT plaintext "1234567890"

### C. Retrieve Form and Verify SSN is Decrypted

```bash
# Get basic info (decrypts automatically)
curl -X GET http://localhost:8080/form/1/basic-info \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Expected Result:**
```json
{
  "id": 1,
  "social_security_number": "1234567890",  ← Decrypted!
  "gender": 1,
  "birth_date": "1990-01-01",
  ...
}
```

✅ API should return decrypted SSN to authorized user

---

## 6️⃣ Test Address Encryption

### A. Update Contact Info with Address

```bash
# Update contact info with address
curl -X PUT http://localhost:8080/form/1/contact \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "form_id": 1,
    "user_id": 1,
    "name": "Test User",
    "address": "123 Main Street, Tehran, Iran",
    "postal_code": "1234567890",
    "test_gen": false,
    "fm_test_gen": false
  }'
```

### B. Verify Address is Encrypted in Database

```bash
# Connect to PostgreSQL
psql -h localhost -p 5433 -U famcan0541 -d famcan

# Check the encrypted address
SELECT id, form_id, LEFT(address, 50) as address_encrypted
FROM contact_infos;
```

**Expected Result:**
```
 id | form_id |           address_encrypted
----+---------+----------------------------------------
  1 |       1 | Ab7xY2z3K9mN5qP8r... (base64 ciphertext)
```

✅ Address should be encrypted, NOT "123 Main Street, Tehran, Iran"

### C. Retrieve Contact Info and Verify Address is Decrypted

```bash
# Get contact info
curl -X GET http://localhost:8080/form/1/contact \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Expected Result:**
```json
{
  "id": 1,
  "name": "Test User",
  "address": "123 Main Street, Tehran, Iran",  ← Decrypted!
  "postal_code": "1234567890",
  ...
}
```

---

## 7️⃣ Database Verification Checklist

```sql
-- Connect to PostgreSQL
psql -h localhost -p 5433 -U famcan0541 -d famcan

-- 1. Check password hashing
SELECT id, phone,
       CASE
         WHEN password LIKE '$2a$12$%' THEN 'HASHED ✓'
         ELSE 'PLAINTEXT ✗'
       END as password_status
FROM users
LIMIT 5;

-- 2. Check SSN encryption
SELECT id,
       CASE
         WHEN social_security_number ~ '^[0-9]{10}$' THEN 'PLAINTEXT ✗'
         ELSE 'ENCRYPTED ✓'
       END as ssn_status,
       LENGTH(social_security_number) as ssn_length
FROM basic_infos
LIMIT 5;

-- 3. Check address encryption
SELECT id,
       CASE
         WHEN address LIKE '%Street%' OR address LIKE '%Tehran%' THEN 'PLAINTEXT ✗'
         ELSE 'ENCRYPTED ✓'
       END as address_status,
       LENGTH(address) as address_length
FROM contact_infos
LIMIT 5;
```

**Expected Results:**
- ✅ All passwords: `HASHED ✓`
- ✅ All SSNs: `ENCRYPTED ✓` with length > 100
- ✅ All addresses: `ENCRYPTED ✓` with large length

---

## 8️⃣ Quick Test Script

Save this as `test_security.sh`:

```bash
#!/bin/bash

BASE_URL="http://localhost:8080"
echo "🧪 Testing FamCan Security Features..."
echo ""

# Test 1: Health Check
echo "1️⃣ Testing Server Health..."
curl -s $BASE_URL/enums > /dev/null && echo "✅ Server is running" || echo "❌ Server is down"
echo ""

# Test 2: Rate Limiting
echo "2️⃣ Testing Rate Limiting (sending 101 requests)..."
success_count=0
fail_count=0

for i in {1..101}; do
  response=$(curl -s -o /dev/null -w "%{http_code}" $BASE_URL/enums)
  if [ "$response" = "200" ]; then
    ((success_count++))
  else
    ((fail_count++))
  fi
done

echo "   ✅ Successful requests: $success_count"
echo "   🚫 Rate limited requests: $fail_count"

if [ $fail_count -gt 0 ]; then
  echo "   ✅ Rate limiting is working!"
else
  echo "   ❌ Rate limiting NOT working"
fi
echo ""

# Test 3: OTP Backdoor
echo "3️⃣ Testing OTP Backdoor Removal..."
echo "   Attempting login with backdoor OTP '111111'..."
response=$(curl -s -X POST $BASE_URL/auth/verify-otp \
  -H "Content-Type: application/json" \
  -d '{"phone":"09123456789","otp":"111111"}' \
  -w "%{http_code}")

if [[ $response == *"401"* ]] || [[ $response == *"400"* ]]; then
  echo "   ✅ Backdoor OTP rejected - Security working!"
else
  echo "   ❌ Backdoor OTP accepted - SECURITY ISSUE!"
fi
echo ""

echo "🎉 Security testing complete!"
echo ""
echo "📝 Next steps:"
echo "   1. Check database for encrypted data (see TEST_SECURITY.md section 7)"
echo "   2. Test password hashing with actual user login"
echo "   3. Test field encryption with form creation"
```

Run it:
```bash
chmod +x test_security.sh
./test_security.sh
```

---

## 9️⃣ Environment Variables Check

Verify security settings:

```bash
# Check if encryption key is set
grep ENCRYPTION_KEY .env

# Check if rate limiting is configured
grep RATE_LIMIT .env
```

**Expected:**
```
ENCRYPTION_KEY=ymt2wHO9X6q3p+13kTvbS0KEb9bWTOjS
RATE_LIMIT_PER_MINUTE=100
RATE_LIMIT_WINDOW_SECONDS=60
```

---

## 🎯 Summary Checklist

- [ ] **Password Hashing**: Passwords stored as bcrypt hashes starting with `$2a$12$`
- [ ] **Password Login**: Can login with correct password, fails with wrong password
- [ ] **Rate Limiting**: 101st request returns HTTP 429
- [ ] **OTP Security**: Backdoor "111111" rejected
- [ ] **SSN Encryption**: SSN stored as base64 ciphertext in database
- [ ] **SSN Decryption**: SSN returned as plaintext via API
- [ ] **Address Encryption**: Address stored as base64 ciphertext
- [ ] **Address Decryption**: Address returned as plaintext via API
- [ ] **Server Builds**: `go build ./...` succeeds
- [ ] **Tests Pass**: `go test ./...` succeeds

---

## 🚨 Troubleshooting

**Rate limiting not working?**
- Check Redis is running: `redis-cli -h localhost -p 6380 ping`
- Check Redis keys: `redis-cli -h localhost -p 6380 KEYS rate_limit:*`

**Encryption errors?**
- Verify ENCRYPTION_KEY is exactly 32 characters
- Check .env file is loaded: `echo $ENCRYPTION_KEY`

**Password hashing not working?**
- Check user_service has passwordHasher injected
- Verify Wire generation: `wire ./wire`

**Build fails?**
- Run: `go mod tidy`
- Regenerate wire: `wire ./wire`
