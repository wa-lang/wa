struct ret_t {
    int v0, v1, v2, v3;
};

struct ret_t add(int a0, int a1, int a2, int a3, int a4, int a5, int a6, int a7, int a8, int a9) {
    struct ret_t v;
    v.v0 = a0;
    return v;
}

int main() {
    struct ret_t x = add(0, 1, 2, 3, 4, 5, 6, 7, 8, 9);
    return 0;
}
