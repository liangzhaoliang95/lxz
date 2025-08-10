/**
 * @author  zhaoliang.liang
 * @date  2025/7/23 11:29
 */

package main

import (
	"flag"
	"github.com/liangzhaoliang95/lxz/cmd"
	"k8s.io/klog/v2"
)

func init() {
	klog.InitFlags(nil)

	if err := flag.Set("logtostderr", "false"); err != nil {
		panic(err)
	}
	if err := flag.Set("alsologtostderr", "false"); err != nil {
		panic(err)
	}
	if err := flag.Set("stderrthreshold", "fatal"); err != nil {
		panic(err)
	}
	if err := flag.Set("v", "0"); err != nil {
		panic(err)
	}
}

func main() {
	cmd.Execute()
}
