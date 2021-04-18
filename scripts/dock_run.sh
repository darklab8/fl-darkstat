set -x

rep=fldarknet
port=8000

docker run --name ${rep} -t -d -p ${port}:$port --rm ${rep}:latest

