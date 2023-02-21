

test:
	( sleep 3; curl -G http://localhost:7999/set --data-urlencode name=peers \
		--data-urlencode value=c,b,e,f ) &
	go run . -peers a,b,c -config-port 7999
