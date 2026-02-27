#!/bin/bash

set -euo pipefail

ENVIRONMENT=$1
SERVICE_NAME=$2
SERVICE_DIR=/opt/$SERVICE_NAME

sudo mkdir -p "$SERVICE_DIR/logs"

UNIT_FILE=/etc/systemd/system/$SERVICE_NAME.service
sudo bash -c "cat > $UNIT_FILE" <<EOF
[Unit]
Description=P2P Management Service
After=network.target

[Service]
Type=simple
Restart=always
User=ec2-user
Environment=ENVIRONMENT=$ENVIRONMENT
WorkingDirectory=$SERVICE_DIR
ExecStart=$SERVICE_DIR/bin/service >> $SERVICE_DIR/logs/service.log 2>&1

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable "$SERVICE_NAME"
sudo systemctl restart "$SERVICE_NAME"

sleep 5
sudo systemctl status "$SERVICE_NAME" --no-pager || true

echo "$SERVICE_NAME deployment to $ENVIRONMENT completed."