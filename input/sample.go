package input

import (
	c "github.com/zeals-co-ltd/zero-api/generated/go/entities/common"
)

type Tester struct {
	field1 *c.Card `accessor:"getter"`
	field2 int32   `accessor:"getter:GetSecondField"`
	field3 *bool
}
