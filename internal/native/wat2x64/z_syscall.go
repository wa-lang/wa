// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2x64

import (
	_ "embed"

	"wa-lang.org/wa/internal/wat/ast"
)

const (
	kSyscallMalloc = "malloc"
)

func (p *wat2X64Worker) checkSyscallSig(spec *ast.ImportSpec) {
	// TODO: 检查系统调用函数签名类型是否匹配
}

const syscallListWin64 = `
# --- 文件 I/O 相关 ---
.extern _open           # int _open(const char *filename, int oflag, ...)
.extern _close          # int _close(int fd)
.extern _read           # int _read(int fd, void *buffer, unsigned int count)
.extern _write          # int _write(int fd, const void *buffer, unsigned int count)
.extern _lseek          # long _lseek(int fd, long offset, int origin)
.extern _isatty         # int _isatty(int fd)
.extern _commit         # int _commit(int fd) (类似 fsync)
.extern _dup            # int _dup(int fd)
.extern _dup2           # int _dup2(int fd1, int fd2)

# --- 文件系统相关 ---
.extern _mkdir          # int _mkdir(const char *dirname)
.extern _rmdir          # int _rmdir(const char *dirname)
.extern _chdir          # int _chdir(const char *dirname)
.extern _getcwd         # char *_getcwd(char *buffer, int maxlen)
.extern _unlink         # int _unlink(const char *filename) (即删除文件)
.extern _stat64         # int _stat64(const char *path, struct __stat64 *buffer)
.extern _fstat64        # int _fstat64(int fd, struct __stat64 *buffer)

# --- 进程与内存相关 ---
.extern _exit           # void _exit(int status) (不刷新缓冲，最接近 syscall exit)
.extern exit            # void exit(int status)  (标准 C 退出)
.extern abort           # void abort(void)
.extern _getpid         # int _getpid(void)
.extern malloc          # void *malloc(size_t size)
.extern free            # void free(void *ptr)
.extern realloc         # void *realloc(void *ptr, size_t size)
.extern calloc          # void *calloc(size_t nmemb, size_t size)
`

const syscallListLinux = `
# --- 文件 I/O 相关 (Linux/POSIX 标准名称) ---
.extern open            # int open(const char *pathname, int flags, ...)
.extern close           # int close(int fd)
.extern read            # ssize_t read(int fd, void *buf, size_t count)
.extern write           # ssize_t write(int fd, const void *buf, size_t count)
.extern lseek           # off_t lseek(int fd, off_t offset, int whence)
.extern isatty          # int isatty(int fd)
.extern fsync           # int fsync(int fd) (对应 Windows 的 _commit)
.extern dup             # int dup(int oldfd)
.extern dup2            # int dup2(int oldfd, int newfd)

# --- 文件系统相关 ---
.extern mkdir           # int mkdir(const char *pathname, mode_t mode)
.extern rmdir           # int rmdir(const char *pathname)
.extern chdir           # int chdir(const char *path)
.extern getcwd          # char *getcwd(char *buf, size_t size)
.extern unlink          # int unlink(const char *pathname) (删除文件)
.extern stat            # int stat(const char *pathname, struct stat *statbuf)
.extern fstat           # int fstat(int fd, struct stat *statbuf)

# --- 进程与内存相关 ---
.extern _exit           # void _exit(int status) (直接终止，不清理 C 缓存)
.extern exit            # void exit(int status)  (标准 C 退出)
.extern abort           # void abort(void)
.extern getpid          # pid_t getpid(void)
.extern fork            # pid_t fork(void) (Windows 没有对应函数)
.extern malloc          # void *malloc(size_t size)
.extern free            # void free(void *ptr)
.extern mmap            # void *mmap(void *addr, size_t length, ...)
`
