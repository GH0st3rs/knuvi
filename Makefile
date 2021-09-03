BUILD=go build
FLAGS=--ldflags="-s -w"

build: clean
	$(BUILD) $(FLAGS) .
mips:
	GOARCH=mipsle GOOS=linux $(BUILD) $(FLAGS) -o knuvi_mips .
clean:
	rm knuvi_*
