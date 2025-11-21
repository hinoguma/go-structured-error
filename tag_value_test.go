package fault

import "testing"

func TestTagValue_String(t *testing.T) {
	testCases := []struct {
		label    string
		tagVal   TagValue
		expected string
	}{
		{
			label:    "StringTagValue returns correct string",
			tagVal:   StringTagValue("example string"),
			expected: "example string",
		},
		{
			label:    "IntTagValue returns correct string",
			tagVal:   IntTagValue(42),
			expected: "42",
		},
		{
			label:    "BoolTagValue returns correct string",
			tagVal:   BoolTagValue(true),
			expected: "true",
		},
		{
			label:    "BoolTagValue returns correct string for false",
			tagVal:   BoolTagValue(false),
			expected: "false",
		},
		{
			label:    "FloatTagValue returns correct string",
			tagVal:   FloatTagValue(3.14),
			expected: "3.14",
		},
		{
			label:    "NilTagValue returns correct string",
			tagVal:   NilTagValue{},
			expected: "null",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := tc.tagVal.String()
			if got != tc.expected {
				t.Errorf("expected string %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestTagValue_JsonValueString(t *testing.T) {
	testCases := []struct {
		label    string
		tagVal   TagValue
		expected string
	}{
		{
			label:    "StringTagValue",
			tagVal:   StringTagValue("example string"),
			expected: "\"example string\"",
		},
		{
			label:    "string with escapes",
			tagVal:   StringTagValue("line1\nline2\"quote\"\\backslash"),
			expected: "\"line1\\nline2\\\"quote\\\"\\\\backslash\"",
		},
		{
			label:    "IntTagValue",
			tagVal:   IntTagValue(42),
			expected: "42",
		},
		{
			label:    "BoolTagValue true",
			tagVal:   BoolTagValue(true),
			expected: "true",
		},
		{
			label:    "BoolTagValue false",
			tagVal:   BoolTagValue(false),
			expected: "false",
		},
		{
			label:    "FloatTagValue",
			tagVal:   FloatTagValue(3.14),
			expected: "3.14",
		},
		{
			label:    "NilTagValue",
			tagVal:   NilTagValue{},
			expected: "null",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := tc.tagVal.JsonValueString()
			if got != tc.expected {
				t.Errorf("expected JSON string %v, got %v", tc.expected, got)
			}
		})
	}
}
