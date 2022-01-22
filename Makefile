windows:
	@echo "Building for windows"
	GOOS=windows GOARCH=386 go build -o ./bin/windows/gtasks.exe
linux:
	@echo "Building for linux"
	go build -o ./bin/linux/gtasks
	@ cd bin/linux
	gtasks
mac:
	@echo "Building for mac"
	GOOS=darwin GOARCH=amd64 go build -o ./bin/mac/gtasks
all:
	@echo "Building for every OS and Platform"
	GOOS=windows GOARCH=386 GO386=softfloat go build -o ./bin/windows/gtasks.exe
	GOOS=linux GOARCH=386 GO386=softfloat go build -o ./bin/linux/gtasks
	GOOS=freebsd GOARCH=386 GO386=softfloat go build -o ./bin/freebsd/gtasks
	GOOS=darwin GOARCH=amd64 go build -o ./bin/mac/gtasks
	GOOS=darwin GOARCH=arm64 go build -o ./bin/m1/gtasks
	@echo "Zipping for release"
	@tar -czf bin/releases/gtasks_linux_$v.tar.gz LICENSE -C bin/linux gtasks
	@tar -czf bin/releases/gtasks_win_$v.tar.gz LICENSE -C  bin/windows gtasks.exe
	@tar -czf bin/releases/gtasks_mac_amd64_$v.tar.gz LICENSE -C bin/mac gtasks 
	@tar -czf bin/releases/gtasks_mac_m1_arm64_$v.tar.gz LICENSE -C bin/m1 gtasks 
	@tar -czf bin/releases/gtasks_bsd_$v.tar.gz LICENSE -C bin/freebsd gtasks

release:
	gh release create $v 'bin/releases/gtasks_linux_$v.tar.gz' 'bin/releases/gtasks_win_$v.tar.gz' 'bin/releases/gtasks_bsd_$v.tar.gz' 'bin/releases/gtasks_mac_amd64_$v.tar.gz' 'bin/releases/gtasks_mac_m1_arm64_$v.tar.gz'