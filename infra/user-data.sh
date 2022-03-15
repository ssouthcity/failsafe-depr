#!/bin/bash

sudo snap install --classic go
sudo snap install doctl
sudo snap install jq

mkdir /etc/failsafe
echo '${config_file_content}' > /etc/failsafe/config.json
echo '${service_file_content}' > /etc/systemd/system/failsafe.service
