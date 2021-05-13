#!/usr/bin/env bash

clear; echo "Waiting for environment to be setup..."; while [ ! -f /tmp/.env_good ]; do sleep 5; done; echo 'Environment is now configured for the course!';
