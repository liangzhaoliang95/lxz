/**
 * @author  zhaoliang.liang
 * @date  2025/7/31 10:35
 */

package helper

type GitHelper struct{}

// 获取当前文件夹下的子模块分支
func (g *GitHelper) GetSubmoduleBranches() ([]string, error) {
	// 这里应该调用 git 命令获取子模块分支信息
	// 例如使用 "git submodule foreach 'git rev-parse --abbrev-ref HEAD'" 命令
	// 返回一个包含所有子模块分支名称的字符串切片
	return nil, nil // TODO: 实现具体逻辑
}
