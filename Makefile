.PHONY: build

vendor:
	go mod tidy

gen:
	mkdir -p ./proto/go
	rm -rf ./proto/go/*
	protoc --proto_path=proto/. --twirp_out=proto/go/ --go_out=proto/go/ proto/*.proto

deploy:
	sam build
	sam deploy --no-confirm-changeset --profile me
