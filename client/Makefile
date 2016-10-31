BUILD_DIR ?= public
BIN_DIR ?= node_modules/.bin
BINDATA_TARGET = embedded.go
ENV ?= production

dev:
	npm start

clean:
	rm -rf $(BUILD_DIR)/*.{js,css}

dist: clean
	NODE_ENV=$(ENV) $(BIN_DIR)/webpack --progress

serve: dist
	$(BIN_DIR)/static $(BUILD_DIR) --spa

embed: dist
	go-bindata -prefix "public" -o $(BINDATA_TARGET) -pkg client public/...
	gofmt -s -w $(BINDATA_TARGET)
	goimports -w $(BINDATA_TARGET)

.PHONY: dist dev clean serve embed
