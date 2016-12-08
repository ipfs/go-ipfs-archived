# golang utilities

go-pkg-name=$(shell go list ./$(1))
go-main-name=$(notdir $(call go-pkg-name,$(1)))$(exe?)
go-curr-pkg-tgt=$(d)/$(call go-main-name,$(d))

go-flags=$(if $(GOTAGS), -tags $(call join-with,$(GOTAGS),$(comma)))

define go-build=
go build $(GOFLAGS)$(call go-flags) -o "$@" "$(call go-pkg-name,$<)"
endef



