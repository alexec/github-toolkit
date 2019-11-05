build: ght
ght:
	GO111MODULE=on go install ./cmd/ght
test:
	GO111MODULE=on go test `go list ./...`