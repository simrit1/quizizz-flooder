build:
	go build -o bin/main main.go

compile:
	echo "Compiling for every OS and Platform"
	GOOS=freebsd GOARCH=386 go build -o bin/main-freebsd-386 .
	GOOS=linux GOARCH=386 go build -o bin/main-linux-386 .
	GOOS=windows GOARCH=386 go build -o bin/main-windows-386 .

run:
	go run .