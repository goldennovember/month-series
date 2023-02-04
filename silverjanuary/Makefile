postgres:
	docker compose up -d

python:
	python -m venv venv
	source venv/bin/activate
	pip install -r requirements.txt

go:
	go mod init github.com/silverjanuary/silverjanuary
	go mod download
	go mod tidy

binding:
	mkdir lib
	go build -o ./lib/_silverjanuary.dll -buildmode=c-shared ./_silverjanuary.go

.PHONY.: postgres python go binding
