#!/usr/bin/env bash
echo "border0-example-connector" > /etc/hostname ; hostname -F /etc/hostname
apt-get update && sudo apt-get -y install gpg curl
install -m 0755 -d /etc/apt/keyrings
install -m 0755 -d /etc/sysconfig

echo """
#turning off logs so the API dones't rate limit us
#BORDER0_CONNECTOR_LOG_LEVEL=fatal
#BORDER0_LOG_LEVEL=fatal

#BORDER0_VERY_VERBOSE=true
#BORDER0_LOG_LEVEL=debug
#BORDER0_CONNECTOR_LOG_LEVEL=fatal

# pointing at staging env
BORDER0_TUNNEL=tunnel.staging.border0.com
BORDER0_CONNECTOR_SERVER=capi.staging.border0.com:443
BORDER0_API=https://api.staging.border0.com/api/v1
BORDER0_WEB_URL=https://portal.staging.border0.com
BORDER0_RELAY_URL=wss://relay.staging.border0.com

# restirc comms to turtle only
# BORDER0_ALLOWED_METHODS=udp4
""" > /etc/sysconfig/border0

curl -fsSL https://download.border0.com/deb/gpg | gpg --dearmor -o /etc/apt/keyrings/border0.gpg
echo "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/border0.gpg] https://download.border0.com/deb/ stable main" > /etc/apt/sources.list.d/border0.list
apt-get -y update
BORDER0_TOKEN=${border0_connector_token_path} apt-get -y install border0
border0 --version
#eof