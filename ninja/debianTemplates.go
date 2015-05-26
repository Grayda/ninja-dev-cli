package main

const ninjaDebUpstart = `
description "Ninja Sphere {{.PackageName}}"
author      "http://www.ninjablocks.com"

start on filesystem and started mosquitto and started sphere-homecloud
stop on runlevel [016]

respawn
respawn limit 99 1

env RUN_AS=root
env APP={{.PackageName}}
env LOG=/var/log/{{.PackageName}}.log
env PID=/var/run/{{.PackageName}}.pid

env NINJA_APP_PATH={{.BaseDir}}/{{.TargetDir}}/{{.PackageName}}
env NINJA_APP_DATA_PATH=/data/sphere/appdata/{{.TargetDir}}/{{.PackageName}}

script
    . /etc/profile.d/ninja.sh
    mkdir -p $NINJA_APP_DATA_PATH
    exec start-stop-daemon -d $NINJA_APP_PATH --start --chuid $RUN_AS --make-pidfile --pidfile $PID --exec $APP >> $LOG 2>&1
end script
`

const ninjaDebPostInstall = `
service={{.PackageName}}

if test $(ps -ef | grep -v grep | grep $service | wc -l) -gt 0; then
	service $service stop
	service $service start
fi
`
