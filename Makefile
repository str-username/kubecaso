DEPLOY_LOCAL_INFRA_MANIFESTS_PATH="./deploy/local"

all: cluster

cluster:
	@echo "=== create cluster  ====" \
	&& kind create cluster \
	--config=${DEPLOY_LOCAL_INFRA_MANIFESTS_PATH}/cluster.yaml
	@echo "========================"

clean:
	@echo "=== delete all clusters  ====" \
	&& kind delete clusters --all
	@echo "============================="