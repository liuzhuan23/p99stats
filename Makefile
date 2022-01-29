#STUB_BIN_NAME := 'libp99stats.so'
STUB_BIN_NAME := 'libp99stats.a'

all:
	@mkdir SO
#go build -o ./SO/$(STUB_BIN_NAME) -buildmode=c-shared
	go build -o ./SO/$(STUB_BIN_NAME) -buildmode=c-archive
	@chmod 755 ./SO/$(STUB_BIN_NAME)
clean:
	@rm -rf ./SO
