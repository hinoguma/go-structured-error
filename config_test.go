package fault

import "testing"

func TestSetMaxDepthStackTrace(t *testing.T) {
	testCases := []struct {
		label    string
		depth    int
		expected int
	}{
		{
			label:    "set max depth to 50",
			depth:    50,
			expected: 50,
		},
		{
			label:    "set max depth to 10",
			depth:    10,
			expected: 10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			SetMaxDepthStackTrace(tc.depth)
			got := GetMaxDepthStackTrace()
			if got != tc.expected {
				t.Errorf("expected max depth %v, got %v", tc.expected, got)
			}
		})
	}
}
