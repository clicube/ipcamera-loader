#!/bin/bash

cp ipcamera-loaderd.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable ipcamera-loaderd
systemctl start ipcamera-loaderd
