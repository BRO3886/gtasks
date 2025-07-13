# Load environment variables from .env file
include .env
export

# Version (can be overridden with: make release v=v1.2.3)
v ?= v0.10.0

# Build flags with credentials
LDFLAGS = -X github.com/BRO3886/gtasks/internal/config.ClientID=$(GTASKS_CLIENT_ID) -X github.com/BRO3886/gtasks/internal/config.ClientSecret=$(GTASKS_CLIENT_SECRET)

dev:
	@echo "Building for development"
	go build -ldflags "$(LDFLAGS)" -o ./gtasks .

windows:
	@echo "Building for windows"
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o ./bin/windows-amd64/gtasks.exe
linux:
	@echo "Building for linux"
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o ./bin/linux-amd64/gtasks
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o ./bin/linux-arm64/gtasks
	@ cd bin/linux
	gtasks
mac:
	@echo "Building for mac"
	GOOS=darwin GOARCH=amd64 go build -o ./bin/mac-amd64/gtasks
	GOOS=darwin GOARCH=arm64 go build -o ./bin/mac-arm64/gtasks
all:
	@echo "Building for every OS and Platform"
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o ./bin/windows_amd64/gtasks.exe
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o ./bin/linux_amd64/gtasks
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o ./bin/linux_arm64/gtasks
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o ./bin/mac_amd64/gtasks
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o ./bin/mac_arm64/gtasks
	@echo "Zipping for release"
	@mkdir -p bin/releases
	@tar -czf bin/releases/gtasks_linux_amd64_$(v).tar.gz LICENSE -C bin/linux_amd64 gtasks
	@tar -czf bin/releases/gtasks_linux_arm64_$(v).tar.gz LICENSE -C bin/linux_arm64 gtasks
	@tar -czf bin/releases/gtasks_win_$(v).tar.gz LICENSE -C bin/windows_amd64 gtasks.exe
	@tar -czf bin/releases/gtasks_mac_amd64_$(v).tar.gz LICENSE -C bin/mac_amd64 gtasks 
	@tar -czf bin/releases/gtasks_mac_arm64_$(v).tar.gz LICENSE -C bin/mac_arm64 gtasks 

release:
	gh release create $v 'bin/releases/gtasks_linux_amd64_$v.tar.gz' 'bin/releases/gtasks_linux_arm64_$v.tar.gz' 'bin/releases/gtasks_win_$v.tar.gz' 'bin/releases/gtasks_mac_amd64_$v.tar.gz' 'bin/releases/gtasks_mac_arm64_$v.tar.gz'