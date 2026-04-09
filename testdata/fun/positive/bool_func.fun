fun nonneg(x) {
  return not x < 0;
}

main {
  return nonneg(5) and not nonneg(0 - 1);
}
