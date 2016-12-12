include mk/header.mk

T_$(d) = $(sort $(wildcard t[0-9][0-9][0-9][0-9]-*.sh))

DEPS_$(d) = test/bin/random test/bin/multihash test/bin/pollEndpoint \
	   test/bin/iptb test/bin/go-sleep test/bin/random-files \
	   test/bin/go-timeout test/bin/hang-fds
DEPS_$(d) += cmd/ipfs/ipfs

$(T_$(d)): $(DEPS_$(d))
	@echo "*** $@ ***"
	./$@
.PHONY: $(T_$(d))




sharness_deps: $(DEPS_$(d))




$(d)/chain.mk:



include mk/footer.mk
