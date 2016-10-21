#! /bin/bash

for((i=8000;i>=500;i-=500)) do
	output=$i.output
	./gspan $i >$output
done
