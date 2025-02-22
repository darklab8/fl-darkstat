# it can make json in grpc api gateway equal to json serializing in json structs

with open("darkapis/darkgrpc/statproto/darkstat.pb.go", "r") as file:
    data = file.read()

lines = data.split("\n")
print(len(lines), type(lines))

for index, line in enumerate(lines):
    if "json=" not in line:
        continue

    should_be = line.split("json=")[1].split(",")[0]

    it_is = line.split("json:\"")[1].split(",")[0]

    if it_is == should_be:
        continue
    print(f"{index=}", should_be, it_is)

    # with open("darkapis/darkgrpc/statproto/darkstat.proto", "r") as file:
    #     data = file.read()

    # with open("darkapis/darkgrpc/statproto/darkstat.proto", "w") as file:
    #     file.write(data.replace(f"{it_is} =", f"{should_be} ="))
