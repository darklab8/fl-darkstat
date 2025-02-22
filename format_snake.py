# it can make json in grpc api gateway equal to json serializing in json structs
import subprocess

data = str(subprocess.run("buf lint darkapis/darkgrpc/statproto/darkstat.proto", shell=True, capture_output=True).stdout)

lines = data.split(r"\n")
print(len(lines), type(lines))

for index, line in enumerate(lines):
    if "should be lower_snake_case, such as" not in line:
        continue
    
    words = line.split("\"")
    should_be = words[3]
    it_is = words[1]

    if it_is == should_be:
        continue
    print(f"{index=}", should_be, it_is)

    with open("darkapis/darkgrpc/statproto/darkstat.proto", "r") as file:
        data = file.read()

    with open("darkapis/darkgrpc/statproto/darkstat.proto", "w") as file:
        file.write(data.replace(f"{it_is} =", f"{should_be} ="))
