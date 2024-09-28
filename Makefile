.PHONY: bench test

bench:
	 go test -bench=. -benchmem -count 3 -benchtime 5s

test:
	go test -v -count 1 .