PROJECT := github.com/eurielec/bulkerrs

all: update-mod check-license check-fmt check
.PHONY: all

check: check-license check-fmt
	@(echo " -> Go Test")
	go test $(PROJECT)/...

.PHONY: update-mod
update-mod:
	@(echo " -> Go Mod")
	@(cd $(GOPATH)/src/$(PROJECT) && dep ensure -update)
	@(go mod vendor)
	@(go mod tidy)

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

.PHONY: go-fmt
go-fmt:
	@(echo " -> Go Format")
	@(find ./ -path "./vendor/*" -prune -o -name "*.go" -exec gofmt -l -w {} \;) | \
		sed -e "s/^./  - fixed: /g"
