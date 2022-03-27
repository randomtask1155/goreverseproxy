#!/usr/bin/bash

export BACKEND_ADDRESS=ip:port
export LISTEN_PORT=8080 
nohup /home/pi/bin/goreverseproxy &> /tmp/proxy.log &
