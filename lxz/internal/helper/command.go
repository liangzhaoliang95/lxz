/**
 * @author  zhaoliang.liang
 * @date  2025/7/31 10:31
 */

package helper

import (
	"fmt"
	"os"
	"os/exec"
)

func Command(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Env = append(os.Environ(), "LC_ALL=C") // 设置环境变量，避免中文乱码
	//cmd.Args = append(cmd.Args, arg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("cmd %s,args %s, call failed %s", command, args, err.Error())
	}
	return fmt.Sprintf("cmd %s,args %s, call success", command, args), nil
}
