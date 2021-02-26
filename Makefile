build: clean
	go build -o dist/bin/telegram-bot -ldflags="-s -w" main.go
	cp -r ./local/config/ ./dist/config/

run: build
	./dist/bin/telegram-bot

clean:
	rm -rf dist/
