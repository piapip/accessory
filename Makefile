test:
	go build main.go
	./main -type Tester -output result.go ./input
