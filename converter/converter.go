package converter

import (
	"fmt"
	"regexp"
	"strings"
)

type CurlRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Cookies map[string]string
	Body    string
}

// CurlToGo converts a curl command to Go code
func CurlToGo(curlCommand string) (string, error) {
	req, err := ParseCurl(curlCommand)
	if err != nil {
		return "", err
	}

	return GenerateGoCode(req), nil
}

// ParseCurl parses a curl command and extracts request details
func ParseCurl(curlCommand string) (*CurlRequest, error) {
	req := &CurlRequest{
		Method:  "GET",
		Headers: make(map[string]string),
		Cookies: make(map[string]string),
	}

	// Remove leading/trailing whitespace and newlines
	curlCommand = strings.TrimSpace(curlCommand)

	// Remove line continuations (backslash + newline)
	curlCommand = regexp.MustCompile(`\\\s*\n\s*`).ReplaceAllString(curlCommand, " ")

	// Remove extra whitespace
	curlCommand = regexp.MustCompile(`\s+`).ReplaceAllString(curlCommand, " ")

	// Extract URL - look for http(s):// patterns
	// This handles both quoted and unquoted URLs
	urlPattern := regexp.MustCompile(`(?:'(https?://[^']+)'|"(https?://[^"]+)"|(https?://\S+))`)
	if matches := urlPattern.FindStringSubmatch(curlCommand); len(matches) > 0 {
		for i := 1; i < len(matches); i++ {
			if matches[i] != "" {
				req.URL = matches[i]
				break
			}
		}
	}

	if req.URL == "" {
		return nil, fmt.Errorf("no URL found in curl command")
	}

	// Extract method (-X or --request)
	methodRegex := regexp.MustCompile(`(?:-X|--request)\s+(?:'([^']+)'|"([^"]+)"|(\S+))`)
	if matches := methodRegex.FindStringSubmatch(curlCommand); len(matches) > 0 {
		for i := 1; i < len(matches); i++ {
			if matches[i] != "" {
				req.Method = strings.ToUpper(matches[i])
				break
			}
		}
	}

	// Extract headers (-H or --header)
	headerRegex := regexp.MustCompile(`(?:-H|--header)\s+(?:'([^']+)'|"([^"]+)")`)
	headerMatches := headerRegex.FindAllStringSubmatch(curlCommand, -1)
	for _, match := range headerMatches {
		headerValue := ""
		for i := 1; i < len(match); i++ {
			if match[i] != "" {
				headerValue = match[i]
				break
			}
		}

		// Parse header into key: value
		parts := strings.SplitN(headerValue, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Check if this is a cookie header
			if strings.ToLower(key) == "cookie" {
				parseCookies(value, req.Cookies)
			} else {
				req.Headers[key] = value
			}
		}
	}

	// Extract cookies (--cookie)
	cookieRegex := regexp.MustCompile(`(?:--cookie)\s+(?:'([^']+)'|"([^"]+)")`)
	if matches := cookieRegex.FindStringSubmatch(curlCommand); len(matches) > 0 {
		cookieValue := ""
		for i := 1; i < len(matches); i++ {
			if matches[i] != "" {
				cookieValue = matches[i]
				break
			}
		}
		parseCookies(cookieValue, req.Cookies)
	}

	// Extract data (-d, --data, --data-raw, --data-binary)
	dataRegex := regexp.MustCompile(`(?:-d|--data|--data-raw|--data-binary)\s+(?:'([^']+)'|"([^"]+)"|(\S+))`)
	if matches := dataRegex.FindStringSubmatch(curlCommand); len(matches) > 0 {
		for i := 1; i < len(matches); i++ {
			if matches[i] != "" {
				req.Body = matches[i]
				// If data is present and method is GET, change to POST
				if req.Method == "GET" {
					req.Method = "POST"
				}
				break
			}
		}
	}

	return req, nil
}

// parseCookies parses a cookie string and adds them to the cookies map
func parseCookies(cookieStr string, cookies map[string]string) {
	parts := strings.Split(cookieStr, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		kvParts := strings.SplitN(part, "=", 2)
		if len(kvParts) == 2 {
			cookies[strings.TrimSpace(kvParts[0])] = strings.TrimSpace(kvParts[1])
		}
	}
}

// GenerateGoCode generates Go code from a CurlRequest
func GenerateGoCode(req *CurlRequest) string {
	var sb strings.Builder

	sb.WriteString("package main\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("\t\"fmt\"\n")
	sb.WriteString("\t\"io\"\n")
	sb.WriteString("\t\"net/http\"\n")

	if req.Body != "" {
		sb.WriteString("\t\"strings\"\n")
	}

	sb.WriteString(")\n\n")

	// Write constants
	sb.WriteString("func main() {\n")
	sb.WriteString(fmt.Sprintf("\tconst method = %q\n", req.Method))
	sb.WriteString(fmt.Sprintf("\tconst url = %q\n", req.URL))

	// Write body constant if present
	if req.Body != "" {
		sb.WriteString(fmt.Sprintf("\tconst body = %q\n", req.Body))
	}

	sb.WriteString("\n")

	// Create request
	if req.Body != "" {
		sb.WriteString("\treq, err := http.NewRequest(method, url, strings.NewReader(body))\n")
	} else {
		sb.WriteString("\treq, err := http.NewRequest(method, url, nil)\n")
	}
	sb.WriteString("\tif err != nil {\n")
	sb.WriteString("\t\tpanic(err)\n")
	sb.WriteString("\t}\n\n")

	// Add headers
	if len(req.Headers) > 0 {
		sb.WriteString("\t// Headers\n")
		for key, value := range req.Headers {
			sb.WriteString(fmt.Sprintf("\treq.Header.Set(%q, %q)\n", key, value))
		}
		sb.WriteString("\n")
	}

	// Add cookies
	if len(req.Cookies) > 0 {
		sb.WriteString("\t// Cookies\n")
		for name, value := range req.Cookies {
			sb.WriteString(fmt.Sprintf("\treq.AddCookie(&http.Cookie{Name: %q, Value: %q})\n", name, value))
		}
		sb.WriteString("\n")
	}

	// Make request
	sb.WriteString("\tclient := &http.Client{}\n")
	sb.WriteString("\tresp, err := client.Do(req)\n")
	sb.WriteString("\tif err != nil {\n")
	sb.WriteString("\t\tpanic(err)\n")
	sb.WriteString("\t}\n")
	sb.WriteString("\tdefer resp.Body.Close()\n\n")

	// Read response
	sb.WriteString("\trespBody, err := io.ReadAll(resp.Body)\n")
	sb.WriteString("\tif err != nil {\n")
	sb.WriteString("\t\tpanic(err)\n")
	sb.WriteString("\t}\n\n")

	sb.WriteString("\tfmt.Println(\"Status:\", resp.Status)\n")
	sb.WriteString("\tfmt.Println(\"Body:\", string(respBody))\n")
	sb.WriteString("}\n")

	return sb.String()
}
