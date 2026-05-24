import base64
import os

with open("ssh_key", "rb") as file:
    data = file.read()

env_var = str(base64.b64encode(data),encoding='utf8')
os.environ["ID_RSA_FILES_FREELANCER_VANILLA"] =env_var

print(env_var)
print("\n\n")

with open("ssh_key.out", "wb", 0o600) as file:
    file.write(base64.b64decode(bytes(os.environ["ID_RSA_FILES_FREELANCER_VANILLA"], encoding='utf8')))
os.chmod("ssh_key.out", 0o600)