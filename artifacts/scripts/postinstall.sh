#!/bin/sh
echo "Executing postinstall script"

dv_user="${DV_USERNAME}"
if [ -z "$dv_user" ]; then
  dv_user="dv"
fi

id -u "$dv_user" &>/dev/null || useradd -r -s /bin/bash -m "$dv_user"

usermod -aG $dv_user dv

if [ -e /home/dv/environment/updater.config.yaml ] && ! [ -e /home/dv/updater/config.yaml ]; then
   echo "Found dv-environment config. Copying..."
   cp /home/dv/environment/updater.config.yaml /home/dv/updater/config.yaml
   chown "$dv_user":"$dv_user" /home/dv/updater/config.yaml
fi

if [ -e /home/dv/updater/dv-updater.service ] && ! [ -e /etc/systemd/system/dv-updater.service ]
 then
   echo "Unit file not exists. Copying..."
   cp /home/dv/updater/dv-updater.service /etc/systemd/system/dv-updater.service
fi

echo "Configuring sudoers for $dv_user..."
cat > /etc/sudoers.d/dv-updater << EOF
$dv_user ALL=(ALL) NOPASSWD: /usr/bin/apt *, /usr/bin/yum *, /usr/bin/systemctl stop dv-updater.service, /usr/bin/systemctl start dv-updater.service, /usr/bin/systemctl enable dv-updater.service, /usr/bin/systemctl restart dv-updater.service, /usr/bin/systemctl try-restart dv-updater.service, /usr/bin/dpkg *
EOF
chmod 440 /etc/sudoers.d/dv-updater

chmod +x /home/dv/updater/dv-updater

echo "Enabling and restarting dv-updater.service..."
systemctl enable dv-updater.service
systemctl restart dv-updater.service

echo "Postinstall script done"
exit 0