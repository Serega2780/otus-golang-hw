#!/usr/bin/env bash

./test_ctrl_d.sh < command.sh  2> /tmp/out.txt

if grep -q "...EOF" /tmp/out.txt; then
    echo FOUND
else
    echo NOT FOUND
fi

rm -f /tmp/out.txt
