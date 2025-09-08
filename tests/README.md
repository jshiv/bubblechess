# A2A Server Test Suite

This folder contains tests for the JSON-RPC A2A chess server.

## Test Files

- `test_jsonrpc_a2a.json` - Test request for the `message/send` method
- `test_tasks_send.json` - Test request for the `tasks/send` method
- `run_tests.sh` - Main test runner script

## Prerequisites

1. **Server Running**: The A2A server must be running on `http://localhost:8080`
2. **jq**: JSON processor for pretty-printing responses (install with `brew install jq` on macOS)
3. **curl**: HTTP client (usually pre-installed)

## Quick Start

### 1. Start the A2A Server

```bash
# From the project root
./jsonrpc_a2a_server -ollama-url http://localhost:11434 -model gpt-oss:20b
```

### 2. Run All Tests

```bash
# From the project root
./tests/run_tests.sh
```

## Test Script Usage

### Run All Tests
```bash
./tests/run_tests.sh
```

### Run Specific Test Categories
```bash
./tests/run_tests.sh -b          # Basic endpoints only
./tests/run_tests.sh -a          # A2A protocol endpoints only
./tests/run_tests.sh -e          # Error handling only
./tests/run_tests.sh -s          # Check server status only
```

### Get Help
```bash
./tests/run_tests.sh -h
```

## Test Categories

### 1. Basic Endpoints
- **Root** (`/`) - Server information
- **Agent Card** (`/.well-known/agent.json`) - A2A agent discovery

### 2. A2A Protocol Endpoints
- **Message Send** (`/a2a` with `message/send` method) - Send chess move request
- **Tasks Send** (`/a2a` with `tasks/send` method) - Send chess task request

### 3. Error Handling
- **Invalid JSON** - Test malformed request handling
- **Invalid Endpoint** - Test 404 handling

## Test Data

### Message Send Request
```json
{
  "jsonrpc": "2.0",
  "method": "message/send",
  "params": {
    "message": {
      "kind": "message",
      "role": "user",
      "messageId": "msg_1",
      "parts": [
        {
          "kind": "text",
          "text": "{\"board_state\":\"...\",\"player_color\":\"black\",\"game_history\":[\"e2e4\",\"e7e5\"]}"
        }
      ]
    }
  },
  "id": 1
}
```

### Tasks Send Request
```json
{
  "jsonrpc": "2.0",
  "method": "tasks/send",
  "params": {
    "message": {
      "kind": "message",
      "role": "user",
      "messageId": "task_1",
      "parts": [
        {
          "kind": "text",
          "text": "{\"board_state\":\"...\",\"player_color\":\"white\",\"game_history\":[\"e2e4\"]}"
        }
      ]
    }
  },
  "id": 2
}
```

## Expected Responses

### Successful A2A Response
```json
{
  "id": 1,
  "jsonrpc": "2.0",
  "result": {
    "kind": "message",
    "messageId": "msg_1234567890",
    "role": "agent",
    "parts": [
      {
        "kind": "text",
        "text": "Generated move: Nf6"
      }
    ]
  }
}
```

### Error Response
```json
{
  "jsonrpc": "2.0",
  "error": {
    "code": -32602,
    "message": "Invalid params",
    "data": "Error description"
  },
  "id": 1
}
```

## Troubleshooting

### Server Not Running
If you get "Server is not running" errors:
1. Check if the server process is active: `ps aux | grep jsonrpc_a2a_server`
2. Start the server: `./jsonrpc_a2a_server -ollama-url http://localhost:11434 -model gpt-oss:20b`
3. Wait a few seconds for startup

### Ollama Connection Issues
If you get Ollama-related errors:
1. Check if Ollama is running: `curl http://localhost:11434/api/tags`
2. Verify the model exists: `ollama list`
3. Pull the model if needed: `ollama pull gpt-oss:20b`

### Test Failures
- Check server logs for detailed error information
- Verify test JSON files are valid
- Ensure the server is responding on the expected port

## Adding New Tests

To add new tests:

1. **Create test data file**: Add new JSON files for different scenarios
2. **Update test script**: Add new test functions in `run_tests.sh`
3. **Update README**: Document new test cases and expected responses

## Test Output

The test script provides:
- ✅ **Green checkmarks** for successful tests
- ❌ **Red X marks** for failed tests
- 🧪 **Test details** including endpoints and methods
- 📊 **Response data** formatted with jq
- 🎨 **Colored output** for better readability
