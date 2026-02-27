.PHONY: build web-build deploy dev clean test release

build: web-build
	go build -o bin/fpv-ground-station ./cmd/fpv-ground-station

web-build:
	cd web-ui && bun install && bun run build
	rm -rf cmd/fpv-ground-station/dist
	cp -r web-ui/dist cmd/fpv-ground-station/dist

deploy: build
	./bin/fpv-ground-station

test:
	go test ./internal/... -v -race

release: web-build
	GOOS=darwin GOARCH=arm64 go build -o bin/fpv-ground-station-darwin-arm64 ./cmd/fpv-ground-station
	GOOS=darwin GOARCH=amd64 go build -o bin/fpv-ground-station-darwin-amd64 ./cmd/fpv-ground-station
	GOOS=linux GOARCH=amd64 go build -o bin/fpv-ground-station-linux-amd64 ./cmd/fpv-ground-station
	GOOS=windows GOARCH=amd64 go build -o bin/fpv-ground-station-windows-amd64.exe ./cmd/fpv-ground-station
