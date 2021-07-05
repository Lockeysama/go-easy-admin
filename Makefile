.PHONY: build update

export VERSION = 0.1.3
export REGISTRY = registry.cn-shanghai.aliyuncs.com/entertech

ifeq (${ENV}, test)
	export ENVFLAG = -test
else
	ifneq (${ENV}, release)
		$(error "Error: ENV(${ENV}) undefined.")
	endif
endif


_echo:
	@echo "Env: " ${ENV} "(" ${ENVFLAG} ")"
	@echo "Version: " ${VERSION}
	@echo "Registry: " ${REGISTRY}
	@echo "Project: " && pwd
	@echo "Git Branch: " && git branch | grep "*" | awk '{print $2}'
	@read -p "按任意键继续..."

build: _echo
	docker-compose build --build-arg ENVFLAG=${ENVFLAG} --build-arg VERSION=${VERSION}

update: build
	docker push ${REGISTRY}/tdcs${ENVFLAG}:${VERSION}
