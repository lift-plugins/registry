NAME := lift-registry
VERSION := v1.0.0
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.AppName=$(NAME)"
BLDTAGS := -tags "bleve"

test:
	go test $(BLDTAGS) -parallel 2 ./...

testrace:
	go test $(BLDTAGS) -parallel 2 -race ./...

generate:
	go generate ui/webapp.go

build:
	go build $(BLDTAGS) $(LDFLAGS)

clean:
	go clean $(BLDTAGS) $(LDFLAGS)

compile: swagger
	@rm -rf build/
	@gox $(BLDTAGS) $(LDFLAGS) \
	-os="darwin" \
	-os="linux" \
	-output "build/{{.Dir}}_$(VERSION)_{{.OS}}_{{.Arch}}/$(NAME)" \
	./...

install:
	go install $(BLDTAGS) $(LDFLAGS)

dist: compile
	$(eval FILES := $(shell ls build))
	@rm -rf dist && mkdir dist
	@for f in $(FILES); do \
		(cd $(shell pwd)/build/$$f && tar -cvzf ../../dist/$$f.tar.gz *); \
		(cd $(shell pwd)/dist && shasum -a 512 $$f.tar.gz > $$f.sha512); \
		echo $$f; \
	done

swagger:
	cp $(GOPATH)/src/github.com/hooklift/apis/browser/lift-registry.swagger.json ui/public/lib/api.swagger.json

release: test dist
	@latest_tag=$$(git describe --tags `git rev-list --tags --max-count=1`); \
	comparison="$$latest_tag..HEAD"; \
	if [ -z "$$latest_tag" ]; then comparison=""; fi; \
	changelog=$$(git log $$comparison --oneline --no-merges --reverse); \
	github-release hooklift/$(NAME) $(VERSION) "$$(git rev-parse --abbrev-ref HEAD)" "**Changelog**<br/>$$changelog" 'dist/*'; \
	git pull

devcerts:
	openssl ecparam -genkey -name secp384r1 -out certs/server-key.pem && \
	openssl req -new -x509 -key certs/server-key.pem -out certs/server.pem -days 90


.PHONY: build compile protoc install deps dist release
