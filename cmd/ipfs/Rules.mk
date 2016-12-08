p 				:= $(sp).x
dirstack_$(sp)	:= $(d)
d				:= $(dir)# keep track of dirs


TGTS_$(d) := $(call go-curr-pkg-tgt)

TGT_BIN := $(TGT_BIN) $(TGTS_$(d))

$(TGTS_$(d)): $(d) $(GOSRC)
	$(go-build)

d		:= $(dirstack_$(sp))
sp		:= $(basename $(sp))
