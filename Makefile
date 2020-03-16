PROJECT := github.com/eurielec/bulkerrs

all: update-mod check-license check-go check-doc check
.PHONY: all

check: check-license check-fmt check-doc
	@(echo " -> Go Test")
	go test $(PROJECT)/...

.PHONY: update-mod
update-mod:
	@(echo " -> Go Mod")
	@(cd $(GOPATH)/src/$(PROJECT) && dep ensure -update)
	@(go mod vendor -v)
	@(go mod tidy -v)

.PHONY: check-license
check-license:
	@(echo " -> Go Check License")
	@(grep -rl "Licensed under the GPLv3" --include="*.go" --exclude-dir="./vendor" .;\
		find ./ -path "./vendor/*" -prune -o -name "*.go" -print) | sed -e 's,\./,,' | sort | uniq -u | \
		xargs -I {} echo FAIL: licence missed: {}

.PHONY: check-fmt
check-fmt:
	@(echo " -> Go Check Format")
	$(eval GOFMT := $(strip $(shell gofmt -l .| grep -v "^vendor/" | sed -e "s/^//g")))
	@(if [ "x$(GOFMT)" != "x" ]; then \
		echo "  detected wrongly formatted files: $(GOFMT)"; \
		echo '  please run "make go-fmt"'; \
		exit 1; \
	fi)
	@(go vet -all -composites=false -copylocks=false .)

.PHONY: check-doc
check-doc:
	@(echo " -> Go Check Doc")
	$(eval GODOC := $(shell gomarkdoc -c . 2>&1))
	@(if [ "$(GODOC)" != "" ]; then \
		echo "  Fail! Did you forget to run gomarkdoc?"; \
		echo '  please run "make doc"'; \
		exit 1; \
	else \
	  echo "  OK!"; \
	fi)

.PHONY: go-fmt
go-fmt:
	@(echo " -> Go Format")
	@(find ./ -path "./vendor/*" -prune -o -name "*.go" -exec gofmt -l -w {} \;) | \
		sed -e "s/^./  - fixed: /g"

.PHONY: doc
doc:
	@(echo " -> Go Doc")
	@(gomarkdoc .)

.PHONY: tag
tag:
	@(echo " -> Git make tag")
	@(git tag "v$(shell cat VERSION)")
