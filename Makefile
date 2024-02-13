

test:
	# In CI it takes quite a while for the go process to run
	# prime it by compiling/running:
	go run . help
	( sleep 10; pkill -int memstore) & # make sure this happens even if the curl below fails
	( sleep 4; curl -G http://localhost:7999/set \
		--data-urlencode name=peers \
		--data-urlencode value=c,b,e,f ; echo) &
	@echo 'Expect to see: Success "peers" -> "c,b,e,f"'
	go run . -peers a,b,c -config-port 7999

local-k8s:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" .
	-kubectl delete statefulset -n memstore memstore # so it'll reload the image
	docker buildx build --load --tag fortio/memstore:latest .
	kubectl apply -f deploy

debug-pod:
	kubectl run debug --image=ubuntu --restart=Never -- /bin/sleep infinity
