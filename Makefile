.PHONY: test
test:
	ginkgo ./...

.PHONY: gen
gen:
	buf generate