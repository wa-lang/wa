// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

extern int add_f(int a, int b, int c, int d);
extern int write(int fd, void* p, int size);

void printch(char c) {
    write(1, &c, 1);
}

void printint(int n) {
    if(n >= 10) {
        printint(n/10);
    }
    printch((n%10)+'0');
}

int main() {
    int ret = add_f(1, 2, 3, 4);
    printint(ret);
    printch('\n');
    return 0;
}
