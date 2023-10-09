build:
	go build -o ./bin/toolbox

run: build
	./bin/toolbox

test:
	go test -v ./...

download-meta:
	./bin/k0yote3web download meta -k${GO_PIVATE_KEY} -u${GO_ALCHEMY_RPC} -n${GO_NODE_PROVIDER} -a${GO_API_KEY} -s${GO_START_TOKENID} -e${GO_END_TOKENID} -b${GO_META_URL}

download-image:
	./bin/k0yote3web download image

cmd: FORCE
	cd cmd/k0yote3web && go build -o ../../bin/k0yote3web && cd -	

FORCE: ;