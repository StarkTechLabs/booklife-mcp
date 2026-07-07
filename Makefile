BINARY := booklife
DIST   := dist
CMD    := ./cmd/booklife

.PHONY: build clean

build:
	mkdir -p $(DIST)
	cd booklife-mcp && go build -o ../$(DIST)/$(BINARY) $(CMD)

clean:
	rm -rf $(DIST)
