ll: chart_service scan_server seele_syncer node_service
chart_service:
	go build -o ./build/chart/chart_service ./cmd/chart_service
	cp ./cmd/chart_service/cmd/server1.json ./build/chart/
	cp ./cmd/chart_service/cmd/server2.json ./build/chart/
	@echo "Done chart_service building"

scan_server:
	go build -o ./build/server/scan_server ./cmd/scan_server 
	cp ./cmd/scan_server/cmd/server1.json ./build/server/
	cp ./cmd/scan_server/cmd/server2.json ./build/server/
	@echo "Done scan_server building"

node_service:
	go build -o ./build/node/node_service ./cmd/node_service 
	cp ./cmd/node_service/cmd/server1.json ./build/node/
	cp ./cmd/node_service/cmd/server2.json ./build/node/
	@echo "Done node_service building"

seele_syncer:
	go build -o ./build/syncer/seele_syncer ./cmd/seele_syncer 
	cp ./cmd/seele_syncer/cmd/server1.json ./build/syncer/
	cp ./cmd/seele_syncer/cmd/server2.json ./build/syncer/
	@echo "Done seele_syncer building"

.PHONY: chart_service scan_server node_service seele_syncer
