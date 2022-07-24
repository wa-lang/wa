// 版权 @2021 凹语言 作者。保留所有权利。

#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// ----------------------------------------------------------------------------
// bool
// ----------------------------------------------------------------------------

typedef uint8_t ugo_bool_t;

// ----------------------------------------------------------------------------
// ugo_string_t
// ----------------------------------------------------------------------------

typedef struct {
	char* data;
	int size;
} ugo_string_t;

ugo_string_t* ugo_string_new(int size, const char* data) {
	ugo_string_t* p = malloc(sizeof(*p)+size+1);
	p->data = (char*)(p+1);
	p->size = size;

	if(data != NULL) {
		memcpy(p->data, data, size);
		p->data[size] = '\0';
	} else {
		memset(p->data, 0, size+1);
	}

	return p;
}

ugo_string_t* ugo_string_clone(ugo_string_t* s) {
	ugo_string_t* p = malloc(sizeof(*p)+s->size+1);
	memcpy(p->data, s->data, s->size);
	p->data[s->size] = '\0';
	p->size = s->size;
	return p;
}

void ugo_string_free(ugo_string_t* s) {
	free(s);
}

// ----------------------------------------------------------------------------
// cstring
// ----------------------------------------------------------------------------

char* ugo_cstring_join(const char* s1, const char* s2) {
	int len1 = strlen(s1);
	int len2 = strlen(s2);
	char* dst = malloc(len1+len2+1);
	strcpy(dst, s1);
	strcat(dst, s2);
	return dst;
}

char* ugo_cstring_slice(const char* s, int low, int high) {
	int src_len = strlen(s);
	int dst_len = 0;
	char* dst = NULL;

	if(low < 0) { low = 0; }
	if(high < 0) { high = src_len; }

	if(low >= src_len) { low = src_len-1; }
	if(high > src_len) { low = src_len; }

	dst_len = high-low;
	dst = malloc(dst_len+1);

	memcpy(dst, s+low, dst_len);
	dst[dst_len] = '\0';

	return dst;
}

int ugo_cstring_index(const char* s, int idx) {
	return s[idx];
}

int ugo_cstring_cmp(const char* s1, const char* s2) {
	return strcmp(s1, s2);
}

// ----------------------------------------------------------------------------
// ugo_print_xxx
// ----------------------------------------------------------------------------

int ugo_print_rune(int ch) {
	return printf("%c", ch);
}
int ugo_print_cstring(const char* s) {
	return printf("%s", s);
}
int ugo_print_cstring_len(const char* s, int n) {
	while(n-- > 0) {
		printf("%c", *s++);
	}
	return 1;
}
int ugo_print_string(const ugo_string_t* s) {
	return printf("%.*s", s->size, s->data);
}

int ugo_print_bool(ugo_bool_t x) {
	return printf("%s", x? "true": "false");
}
int ugo_print_int(int x) {
	return printf("%d", x);
}
int ugo_print_int64(int64_t x) {
	return printf("%lld", x);
}
int ugo_print_ptr(void *p) {
	return printf("0x%x", p, p);
}

int ugo_printf(const char* format, ...) {
	va_list args;
	int n;

	va_start(args, format);
	n = vprintf(format, args);
	va_end(args);
	return n;
}

// ----------------------------------------------------------------------------
// builtin
// ----------------------------------------------------------------------------

int ugo_builtin_println(int x) {
	return printf("%d\n", x);
}

int ugo_builtin_exit(int x) {
	exit(x);
	return 0;
}

// ----------------------------------------------------------------------------
// END
// ----------------------------------------------------------------------------
