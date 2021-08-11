#!/usr/bin/env bash
clear
while [ ! -f /tmp/yot-installed ]
do
  clear
  sleep .2
done
clear && echo "YAML Overlay Tool lab environment successfully configured"
