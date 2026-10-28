# The ACP spec repo and ref the vendored schemas come from. A release tag (e.g.
# SCHEMA_REF=schema-v1.26.0) pins them; main follows the latest published spec.
SPEC ?= https://raw.githubusercontent.com/agentclientprotocol/agent-client-protocol
SCHEMA_REF ?= main

.PHONY: schema test lint

# schema refreshes schema/v1 and schema/v2 from the public spec. types.go follows
# v1 by hand: review the diff and update the types to match.
schema:
	@for v in v1 v2; do mkdir -p schema/$$v; for f in schema.json meta.json; do \
		curl -fsSL $(SPEC)/$(SCHEMA_REF)/schema/$$v/$$f -o schema/$$v/$$f || exit 1; done; done
	@git diff --stat -- schema

test:
	go test -race ./...

lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0 run ./...
