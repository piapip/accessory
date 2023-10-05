package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	// caser is for converting string to UpperCase title.
	caser = cases.Title(language.AmericanEnglish)
)

const (
	dir          = "./input-enum/input.proto"
	protoPackage = "delivery_settings_entities"
	modelPackage = "models"
)

type Enum struct {
	Title    string
	Receiver string
	Values   []*Value
}

// GetTitle returns the Enum's Title.
func (e *Enum) GetTitle() string {
	if e == nil {
		return ""
	}

	return e.Title
}

// GetTitle returns the Enum's Title.
func (e *Enum) GetReceiver() string {
	if e == nil {
		return ""
	}

	if e.Receiver == "" {
		e.Receiver = string(strings.ToLower(e.Title)[0])
	}

	return e.Receiver
}

// GetValues returns the Enum's Values.
func (e *Enum) GetValues() []*Value {
	if e == nil {
		return nil
	}

	return e.Values
}

func (e *Enum) ToString() string {
	if e == nil {
		return ""
	}

	return fmt.Sprintf(`
type %s int32

const (
	%s)`, e.Title, ConvertValuesToStruct(e.Values, e.Title))
}

func (e *Enum) ToProto() string {
	if e == nil {
		return ""
	}

	return fmt.Sprintf(`
// ToProto converts the %s to Protobuf version.
func (%s %s) ToProto() %s.%s {
	switch %s {%s
	}
}`, e.Title, e.GetReceiver(), e.Title, protoPackage, e.Title, e.GetReceiver(), ConvertValuesToProtos(e.Values, e.Title))
}

func (e *Enum) ProtoToEnum() string {
	if e == nil {
		return ""
	}

	return fmt.Sprintf(`
// ProtoTo%s converts from Protobuf version to the %s.
func ProtoTo%s(%s %s.%s) %s {
	switch %s {%s
	}
}`, e.Title, e.Title, e.Title, e.GetReceiver(), protoPackage, e.Title, e.Title, e.GetReceiver(), ConvertValuesToProtoToEnum(e.Values, e.Title))
}

func (e *Enum) GenerateTest() string {
	return fmt.Sprintf(`
func Test%s_Convert(t *testing.T) {
	type want struct {
		args      %s.%s
		wantProto %s.%s
	}

	type Context struct {
		testData *want
	}

	contextInitiateFunction := func(t *testing.T) *Context {
		return &Context{}
	}

	gt.Begin(t,
		contextInitiateFunction,
		gt.Run("Convert from model to Proto and then convert back to model", func(t *testing.T, ctx *Context) {
			// Convert from model to Proto.
			gotProto := ctx.testData.args.ToProto()
			assert.Equal(t, ctx.testData.wantProto, gotProto)

			// Then convert from Proto back to model
			gotModel := %s.ProtoToTimeUnit(gotProto)
			assert.Equal(t, ctx.testData.args, gotModel)
		}).
		%s
	)
}	
`, e.GetTitle(), modelPackage, e.GetTitle(), protoPackage, e.GetTitle(), modelPackage, AssembleTestData(e.GetValues()))
}

type Value struct {
	OriginalStringValue string
	StringValue         string
	NumberValue         int
	Comment             string
}

// GetOriginalStringValue returns the Value's OriginalStringValue.
func (v *Value) GetOriginalStringValue() string {
	if v == nil {
		return ""
	}

	return v.OriginalStringValue
}

// GetStringValue returns the Value's StringValue.
func (v *Value) GetStringValue() string {
	if v == nil {
		return ""
	}

	return v.StringValue
}

// GetNumberValue returns the Value's IntValue.
func (v *Value) GetNumberValue() int {
	if v == nil {
		return 0
	}

	return v.NumberValue
}

// GetComment returns the Value's Comment.
func (v *Value) GetComment() string {
	if v == nil {
		return ""
	}

	return v.Comment
}

func (v *Value) ToStruct(enumTitle string) string {
	if v == nil {
		return ""
	}

	if v.NumberValue == 0 {
		return fmt.Sprintf(`%s
	%s %s = iota
`, v.Comment, v.StringValue, enumTitle)
	}

	return fmt.Sprintf(`%s
	%s
`, v.Comment, v.StringValue)
}

func (v *Value) ToProto(enumTitle string) string {
	if v == nil {
		return ""
	}

	if v.NumberValue == 0 {
		return fmt.Sprintf(`
	default:
		return %s.%s_%s`, protoPackage, enumTitle, v.OriginalStringValue)
	}

	return fmt.Sprintf(`
	case %s:
		return %s.%s_%s`, v.StringValue, protoPackage, enumTitle, v.OriginalStringValue)
}

func (v *Value) ProtoToEnum(enumTitle string) string {
	if v == nil {
		return ""
	}

	if v.NumberValue == 0 {
		return fmt.Sprintf(`
	default:
		return %s`, v.StringValue)
	}

	return fmt.Sprintf(`
	case %s.%s_%s:
		return %s`, protoPackage, enumTitle, v.OriginalStringValue, v.StringValue)
}

func (v *Value) ToTestData() string {
	if v == nil {
		return ""
	}

	return fmt.Sprintf(`
		Using("given %s value", func(t *testing.T, ctx *Context) {
			ctx.testData = &want{
				args:      %s.%s,
				wantProto: %s.%s,
			}
		})`, v.StringValue, modelPackage, v.StringValue, protoPackage, v.OriginalStringValue)
}

func ConvertValuesToProtos(values []*Value, enumTitle string) string {
	if len(values) == 0 {
		return ""
	}

	result := ""

	for i, value := range values {
		if i != 0 {
			result = result + value.ToProto(enumTitle)
		}
	}

	result = result + values[0].ToProto(enumTitle)

	return result
}

func ConvertValuesToProtoToEnum(values []*Value, enumTitle string) string {
	if len(values) == 0 {
		return ""
	}

	result := ""

	for i, value := range values {
		if i != 0 {
			result = result + value.ProtoToEnum(enumTitle)
		}
	}

	result = result + values[0].ProtoToEnum(enumTitle)

	return result
}

func ConvertValuesToStruct(values []*Value, enumTitle string) string {
	result := ""

	for _, value := range values {
		result = result + value.ToStruct(enumTitle)
	}

	return result
}

func AssembleTestData(values []*Value) string {
	result := ""

	for i, value := range values {
		result = result + value.ToTestData()

		// if not the final piece of test data, then Using(...).
		if i != len(values)-1 {
			result = result + "."
		} else {
			// otherwise, it should be Using(...),
			// The comma "," marks as this is the final item.
			result = result + ","
		}
	}

	return result
}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	log.SetFlags(0 | log.Lshortfile)

	bytes, err := os.ReadFile(dir)
	if err != nil {
		fmt.Println(err)
	}

	enums, err := extractEnum(string(bytes))
	if err != nil {
		panic(err)
	}

	// Export the extracted data with the formatted template.
	for _, enum := range enums {
		f, err := os.Create(fmt.Sprintf("./input-enum/%s.go", enum.GetTitle()))
		checkErr(err)
		defer f.Close()

		w := bufio.NewWriter(f)

		result := fmt.Sprintf(`package input_enum
		%s
		%s
		%s
		%s`, enum.ToString(), enum.ToProto(), enum.ProtoToEnum(), enum.GenerateTest())

		_, err = w.WriteString(result)
		checkErr(err)

		w.Flush()
	}
}

// func loggingStuff(v interface{}) {
// 	reqJSON, err := json.MarshalIndent(v, "", "  ")
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("%s\n", reqJSON)
// }

// extractEnum takes the content of the file, and converts the content to a bunch of Enum.
func extractEnum(input string) ([]*Enum, error) {
	// formatRegex will separate the enum title and values based on this pattern:/
	// enum {{Title}} {
	//   {{Content over here}}
	// }
	formatRegex := regexp.MustCompile(`enum\s+(\w+)\s+{([^}]*)}`)

	matches := formatRegex.FindAllStringSubmatch(input, -1)

	enums := make([]*Enum, 0)

	for _, match := range matches {
		// For each title, finds its values.
		enum := &Enum{
			Title: match[1],
		}

		// Extract for the value of the enum
		values, err := extractValues(match[2])
		if err != nil {
			return nil, fmt.Errorf("error occurred while extracking values for title (%s): %w", match[1], err)
		}

		enum.Values = values

		sort.Slice(values, func(i, j int) bool {
			return values[i].NumberValue < values[j].NumberValue
		})

		enums = append(enums, enum)
	}

	return enums, nil
}

// extractValues tries to extract
//
//	// TIME_UNIT_UNSPECIFIED.
//	TIME_UNIT_UNSPECIFIED = 0;
//	// TIME_UNIT_SECOND.
//	TIME_UNIT_SECOND = 1;
//
// to a bunch of Value format.
func extractValues(input string) ([]*Value, error) {
	commentPattern := regexp.MustCompile(`^// (.+)`)
	commentTracker := ""
	// Sometimes, the comment is formed by combining multiple lines. Sometimes, there's no comment.
	// So just to be sure, we'll traverse line by line.
	lines := strings.Split(input, "\n")

	result := make([]*Value, 0)

	for _, line := range lines {
		// Trim the line first, it's usually starts with a bunch of whitespace.
		line = strings.Trim(line, " ")

		if len(line) == 0 {
			continue
		}

		// if the line starts wih '//', it's a comment.
		// we store the comment in the commentTracker then continue to another line
		if commentMatch := commentPattern.MatchString(line); commentMatch {
			if commentTracker == "" {
				commentTracker = line
			} else {
				commentTracker = commentTracker + "\n" + line
			}

			continue
		}

		// otherwise, it's actual content, extract the string and the number value.
		// 1. Extract OriginalStringValue, String and NumberValue
		// 2. Store the commentTracker in a Value object,
		// 3. and reset that commentTracker to start recording the other Value Objects

		// 1. Extract String and NumberValue

		// Remove all the white space and potential ";" at the end.
		line = strings.Replace(line, " ", "", -1)
		line = strings.Replace(line, ";", "", -1)

		contents := strings.Split(line, "=")
		if len(contents) != 2 {
			return nil, fmt.Errorf("content value for enum protobuf is expected to have this format: (StringValue) = (NumberValue), received content: %s",
				line)
		}

		// We go with this format:
		// (StringValue) = (NumberValue), so we'll always expect StringValue to be on the left [0] and NumberValue to be on the right [1].
		// Also, for the getter in go, we don't want all CAPS, we want PascalCase instead.
		originalStringValue := contents[0]
		stringValue := convertSnekToPascalCase(contents[0])
		numberValue, err := strconv.Atoi(contents[1])
		if err != nil {
			return nil, fmt.Errorf("failed to convert int for %s, content value for enum protobuf is expected to have this format: (StringValue) = (NumberValue): %w",
				line, err)
		}

		// 2. Store the commentTracker in a Value object,
		value := &Value{
			OriginalStringValue: originalStringValue,
			Comment:             commentTracker,
			StringValue:         stringValue,
			NumberValue:         numberValue,
		}

		// 3. Reset.
		commentTracker = ""

		result = append(result, value)
	}

	return result, nil
}

// convertSnekToPascalCase converts UPPER_CASE_SNAKE to UpperCaseSnake.
func convertSnekToPascalCase(input string) string {
	// Split the input string into words using underscores
	words := strings.Split(input, "_")

	// Capitalize the first letter of each word
	for i := range words {
		words[i] = caser.String(words[i])
	}

	// Join the words together to form the UpperCamelCase string
	upperCamel := strings.Join(words, "")

	return upperCamel
}
