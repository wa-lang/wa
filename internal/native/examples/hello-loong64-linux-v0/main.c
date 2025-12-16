// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

extern void printchar(int c);
extern void printstring(const char* s, int len);

static const char msg[] = "Hello, LoongArch 64!\n";

int main() {
    printstring(msg, sizeof(msg));
    printchar('a');
    printchar('b');
    printchar('c');
    printchar('\n');
    return 0;
}
