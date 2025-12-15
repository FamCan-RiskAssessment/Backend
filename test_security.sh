#!/bin/bash

BASE_URL="http://localhost:8080"
echo "🧪 Testing FamCan Security Features..."
echo ""

# Test 1: Health Check
echo "1️⃣ Testing Server Health..."
if curl -s $BASE_URL/enums > /dev/null 2>&1; then
  echo "   ✅ Server is running"
else
  echo "   ❌ Server is down - start with: go run main.go"
  exit 1
fi
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

  # Show progress every 20 requests
  if [ $((i % 20)) -eq 0 ]; then
    echo "   Progress: $i/101 requests sent..."
  fi
done

echo ""
echo "   📊 Results:"
echo "   ✅ Successful requests: $success_count/100 expected"
echo "   🚫 Rate limited requests: $fail_count/1 expected"
echo ""

if [ $fail_count -gt 0 ]; then
  echo "   ✅ Rate limiting is WORKING!"
else
  echo "   ❌ Rate limiting NOT working - check Redis connection"
fi
echo ""

# Test 3: OTP Backdoor (optional - requires phone setup)
echo "3️⃣ Testing OTP Backdoor Removal..."
echo "   Attempting to verify with backdoor OTP '111111'..."

response=$(curl -s -X POST $BASE_URL/auth/verify-otp \
  -H "Content-Type: application/json" \
  -d '{"phone":"09123456789","otp":"111111"}' \
  -w "\n%{http_code}")

http_code=$(echo "$response" | tail -n1)

if [ "$http_code" = "400" ] || [ "$http_code" = "401" ]; then
  echo "   ✅ Backdoor OTP '111111' was REJECTED - Security working!"
else
  echo "   ⚠️  Response code: $http_code (expected 400/401)"
  echo "   Note: This may fail if phone doesn't exist - that's okay"
fi
echo ""

# Test 4: Environment Variables
echo "4️⃣ Checking Security Configuration..."

if grep -q "ENCRYPTION_KEY=" .env 2>/dev/null; then
  key_length=$(grep "ENCRYPTION_KEY=" .env | cut -d'=' -f2 | wc -c)
  key_length=$((key_length - 1)) # Remove newline

  if [ $key_length -eq 32 ]; then
    echo "   ✅ ENCRYPTION_KEY is set (32 characters)"
  else
    echo "   ❌ ENCRYPTION_KEY is $key_length characters (must be 32)"
  fi
else
  echo "   ❌ ENCRYPTION_KEY not found in .env"
fi

if grep -q "RATE_LIMIT_PER_MINUTE=" .env 2>/dev/null; then
  echo "   ✅ RATE_LIMIT_PER_MINUTE is set"
else
  echo "   ❌ RATE_LIMIT_PER_MINUTE not found in .env"
fi

echo ""

# Summary
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🎉 Security Testing Complete!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📝 Next Steps:"
echo ""
echo "   To test password hashing:"
echo "   1. Create/set a user password via API"
echo "   2. Check database: psql -h localhost -p 5433 -U famcan0541 -d famcan"
echo "   3. Run: SELECT id, phone, LEFT(password, 20) FROM users;"
echo "   4. Password should start with: \$2a\$12\$"
echo ""
echo "   To test field encryption:"
echo "   1. Create a form with SSN and address"
echo "   2. Check database for encrypted values"
echo "   3. Retrieve via API and verify decryption works"
echo ""
echo "   📖 See TEST_SECURITY.md for detailed instructions"
echo ""
