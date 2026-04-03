fun fact(n) {
  if n < 2 {
    n = 1;
  } else {
    n = n * fact(n - 1);
  }
  return n;
}

main {
  return fact(5);
}
