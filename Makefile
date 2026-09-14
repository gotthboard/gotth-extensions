.PHONY: build fmt-check vet test race coverage proto proto-contract consumer verify

build:
	go build ./...

fmt-check:
	@test -z "$$(gofmt -l pkg test)"

vet:
	go vet -mod=readonly ./...

test:
	go test -mod=readonly ./...

race:
	go test -mod=readonly -race ./...

coverage:
	go test -mod=readonly -coverprofile=coverage.out ./pkg/extensions
	go tool cover -func=coverage.out

proto:
	@descriptor="$$(mktemp)"; \
	trap 'rm -f "$$descriptor"' EXIT; \
	protoc --proto_path=. --include_imports --descriptor_set_out="$$descriptor" proto/gotth/extensions/v1/control.proto; \
	test -s "$$descriptor"

proto-contract:
	@if rg -n 'rpc[[:space:]]+(Invoke|Execute|Configure|Start|Stop|Restart)|map<|secret_value|bytes[[:space:]]+payload|google\.protobuf\.(Any|Struct)' proto; then \
		echo "forbidden generic, secret, or process-control field in control schema" >&2; \
		exit 1; \
	fi

consumer:
	cd test/consumer && go test -mod=readonly ./...

verify: fmt-check vet test race coverage proto proto-contract consumer
	git diff --check
