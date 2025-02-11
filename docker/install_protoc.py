import subprocess

arch_info=str(subprocess.run("lscpu", shell=True, check=True, capture_output=True).stdout)
if "x86_64" in arch_info:
    arch = "x86_64"
else:
    arch = "aarch64"

if arch not in ["x86_64","aarch64"]:
    raise Exception(f"not support cpu architecture {arch=}")

if arch == "x86_64":
    subprocess.run("wget https://github.com/protocolbuffers/protobuf/releases/download/v29.3/protoc-29.3-linux-x86_64.zip -O protoc.zip", shell=True, check=True)

if arch == "aarch64":
    subprocess.run("wget https://github.com/protocolbuffers/protobuf/releases/download/v29.3/protoc-29.3-linux-aarch_64.zip -O protoc.zip", shell=True, check=True)

subprocess.run("unzip -p protoc.zip bin/protoc > /usr/bin/protoc", shell=True, check=True)
subprocess.run("chmod 777 /usr/bin/protoc", shell=True, check=True)
subprocess.run("rm protoc.zip", shell=True, check=True)
