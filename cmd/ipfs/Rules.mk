p 				:= $(sp).x
dirstack_$(sp)	:= $(d)
d				:= $(dir)# keep track of dirs


TGTS_$(d) := $(call go-curr-pkg-tgt)

TGT_BIN += $(TGTS_$(d))
CLEAN += $(TGTS_$(d))


$(TGTS_$(d)): $(d) ALWAYS
	$(go-build)

d		:= $(dirstack_$(sp))
sp		:= $(basename $(sp))
