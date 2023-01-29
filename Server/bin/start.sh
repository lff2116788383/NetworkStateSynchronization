#!/bin/sh
ulimit -c unlimited
./Server ../config/config.json -d
