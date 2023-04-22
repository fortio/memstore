

test:
	( sleep 3; curl -G http://localhost:7999/set \
		--data-urlencode name=peers \
		--data-urlencode value=c,b,e,f ; echo ; sleep 5; pkill -int memstore) &
	go run . -peers a,b,c -config-port 7999
