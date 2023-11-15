#!/usr/bin/bash

export BACKEND_ADDRESS=ip:port
export LISTEN_PORT=8080 
export SERVER_CERT="cert.crt"
export SERVER_KEY="cert.key"
nohup /home/pi/bin/goreverseproxy &> /tmp/proxy.log &
