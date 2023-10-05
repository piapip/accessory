# Example:
#   make test STRUCT_NAME=DeliverySetting
#  // The error "'expected operand, found '=='" may occur if the receiver is empty
#  Look at the write.go for this
test:
	go build ./generator/struct/main.go
	./main -type ${STRUCT_NAME} -output result_${STRUCT_NAME}.go ./input

# Example:
#   make test-enum
#   Fill in the ./input-enum/input.proto file with proto definition stuff
test-enum:
	go build ./generator/enum/main.go
	./main
