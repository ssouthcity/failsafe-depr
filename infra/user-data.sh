#!/bin/bash

sudo snap install --classic go
sudo snap install doctl

# echo
# [Unit]
# Description=Failsafe Discord bot
# 
# [Service]
# ExecStart=/root/go/bin/failsafe
# Restart=on-failure
# 
# [Install]
# WantedBy=multi-user.target