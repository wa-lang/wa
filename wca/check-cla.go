/* ！！！仔细阅读本注释！！！
本文件记录了所有签署过 WCA 协议的邮箱地址；未签署 WCA 的邮箱发起的 PR 无法被合并。
*/

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var wcaMap = map[string]bool{
	"406134592@qq.com": true,
	"476582@qq.com":    true,
	"50434393+xxxDeveloper@users.noreply.github.com": true,
	"547551933@qq.com": true,
	"64215+codefromthecrypt@users.noreply.github.com": true,
	"704566072@qq.com":            true,
	"ben.shi@streamcomputing.com": true,
	"chaishushan@gmail.com":       true,
	"humengmingxx@gmail.com":      true,
	"imcusg@gmail.com":            true,
	"powerman1st@163.com":         true,
	"visus@qq.com":                true,
	"wuxuan.ecios@gmail.com":      true,
	"zhanluxianshen@163.com":      true,
	"dragon-fly@qq.com":           true,
}

func main() {
	gitLostEmails := getGitLogEmails()

	var sortedNames []string
	for s := range gitLostEmails {
		sortedNames = append(sortedNames, s)
	}
	for _, s := range sortedNames {
		if wcaMap[s] {
			fmt.Println("check cla:", s, "ok")
		} else {
			fmt.Println("check cla:", s, "failed")
			os.Exit(1)
		}
	}
	if len(gitLostEmails) != len(wcaMap) {
		for k := range wcaMap {
			if !gitLostEmails[k] {
				fmt.Println("check cla:", k, "missing")
			}
		}
		os.Exit(1)
	}
}

// git log --pretty=format:"%ae"
func getGitLogEmails() map[string]bool {
	cmd := exec.Command("git", "log", `--pretty=format:"%ae"`)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}

	m := make(map[string]bool)
	for _, s := range strings.Split(string(stdoutStderr), "\n") {
		email := strings.TrimSpace(s)
		if !strings.Contains(email, "@") {
			continue
		}

		email = strings.TrimFunc(email, func(r rune) bool { return r == '"' })
		m[email] = true
	}

	return m
}
