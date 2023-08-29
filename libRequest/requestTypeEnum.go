// Code generated by "enumer -type=Type -json -output requestTypeEnum.go"; DO NOT EDIT.

package libRequest

import (
	"encoding/json"
	"fmt"
	"strings"
)

const _TypeName = "JSONJSONWithUriQuery"

var _TypeIndex = [...]uint8{0, 4, 15, 20}

const _TypeLowerName = "jsonjsonwithuriquery"

func (i Type) String() string {
	if i < 0 || i >= Type(len(_TypeIndex)-1) {
		return fmt.Sprintf("Type(%d)", i)
	}
	return _TypeName[_TypeIndex[i]:_TypeIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _TypeNoOp() {
	var x [1]struct{}
	_ = x[JSON-(0)]
	_ = x[JSONWithUri-(1)]
	_ = x[Query-(2)]
}

var _TypeValues = []Type{JSON, JSONWithUri, Query}

var _TypeNameToValueMap = map[string]Type{
	_TypeName[0:4]:        JSON,
	_TypeLowerName[0:4]:   JSON,
	_TypeName[4:15]:       JSONWithUri,
	_TypeLowerName[4:15]:  JSONWithUri,
	_TypeName[15:20]:      Query,
	_TypeLowerName[15:20]: Query,
}

var _TypeNames = []string{
	_TypeName[0:4],
	_TypeName[4:15],
	_TypeName[15:20],
}

// TypeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func TypeString(s string) (Type, error) {
	if val, ok := _TypeNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _TypeNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Type values", s)
}

// TypeValues returns all values of the enum
func TypeValues() []Type {
	return _TypeValues
}

// TypeStrings returns a slice of all String values of the enum
func TypeStrings() []string {
	strs := make([]string, len(_TypeNames))
	copy(strs, _TypeNames)
	return strs
}

// IsAType returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Type) IsAType() bool {
	for _, v := range _TypeValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for Type
func (i Type) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for Type
func (i *Type) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("Type should be a string, got %s", data)
	}

	var err error
	*i, err = TypeString(s)
	return err
}
