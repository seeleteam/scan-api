# Block APIs
#### 获取区块列表
	
	https://api.seelescan.io/api/v1/blocks

#### 参数 
1. p:要显示的页码,默认值为1
2. ps: 每页显示数量,默认值为25

#### 返回
1. code: 错误码,0为正常,非0为错误
2. message: 错误提示,正确执行会空
3. data: 返回一个指定页码的按照区块高度降序排序的区块列表
	- lsit: 区块列表
	- pageInfo: 分页详情信息

#### 例子
	//Request
	http://api.seelescan.io/api/v1/blocks?p=1&ps=10
	
	//Return
	{
		"code":0,
		"data":{
			"list":[
				{
					"height":5957,
					"age":"1 hours ago",
					"txn":1,
					"miner":"0x4dd6881d13ab5152127533c5954e4e062eb4bb2dcd93becf4f4e9b1d2d69f1363eea0395e8e76a2716b033d1e3cc8da2bf24811b1e31a86ac8bcacca4c4b29bd"
				},
				{
					"height":5958,
					"age":"1 hours ago",
					"txn":1,
					"miner":"0x4dd6881d13ab5152127533c5954e4e062eb4bb2dcd93becf4f4e9b1d2d69f1363eea0395e8e76a2716b033d1e3cc8da2bf24811b1e31a86ac8bcacca4c4b29bd"
				}
			],
			"pageInfo":{
				"begin":5956,
				"curPage":0,
				"end":5976,
				"totalCount":5976
			}
		},
		"message":""
	}

#### 获取区块详细信息

	http://api.seelescan.io/api/v1/block
	
#### 参数 
1. height:待查询的区块高度
2. hash: 待查询的区块的Hash值

#### 返回
1. code: 错误码,0为正常,非0为错误
2. message: 错误提示,正确执行会空
3. data: 返回要查询的区块的详细信息

#### 例子
	//Request
		//By Height
		http://api.seelescan.io/api/v1/block?height=5567
		//By Hash
		http://api.seelescan.io/api/v1/block?hash=0x00000057df238881381bb218a5d5f6b1589d969e6c6fb0aa50129dd85786e69d
	
	//Return
	{
		"code": 0, 
		"data": {
			"headHash": "0x000000a830505c2df9ff542d2fe70f72efeb8ced3927460b44c64321159a2ec0", 
			"preBlockHash": "0x0000019d36b3c399a297c68540ff1a0bca75321c3d115ec7bb454ae4e7ea1195", 
			"height": 4, 
			"age": "14 days ago", 
			"difficulty": 10000000, 
			"miner": "0x1cba7cc4097c34ef9d90c0bf1fa9babd7e2fb26db7b49d7b1eb8f580726e3a99d3aec263fc8de535e74a79138622d320b3765b0a75fabd084985c456c6fe65bb", 
			"nonce": "13260572831091132416", 
			"txcount": 1
		}, 
		"message": ""
	}

# Transaction APIs
#### 获取交易列表
    
	https://api.seelescan.io/api/v1/txs?p=1&ps=10&block=5567
	
#### 参数 
1. p:要显示的页码 
2. ps: 每页显示数量
3. block:区块的高度

#### 返回
1. code: 错误码,0为正常,非0为错误
2. message: 错误提示,正确执行会空
3. data: 返回一个指定页码的按照交易序号降序排序的交易列表
	- lsit: 交易列表
	- pageInfo: 分页详情信息

#### 例子
	//Request
	//By Default
	https://api.seelescan.io/api/v1/txs?p=1&ps=10
	//By Block
	https://api.seelescan.io/api/v1/txs?p=1&ps=10&block=5567
	
	//Return
	{
		"code": 0, 
		"data": {
			"list": [
				{
					"Hash": "0x26d8ecfb5b75e3f6da5750b072ac6b8bd969d1ce453206f7f33062cad89397eb", 
					"From": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", 
					"To": "0x1cba7cc4097c34ef9d90c0bf1fa9babd7e2fb26db7b49d7b1eb8f580726e3a99d3aec263fc8de535e74a79138622d320b3765b0a75fabd084985c456c6fe65bb", 
					"Amount": "10", 
					"Timestamp": "1525425890093552640"
				}
			], 
			"pageInfo": {
				"begin": 0, 
				"curPage": 0, 
				"end": 1, 
				"totalCount": 1
			}
		}, 
		"message": ""
	}

#### 获取交易详细
    
	https://api.seelescan.io/api/v1/tx
	
#### 参数 
1. txhash: 交易的哈希值

#### 返回
1. code: 错误码,0为正常,非0为错误
2. message: 错误提示,正确执行会空
3. data: 返回一个指定交易的详细信息

#### 例子
	//Request
	https://api.seelescan.io/api/v1/tx?txhash=0x4d58d1edcbdb91f9942186b3db4d0214c5d2ab9fff5c79766d7beb46cac7881f
	
	//Return
	{
		"code": 0, 
		"data": {
			"txHash": "0x649b7ab12c0bf721e9a5bda7fa19f1029e3f70ed2d6fd49eafe066149e7cbf98", 
			"block": 4, 
			"age": "14 days ago", 
			"from": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", 
			"to": "0x1cba7cc4097c34ef9d90c0bf1fa9babd7e2fb26db7b49d7b1eb8f580726e3a99d3aec263fc8de535e74a79138622d320b3765b0a75fabd084985c456c6fe65bb", 
			"value": "10"
		}, 
		"message": ""
	}

# Account APIs
#### 获取账户列表
	
	https://api.seelescan.io/api/v1/accounts

#### 参数 
1. p:要显示的页码,默认值为1
2. ps: 每页显示数量,默认值为25

#### 返回
1. code: 错误码,0为正常,非0为错误
2. message: 错误提示,正确执行会空
3. data: 返回一个指定页码的按照账户余额降序排序的账户列表
	- lsit: 账户列表
	- pageInfo: 分页详情信息

#### 例子
	//Request
	http://api.seelescan.io/api/v1/blocks?p=1&ps=10
	
	//Return
	{
			"code": 0, 
			"data": {
					"list": [
							{
									"rank": 1, 
									"address": "0x4dd6881d13ab5152127533c5954e4e062eb4bb2dcd93becf4f4e9b1d2d69f1363eea0395e8e76a2716b033d1e3cc8da2bf24811b1e31a86ac8bcacca4c4b29bd", 
									"balance": 318600, 
									"percentage": 0, 
									"txcount": 1593
							}, 
							{
									"rank": 2, 
									"address": "0x1cba7cc4097c34ef9d90c0bf1fa9babd7e2fb26db7b49d7b1eb8f580726e3a99d3aec263fc8de535e74a79138622d320b3765b0a75fabd084985c456c6fe65bb", 
									"balance": 54910, 
									"percentage": 0, 
									"txcount": 5491
							}
					], 
					"pageInfo": {
							"begin": 0, 
							"curPage": 1, 
							"end": 2, 
							"totalCount": 2
					}
			}, 
			"message": ""
	}

#### 获取账户详细
	
	https://api.seelescan.io/api/v1/account

#### 参数 
1. address: 账户的地址

#### 返回
返回一个指定账户的详细信息

#### 例子
	//Request
	https://api.seelescan.io/api/v1/account?address=0x4dd6881d13ab5152127533c5954e4e062eb4bb2dcd93becf4f4e9b1d2d69f1363eea0395e8e76a2716b033d1e3cc8da2bf24811b1e31a86ac8bcacca4c4b29bd
	
	//Return
	{
        "code": 0, 
        "data": {
                "address": "0x4dd6881d13ab5152127533c5954e4e062eb4bb2dcd93becf4f4e9b1d2d69f1363eea0395e8e76a2716b033d1e3cc8da2bf24811b1e31a86ac8bcacca4c4b29bd", 
                "balance": 318600, 
                "percentage": 0, 
                "txcount": 1593, 
                "txs": [
                        {
                                "hash": "0x199de9e63d8f986cb26c52ebc553cd2f020d08e15ac842ed7669310a036d5eca", 
                                "block": 7084, 
                                "from": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", 
                                "to": "0x4dd6881d13ab5152127533c5954e4e062eb4bb2dcd93becf4f4e9b1d2d69f1363eea0395e8e76a2716b033d1e3cc8da2bf24811b1e31a86ac8bcacca4c4b29bd", 
                                "amount": 200, 
                                "age": "1 secs ago", 
                                "txfee": 0, 
                                "inorout": true
                        }
                ]
        }, 
        "message": ""
	}
	
# Transaction APIs
#### 获取交易列表
    
	https://api.seelescan.io/api/v1/txs?p=1&ps=10&block=5567
	
#### 参数 
1. p:要显示的页码 
2. ps: 每页显示数量
3. block:区块的高度

#### 返回
1. code: 错误码,0为正常,非0为错误
2. message: 错误提示,正确执行会空
3. data: 返回一个指定页码的按照交易序号降序排序的交易列表
	- lsit: 交易列表
	- pageInfo: 分页详情信息

#### 例子
	//Request
	//By Default
	https://api.seelescan.io/api/v1/txs?p=1&ps=10
	//By Block
	https://api.seelescan.io/api/v1/txs?p=1&ps=10&block=5567
	
	//Return
	{
		"code": 0, 
		"data": {
			"list": [
				{
					"Hash": "0x26d8ecfb5b75e3f6da5750b072ac6b8bd969d1ce453206f7f33062cad89397eb", 
					"From": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", 
					"To": "0x1cba7cc4097c34ef9d90c0bf1fa9babd7e2fb26db7b49d7b1eb8f580726e3a99d3aec263fc8de535e74a79138622d320b3765b0a75fabd084985c456c6fe65bb", 
					"Amount": "10", 
					"Timestamp": "1525425890093552640"
				}
			], 
			"pageInfo": {
				"begin": 0, 
				"curPage": 0, 
				"end": 1, 
				"totalCount": 1
			}
		}, 
		"message": ""
	}

#### 获取交易详细
    
	https://api.seelescan.io/api/v1/tx
	
#### 参数 
1. txhash: 交易的哈希值

#### 返回
1. code: 错误码,0为正常,非0为错误
2. message: 错误提示,正确执行会空
3. data: 返回一个指定交易的详细信息

#### 例子
	//Request
	https://api.seelescan.io/api/v1/tx?txhash=0x4d58d1edcbdb91f9942186b3db4d0214c5d2ab9fff5c79766d7beb46cac7881f
	
	//Return
	{
		"code": 0, 
		"data": {
			"txHash": "0x649b7ab12c0bf721e9a5bda7fa19f1029e3f70ed2d6fd49eafe066149e7cbf98", 
			"block": 4, 
			"age": "14 days ago", 
			"from": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", 
			"to": "0x1cba7cc4097c34ef9d90c0bf1fa9babd7e2fb26db7b49d7b1eb8f580726e3a99d3aec263fc8de535e74a79138622d320b3765b0a75fabd084985c456c6fe65bb", 
			"value": "10"
		}, 
		"message": ""
	}

# Node APIs
#### 获取节点列表
		https://api.seelescan.io/api/v1/nodes
	
#### 参数 
1. p:要显示的页码,默认值为1
2. ps: 每页显示数量,默认值为25

#### 返回
1. code: 错误码,0为正常,非0为错误
2. message: 错误提示,正确执行会空
3. data: 返回一个指定页码的节点列表
	- lsit: 节点列表
	- pageInfo: 分页详情信息

#### 例子
	//Request
	https://api.seelescan.io/api/v1/nodes?p=1&ps=10
	
	//Return
	{
        "code": 0, 
        "data": {
                "list": [
                        {
                                "ID": "d608043ca78cae6074deccaa0290320277abf9ffc1b004badf5e85fc648a0cd597de02338626729b60b1721a352ca621bd13ed378819e29c60d1bff3ad5cabef", 
                                "Host": "116.24.64.70", 
                                "Port": "51563", 
                                "City": "Shenzhen", 
                                "Region": "Guangdong", 
                                "Country": "China", 
                                "Client": "Geth/v1.8.8-stable-2688dab4/windows-amd64/go1.10.1", 
                                "Caps": "eth/63", 
                                "LastSeen": 1527493747, 
                                "LongitudeAndLatitude": "[114.1333,22.5333]"
                        }
                ], 
                "pageInfo": {
                        "begin": 1306, 
                        "curPage": 1, 
                        "end": 1316, 
                        "totalCount": 1316
                }
        }, 
        "message": ""
	}
	
#### 获取节点详情
	
	https://api.seelescan.io/api/v1/node

#### 参数 
1. id: 节点的node ID

#### 返回
1. code: 错误码,0为正常,非0为错误
2. message: 错误提示,正确执行会空
3. data: 返回节点的详细信息

#### 例子
	//Request
	https://api.seelescan.io/api/v1/node?id=23ddfb54a488f906cdb9cbd257eac5663a4c74ba25619bb902651602a4491be4ce437907fcc567b31be6746a014931f4670ac116c0010e5beb28b0dce2c6eaad
	
	//Return
	{
        "code": 0, 
        "data": {
                "ID": "b9fc9dc30a707b90255e0abc728ccc8eef443f71f9bfee1683db17b61f0f69fd0ec88616b4ffadc527343ec3940eb69e875a5cee71932f7f4a197393ca1a2f93", 
                "Host": "180.167.100.186", 
                "Port": "49736", 
                "City": "Shanghai", 
                "Region": "Shanghai", 
                "Country": "China", 
                "Client": "Parity/v1.10.3-stable-b9ceda3-20180507/x86_64-macos/rustc1.25.0", 
                "Caps": "eth/62|eth/63|par/1|par/2|pip/1", 
                "LastSeen": 1527495787, 
                "LongitudeAndLatitude": "[121.3997,31.0456]"
        }, 
        "message": ""
	}

#### 获取节点地图
	
	https://api.seelescan.io/api/v1/nodemap

#### 返回
1. code: 错误码,0为正常,非0为错误
2. message: 错误提示,正确执行会空
3. data: 返回所有的全节点列表

#### 例子
	//Request
	https://api.seelescan.io/api/v1/nodemap
	
	//Return
	{
		"code": 0, 
		"data": [
			{
				"ID": "60c2425e1eaf7ef60a9770662bf28199b37092ba332800dffacbfe23096d02219d8d7b84ca4640fe89b9feffc81d5c1c6b09311a4c888c363b24aed75246b0f9", 
				"Host": "118.25.66.79", 
				"Port": "30303", 
				"City": "Beijing", 
				"Region": "Beijing", 
				"Country": "China", 
				"Client": "Geth/v1.8.8-stable/linux-amd64/go1.9.2", 
				"Caps": "eth/63", 
				"LastSeen": 1527157086, 
				"LongitudeAndLatitude": "[116.3883,39.9289]"
			}
		], 
		"message": ""
	}
	
# Stat APIs
#### 获取最新区块高度
	https://api.seelescan.io/api/v1/lastblock

#### 返回
1. code: 错误码,0为正常,非0为错误
2. message: 错误提示,正确执行会空
3. data: 最新区块的高度

#### 例子
	//Request
	https://api.seelescan.io/api/v1/lastblock
	
	//Return
	{
			"code": 0, 
			"data": 44825, 
			"message": ""
	}
	
#### 获取最新区块生成时间
	https://api.seelescan.io/api/v1/bestblock
	
#### 返回
1. code: 错误码,0为正常,非0为错误
2. message: 错误提示,正确执行会空
3. data: 最新区块的生成时间

#### 例子
	//Request
	https://api.seelescan.io/api/v1/bestblock
	
	//Return
	{
        "code": 0, 
        "data": "31 secs ago", 
        "message": ""
	}
	
#### 获取平均区块生成时间
	https://api.seelescan.io/api/v1/avgblocktime
	
#### 返回
1. code: 错误码,0为正常,非0为错误
2. message: 错误提示,正确执行会空
3. data: 平均区块的生成时间,单位是秒

#### 例子
	//Request
	https://api.seelescan.io/api/v1/avgblocktime
	
	//Return
	{
			"code": 0, 
			"data": 14, 
			"message": ""
	}
	

#### 获取交易总数量
	https://api.seelescan.io/api/v1/txcount
	
#### 返回
1. code: 错误码,0为正常,非0为错误
2. message: 错误提示,正确执行会空
3. data: 平均区块的生成时间,单位是秒

#### 例子
	//Request
	https://api.seelescan.io/api/v1/txcount
	
	//Return
	{
			"code": 0, 
			"data": 44842, 
			"message": ""
	}
	
#### 获取平均区块难度
	https://api.seelescan.io/api/v1/difficulty
	
#### 返回
1. code: 错误码,0为正常,非0为错误
2. message: 错误提示,正确执行会空
3. data: 区块平均难度

#### 例子
	//Request
	https://api.seelescan.io/api/v1/txcount
	
	//Return
	{
			"code": 0, 
			"data": 16628715, 
			"message": ""
	}
	
#### 获取平均哈希速率
	https://api.seelescan.io/api/v1/hashrate
	
#### 返回
1. code: 错误码,0为正常,非0为错误
2. message: 错误提示,正确执行会空
3. data: 平均哈希速率

#### 例子
	//Request
	https://api.seelescan.io/api/v1/txcount
	
	//Return
	{
			"code": 0, 
			"data": 1217636.8141435185, 
			"message": ""
	}
	
# Search APIs
#### 查询账户，区块或者交易的详细信息

	https://api.seelescan.io/api/v1/search

#### 参数 
1. content: 区块的高度，交易的哈希值，账户的地址之一

#### 返回
1. code: 错误码,0为正常,非0为错误
2. message: 错误提示,正确执行会空
3. data: 返回搜索的区块,交易或者账户的详细信息
	- info: 区块,交易,账户的详细信息
	- type: 返回的数据类型:block, transaction, account

#### 例子
查询交易
	
	//Request
	https://api.seelescan.io/api/v1/search?content=0x4d58d1edcbdb91f9942186b3db4d0214c5d2ab9fff5c79766d7beb46cac7881f
	
	//Return
	{
			"code": 0, 
			"data": {
					"info": {
							"txHash": "0x649b7ab12c0bf721e9a5bda7fa19f1029e3f70ed2d6fd49eafe066149e7cbf98", 
							"block": 4, 
							"age": "15 days ago", 
							"from": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", 
							"to": "0x1cba7cc4097c34ef9d90c0bf1fa9babd7e2fb26db7b49d7b1eb8f580726e3a99d3aec263fc8de535e74a79138622d320b3765b0a75fabd084985c456c6fe65bb", 
							"value": "10"
					}, 
					"type": "transaction"
			}, 
			"message": ""
	}

查询区块
	
	//Request
	http://106.75.171.117:3003/api/v1/search?content=1000
	
	//Return
	{
			"code": 0, 
			"data": {
					"info": {
							"headHash": "0x0000002485f23e29bbc541e0b2e3be07f41d91b0ccbf2d1faeb9c5aa8aa35328", 
							"preBlockHash": "0x000000b844bc2a667ecac99f88127d8affd40fd6ef7d6a5b989411d381fd0c2a", 
							"height": 1000, 
							"age": "7 days ago", 
							"difficulty": 12570123, 
							"miner": "0x4dd6881d13ab5152127533c5954e4e062eb4bb2dcd93becf4f4e9b1d2d69f1363eea0395e8e76a2716b033d1e3cc8da2bf24811b1e31a86ac8bcacca4c4b29bd", 
							"nonce": "11984700971403517952", 
							"txcount": 1, 
							"maxheight": 44802, 
							"minheight": 0
					}, 
					"type": "block"
			}, 
			"message": ""
	}
	
查询账户
	
	//Request
	http://106.75.171.117:3003/api/v1/search?content=0x4dd6881d13ab5152127533c5954e4e062eb4bb2dcd93becf4f4e9b1d2d69f1363eea0395e8e76a2716b033d1e3cc8da2bf24811b1e31a86ac8bcacca4c4b29bd
	
	//Return
	{
			"code": 0, 
			"data": {
					"info": {
							"address": "0x4dd6881d13ab5152127533c5954e4e062eb4bb2dcd93becf4f4e9b1d2d69f1363eea0395e8e76a2716b033d1e3cc8da2bf24811b1e31a86ac8bcacca4c4b29bd", 
							"balance": 467080000000000, 
							"percentage": 0, 
							"txcount": 23354, 
							"txs": [
									{
											"hash": "0xe078976326b60829e6a2d91de5402cdc0f8440158869c6c9a9fddfe38871bbb7", 
											"block": 44805, 
											"from": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", 
											"to": "0x4dd6881d13ab5152127533c5954e4e062eb4bb2dcd93becf4f4e9b1d2d69f1363eea0395e8e76a2716b033d1e3cc8da2bf24811b1e31a86ac8bcacca4c4b29bd", 
											"amount": 20000000000, 
											"age": "1 secs ago", 
											"txfee": 0, 
											"inorout": true
									}
							]
					}, 
					"type": "account"
			}, 
			"message": ""
	}
	
# Chart APIs
#### 获取交易历史图表
	https://api.seelescan.io/api/v1/chart/tx
	
#### 返回
1. code: 错误码,0为正常,非0为错误
2. message: 错误提示,正确执行会空
3. data: 返回按日期排序的每天交易信息列表

#### 例子
	//Request
	https://api.seelescan.io/api/v1/chart/tx
	
	//Return
	{
		"code": 0, 
		"data": [
			{
				"TotalTxs": 6095, 
				"TotalBlocks": 6095, 
				"HashRate": 1216530.3790393518, 
				"Difficulty": 17244991.75537326, 
				"AvgTime": 14.175553732567678, 
				"Rewards": 121900000000000, 
				"TotalAddresss": 2, 
				"TodayIncrease": 0, 
				"TimeStamp": 1527350400
			}
		], 
		"message": ""
	}

#### 获取区块难度增长图表
	https://api.seelescan.io/api/v1/chart/difficulty
	
#### 返回
1. code: 错误码,0为正常,非0为错误
2. message: 错误提示,正确执行会空
3. data: 返回按日期排序的每天区块难度增长列表

#### 例子
	//Request
	https://api.seelescan.io/api/v1/chart/difficulty
	
	//Return
	{
		"code": 0, 
		"data": [
			{
				"Difficulty": 17244991.75537326, 
				"TimeStamp": 1527350400
			}
		], 
		"message": ""
	}

#### 获取地址增长图表
	https://api.seelescan.io/api/v1/chart/address
	
#### 返回
1. code: 错误码,0为正常,非0为错误
2. message: 错误提示,正确执行会空
3. data: 返回按日期排序的每天地址增长列表

#### 例子
	//Request
	https://api.seelescan.io/api/v1/chart/address
	
	//Return
	{
		"code": 0, 
		"data": [
			{
				"TotalAddresss": 2, 
				"TodayIncrease": 0, 
				"TimeStamp": 1527350400
			}
		], 
		"message": ""
	}

#### 获取区块数量和奖励图表
	https://api.seelescan.io/api/v1/chart/blocks
	
#### 返回
1. code: 错误码,0为正常,非0为错误
2. message: 错误提示,正确执行会空
3. data: 返回按日期排序的每天区块数量和奖励列表

#### 例子
	//Request
	https://api.seelescan.io/api/v1/chart/blocks
	
	//Return
	{
		"code": 0, 
		"data": [
			{
				"TotalBlocks": 6095, 
				"Rewards": 121900000000000, 
				"TimeStamp": 1527350400
			}
		], 
		"message": ""
	}

#### 获取哈希速率图表
	https://api.seelescan.io/api/v1/chart/hashrate
	
#### 返回
1. code: 错误码,0为正常,非0为错误
2. message: 错误提示,正确执行会空
3. data: 返回按日期排序的哈希速率列表

#### 例子
	//Request
	https://api.seelescan.io/api/v1/chart/hashrate
	
	//Return
	{
		"code": 0, 
		"data": [
			{
				"TotalBlocks": 6095, 
				"Rewards": 121900000000000, 
				"TimeStamp": 1527350400
			}
		], 
		"message": ""
	}
#### 获取区块产生时间图表
	https://api.seelescan.io/api/v1/chart/blocktime
	
#### 返回
1. code: 错误码,0为正常,非0为错误
2. message: 错误提示,正确执行会空
3. data: 返回按日期排序的平均出块时间列表

#### 例子
	//Request
	https://api.seelescan.io/api/v1/chart/blocktime
	
	//Return
	{
		"code": 0, 
		"data": [
			{
				"AvgTime": 14.175553732567678, 
				"TimeStamp": 1527350400
			}
		], 
		"message": ""
	}
	
#### 获取挖矿排行图表
	https://api.seelescan.io/api/v1/chart/miner
	
#### 返回
1. code: 错误码,0为正常,非0为错误
2. message: 错误提示,正确执行会空
3. data: 返回按日期排序的平均出块时间列表

#### 例子
	//Request
	https://api.seelescan.io/api/v1/chart/miner
	
	//Return
	{
		"code": 0, 
		"data": [
			{
				"Rank": [
					{
						"Address": "0x4dd6881d13ab5152127533c5954e4e062eb4bb2dcd93becf4f4e9b1d2d69f1363eea0395e8e76a2716b033d1e3cc8da2bf24811b1e31a86ac8bcacca4c4b29bd", 
						"Mined": 21215, 
						"Percentage": 0.5217658632562715
					}, 
					{
						"Address": "0x23ddfb54a488f906cdb9cbd257eac5663a4c74ba25619bb902651602a4491be4ce437907fcc567b31be6746a014931f4670ac116c0010e5beb28b0dce2c6eaad", 
						"Mined": 19445, 
						"Percentage": 0.4782341367437285
					}
				]
			}
		], 
		"message": ""
	}
