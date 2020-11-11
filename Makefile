windows:
	echo "Building for windows"
	GOOS=windows GOARCH=386 go build -o ./bin/windows/gtasks.exe
linux:
	echo "Building for linux"
	go build -o ./bin/linux/gtasks
all:
	echo "Building for every OS and Platform"
	GOOS=windows GOARCH=386 go build -o ./bin/windows/gtasks.exe
	GOOS=linux GOARCH=386 go build -o ./bin/linux/gtasks
	GOOS=freebsd GOARCH=386 go build -o ./bin/freebsd/gtasks
run:
	go run .
