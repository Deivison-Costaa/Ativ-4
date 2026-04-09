fun dup(x) {
  return x + x;
}

fun quad(x) {
  return dup(dup(x));
}

main {
  return quad(7);
}
