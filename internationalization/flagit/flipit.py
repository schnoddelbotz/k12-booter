#!/usr/bin/env python3

import argparse
import sys

# read SysV printerbanner from stdin and rotate, somehow

parser = argparse.ArgumentParser()
parser.add_argument('infile',
                    default=sys.stdin,
                    type=argparse.FileType('r'),
                    nargs='?')
args = parser.parse_args()
data = args.infile.readlines()

# analyse input - get number of rows and longest line
maxLineLength = 0
inputRows = 0
for x in data:
	inputRows += 1
	lineLength = len(x)
	if lineLength > maxLineLength:
		maxLineLength = lineLength

# trim
for x in range(0, inputRows):
	data[x] = data[x].replace('\n', '') # chomp
maxLineLength = maxLineLength -1

# print(f'maxLineLength {maxLineLength}')
# print(f'inputRows {inputRows}')

# pad all lines using whitespace to match maxLineLngth "matrix"
for x in range(0, inputRows):
	diff = maxLineLength - len(data[x])
	if diff > 0:
		data[x] = data[x] + diff * ' '
	# print (f'R{x:3} : {data[x]}')

# init output
output = []
for y in range(0, maxLineLength):
	output.append(list(inputRows * "z"))

# flip!
for x in range(0, inputRows):
	for y in range(0, maxLineLength):
		output[y][x] = data[x][y]

for l in reversed(output):
	for c in l:
		print(c, end='')
	print()
