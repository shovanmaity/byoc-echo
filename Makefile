.PHONY: dep
dep:
	go mod tidy

.PHONY: bin
bin: dep
	mkdir -p bin
	rm -rf bin/app
	CGO_ENABLED=0 go build -o bin/app .