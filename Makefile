RELEASE_PATH=root/etc/release
NETWORK_INTERFACES_PATH=root/etc/network/interfaces

deps:
	sudo apt install -y golang-go


build:
	go build .


run:
	NETWORK_INTERFACES_PATH=${NETWORK_INTERFACES_PATH} go run .

dev:
	NETWORK_INTERFACES_PATH=${NETWORK_INTERFACES_PATH} nodemon --exec go run system.go --signal SIGTERM
