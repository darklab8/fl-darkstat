#!/bin/sh
set -x
OLD_PASSWORD="smth"
while :
do
	echo "attempting to clonse freelancer folder"
    git clone gogs@git.discoverygc.com:DGCRepository/game-repository.git freelancer_folder
    $(cd freelancer_folder && git pull)
    NEW_PASSWORD=$(cd freelancer_folder && git rev-parse HEAD)
    if [ "$OLD_PASSWORD" != "$NEW_PASSWORD" ]
    then
        echo "Detected New Password $NEW_PASSWORD"
        docker service update --env-add DARKCORE_PASSWORD=$NEW_PASSWORD --image darkwind8/darkstat:production dev-darkstat-app
        OLD_PASSWORD=$NEW_PASSWORD
    fi
	sleep 30
done
