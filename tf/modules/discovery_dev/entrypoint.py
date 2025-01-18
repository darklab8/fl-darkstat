import os
import time
import subprocess
import requests
import json

OLD_PASSWORD="smth"
while True:
    print("attempting to clonse freelancer folder")
    os.system("git clone gogs@git.discoverygc.com:DGCRepository/game-repository.git freelancer_folder")
    os.system("cd freelancer_folder && git pull")
    NEW_PASSWORD=subprocess.run("cd freelancer_folder && git rev-parse HEAD", shell=True, capture_output=True).stdout.decode("utf8").replace("\n","")
    if OLD_PASSWORD != NEW_PASSWORD:
        print(f"Detected New Password {NEW_PASSWORD}")
        os.system(f"docker service update --env-add DARKCORE_PASSWORD={NEW_PASSWORD} --image darkwind8/darkstat:production dev-darkstat-app")
        OLD_PASSWORD=NEW_PASSWORD
        STATE=json.loads(subprocess.run("docker service inspect dev-darkstat-app", shell=True, capture_output=True).stdout.decode("utf8").replace("\n",""))[0]["UpdateStatus"]["State"]
        data = dict(username="Darkstat",content=f"https://darkstat-dev.dd84ai.com/?password={NEW_PASSWORD} state={STATE}")
        result = requests.post(os.environ["DISCO_DEV_WEBHOOK"], json = data)
    time.sleep(30)
