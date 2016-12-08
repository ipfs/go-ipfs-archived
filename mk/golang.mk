# golang utilities

GOTAGS:=
GOFLAGS:=
GOTFLAGS:=

go-pkg-name=$(shell go list ./$(1))
go-main-name=$(notdir $(call go-pkg-name,$(1)))$(exe?)
go-curr-pkg-tgt=$(d)/$(call go-main-name,$(d))

go-flag-from-tags=$(if $(GOTAGS), -tags $(call join-with,$(comma),$(GOTAGS)))

define go-build=
go build $(GOFLAGS)$(call go-flag-from-tags) -o "$@" "$(call go-pkg-name,$<)"
endef


