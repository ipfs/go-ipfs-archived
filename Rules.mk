include mk/util.mk
include mk/golang.mk

GOSRC := $(wildcard *.go)

all: targets


.PHONY: targets
targets: $(TGT_BIN)



dir := cmd/ipfs
include $(dir)/Rules.mk

