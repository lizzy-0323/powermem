# PowerMem Go Client Example

A simple Go client example demonstrating how to integrate PowerMem's intelligent memory capabilities into Go applications using the HTTP API Server.

## Prerequisites

- Go 1.21 or higher
- PowerMem HTTP API Server running (default: <http://localhost:8000>)

## Installation

```bash
cd examples/go
go mod download
```

## Configuration

Configure the client using environment variables:

```bash
export POWERMEM_BASE_URL="http://localhost:8000"
export POWERMEM_API_KEY="your-api-key-here"
```

Or use the defaults:

- Base URL: `http://localhost:8000`
- API Key: empty (if authentication is disabled)

## Project Structure

```
examples/go/
├── main.go       # Main example demonstrating all operations
├── client.go     # PowerMem client wrapper
├── models.go     # Request/response structs
├── go.mod        # Go module file
└── README.md     # This file
```

## Usage

### Basic Example

```bash
go run .
```

### With Custom Configuration

```bash
POWERMEM_BASE_URL="http://localhost:8000" \
POWERMEM_API_KEY="test-api-key-123" \
go run .
```

## Examples

### Initialize Client

```go
client := NewClient("http://localhost:8000", "your-api-key")
```

### Create Memory

```go
memories, err := client.CreateMemory(CreateMemoryRequest{
    Content: "User likes coffee",
    UserID:  "user-123",
    AgentID: "agent-456",
    Metadata: map[string]interface{}{
        "source": "conversation",
    },
    Infer: true,
})
```

### Search Memories

```go
results, err := client.SearchMemories(SearchMemoryRequest{
    Query:   "What beverages does the user like",
    UserID:  "user-123",
    AgentID: "agent-456",
    Limit:   10,
})
```

### List Memories

```go
memoryList, err := client.ListMemories(ListMemoriesRequest{
    UserID:  "user-123",
    AgentID: "agent-456",
    Limit:   10,
    Offset:  0,
})
```

### Update Memory

```go
updatedMemory, err := client.UpdateMemory(memoryID, UpdateMemoryRequest{
    Content: "User loves latte coffee",
    UserID:  "user-123",
    AgentID: "agent-456",
})
```

### Delete Memory

```go
err := client.DeleteMemory(memoryID, "user-123", "agent-456")
```

## API Reference

### Client Methods

- `NewClient(baseURL, apiKey string) *Client` - Create a new client instance
- `CreateMemory(req CreateMemoryRequest) ([]Memory, error)` - Create a new memory
- `ListMemories(req ListMemoriesRequest) (*MemoryListData, error)` - List memories with filtering
- `SearchMemories(req SearchMemoryRequest) (*SearchData, error)` - Search memories semantically
- `UpdateMemory(memoryID int64, req UpdateMemoryRequest) (*Memory, error)` - Update a memory
- `DeleteMemory(memoryID int64, userID, agentID string) error` - Delete a memory
