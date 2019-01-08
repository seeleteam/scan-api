# Scan API

## Env
>recommend
```text
go 1.10+
mongo 3.6.3
```

## Project structure
```text
┌── api: api interface
│   ├── handlers: router handler
│   └── routers:  the http router
├── chart: chart data processor
│   ├── address: address chart processor
│   ├── block: block count and reward chart processor
│   ├── blockdifficulty: block difficulty chart processor
│   ├── blocktime: block average time chart processor
│   ├── hashrate: hashrate chart processor
│   ├── topminers: topminers chart processor
│   ├── tx_history: transaction hsitory chart processor
├── cmd: app entrance
|   ├── chart_service: chart service entrance
|   ├── node_service: node service entrance
|   ├── seele_syncer: seele syncer entrance
│   └── scan_server:  http service entrance
├── database: mongodb database
├── log: third logger warpper
├── node: node service
├── rpc:  json rpc
├── server:  scan server
└── vendor: third dependencies

```

## Start
```
cd somewhere/you/want/to/download/
go get -u -v github.com/seeleteam/scan-api
cd scan-api

# generate the executable file
make

# start seele_syncer
cd build/syncer/
./seele_syncer -c server.json

# start scan_server
cd build/server/
./scan_server -c server.json

# start chart_service
 cd build/chart
 ./chart_service -c server.json
 
# start node_service
cd build/node
./node_service -c server.json
```

## Config
```text

"GinMode":"debug"
# gin run mode, format ip:port

"Addr": ":8888"
# server listen address and port

"LimitConnection": 0
# connection limit number

"DefaultHammerTime": 30
# connection limit number

"RpcURL": "127.0.0.1:55028"
# seele node rpc address and port

"WriteLog": true
# enable write log out

"LogLevel": "debug"
# log level

"LogFile": "scan-api.log"
# log filename

"MaxHeaderBytes": 20
"ReadTimeout":300,
"IdleTimeout": 0,
"WriteTimeout":120,
# gin settting

"DataBaseConnUrl":"127.0.0.1:27017",
"DataBaseName":"seele",
# mongodb name address and port 

"Interval":30
# sync interval

```
