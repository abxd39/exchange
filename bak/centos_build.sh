
#cd ../gateway
#GOOS=linux GOARCH=amd64 go build
#mv gateway ../bin
#cd ../user_service
#GOOS=linux GOARCH=amd64 go build
#mv user_service ../bin
cd ../currency_service
GOOS=linux GOARCH=amd64 go build
mv currency_service ../bin
#cd ../public_service
#GOOS=linux GOARCH=amd64 go build
#mv public_service ../bin
#cd ../token_service
#GOOS=linux GOARCH=amd64 go build
#mv token_service ../bin
#cd ../price_service
#GOOS=linux GOARCH=amd64 go build
#mv price_service ../bin
#cd ../bin

scp -r -2 /home/wood/gopath/src/digicon/bin/* root@47.106.136.96:/root/go/src/dig
ssh root@47.106.136.96   "cd /root/go/src/dig/ && sh del.sh"
ssh root@47.106.136.96   "cd /root/go/src/dig/ && sh rb.sh"

