/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package node

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/seeleteam/scan-api/database"
	"github.com/seeleteam/scan-api/log"
	"github.com/seeleteam/scan-api/rpc"
	mgo "gopkg.in/mgo.v2"
)

const (
	geoURL = "http://www.geoplugin.net/json.gp?ip="
)

//geoPluginRet the reulst return by the geoplugin website
type geoPluginRet struct {
	GeopluginRequest                string `json:"Geoplugin_request"`
	GeopluginStatus                 int    `json:"Geoplugin_status"`
	GeopluginDelay                  string `json:"Geoplugin_delay"`
	GeopluginCredit                 string `json:"Geoplugin_credit"`
	GeopluginCity                   string `json:"Geoplugin_city"`
	GeopluginRegion                 string `json:"Geoplugin_region"`
	GeopluginRegionCode             string `json:"Geoplugin_regionCode"`
	GeopluginRegionName             string `json:"Geoplugin_regionName"`
	GeopluginAreaCode               string `json:"Geoplugin_areaCode"`
	GeopluginDmaCode                string `json:"Geoplugin_dmaCode"`
	GeopluginCountryCode            string `json:"Geoplugin_countryCode"`
	GeopluginCountryName            string `json:"Geoplugin_countryName"`
	GeopluginInEU                   int    `json:"Geoplugin_inEU"`
	GeopluginContinentCode          string `json:"Geoplugin_continentCode"`
	GeopluginContinentName          string `json:"Geoplugin_continentName"`
	GeopluginLatitude               string `json:"Geoplugin_latitude"`
	GeopluginLongitude              string `json:"Geoplugin_longitude"`
	GeopluginLocationAccuracyRadius string `json:"Geoplugin_locationAccuracyRadius"`
	GeopluginTimezone               string `json:"Geoplugin_timezone"`
}

//NodeService is the find node service
type NodeService struct {
	nodeMap     map[string]database.DBNodeInfo
	nodeMapLock sync.Mutex

	nodeDB NodeDB
	cfg    *Config
}

func New(cfg *Config, nodeDB NodeDB) *NodeService {
	return &NodeService{
		nodeDB:  nodeDB,
		cfg:     cfg,
		nodeMap: make(map[string]database.DBNodeInfo),
	}
}

//getGeoInfoByHTTP get location information by node ip
func getGeoInfoByHTTP(ip string) (*geoPluginRet, error) {
	getURL := geoURL + ip
	resp, err := http.Get(getURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ret geoPluginRet
	err = json.Unmarshal(body, &ret)
	return &ret, nil
}

//ProcessSinglePeer convert a rpc.peerinfo to nodeinfo
func (n *NodeService) ProcessSinglePeer(peer *rpc.PeerInfo, c chan int) {
	defer func() {
		c <- 1
	}()
	var nodeInfo database.DBNodeInfo
	nodeInfo.ID = peer.ID
	var caps string
	for i := 0; i < len(peer.Caps); i++ {
		if i != 0 {
			caps += "|"
		}
		caps += peer.Caps[i]
	}
	nodeInfo.Caps = caps
	fmt.Println(peer.RemoteAddress)
	ipAndPort := strings.Split(peer.RemoteAddress, ":")
	nodeInfo.Host = ipAndPort[0]
	nodeInfo.Port = ipAndPort[1]
	geo, err := getGeoInfoByHTTP(nodeInfo.Host)
	if err != nil {
		return
	}
	nodeInfo.Country = geo.GeopluginCountryName
	nodeInfo.Region = geo.GeopluginRegionName
	nodeInfo.City = geo.GeopluginCity
	nodeInfo.LastSeen = time.Now().Unix()
	nodeInfo.LongitudeAndLatitude = string('[') + geo.GeopluginLongitude + string(',') + geo.GeopluginLatitude + string(']')
	nodeInfo.ShardNumber = peer.ShardNumber
	if nodeInfo.ShardNumber <= 0 {
		nodeInfo.ShardNumber = 1
	}

	n.nodeMapLock.Lock()
	n.nodeMap[nodeInfo.Host] = nodeInfo
	n.nodeMapLock.Unlock()
	_, err = n.nodeDB.GetNodeInfoByID(nodeInfo.ID)
	if err == mgo.ErrNotFound {
		n.nodeDB.AddNodeInfo(&nodeInfo)
	}
}

//DeleteExpireNode if an node does not appear for a long time, remove ti from the database and nodemap
func (n *NodeService) DeleteExpireNode() {
	now := time.Now().Unix()
	for k, v := range n.nodeMap {
		if now-v.LastSeen > n.cfg.ExpireTime {
			n.nodeDB.DeleteNodeInfo(&v)
			delete(n.nodeMap, k)
		}
	}
}

//FindNode get all peers info and store them into database
func (n *NodeService) FindNode() {

	var allPeerInfos []rpc.PeerInfo
	for i := 0; i < len(n.cfg.RPCNodes); i++ {
		rpcURL := n.cfg.RPCNodes[i]

		rpc := rpc.NewRPC(rpcURL)
		defer func() {
			if rpc != nil {
				rpc.Release()
			}
		}()

		if rpc == nil {
			continue
		}

		if err := rpc.Connect(); err != nil {
			fmt.Printf("rpc init failed, connurl:%v\n", rpcURL)
			continue
		}

		peerInfos, err := rpc.GetPeersInfo()
		if err != nil {
			log.Fatal(err)
			continue
		}

		allPeerInfos = append(allPeerInfos, peerInfos...)
	}

	if len(allPeerInfos) == 0 {
		return
	}

	cnum := make(chan int, len(allPeerInfos))
	for i := 0; i < len(allPeerInfos); i++ {
		peer := allPeerInfos[i]
		if v, ok := n.nodeMap[peer.ID]; ok {
			v.LastSeen = time.Now().Unix()
			cnum <- 1
		} else {
			go n.ProcessSinglePeer(&peer, cnum)
		}
	}

	for i := 0; i < len(allPeerInfos); i++ {
		<-cnum
	}
}

//RestoreNodeFromDB restore data from database into nodemap
func (n *NodeService) RestoreNodeFromDB() {
	nodes, err := n.nodeDB.GetNodeInfos()
	if err != nil {
		return
	}

	now := time.Now().Unix()
	for i := 0; i < len(nodes); i++ {
		nodes[i].LastSeen = now
		n.nodeMap[nodes[i].ID] = *nodes[i]
	}
}

//StartFindNodeService start the node map service
func (n *NodeService) StartFindNodeService() {
	n.RestoreNodeFromDB()
	n.FindNode()

	ticks := time.NewTicker(n.cfg.Interval * time.Second)
	tick := ticks.C
	go func() {
		for range tick {
			n.FindNode()
			_, ok := <-tick
			if !ok {
				break
			}
		}
	}()
}
