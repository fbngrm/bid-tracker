proto_source_files := $(wildcard ./apis/auction/v1/*.proto apis/*.yaml)

gen/proto/go/auction/v1/*.go: $(proto_source_files)
	just buf generate

