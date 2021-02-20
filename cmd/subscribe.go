// //go:generate goversioninfo -icon=resource/icon.ico -manifest=resource/goversioninfo.exe.manifest -64
// https://www.freebuf.com/sectool/246285.html
// https://docs.microsoft.com/zh-cn/dotnet/framework/tools/signtool-exe
// https://blog.csdn.net/wangshubo1989/article/details/50849914
package main

import (
	"fmt"
	"os"

	"subscribe"
)

func main() {
	fmt.Println(os.UserConfigDir())
	fmt.Println(os.UserCacheDir())
	fmt.Println()

	return
	subscribe.Start()
}
