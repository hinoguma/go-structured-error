package fault

import (
	"reflect"
	"testing"
)

func TestTags_SetValueSafe(t *testing.T) {
	testCases := []struct {
		label    string
		initial  Tags
		setFunc  func(attrs *Tags)
		expected Tags
	}{
		{
			label:   "add new attribute on initially empty attributes",
			initial: Tags{},
			setFunc: func(attrs *Tags) {
				attrs.SetValueSafe("key1", StringTagValue("value1"))
			},
			expected: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
				},
				keyMap: map[string]int{
					"key1": 0,
				},
			},
		},
		{
			label: "add new attribute on not empty attributes",
			initial: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
				},
				keyMap: map[string]int{
					"key1": 0,
				},
			},
			setFunc: func(attrs *Tags) {
				attrs.SetValueSafe("key2", StringTagValue("value2"))
			},
			expected: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
					{Key: "key2", Value: StringTagValue("value2")},
				},
				keyMap: map[string]int{
					"key1": 0,
					"key2": 1,
				},
			},
		},
		{
			label: "add new attribute with duplicate key",
			initial: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
				},
				keyMap: map[string]int{
					"key1": 0,
				},
			},
			setFunc: func(attrs *Tags) {
				attrs.SetValueSafe("key1", StringTagValue("updatedValue1"))
			},
			expected: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("updatedValue1")},
				},
				keyMap: map[string]int{
					"key1": 0,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			attrs := tc.initial
			tc.setFunc(&attrs)
			if !reflect.DeepEqual(attrs, tc.expected) {
				t.Errorf("expected attributes %v, got %v", tc.expected, attrs)
			}
		})
	}

}

func TestTags_GetValue(t *testing.T) {
	testCases := []struct {
		label         string
		attributes    Tags
		getKey        string
		expectedValue TagValue
		expectedOk    bool
	}{
		{
			label: "get existing attribute",
			attributes: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
				},
				keyMap: map[string]int{
					"key1": 0,
				},
			},
			getKey:        "key1",
			expectedValue: StringTagValue("value1"),
			expectedOk:    true,
		},
		{
			label: "get non-existing attribute",
			attributes: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
				},
				keyMap: map[string]int{
					"key1": 0,
				},
			},
			getKey:        "key2",
			expectedValue: nil,
			expectedOk:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			value, ok := tc.attributes.GetValue(tc.getKey)
			if ok != tc.expectedOk {
				t.Errorf("expected ok %v, got %v", tc.expectedOk, ok)
			}
			if !reflect.DeepEqual(value, tc.expectedValue) {
				t.Errorf("expected value %v, got %v", tc.expectedValue, value)
			}
		})
	}
}

func TestTags_Delete(t *testing.T) {
	testCases := []struct {
		label     string
		initial   Tags
		deleteKey string
		expected  Tags
	}{
		{
			label: "delete existing attribute only one attribute",
			initial: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
				},
				keyMap: map[string]int{
					"key1": 0,
				},
			},
			deleteKey: "key1",
			expected: Tags{
				tags:   []Tag{},
				keyMap: map[string]int{},
			},
		},
		{
			label: "delete existing attribute in beginning",
			initial: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
					{Key: "key2", Value: StringTagValue("value2")},
				},
				keyMap: map[string]int{
					"key1": 0,
					"key2": 1,
				},
			},
			deleteKey: "key1",
			expected: Tags{
				tags: []Tag{
					{Key: "key2", Value: StringTagValue("value2")},
				},
				keyMap: map[string]int{
					"key2": 0,
				},
			},
		},
		{
			label: "delete existing attribute in end",
			initial: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
					{Key: "key2", Value: StringTagValue("value2")},
				},
				keyMap: map[string]int{
					"key1": 0,
					"key2": 1,
				},
			},
			deleteKey: "key2",
			expected: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
				},
				keyMap: map[string]int{
					"key1": 0,
				},
			},
		},
		{
			label: "delete existing attribute in middle",
			initial: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
					{Key: "key2", Value: StringTagValue("value2")},
					{Key: "key3", Value: StringTagValue("value3")},
					{Key: "key4", Value: StringTagValue("value4")},
					{Key: "key5", Value: StringTagValue("value5")},
				},
				keyMap: map[string]int{
					"key1": 0,
					"key2": 1,
					"key3": 2,
					"key4": 3,
					"key5": 4,
				},
			},
			deleteKey: "key3",
			expected: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
					{Key: "key2", Value: StringTagValue("value2")},
					{Key: "key4", Value: StringTagValue("value4")},
					{Key: "key5", Value: StringTagValue("value5")},
				},
				keyMap: map[string]int{
					"key1": 0,
					"key2": 1,
					"key4": 2,
					"key5": 3,
				},
			},
		},

		{
			label: "delete non-existing attribute",
			initial: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
					{Key: "key2", Value: StringTagValue("value2")},
				},
				keyMap: map[string]int{
					"key1": 0,
					"key2": 1,
				},
			},
			deleteKey: "key3",
			expected: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
					{Key: "key2", Value: StringTagValue("value2")},
				},
				keyMap: map[string]int{
					"key1": 0,
					"key2": 1,
				},
			},
		},
		{
			label: "delete attribute from empty attributes",
			initial: Tags{
				tags:   []Tag{},
				keyMap: map[string]int{},
			},
			deleteKey: "key1",
			expected: Tags{
				tags:   []Tag{},
				keyMap: map[string]int{},
			},
		},
		{
			label:     "delete attribute from nil attributes",
			initial:   Tags{},
			deleteKey: "key1",
			expected:  Tags{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			attrs := tc.initial
			attrs.Delete(tc.deleteKey)
			if !reflect.DeepEqual(attrs, tc.expected) {
				t.Errorf("expected attributes %v, got %v", tc.expected, attrs)
			}
		})
	}
}

func TestTags_saveKey(t *testing.T) {

	testCases := []struct {
		label    string
		initial  Tags
		key      string
		index    int
		expected Tags
	}{
		{
			label: "save key in empty keyMap",
			initial: Tags{
				tags:   []Tag{},
				keyMap: nil,
			},
			key:   "key1",
			index: 0,
			expected: Tags{
				tags: []Tag{},
				keyMap: map[string]int{
					"key1": 0,
				},
			},
		},
		{
			label: "save key in non-empty keyMap",
			initial: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
				},
				keyMap: map[string]int{
					"key1": 0,
				},
			},
			key:   "key1",
			index: 0,
			expected: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
				},
				keyMap: map[string]int{
					"key1": 0,
				},
			},
		},
		{
			label: "save another key in non-empty keyMap",
			initial: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
				},
				keyMap: map[string]int{
					"key1": 0,
				},
			},
			key:   "key2",
			index: 1,
			expected: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
				},
				keyMap: map[string]int{
					"key1": 0,
					"key2": 1,
				},
			},
		},
		{
			label: "overwrite existing key in keyMap",
			initial: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
				},
				keyMap: map[string]int{
					"key1": 0,
				},
			},
			key:   "key1",
			index: 2,
			expected: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
				},
				keyMap: map[string]int{
					"key1": 2,
				},
			},
		},
		{
			label: "save empty key",
			initial: Tags{
				tags:   []Tag{},
				keyMap: map[string]int{},
			},
			key:   "",
			index: 0,
			expected: Tags{
				tags:   []Tag{},
				keyMap: map[string]int{"": 0},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			attrs := tc.initial
			attrs.saveKey(tc.key, tc.index)
			if !reflect.DeepEqual(attrs, tc.expected) {
				t.Errorf("expected attributes %v, got %v", tc.expected, attrs)
			}
		})
	}

}

func TestTags_GetIndexByKey(t *testing.T) {
	testCases := []struct {
		label          string
		attributes     Tags
		key            string
		expectedIdx    int
		expectedExists bool
	}{
		{
			label: "key exists",
			attributes: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
					{Key: "key2", Value: StringTagValue("value2")},
				},
				keyMap: map[string]int{
					"key1": 0,
					"key2": 1,
				},
			},
			key:            "key2",
			expectedIdx:    1,
			expectedExists: true,
		},
		{
			label: "key does not exist",
			attributes: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
					{Key: "key2", Value: StringTagValue("value2")},
				},
				keyMap: map[string]int{
					"key1": 0,
					"key2": 1,
				},
			},
			key:            "key3",
			expectedIdx:    0,
			expectedExists: false,
		},
		{
			label: "empty attributes",
			attributes: Tags{
				tags:   []Tag{},
				keyMap: map[string]int{},
			},
			key:            "key1",
			expectedIdx:    0,
			expectedExists: false,
		},
		{
			label:          "nil attributes",
			attributes:     Tags{},
			key:            "key1",
			expectedIdx:    0,
			expectedExists: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			idx, exists := tc.attributes.GetIndexByKey(tc.key)
			if idx != tc.expectedIdx {
				t.Errorf("expected index %v, got %v", tc.expectedIdx, idx)
			}
			if exists != tc.expectedExists {
				t.Errorf("expected exists %v, got %v", tc.expectedExists, exists)
			}
		})
	}
}

func TestTags_deleteKey(t *testing.T) {
	testCases := []struct {
		label     string
		initial   Tags
		deleteKey string
		expected  Tags
	}{
		{
			label: "delete existing key only one attribute",
			initial: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
				},
				keyMap: map[string]int{
					"key1": 0,
				},
			},
			deleteKey: "key1",
			expected: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
				},
				keyMap: map[string]int{},
			},
		},
		{
			label: "delete existing key in beginning",
			initial: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
					{Key: "key2", Value: StringTagValue("value2")},
				},
				keyMap: map[string]int{
					"key1": 0,
					"key2": 1,
				},
			},
			deleteKey: "key1",
			expected: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
					{Key: "key2", Value: StringTagValue("value2")},
				},
				keyMap: map[string]int{
					"key2": 0,
				},
			},
		},
		{
			label: "delete existing key in end",
			initial: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
					{Key: "key2", Value: StringTagValue("value2")},
				},
				keyMap: map[string]int{
					"key1": 0,
					"key2": 1,
				},
			},
			deleteKey: "key2",
			expected: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
					{Key: "key2", Value: StringTagValue("value2")},
				},
				keyMap: map[string]int{
					"key1": 0,
				},
			},
		},
		{
			label: "delete existing key in middle",
			initial: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
					{Key: "key2", Value: StringTagValue("value2")},
					{Key: "key3", Value: StringTagValue("value3")},
					{Key: "key4", Value: StringTagValue("value4")},
					{Key: "key5", Value: StringTagValue("value5")},
				},
				keyMap: map[string]int{
					"key1": 0,
					"key2": 1,
					"key3": 2,
					"key4": 3,
					"key5": 4,
				},
			},
			deleteKey: "key2",
			expected: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
					{Key: "key2", Value: StringTagValue("value2")},
					{Key: "key3", Value: StringTagValue("value3")},
					{Key: "key4", Value: StringTagValue("value4")},
					{Key: "key5", Value: StringTagValue("value5")},
				},
				keyMap: map[string]int{
					"key1": 0,
					"key3": 1,
					"key4": 2,
					"key5": 3,
				},
			},
		},
		{
			label: "delete non-existing key",
			initial: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
					{Key: "key2", Value: StringTagValue("value2")},
				},
				keyMap: map[string]int{
					"key1": 0,
					"key2": 1,
				},
			},
			deleteKey: "key3",
			expected: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
					{Key: "key2", Value: StringTagValue("value2")},
				},
				keyMap: map[string]int{
					"key1": 0,
					"key2": 1,
				},
			},
		},
		{
			label: "delete key from empty attributes",
			initial: Tags{
				tags:   []Tag{},
				keyMap: map[string]int{},
			},
			deleteKey: "key1",
			expected: Tags{
				tags:   []Tag{},
				keyMap: map[string]int{},
			},
		},
		{
			label:     "delete key from nil attributes",
			initial:   Tags{},
			deleteKey: "key1",
			expected:  Tags{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			attrs := tc.initial
			attrs.deleteKey(tc.deleteKey)
			if !reflect.DeepEqual(attrs, tc.expected) {
				t.Errorf("expected attributes %v, got %v", tc.expected, attrs)
			}
		})
	}

}

func TestTags_JsonString(t *testing.T) {
	testCases := []struct {
		label    string
		tags     Tags
		expected string
	}{
		{
			label: "empty tags",
			tags: Tags{
				tags:   []Tag{},
				keyMap: map[string]int{},
			},
			expected: "{}",
		},
		{
			label: "single tag",
			tags: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
				},
				keyMap: map[string]int{
					"key1": 0,
				},
			},
			expected: `{"key1":"value1"}`,
		},
		{
			label: "multiple tags",
			tags: Tags{
				tags: []Tag{
					{Key: "key1", Value: StringTagValue("value1")},
					{Key: "key2", Value: IntTagValue(42)},
					{Key: "key3", Value: BoolTagValue(true)},
				},
				keyMap: map[string]int{
					"key1": 0,
					"key2": 1,
					"key3": 2,
				},
			},
			expected: `{"key1":"value1","key2":42,"key3":true}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := tc.tags.JsonString()
			if got != tc.expected {
				t.Errorf("expected JSON string %v, got %v", tc.expected, got)
			}
		})
	}
}
