#!/bin/bash
while :
do
	echo "attempting to clonse freelancer folder"
    git clone git@github.com:darklab8/fl-files-discovery.git freelancer_folder
    fl-data-discovery -wd /code/freelancer_folder
	sleep 30
done
