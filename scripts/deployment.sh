set -x

rep=fldarknet
port=$1
rm -r $rep
git clone https://github.com/dd84ai/fldarknet.git
#cp ${rep}.df $rep/dockerfile

docker ps
d=$(docker ps -a -q --filter="name=${rep}")
docker stop $d
docker rmi $(docker images -a -q)

docker build -t ${rep}:latest ${rep}

docker run --name ${rep} -t -d -p ${port}:$port --rm ${rep}:latest
rm -r $rep
