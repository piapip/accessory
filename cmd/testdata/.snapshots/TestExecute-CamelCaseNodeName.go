// Code generated by accessory; DO NOT EDIT.

package test

func (t *Tester) FirstField() string {
	if t == nil {
		return ""
	}
	return t.firstField
}

func (t *Tester) SetSecondField(val int32) {
	if t == nil {
		return
	}
	t.secondField = val
}

func (t *Tester) ThirdField() int32 {
	if t == nil {
		return 0
	}
	return t.thirdField
}

func (t *Tester) SetThirdField(val int32) {
	if t == nil {
		return
	}
	t.thirdField = val
}
