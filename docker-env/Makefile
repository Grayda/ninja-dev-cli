build:
	docker build -t ninjasphere-developer-partial .
	docker run --privileged -it ninjasphere-developer-partial /scripts-cross/privileged-prep
	sh -c "docker commit -a 'Builder' $$(docker ps -l -q) ninjasphere/ninjasphere-developer-cross-base"
