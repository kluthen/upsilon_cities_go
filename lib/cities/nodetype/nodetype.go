package nodetype

import (
	"bytes"
	"encoding/json"
)

type NodeType int

const (
	None         NodeType = 0
	Accessible   NodeType = 1
	Inaccessible NodeType = 2
	Filled       NodeType = 3
)

type GroundType int

const (
	NoGround GroundType = 0
	Plain    GroundType = 1
	Desert   GroundType = 2
	Sea      GroundType = 3
)

type LandscapeType int

const (
	NoLandscape LandscapeType = 0
	Mountain    LandscapeType = 1
	Forest      LandscapeType = 2
	River       LandscapeType = 3
)

type ChangeType int

const (
	Any       ChangeType = 0
	Type      ChangeType = 1
	Ground    ChangeType = 2
	Landscape ChangeType = 3
	Structure ChangeType = 4
	Road      ChangeType = 5
)

var ntToEnum = map[string]NodeType{
	"None":         None,
	"Accessible":   Accessible,
	"Inaccessible": Inaccessible,
	"Filled":       Filled,
}

var gtToEnum = map[string]GroundType{
	"NoGround": NoGround,
	"Plain":    Plain,
	"Sea":      Sea,
	"Desert":   Desert,
}
var ltToEnum = map[string]LandscapeType{
	"NoLandscape": NoLandscape,
	"Forest":      Forest,
	"River":       River,
	"Mountain":    Mountain,
}

var ntNames = [...]string{
	"None",
	"Accessible",
	"Inaccessible",
	"Filled",
}

var ntShortnames = [...]string{
	".",
	".",
	"X",
	"O",
}

var gtNames = [...]string{
	"NoGround",
	"Plain",
	"Desert",
	"Sea",
}

var gtShortnames = [...]string{
	".",
	"P",
	"D",
	"S",
}

var ltNames = [...]string{
	"NoLandscape",
	"Mountain",
	"Forest",
	"River",
}

var ltShortnames = [...]string{
	".",
	"M",
	"F",
	"r",
}

func (node NodeType) String() string {
	return ntNames[node]
}

//Short short name of the node for display.
func (node NodeType) Short() string {
	return ntShortnames[node]
}

func (node GroundType) String() string {
	return gtNames[node]
}

//Short short name of the node for display.
func (node GroundType) Short() string {
	return gtShortnames[node]
}

func (node LandscapeType) String() string {

	return ltNames[node]
}

//Short short name of the node for display.
func (node LandscapeType) Short() string {

	return ltShortnames[node]
}

// MarshalJSON marshals the enum as a quoted json string
func (node LandscapeType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(node.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (node *LandscapeType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*node = ltToEnum[j]
	return nil
}

// MarshalJSON marshals the enum as a quoted json string
func (node GroundType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(node.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (node *GroundType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*node = gtToEnum[j]
	return nil
}
