static char *itoa_helper(char *dest, int i) {
  if (i <= -10) {
    dest = itoa_helper(dest, i/10);
  }
  *dest++ = '0' - i%10;
  return dest;
}

char *itoa(char *dest, int i) {
  char *s = dest;
  if (i < 0) {
    *s++ = '-';
  } else {
    i = -i;
  }
  *itoa_helper(s, i) = '\0';
  return dest;
}

static unsigned long long seed;
int rand () {
    seed = 6364136223846793005ULL*seed + 1;
    return seed >> 33;
}
