pwd
ls
ls /build/
ls /build/internal/
goreleaser --skip-validate --rm-dist --debug --config ./.goreleaser.yml
