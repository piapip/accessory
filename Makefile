# Example:
#   make test STRUCT_NAME=RetentionAvoid
test:
	go build main.go
	./main -type ${STRUCT_NAME} -output result.go ./input
