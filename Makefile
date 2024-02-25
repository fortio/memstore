

test:
	# In CI it takes quite a while for the go process to run
	# prime it by compiling/running:
	go run . help
	( sleep 10; pkill -term memstore) & # make sure this happens even if the curl below fails
	( sleep 4; curl -G http://localhost:7999/set \
		--data-urlencode name=peers \
		--data-urlencode value=c,b,e,f ; echo) &
	@echo 'Expect to see: Success "peers" -> "c,b,e,f"'
	go run . -peers a,b,c -config-port 7999

# Works with docker-desktop for instance:

LOCAL_HELM_OVERRIDES:=--set image.pullPolicy=Never --set debug=true --set epoch=$(shell date +%s)
HELM:=helm
CHART_NAME:=memstore
CHART_DIR:=chart/
HELM_INSTALL_ARGS:=upgrade --install $(CHART_NAME) $(CHART_DIR) $(LOCAL_HELM_OVERRIDES)

local-k8s:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" .
	# -kubectl delete statefulset -n memstore memstore # so it'll reload the image
	docker buildx build --load --tag fortio/memstore:latest .
	$(HELM) $(HELM_INSTALL_ARGS)

# Needs helm plugin install https://github.com/databus23/helm-diff
local-diff:
	$(HELM) diff $(HELM_INSTALL_ARGS)

# Logs of first pod, colorized with logc (go install fortio.org/logc@latest)
tail-log:
	kubectl logs -f -n memstore memstore-0 | logc

debug-pod:
	kubectl run debug --image=ubuntu -- /bin/sleep infinity
