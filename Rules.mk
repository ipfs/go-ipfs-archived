TGT_BIN:=
CLEAN:=
DISTCLEAN:=
TEST:=

include mk/util.mk
include mk/golang.mk

# -------------------- #
#       sub-files      #
# -------------------- #
dir := bin
include $(dir)/Rules.mk

dir := cmd/ipfs
include $(dir)/Rules.mk

dir := namesys/pb
include $(dir)/Rules.mk

# -------------------- #
#   universal rules    #
# -------------------- #

%.pb.go: %.proto
	$(PROTOC)

# -------------------- #
#     core targets     #
# -------------------- #

all: help
.PHONY: all

build: $(TGT_BIN)
.PHONY: build

clean:
	rm -f $(CLEAN)
.PHONY: clean

distclean: clean
	rm -f $(DISTCLEAN)
.PHONY: distclean

test: $(TEST)
.PHONY: test


