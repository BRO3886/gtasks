# Load environment variables from .env file (optional)
-include .env
export

# Version (can be overridden with: make release v=v1.2.3)
v ?= v0.12.0

# Set EMBED_CREDS=1 to embed credentials in binary
# Example: make dev EMBED_CREDS=1
ifdef EMBED_CREDS
ifdef GTASKS_CLIENT_ID
ifdef GTASKS_CLIENT_SECRET
LDFLAGS = -ldflags "-X github.com/BRO3886/gtasks/internal/config.ClientID=$(GTASKS_CLIENT_ID) -X github.com/BRO3886/gtasks/internal/config.ClientSecret=$(GTASKS_CLIENT_SECRET)"
endif
endif
endif

dev:
	@echo "Building for development"
ifdef LDFLAGS
	@echo "  (with embedded credentials)"
else
	@echo "  (without embedded credentials - use EMBED_CREDS=1 to embed)"
endif
	go build $(LDFLAGS) -o ./gtasks .

windows:
	@echo "Building for windows"
ifdef LDFLAGS
	@echo "  (with embedded credentials)"
endif
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o ./bin/windows-amd64/gtasks.exe

linux:
	@echo "Building for linux"
ifdef LDFLAGS
	@echo "  (with embedded credentials)"
endif
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o ./bin/linux-amd64/gtasks
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o ./bin/linux-arm64/gtasks

mac:
	@echo "Building for mac"
ifdef LDFLAGS
	@echo "  (with embedded credentials)"
endif
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o ./bin/mac-amd64/gtasks
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o ./bin/mac-arm64/gtasks

all:
	@echo "Building for every OS and Platform"
ifdef LDFLAGS
	@echo "  (with embedded credentials)"
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o ./bin/windows_amd64/gtasks.exe
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o ./bin/linux_amd64/gtasks
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o ./bin/linux_arm64/gtasks
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o ./bin/mac_amd64/gtasks
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o ./bin/mac_arm64/gtasks
else
	@echo "  (without embedded credentials)"
	GOOS=windows GOARCH=amd64 go build -o ./bin/windows_amd64/gtasks.exe
	GOOS=linux GOARCH=amd64 go build -o ./bin/linux_amd64/gtasks
	GOOS=linux GOARCH=arm64 go build -o ./bin/linux_arm64/gtasks
	GOOS=darwin GOARCH=amd64 go build -o ./bin/mac_amd64/gtasks
	GOOS=darwin GOARCH=arm64 go build -o ./bin/mac_arm64/gtasks
endif
	@echo "Zipping for release"
	@mkdir -p bin/releases
	@tar -czf bin/releases/gtasks_linux_amd64_$(v).tar.gz LICENSE -C bin/linux_amd64 gtasks
	@tar -czf bin/releases/gtasks_linux_arm64_$(v).tar.gz LICENSE -C bin/linux_arm64 gtasks
	@tar -czf bin/releases/gtasks_win_$(v).tar.gz LICENSE -C bin/windows_amd64 gtasks.exe
	@tar -czf bin/releases/gtasks_mac_amd64_$(v).tar.gz LICENSE -C bin/mac_amd64 gtasks
	@tar -czf bin/releases/gtasks_mac_arm64_$(v).tar.gz LICENSE -C bin/mac_arm64 gtasks

release:
	gh release create $(v) 'bin/releases/gtasks_linux_amd64_$(v).tar.gz' 'bin/releases/gtasks_linux_arm64_$(v).tar.gz' 'bin/releases/gtasks_win_$(v).tar.gz' 'bin/releases/gtasks_mac_amd64_$(v).tar.gz' 'bin/releases/gtasks_mac_arm64_$(v).tar.gz'

clean:
	rm -rf bin/ gtasks

.PHONY: dev windows linux mac all release clean
