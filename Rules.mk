TGT_BIN:=
CLEAN:=
DISTCLEAN:=
TEST:=

include mk/util.mk
include mk/golang.mk
include mk/gx.mk

# -------------------- #
#       sub-files      #
# -------------------- #
dir := bin
include $(dir)/Rules.mk

dir := test
include $(dir)/Rules.mk

dir := cmd/ipfs
include $(dir)/Rules.mk

dir := namesys/pb
include $(dir)/Rules.mk

dir := unixfs/pb
include $(dir)/Rules.mk

dir := merkledag/pb
include $(dir)/Rules.mk

dir := exchange/bitswap/message/pb
include $(dir)/Rules.mk

dir := diagnostics/pb
include $(dir)/Rules.mk

dir := pin/internal/pb
include $(dir)/Rules.mk
#
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

deps: gx-deps
.PHONY: deps

install: $$(DEPS_GO)
	go install $(go-flags-with-tags) ./cmd/ipfs

# TEMP
coverage: test
