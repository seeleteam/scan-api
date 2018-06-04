ll: chart_service scan_server node_service
chart_service:
	go build -o ./build/chart/chart_service ./cmd/chart_service
	cp ./cmd/chart_service/cmd/server.json ./build/chart/
	@echo "Done chart_service building"

scan_server:
	go build -o ./build/server/scan_server ./cmd/scan_server 
	cp ./cmd/scan_server/cmd/server.json ./build/server/
	@echo "Done scan_server building"

node_service:
	go build -o ./build/node/node_service ./cmd/node_service 
	cp ./cmd/node_service/cmd/server.json ./build/node/
	@echo "Done node_service building"

.PHONY: chart_service scan_server node_service
