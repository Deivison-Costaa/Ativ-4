var x = 10;

main {
  x %= 4;
  if (x <= 2) and (x != 0) {
    x += 3;
  } else {
    x += 100;
  }
  return x >= 5;
}
