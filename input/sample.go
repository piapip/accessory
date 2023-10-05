package input

type Value struct {
	Value   string
	Comment string
}

type Enums struct {
	Title  string
	Values []*Value
}
