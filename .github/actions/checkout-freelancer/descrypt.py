import base64
import os
with open("~/.ssh/id_rsa.files.freelancer.vanilla.out", "wb") as file:
    file.write(base64.b64decode(bytes(os.environ["ID_RSA_FILES_FREELANCER_VANILLA"], encoding='utf8')))
