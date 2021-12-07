#!/usr/bin/env bash

if [ ! -d go ]; then
  echo "[!] Creating go output directory"
  mkdir go;
fi

if [ ! -d python ]; then
  echo "[!] Creating python output directory"
  mkdir python;
fi

if [ ! -d csharp ]; then
  echo "[!] Creating csharp output directory"
  mkdir csharp;
fi

echo "[+] Building docker container"
docker image build -t tgams-grpc-build:1.0 .
docker container run --detach --name tgamsgrpc tgams-grpc-build:1.0
docker cp tgamsgrpc:/go/src/github.com/BGrewell/tgams/api/go/github.com/BGrewell/tgams/api/go/tgams.pb.go ./go/.
echo "[+] Updating of go library complete"

docker cp tgamsgrpc:/go/src/github.com/BGrewell/tgams/api/csharp/Tgams.cs ./csharp/.
echo "[+] Updating of csharp library complete"

docker cp tgamsgrpc:/go/src/github.com/BGrewell/tgams/api/python/tgams_pb2.py ./python/.
docker cp tgamsgrpc:/go/src/github.com/BGrewell/tgams/api/python/tgams_pb2_grpc.py ./python/.
echo "[+] Updating of python library complete"

echo "[+] Removing docker container"
docker rm tgamsgrpc

echo "[+] Adding new files to source control"
git add go/tgams.pb.go
git add csharp/Tgams.cs
git add python/tgams_pb2.py
git add python/tgams_pb2_grpc.py
git commit -m "regenerated grpc libraries"
git push

echo "[+] Done. Everything has been rebuilt and the repository has been updated and pushed"