TESTS := ${wildcard ./*_test.go}
.PHONY: test 
test: $(TESTS)
./*_test.go:
	go test $<
giv: 
	go test -v giv/givsoft

