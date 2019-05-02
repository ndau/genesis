#! /usr/local/python3

import itertools
a = open("key1.csv").readlines()
b = open("key2.csv").readlines()
result = itertools.product(a, b)

f3 = open("keys.csv", "w")
for a, b in result:
  f3.write(f"{a[:-1]}, {b[:-1]}\n")