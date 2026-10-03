start:
	go run . -rod="show,trace"

dev:
	go run . -rod="show,slow=0.3s,trace"