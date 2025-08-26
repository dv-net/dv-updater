#!/bin/sh

echo "Executing postremove script"

if ! [ -e /home/dv/updater/dv-updater ]
 then
   echo "Dv updater removed. Disabling..."
   if systemctl list-unit-files | grep "dv-updater.service"
    then
       systemctl disable dv-updater.service
       systemctl stop dv-updater.service
   fi
fi

echo "Postremove script done"
