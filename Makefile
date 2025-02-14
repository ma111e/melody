.PHONY: docs

default: help

## certs : Create TLS certs used by the HTTPS dummy server in "var/https/certs"
certs:
	-mkdir -p var/https/certs
	openssl req -x509 -subj "/C=AU/ST=Some-State/O=Internet Widgits Pty Ltd/CN=localhost" -newkey rsa:4096 -keyout var/https/certs/key.pem -out var/https/certs/cert.pem -days 3650 -nodes

## enable_all_rules : Enable all the rule files present in ./rules/rules-available/
enable_all_rules:
	ln -rs ./rules/rules-available/*.yml ./rules/rules-enabled/

## docker_build : Build Docker image
docker_build:
	docker build . -t melody

## docker_run : Run local Docker image
docker_run:
	@echo "This will only run the locally built image."
	@echo "Swap 'melody' with bonjourmalware/melody if you want to run the latest release."
	@echo "MELODY_CLI example : 'export MELODY_CLI=\"-s -i lo -F 'dst port 10080'\"'"
	@-mkdir logs
	docker run \
		--net=host \
		-e "MELODY_CLI=${MELODY_CLI}" \
		--mount type=bind,source="$(shell pwd)/filter.bpf",target=/app/filter.bpf,readonly \
		--mount type=bind,source="$(shell pwd)/config.yml",target=/app/config.yml,readonly \
		--mount type=bind,source="$(shell pwd)/var",target=/app/var,readonly \
		--mount type=bind,source="$(shell pwd)/rules",target=/app/rules,readonly \
		--mount type=bind,source="$(shell pwd)/logs",target=/app/logs/ \
		melody

## docker : Build and run Docker image
docker: docker_build docker_run

## docs : Deploy documentation
docs:
	cd docs/; mkdocs gh-deploy

## run_local_stdout : Start Melody and log to stdout
run_local_stdout: build
	./melody -s

## build : Build and set network capabilities to start Melody without elevated privileges
build:
	go build -ldflags="-s -w" -o melody
	sudo setcap cap_net_raw,cap_setpcap=ep ./melody

## cap : Set network capabilities to start Melody without elevated privileges
cap:
	sudo setcap cap_net_raw,cap_setpcap=ep ./melody

## setup : Patch listen.interface config key with the current default interface
setup:
	@echo "Ensure net-tools is installed in order to use the 'route' command"
	sudo apt install net-tools
	@echo "> Setting listening interface to \"$(shell route | grep '^default' | awk '{print $$8; exit}')\""
	sed -i "s/# listen.interface: \"lo\"/listen.interface: \"$(shell route | grep '^default' | awk '{print $$8; exit}')\"/g" ./config.yml
	@echo
	@echo -n "Current listening interface :\n\t"
	@grep listen.interface ./config.yml

	@echo -n "Current BPF is :\n\t"
	@cat ./filter.bpf

	# Don't forget to filter the noise by editing ./filter.bpf

## supervisor : Create Melody's supervisor configuration
supervisor:
	sudo ln -s $(shell pwd)/etc/melody.conf /etc/supervisor/conf.d/
	sudo supervisorctl reload
	sudo supervisorctl status all

## service : Create Melody's systemd configuration and enable start at boot
service:
	sudo ln -s $(shell pwd)/etc/melody.service /etc/systemd/system/melody.service
	sudo systemctl daemon-reload && sudo systemctl enable melody
	sudo systemctl stop melody && sudo systemctl status melody

## help : Show this help
help: Makefile
	@printf "\n Melody helpers\n\n"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@printf ""
