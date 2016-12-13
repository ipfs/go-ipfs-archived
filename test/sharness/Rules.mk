include mk/header.mk

T_$(d) = $(sort $(wildcard $(d)/t[0-9][0-9][0-9][0-9]-*.sh))

DEPS_$(d) = test/bin/random test/bin/multihash test/bin/pollEndpoint \
	   test/bin/iptb test/bin/go-sleep test/bin/random-files \
	   test/bin/go-timeout test/bin/hang-fds
DEPS_$(d) += cmd/ipfs/ipfs

$(T_$(d)): $(DEPS_$(d))
	@echo "*** $@ ***"
	(cd $(dir $@) && ./$(notdir $@)) 2>&1
.PHONY: $(T_$(d))


sharness: $(T_$(d))
.PHONY: sharness



sharness_deps: $(DEPS_$(d))



include mk/footer.mk
