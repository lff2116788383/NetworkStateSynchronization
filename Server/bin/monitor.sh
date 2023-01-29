#!/bin/sh
ulimit -c unlimited
./monitor ./main.pid main  -d
