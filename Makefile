download-meta:
	./bin/k0yote3web download meta -k${GO_PIVATE_KEY} -u${GO_ALCHEMY_RPC} -n${GO_NODE_PROVIDER} -a${GO_API_KEY} -s${GO_START_TOKENID} -e${GO_END_TOKENID} -b${GO_META_URL}

download-image:
	./bin/k0yote3web download image

download-image:
	./bin/k0yote3web rewrite replace-meta	

cmd: FORCE
	cd cmd/k0yote3web && go build -o ../../bin/k0yote3web && cd -	

start-hardhat:
	docker build . -t hardhat-mainnet-fork
	docker start hardhat-node || docker run --name hardhat-node -d -p 8545:8545 -e "SDK_ALCHEMY_KEY=${SDK_ALCHEMY_KEY}" hardhat-mainnet-fork
	sudo bash ./scripts/test/await-hardhat.sh

stop-hardhat:
	docker stop hardhat-node
	docker rm hardhat-node

FORCE: ;