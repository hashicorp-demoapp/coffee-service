# Assumes kind installed locally
kind delete cluster --name coffee-service-debug

make build_linux

docker build -t coffee-service:debug .

kind create cluster --name coffee-service-debug

kind load docker-image coffee-service:debug --name coffee-service-debug

kubectl apply -f ./deployments/debug/coffee-service.yaml

kubectl get pods