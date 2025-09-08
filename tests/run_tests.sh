#!/bin/bash

# A2A Server Test Suite
# This script tests all endpoints of the JSON-RPC A2A chess server

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SERVER_URL="http://localhost:8080"
TIMEOUT=5

echo -e "${BLUE}🧪 A2A Server Test Suite${NC}"
echo "=================================="
echo "Server URL: $SERVER_URL"
echo "Timeout: ${TIMEOUT}s"
echo ""

# Function to check if server is running
check_server() {
    echo -e "${YELLOW}🔍 Checking if server is running...${NC}"
    if curl -s --max-time $TIMEOUT "$SERVER_URL/" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Server is running${NC}"
        return 0
    else
        echo -e "${RED}❌ Server is not running on $SERVER_URL${NC}"
        echo "Please start the server first with: ./jsonrpc_a2a_server -ollama-url http://localhost:11434 -model gpt-oss:20b"
        return 1
    fi
}

# Function to run a test
run_test() {
    local test_name="$1"
    local endpoint="$2"
    local method="$3"
    local data_file="$4"
    
    echo -e "\n${YELLOW}🧪 Running: $test_name${NC}"
    echo "Endpoint: $endpoint"
    echo "Method: $method"
    
    if [ -n "$data_file" ]; then
        echo "Data: $data_file"
        # Get the directory where this script is located
        local script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
        local full_data_path="$script_dir/$data_file"
        if [ ! -f "$full_data_path" ]; then
            echo -e "${RED}❌ Data file not found: $full_data_path${NC}"
            return 1
        fi
    fi
    
    # Run the test
    if [ "$method" = "GET" ]; then
        response=$(curl -s --max-time $TIMEOUT "$SERVER_URL$endpoint")
    else
        response=$(curl -s --max-time $TIMEOUT -X "$method" -H "Content-Type: application/json" -d "@$full_data_path" "$SERVER_URL$endpoint")
    fi
    
    # Check if request was successful
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Request successful${NC}"
        echo "Response:"
        echo "$response" | jq . 2>/dev/null || echo "$response"
    else
        echo -e "${RED}❌ Request failed${NC}"
        return 1
    fi
}

# Function to test basic endpoints
test_basic_endpoints() {
    echo -e "\n${BLUE}📋 Testing Basic Endpoints${NC}"
    echo "--------------------------------"
    
    # Test root endpoint
    run_test "Root Endpoint" "/" "GET"
    
    # Test agent card endpoint
    run_test "Agent Card" "/.well-known/agent.json" "GET"
}

# Function to test A2A protocol endpoints
test_a2a_endpoints() {
    echo -e "\n${BLUE}🤖 Testing A2A Protocol Endpoints${NC}"
    echo "----------------------------------------"
    
    # Test message/send endpoint
    run_test "Message Send" "/a2a" "POST" "test_jsonrpc_a2a.json"
    
    # Test tasks/send endpoint
    run_test "Tasks Send" "/a2a" "POST" "test_tasks_send.json"
}

# Function to test error handling
test_error_handling() {
    echo -e "\n${BLUE}⚠️  Testing Error Handling${NC}"
    echo "----------------------------"
    
    # Test invalid method
    echo -e "\n${YELLOW}🧪 Testing: Invalid Method${NC}"
    echo "Endpoint: /a2a"
    echo "Method: POST"
    echo "Data: Invalid JSON"
    
    response=$(curl -s --max-time $TIMEOUT -X POST -H "Content-Type: application/json" -d '{"invalid": "json"' "$SERVER_URL/a2a")
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Request completed${NC}"
        echo "Response:"
        echo "$response" | jq . 2>/dev/null || echo "$response"
    else
        echo -e "${RED}❌ Request failed${NC}"
    fi
    
    # Test invalid endpoint
    echo -e "\n${YELLOW}🧪 Testing: Invalid Endpoint${NC}"
    echo "Endpoint: /invalid"
    echo "Method: GET"
    
    response=$(curl -s --max-time $TIMEOUT "$SERVER_URL/invalid")
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Request completed${NC}"
        echo "Response:"
        echo "$response" | jq . 2>/dev/null || echo "$response"
    else
        echo -e "${RED}❌ Request failed${NC}"
    fi
}

# Function to run all tests
run_all_tests() {
    echo -e "\n${BLUE}🚀 Starting Test Suite${NC}"
    echo "========================"
    
    # Check if server is running
    if ! check_server; then
        exit 1
    fi
    
    # Run all test categories
    test_basic_endpoints
    test_a2a_endpoints
    test_error_handling
    
    echo -e "\n${GREEN}🎉 All tests completed!${NC}"
    echo "========================"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --help     Show this help message"
    echo "  -b, --basic    Test only basic endpoints"
    echo "  -a, --a2a      Test only A2A protocol endpoints"
    echo "  -e, --errors   Test only error handling"
    echo "  -s, --server   Check server status only"
    echo ""
    echo "Examples:"
    echo "  $0              # Run all tests"
    echo "  $0 -b          # Test basic endpoints only"
    echo "  $0 -a          # Test A2A endpoints only"
    echo "  $0 -s          # Check server status only"
}

# Main execution
case "${1:-}" in
    -h|--help)
        show_usage
        exit 0
        ;;
    -b|--basic)
        check_server && test_basic_endpoints
        ;;
    -a|--a2a)
        check_server && test_a2a_endpoints
        ;;
    -e|--errors)
        check_server && test_error_handling
        ;;
    -s|--server)
        check_server
        ;;
    "")
        run_all_tests
        ;;
    *)
        echo -e "${RED}❌ Unknown option: $1${NC}"
        show_usage
        exit 1
        ;;
esac
