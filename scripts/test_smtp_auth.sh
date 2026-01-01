#!/bin/bash

# SMTP Authentication Test Script
# Tests SMTP authentication on different ports with various scenarios

set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
SMTP_HOST="${SMTP_HOST:-localhost}"
SUBMISSION_PORT="${SUBMISSION_PORT:-587}"
SMTPS_PORT="${SMTPS_PORT:-465}"
RELAY_PORT="${RELAY_PORT:-25}"

# Test users
TEST_USER_1="test@localhost"
TEST_PASS_1="testpass123"
TEST_USER_2="alice@localhost"
TEST_PASS_2="alice123"

echo "======================================"
echo "SMTP Authentication Test Suite"
echo "======================================"
echo ""

# Test function
test_smtp_auth() {
    local port=$1
    local user=$2
    local pass=$3
    local tls_mode=$4
    local test_name=$5

    echo -n "Testing: $test_name ... "

    # Build swaks command
    cmd="swaks --to test@example.com --from $user --server $SMTP_HOST:$port"

    if [ "$tls_mode" = "starttls" ]; then
        cmd="$cmd --tls"
    elif [ "$tls_mode" = "smtps" ]; then
        cmd="$cmd --tls-on-connect"
    fi

    if [ -n "$user" ]; then
        cmd="$cmd --auth PLAIN --auth-user $user --auth-password $pass"
    fi

    cmd="$cmd --quit-after AUTH 2>&1"

    # Run test
    output=$(eval $cmd || true)

    if echo "$output" | grep -q "Authentication succeeded\|AUTH succeeded\|235 "; then
        echo -e "${GREEN}✓ PASS${NC}"
        return 0
    elif echo "$output" | grep -q "530 Authentication required\|530"; then
        if [ -z "$user" ]; then
            echo -e "${GREEN}✓ PASS${NC} (correctly rejected - no auth)"
            return 0
        else
            echo -e "${RED}✗ FAIL${NC} - Auth required but credentials provided"
            echo "$output" | grep -i "authentication\|error\|fail" | head -5
            return 1
        fi
    elif echo "$output" | grep -q "535 Authentication failed\|535"; then
        if [ "$pass" = "wrongpass" ]; then
            echo -e "${GREEN}✓ PASS${NC} (correctly rejected - bad password)"
            return 0
        else
            echo -e "${RED}✗ FAIL${NC} - Valid credentials rejected"
            echo "$output" | grep -i "authentication\|error\|fail" | head -5
            return 1
        fi
    else
        echo -e "${YELLOW}? UNKNOWN${NC}"
        echo "$output" | head -10
        return 1
    fi
}

# Check if swaks is installed
if ! command -v swaks &> /dev/null; then
    echo -e "${RED}Error: swaks is not installed${NC}"
    echo "Install with: sudo apt-get install swaks"
    exit 1
fi

# Check if server is running
if ! nc -z $SMTP_HOST $SUBMISSION_PORT 2>/dev/null; then
    echo -e "${RED}Error: SMTP server not running on $SMTP_HOST:$SUBMISSION_PORT${NC}"
    echo "Start the server with: ./build/gomailserver run"
    exit 1
fi

echo "Server: $SMTP_HOST"
echo "Ports: Submission=$SUBMISSION_PORT, SMTPS=$SMTPS_PORT, Relay=$RELAY_PORT"
echo ""

# Test Suite
PASS_COUNT=0
FAIL_COUNT=0
TOTAL_COUNT=0

run_test() {
    TOTAL_COUNT=$((TOTAL_COUNT + 1))
    if test_smtp_auth "$@"; then
        PASS_COUNT=$((PASS_COUNT + 1))
    else
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi
}

echo "=== Port 587 (Submission with STARTTLS) ==="
run_test $SUBMISSION_PORT "$TEST_USER_1" "$TEST_PASS_1" "starttls" "Valid credentials"
run_test $SUBMISSION_PORT "$TEST_USER_1" "wrongpass" "starttls" "Invalid password"
run_test $SUBMISSION_PORT "$TEST_USER_2" "$TEST_PASS_2" "starttls" "Second user valid"
run_test $SUBMISSION_PORT "" "" "" "No authentication (should require auth)"
echo ""

echo "=== Port 465 (SMTPS - Implicit TLS) ==="
run_test $SMTPS_PORT "$TEST_USER_1" "$TEST_PASS_1" "smtps" "Valid credentials with TLS"
run_test $SMTPS_PORT "$TEST_USER_1" "wrongpass" "smtps" "Invalid password with TLS"
echo ""

echo "=== Port 25 (Relay - No Auth Required) ==="
# Relay port typically doesn't require auth for receiving
swaks --to test@localhost --from external@example.com --server $SMTP_HOST:$RELAY_PORT --quit-after RCPT &>/dev/null && \
    echo -e "Testing: Relay without auth ... ${GREEN}✓ PASS${NC} (relay accepts mail)" || \
    echo -e "Testing: Relay without auth ... ${RED}✗ FAIL${NC}"
echo ""

# Summary
echo "======================================"
echo "Test Results"
echo "======================================"
echo -e "Total:  $TOTAL_COUNT tests"
echo -e "Passed: ${GREEN}$PASS_COUNT${NC}"
echo -e "Failed: ${RED}$FAIL_COUNT${NC}"
echo ""

if [ $FAIL_COUNT -eq 0 ]; then
    echo -e "${GREEN}✓ All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}✗ Some tests failed${NC}"
    exit 1
fi
