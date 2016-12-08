TGT_BIN:=
CLEAN:=
TEST:=

include mk/util.mk
include mk/golang.mk

# -------------------- #
#       sub-files      #
# -------------------- #
dir := cmd/ipfs
include $(dir)/Rules.mk


# --- core targets --- #

all: help
.PHONY: all

build: $(TGT_BIN)
.PHONY: build

clean:
	rm -f $(CLEAN)
.PHONY: clean

test: $(TEST)
.PHONY: test


