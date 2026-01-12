package converter

import (
	"strings"
	"testing"
)

func TestParseCurl(t *testing.T) {
	tests := []struct {
		name        string
		curlCommand string
		wantMethod  string
		wantURL     string
		wantHeaders int
		wantBody    string
	}{
		{
			name:        "Simple GET request",
			curlCommand: `curl 'https://api.example.com/users'`,
			wantMethod:  "GET",
			wantURL:     "https://api.example.com/users",
			wantHeaders: 0,
			wantBody:    "",
		},
		{
			name: "POST with headers and data",
			curlCommand: `curl -X POST 'https://api.example.com/users' \
				-H 'Content-Type: application/json' \
				-H 'Authorization: Bearer token123' \
				-d '{"name":"John","email":"john@example.com"}'`,
			wantMethod:  "POST",
			wantURL:     "https://api.example.com/users",
			wantHeaders: 2,
			wantBody:    `{"name":"John","email":"john@example.com"}`,
		},
		{
			name: "GET with headers",
			curlCommand: `curl 'https://api.example.com/data' \
				-H 'Accept: application/json' \
				-H 'User-Agent: Mozilla/5.0'`,
			wantMethod:  "GET",
			wantURL:     "https://api.example.com/data",
			wantHeaders: 2,
			wantBody:    "",
		},
		{
			name: "Inferred POST from data",
			curlCommand: `curl 'https://api.example.com/login' \
				-d 'username=admin&password=secret'`,
			wantMethod:  "POST",
			wantURL:     "https://api.example.com/login",
			wantHeaders: 0,
			wantBody:    "username=admin&password=secret",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := ParseCurl(tt.curlCommand)
			if err != nil {
				t.Fatalf("ParseCurl() error = %v", err)
			}

			if req.Method != tt.wantMethod {
				t.Errorf("Method = %v, want %v", req.Method, tt.wantMethod)
			}

			if req.URL != tt.wantURL {
				t.Errorf("URL = %v, want %v", req.URL, tt.wantURL)
			}

			if len(req.Headers) != tt.wantHeaders {
				t.Errorf("Headers count = %v, want %v", len(req.Headers), tt.wantHeaders)
			}

			if req.Body != tt.wantBody {
				t.Errorf("Body = %v, want %v", req.Body, tt.wantBody)
			}
		})
	}
}

func TestCurlToGo(t *testing.T) {
	curlCommand := `curl 'https://api.example.com/users' \
		-H 'Authorization: Bearer token123' \
		-H 'Content-Type: application/json'`

	goCode, err := CurlToGo(curlCommand)
	if err != nil {
		t.Fatalf("CurlToGo() error = %v", err)
	}

	// Check that generated code contains expected elements
	expectedStrings := []string{
		"package main",
		`const method = "GET"`,
		`const url = "https://api.example.com/users"`,
		`req.Header.Set("Authorization", "Bearer token123")`,
		`req.Header.Set("Content-Type", "application/json")`,
		"http.NewRequest",
		"client.Do(req)",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(goCode, expected) {
			t.Errorf("Generated code does not contain: %s", expected)
		}
	}
}

func TestParseCookies(t *testing.T) {
	cookies := make(map[string]string)
	parseCookies("session=abc123; user_id=456; theme=dark", cookies)

	if cookies["session"] != "abc123" {
		t.Errorf("session cookie = %v, want abc123", cookies["session"])
	}

	if cookies["user_id"] != "456" {
		t.Errorf("user_id cookie = %v, want 456", cookies["user_id"])
	}

	if cookies["theme"] != "dark" {
		t.Errorf("theme cookie = %v, want dark", cookies["theme"])
	}
}
