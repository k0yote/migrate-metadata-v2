include .env

download-meta:
	./bin/k0yote3web download meta -s ${GO_START_TOKENID} -e ${GO_END_TOKENID} -b ${GO_META_URL}

download-image:
	./bin/k0yote3web download image

rewrite-meta:
	./bin/k0yote3web rewrite replace-meta -g ${GO_IPFS_IMAGE_URL}

ipfs-upload:
	./bin/k0yote3web ipfs-upload meta -f ${IPFS_UPLOAD_FILES_PATH}

cmd: FORCE
	cd cmd/k0yote3web && go build -o ../../bin/k0yote3web && cd -	

start-hardhat:
	docker build . -t hardhat-mainnet-fork
	docker start hardhat-node || docker run --name hardhat-node -d -p 8545:8545 -e "SDK_ALCHEMY_KEY=${SDK_ALCHEMY_KEY}" hardhat-mainnet-fork
	sh ./scripts/test/await-hardhat.sh

stop-hardhat:
	docker stop hardhat-node
	docker rm hardhat-node

FORCE: ;