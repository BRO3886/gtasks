windows:
	@echo "Building for windows"
	GOOS=windows GOARCH=386 go build -o ./bin/windows/gtasks.exe
linux:
	@echo "Building for linux"
	go build -o ./bin/linux/gtasks
all:
	@echo "Building for every OS and Platform"
	GOOS=windows GOARCH=386 go build -o ./bin/windows/gtasks.exe
	GOOS=linux GOARCH=386 go build -o ./bin/linux/gtasks
	GOOS=freebsd GOARCH=386 go build -o ./bin/freebsd/gtasks-bsd
	GOOS=darwin GOARCH=amd64 go build -o ./bin/mac/gtasks-mac
run:
	go run .
global:
	go install .
push:
	git add .
	git commit -m "$m"
	git push origin master
release:
	gh release create $v './bin/windows/gtasks.exe' './bin/linux/gtasks' './bin/mac/gtasks-mac'