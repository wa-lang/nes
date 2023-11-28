
GO_FILES=$(wildcard ./*.go)
WA_FILES=$(patsubst %.go,./output/%.wa,$(GO_FILES))

default:
	make clean
	mkdir -p output
	make all

all: $(WA_FILES)
	@echo "ok"


clean:
	-rm -rf ./output

./output/%.wa: %.go
	wa go2wa $< >  $@
