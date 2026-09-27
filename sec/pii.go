package sec

import (
	"regexp"
	"strings"
)

const redactedPII = "[REDACTED]"

var (
	piiEmailPattern = regexp.MustCompile(`(?i)\b([a-z0-9][a-z0-9._%+\-]{0,63})@([a-z0-9][a-z0-9.-]*\.[a-z]{2,})\b`)
	piiCardPattern  = regexp.MustCompile(`\b(?:\d[ -]*?){13,19}\b`)
	piiPhonePattern = regexp.MustCompile(`\+?[0-9][0-9 ()-]{6,}[0-9]`)
	piiIPv4Pattern  = regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`)
)

// MaskPII masks common personal data in free-form text. It is intended for
// logs and diagnostics, not for authorization or data-retention enforcement.
// The function preserves enough shape for troubleshooting while preventing
// direct disclosure of email addresses, phone numbers, payment-card numbers,
// and IPv4 addresses.
func MaskPII(s string) string {
	s = piiEmailPattern.ReplaceAllStringFunc(s, func(value string) string {
		parts := strings.SplitN(value, "@", 2)
		local := parts[0]
		if len(local) > 1 {
			local = local[:1] + "***"
		} else {
			local = "***"
		}
		return local + "@" + parts[1]
	})
	s = piiCardPattern.ReplaceAllStringFunc(s, maskCard)
	s = piiPhonePattern.ReplaceAllStringFunc(s, maskPhone)
	s = piiIPv4Pattern.ReplaceAllStringFunc(s, maskIPv4)
	return s
}

// MaskPIIMap returns a copy of fields with common secret-bearing keys
// redacted and all string values passed through MaskPII. The input map is not
// mutated. Nested maps and []any values are handled recursively.
func MaskPIIMap(input map[string]any) map[string]any {
	output := make(map[string]any, len(input))
	for key, value := range input {
		if sensitivePIIKey(key) {
			output[key] = redactedPII
			continue
		}
		output[key] = maskPIIValue(value)
	}
	return output
}

func maskPIIValue(value any) any {
	switch value := value.(type) {
	case string:
		return MaskPII(value)
	case map[string]any:
		return MaskPIIMap(value)
	case []any:
		masked := make([]any, len(value))
		for i, item := range value {
			masked[i] = maskPIIValue(item)
		}
		return masked
	default:
		return value
	}
}

func sensitivePIIKey(key string) bool {
	switch strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(key, "-", ""), "_", "")) {
	case "password", "passwd", "secret", "token", "accesstoken", "refreshtoken", "authorization", "cookie", "setcookie", "clientsecret", "privatekey", "ssn", "socialsecuritynumber":
		return true
	default:
		return false
	}
}

func maskCard(value string) string {
	digits := make([]byte, 0, len(value))
	for i := 0; i < len(value); i++ {
		if value[i] >= '0' && value[i] <= '9' {
			digits = append(digits, value[i])
		}
	}
	if len(digits) < 13 || len(digits) > 19 {
		return value
	}
	return "****" + string(digits[len(digits)-4:])
}

func maskPhone(value string) string {
	digits := make([]byte, 0, len(value))
	for i := 0; i < len(value); i++ {
		if value[i] >= '0' && value[i] <= '9' {
			digits = append(digits, value[i])
		}
	}
	if len(digits) < 8 {
		return value
	}
	return "***-***-" + string(digits[len(digits)-4:])
}

func maskIPv4(value string) string {
	parts := strings.Split(value, ".")
	if len(parts) != 4 {
		return value
	}
	for _, part := range parts {
		if part == "" || len(part) > 3 {
			return value
		}
		for i := 0; i < len(part); i++ {
			if part[i] < '0' || part[i] > '9' {
				return value
			}
		}
	}
	return parts[0] + "." + parts[1] + ".xxx.xxx"
}
