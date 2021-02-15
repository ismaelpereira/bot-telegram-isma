build: clean
	go build -o dist/bin/telegram-bot -ldflags="-s -w" main.go
	upx -9 dist/bin/*
	cp -r ./local/config/ ./dist/config/

run: build
	./dist/bin/telegram-bot

clean:
	rm -rf dist/
