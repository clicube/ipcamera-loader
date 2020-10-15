#!/bin/bash -x

npm install
sudo cp ipcamera-loaderd.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable ipcamera-loaderd
sudo systemctl restart ipcamera-loaderd
