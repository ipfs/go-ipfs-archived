include mk/header.mk
T_$(d) = $(sort $(wildcard t[0-9][0-9][0-9][0-9]-*.sh))
BINS_$(d) = test/bin/random test/bin/multihash test/bin/pollEndpoint \
	   test/bin/iptb test/bin/go-sleep test/bin/random-files \
	   test/bin/go-timeout test/bin/hang-fds
BINS_$(d) += cmd/ipfs/ipfs 

sharness: $(BINS_$(d))
include mk/footer.mk
