gx-deps: bin/gx bin/gx-go $(CHECK_GO)
	gx install --global >/dev/null 2>&1

DEPS_GO += gx-deps
