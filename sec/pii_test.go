package sec

// Tests for PII masking and safe diagnostic redaction.
import (
	"reflect"
	"testing"
)

func TestMaskPII(t *testing.T) {
	input := "email alice@example.com phone +1 (555) 123-4567 card 4111 1111 1111 1111 ip 192.168.1.42"
	want := "email a***@example.com phone ***-***-4567 card ****1111 ip 192.168.xxx.xxx"
	if got := MaskPII(input); got != want {
		t.Fatalf("MaskPII() = %q, want %q", got, want)
	}
}

func TestMaskPIIMapDoesNotMutateAndRedactsNestedValues(t *testing.T) {
	input := map[string]any{
		"email":    "alice@example.com",
		"password": "do-not-log",
		"nested": map[string]any{
			"phone": "+1 555 123 4567",
		},
		"items": []any{"bob@example.com"},
	}
	original := map[string]any{
		"email":    "alice@example.com",
		"password": "do-not-log",
		"nested":   map[string]any{"phone": "+1 555 123 4567"},
		"items":    []any{"bob@example.com"},
	}

	got := MaskPIIMap(input)
	if !reflect.DeepEqual(input, original) {
		t.Fatalf("MaskPIIMap mutated input: %#v", input)
	}
	if got["password"] != "[REDACTED]" {
		t.Fatalf("password = %#v", got["password"])
	}
	if got["email"] != "a***@example.com" {
		t.Fatalf("email = %#v", got["email"])
	}
	nested := got["nested"].(map[string]any)
	if nested["phone"] != "***-***-4567" {
		t.Fatalf("nested phone = %#v", nested["phone"])
	}
}

func TestMaskPIILeavesNonPIIText(t *testing.T) {
	if got := MaskPII("request completed with status 200"); got != "request completed with status 200" {
		t.Fatalf("unexpected masking: %q", got)
	}
}
