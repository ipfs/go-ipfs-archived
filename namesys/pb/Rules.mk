include mk/header.mk

PB_$(d) = $(wildcard $(d)/*.proto)
TGTS_$(d) = $(PB_$(d):.proto=.pb.go)

GODEPS += $(TGTS_$(d))

include mk/footer.mk
