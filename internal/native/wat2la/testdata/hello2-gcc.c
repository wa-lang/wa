struct ret_t {
    int v0, v1, v2, v3;
};

int add2(int a0, int a1) {
    return a0+a1;
}
struct ret_t add10(int a0, int a1, int a2, int a3, int a4, int a5, int a6, int a7, int a8, int a9) {
    struct ret_t v;
    v.v3 = a0+a9;
    return v;
}

int main() {
    int v0 = add2(1000, 200); 
    struct ret_t x = add10(0, 1, 2, 3, 4, 5, 6, 7, 8, 9);
    return 0;
}
