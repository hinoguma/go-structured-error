package serrors

import "fmt"

// Tag represents a key-value pair used to provide additional information about an error or event.
//
// This becomes {"Tag.key": Tag.value} in JSON format.
type Tag struct {
	Key   string
	Value TagValue
}

func (tag Tag) String() string {
	return fmt.Sprintf("key:%s value:%s", tag.Key, tag.Value.String())
}

func NewTags() Tags {
	return Tags{
		tags:   make([]Tag, 0),
		keyMap: make(map[string]int),
	}
}

// Tags are key-value pairs that provide additional information about an error or event.
//
// Tags can be converted to JSON format for easy serialization and logging.
// To ensure the order of tags is preserved, Tags are stored as a slice of Tag structs internally.
// keyMap is used for quick lookup of tag indices by their keys.
type Tags struct {
	tags   []Tag
	keyMap map[string]int
}

func (tags Tags) GetValue(key string) (TagValue, bool) {
	index, exists := tags.GetIndexByKey(key)
	if !exists {
		return nil, false
	}
	return tags.tags[index].Value, true
}

// SetValue() is not implemented yet.
// It`s easy and simple for user to add value to Tags
//func (tags Tags) SetValue(key string, value any) bool {
//	panic("not implemented")
//}

// SetValueSafe sets the value of a tag with the given key.
//
// If the key already exists, its value is updated.
// If the key does not exist, a new tag is added.
func (tags *Tags) SetValueSafe(key string, value TagValue) {
	if tags.tags == nil {
		tags.tags = make([]Tag, 0)
	}

	index, exists := tags.GetIndexByKey(key)

	// update existing tag if key exists
	if exists {
		tags.tags[index].Value = value
		return
	}

	// add new tag
	tags.saveKey(key, len(tags.tags))
	tags.tags = append(tags.tags, Tag{
		Key:   key,
		Value: value,
	})
}

func (tags *Tags) Delete(key string) {
	index, exists := tags.GetIndexByKey(key)
	if !exists {
		return
	}
	tags.deleteKey(key)
	tags.tags = append(tags.tags[:index], tags.tags[index+1:]...)
}

func (tags Tags) GetIndexByKey(key string) (int, bool) {
	if tags.keyMap == nil {
		return 0, false
	}
	index, ok := tags.keyMap[key]
	return index, ok
}

func (tags *Tags) saveKey(key string, index int) {
	if tags.keyMap == nil {
		tags.keyMap = make(map[string]int)
	}
	tags.keyMap[key] = index
}

func (tags *Tags) deleteKey(key string) {
	if tags.keyMap == nil {
		return
	}
	index, exists := tags.GetIndexByKey(key)
	if !exists {
		return
	}
	for tmpKey, tmpIndex := range tags.keyMap {
		if tmpIndex > index {
			tags.keyMap[tmpKey] = tmpIndex - 1
		}
	}
	delete(tags.keyMap, key)
}

func (tags Tags) JsonValueString() string {
	result := "{"
	for i, tag := range tags.tags {
		if i > 0 {
			result += JsonItemSeparator
		}
		result += "\"" + tag.Key + "\":" + tag.Value.JsonValueString()
	}
	result += "}"
	return result
}
