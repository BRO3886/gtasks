windows:
	@echo "Building for windows"
	GOOS=windows GOARCH=386 go build -o ./bin/windows/gtasks.exe
linux:
	@echo "Building for linux"
	go build -o ./bin/linux/gtasks
	@ cd bin/linux
	gtasks
all:
	@echo "Building for every OS and Platform"
	GOOS=windows GOARCH=386 go build -o ./bin/windows/gtasks.exe
	GOOS=linux GOARCH=386 go build -o ./bin/linux/gtasks
	GOOS=freebsd GOARCH=386 go build -o ./bin/freebsd/gtasks
	GOOS=darwin GOARCH=amd64 go build -o ./bin/mac/gtasks
	@echo "Zipping for release"
	@tar -czf bin/releases/gtasks_linux.tar.gz LICENSE -C bin/linux gtasks
	@tar -czf bin/releases/gtasks_win.tar.gz LICENSE -C  bin/windows gtasks.exe
	@tar -czf bin/releases/gtasks_mac_amd64.tar.gz LICENSE -C bin/mac gtasks 
	@tar -czf bin/releases/gtasks_bsd.tar.gz LICENSE -C bin/freebsd gtasks
run:
	go run .
global:
	go install .
push:
	git add .
	git commit -m "$m"
	git push origin master
release:
	gh release create $v 'bin/releases/gtasks_linux.tar.gz' 'bin/releases/gtasks_win.tar.gz' 'bin/releases/gtasks_bsd.tar.gz' 'bin/releases/gtasks_mac_amd64.tar.gz' 