include mk/util.mk
include mk/golang.mk

.PHONY: all
all: help

.PHONY: build
build: targets

.PHONY: targets
targets: $(TGT_BIN)

.PHONY: clean
clean:
	rm -f $(CLEAN)

.PHONY: test
test: $(TEST)

# -------------------- #
#       sub-files      #
# -------------------- #
dir := cmd/ipfs
include $(dir)/Rules.mk

