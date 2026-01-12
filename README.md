# cURL to Go Converter

A web application that converts cURL commands (copied from Firefox's network tab) into executable Go code.

## Features

- Parse cURL commands and extract:
  - HTTP method (GET, POST, etc.)
  - URL
  - Headers
  - Cookies
  - Request body/data
- Generate clean Go code with constants for all request parameters
- Modern, responsive web interface
- Real-time conversion
- Copy-to-clipboard functionality

## Quick Start

```bash
# Run the server
go run main.go

# Visit http://localhost:8080 in your browser
```

## Usage

1. Open Firefox Developer Tools (F12)
2. Go to the Network tab
3. Right-click on any request and select "Copy as cURL"
4. Paste the cURL command into the web interface
5. Click "Convert to Go"
6. Copy the generated Go code

## Example

**Input (cURL command):**
```bash
curl 'https://api.example.com/users' \
  -H 'Authorization: Bearer token123' \
  -H 'Content-Type: application/json' \
  -d '{"name":"John"}'
```

**Output (Go code):**
```go
package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {
	const method = "POST"
	const url = "https://api.example.com/users"
	const body = "{\"name\":\"John\"}"

	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		panic(err)
	}

	// Headers
	req.Header.Set("Authorization", "Bearer token123")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("Status:", resp.Status)
	fmt.Println("Body:", string(respBody))
}
```

## Project Structure

```
curl-to-go/
├── main.go              # Web server with chi router
├── converter/
│   ├── converter.go     # cURL parsing and Go code generation
│   └── converter_test.go # Tests
└── static/
    └── index.html       # Web interface
```

## API

### POST /api/convert

Convert a cURL command to Go code.

**Request:**
```json
{
  "curl_command": "curl 'https://example.com' -H 'Authorization: Bearer token'"
}
```

**Response:**
```json
{
  "go_code": "package main\n\n..."
}
```

**Error Response:**
```json
{
  "error": "no URL found in curl command"
}
```

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test ./converter/... -v
```

## Supported cURL Flags

- `-X, --request` - HTTP method
- `-H, --header` - Headers
- `-d, --data, --data-raw, --data-binary` - Request body
- `--cookie` - Cookies
- Cookie header parsing

## Development

Built with:
- [go-chi/chi](https://github.com/go-chi/chi) - Lightweight router
- Vanilla JavaScript - No frontend frameworks
- Go standard library for HTTP client code generation

## License

MIT
