echo build for windows/386...
GOOS=windows GOARCH=386 go build -o ./bin/subscribe_windows_386.exe -a  -ldflags '-w -s'
echo build for windows/amd64...
GOOS=windows GOARCH=amd64 go build -o ./bin/subscribe_windows_amd64.exe -a  -ldflags '-w -s'
echo build for windows/arm...
GOOS=windows GOARCH=arm go build -o ./bin/subscribe_windows_arm.exe -a  -ldflags '-w -s'

echo build for linux/386...
GOOS=linux GOARCH=386 go build -o ./bin/subscribe_linux_386 -a  -ldflags '-w -s'
echo build for linux/amd64...
GOOS=linux GOARCH=amd64 go build -o ./bin/subscribe_linux_amd64 -a  -ldflags '-w -s'
echo build for linux/arm...
GOOS=linux GOARCH=arm go build -o ./bin/subscribe_linux_arm -a  -ldflags '-w -s'
echo build for linux/arm64...
GOOS=linux GOARCH=arm64 go build -o ./bin/subscribe_linux_arm64 -a  -ldflags '-w -s'
echo build for linux/ppc64...
GOOS=linux GOARCH=ppc64 go build -o ./bin/subscribe_linux_ppc64 -a  -ldflags '-w -s'
echo build for linux/ppc64le...
GOOS=linux GOARCH=ppc64le go build -o ./bin/subscribe_linux_ppc64le -a  -ldflags '-w -s'
echo build for linux/mips...
GOOS=linux GOARCH=mips go build -o ./bin/subscribe_linux_mips -a  -ldflags '-w -s'
echo build for linux/mipsle...
GOOS=linux GOARCH=mipsle go build -o ./bin/subscribe_linux_mipsle -a  -ldflags '-w -s'
echo build for linux/mips64...
GOOS=linux GOARCH=mips64 go build -o ./bin/subscribe_linux_mips64 -a  -ldflags '-w -s'
echo build for linux/mips64le...
GOOS=linux GOARCH=mips64le go build -o ./bin/subscribe_linux_mips64le -a  -ldflags '-w -s'
echo build for linux/riscv64...
GOOS=linux GOARCH=riscv64 go build -o ./bin/subscribe_linux_riscv64 -a  -ldflags '-w -s'
echo build for linux/s390x...
GOOS=linux GOARCH=s390x go build -o ./bin/subscribe_linux_s390x -a  -ldflags '-w -s'

echo build for darwin/386...
GOOS=darwin GOARCH=386 go build -o ./bin/subscribe_darwin_386 -a  -ldflags '-w -s'
echo build for darwin/amd64...
GOOS=darwin GOARCH=amd64 go build -o ./bin/subscribe_darwin_amd64 -a  -ldflags '-w -s'

echo build for aix/ppc64...
GOOS=aix GOARCH=ppc64 go build -o ./bin/subscribe_aix_ppc64 -a  -ldflags '-w -s'

echo build for dragonfly/amd64...
GOOS=dragonfly GOARCH=amd64 go build -o ./bin/subscribe_dragonfly_amd64 -a  -ldflags '-w -s'

echo build for freebsd/386...
GOOS=freebsd GOARCH=386 go build -o ./bin/subscribe_freebsd_386 -a  -ldflags '-w -s'
echo build for freebsd/amd64...
GOOS=freebsd GOARCH=amd64 go build -o ./bin/subscribe_freebsd_amd64 -a  -ldflags '-w -s'
echo build for freebsd/arm...
GOOS=freebsd GOARCH=arm go build -o ./bin/subscribe_freebsd_arm -a  -ldflags '-w -s'
echo build for freebsd/arm64...
GOOS=freebsd GOARCH=arm64 go build -o ./bin/subscribe_freebsd_arm64 -a  -ldflags '-w -s'

echo build for illumos/amd64...
GOOS=illumos GOARCH=amd64 go build -o ./bin/subscribe_illumos_amd64 -a  -ldflags '-w -s'

echo build for js/wasm...
GOOS=js GOARCH=wasm go build -o ./bin/subscribe_js_wasm -a  -ldflags '-w -s'

echo build for netbsd/386...
GOOS=netbsd GOARCH=386 go build -o ./bin/subscribe_netbsd_386 -a  -ldflags '-w -s'
echo build for netbsd/amd64...
GOOS=netbsd GOARCH=amd64 go build -o ./bin/subscribe_netbsd_amd64 -a  -ldflags '-w -s'
echo build for netbsd/arm...
GOOS=netbsd GOARCH=arm go build -o ./bin/subscribe_netbsd_arm -a  -ldflags '-w -s'
echo build for netbsd/arm64...
GOOS=netbsd GOARCH=arm64 go build -o ./bin/subscribe_netbsd_arm64 -a  -ldflags '-w -s'

echo build for openbsd/386...
GOOS=openbsd GOARCH=386 go build -o ./bin/subscribe_openbsd_386 -a  -ldflags '-w -s'
echo build for openbsd/amd64...
GOOS=openbsd GOARCH=amd64 go build -o ./bin/subscribe_openbsd_amd64 -a  -ldflags '-w -s'
echo build for openbsd/arm...
GOOS=openbsd GOARCH=arm go build -o ./bin/subscribe_openbsd_arm -a  -ldflags '-w -s'
echo build for openbsd/arm64...
GOOS=openbsd GOARCH=arm64 go build -o ./bin/subscribe_openbsd_arm64 -a  -ldflags '-w -s'

echo build for plan9/386...
GOOS=plan9 GOARCH=386 go build -o ./bin/subscribe_plan9_386 -a  -ldflags '-w -s'
echo build for plan9/amd64...
GOOS=plan9 GOARCH=amd64 go build -o ./bin/subscribe_plan9_amd64 -a  -ldflags '-w -s'
echo build for plan9/arm...
GOOS=plan9 GOARCH=arm go build -o ./bin/subscribe_plan9_arm -a  -ldflags '-w -s'

echo build for solaris/amd64...
GOOS=solaris GOARCH=amd64 go build -o ./bin/subscribe_solaris_amd64 -a  -ldflags '-w -s'

