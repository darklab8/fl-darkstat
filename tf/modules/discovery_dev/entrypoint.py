import os
import time
import subprocess
import requests
import json
from typing import List

OLD_PASSWORD="smth"
while True:
    print("finding out old password")
    env_list: List[str] = json.loads(subprocess.run("docker service inspect dev-darkstat-app", shell=True, capture_output=True).stdout.decode("utf8").replace("\n",""))[0]["Spec"]["TaskTemplate"]["ContainerSpec"]["Env"]
    environ ={}
    for env in env_list:
        values = env.split("=")
        environ[values[0]] = values[1]
    OLD_PASSWORD=environ["DARKCORE_PASSWORD"]

    print("attempting to clone freelancer folder")
    os.system("git clone gogs@git.discoverygc.com:DGCRepository/game-repository.git freelancer_folder")
    os.system("cd freelancer_folder && git checkout 5.3 && git pull")
    NEW_PASSWORD=subprocess.run("cd freelancer_folder && git rev-parse HEAD", shell=True, capture_output=True).stdout.decode("utf8").replace("\n","")
    if OLD_PASSWORD != NEW_PASSWORD:
        print(f"Detected New Password {NEW_PASSWORD}")
        os.system(f'docker service update --env-add DARKCORE_PASSWORD={NEW_PASSWORD} --env-add "FLDARKSTAT_HEADING=commit:{NEW_PASSWORD}" --image darkwind8/darkstat:production dev-darkstat-app')
        OLD_PASSWORD=NEW_PASSWORD
        STATE=json.loads(subprocess.run("docker service inspect dev-darkstat-app", shell=True, capture_output=True).stdout.decode("utf8").replace("\n",""))[0]["UpdateStatus"]["State"]
        content_msg = f"https://darkstat-dev.dd84ai.com/?password={NEW_PASSWORD} state={STATE}"
        if "state=completed" not in content_msg:
            content_msg += " <@370435997974134785>"
        data = dict(username="Darkstat",content=content_msg)
        result = requests.post(os.environ["DISCO_DEV_WEBHOOK"], json = data)
    time.sleep(30)
