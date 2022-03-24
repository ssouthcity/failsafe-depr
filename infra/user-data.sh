#!/bin/bash

mkdir /etc/failsafe
echo '${failsafe_config_json}' > /etc/failsafe/config.json
echo '${failsafe_service_content}' > /etc/systemd/system/failsafe.service

sudo systemctl enable failsafe.service