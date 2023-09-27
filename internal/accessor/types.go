package accessor

import (
	"go/types"

	"golang.org/x/tools/go/packages"
)

type Package struct {
	*packages.Package
	Dir     string
	Structs []*Struct
}

type Struct struct {
	Name   string
	Fields []*Field
}

// Example:
//
//	 type Tester struct {
//		field1 string  `accessor:"getter"`
//		field2 int32   `accessor:"getter:GetSecondField"`
//		field3 *bool
//		field4 *c.Card `accessor:"getter"`
//	}
//
// Field will be:
//
//   - Name: field1
//     Type: string
//     Tag: {Getter: "", Setter: nil}
//
//   - Name: field2
//     Type: int32
//     Tag: {Getter: "GetSecondField", Setter: nil}
//
//   - Name: field3
//     Type: *bool
//     Tag: {Getter: nil, Setter: nil}
//
//   - Name: field4
//     Type: *github.com/zeals-co-ltd/zero-api/generated/go/entities/common.Card
//     Tag: {Getter: "", Setter: nil}
type Field struct {
	Name string
	Type types.Type
	Tag  *Tag
}

type Tag struct {
	Getter *string
	Setter *string
}
