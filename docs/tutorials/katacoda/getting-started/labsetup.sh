#!/usr/bin/env bash

echo "Installing YAML Overlay Tool"

curl -sL -o /tmp/yot_amd64.deb https://github.com/vmware-tanzu-labs/yaml-overlay-tool/releases/download/v0.6.3/yot_amd64.deb && sudo dpkg -i /tmp/yot_amd64.deb

echo "YAML Overlay Tool installed successfully"

mkdir ~/source-manifests && mv ~/*.yaml ~/source-manifests/ && mv /tmp/yot.yaml /root/yot.yaml
touch /tmp/yot-installed