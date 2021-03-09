#!/usr/bin/env bash

# Ensure we have git and python
sudo apt-get install -y git python3

# Clone the repo
mkdir ~/git && cd ~/git
git clone https://github.com/vmware-tanzu-labs/yot.git

# Install python app deps
cd ~/git/yot
pip3 install -r requirements.txt

# Create symlink
ln -s ~/git/yot/yot /usr/local/bin/yot

# Ensure executable
chmod +x ~/git/yot/yot

# Add semaphore for foreground script to know env is good to go
echo "done" > /tmp/.env_good
clear
