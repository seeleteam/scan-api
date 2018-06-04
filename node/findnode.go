/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package node

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"scan-api/database"
	"scan-api/log"
	"scan-api/rpc"
	"strings"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
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

var (
	nodeMap     map[string]database.DBNodeInfo
	nodeMapLock sync.Mutex
)

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
func ProcessSinglePeer(peer *rpc.PeerInfo, c chan int) {
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
	nodeMapLock.Lock()
	nodeMap[nodeInfo.Host] = nodeInfo
	nodeMapLock.Unlock()
	_, err = database.GetNodeInfo(nodeInfo.Host)
	if err == mgo.ErrNotFound {
		database.AddNodeInfo(&nodeInfo)
	}
}

//DeleteExpireNode if an node does not appear for a long time, remove ti from the database and nodemap
func DeleteExpireNode(cfg *Config) {
	now := time.Now().Unix()
	for k, v := range nodeMap {
		if now-v.LastSeen > cfg.ExpireTime {
			database.DeleteNodeInfo(&v)
			delete(nodeMap, k)
		}
	}
}

//FindNode get all peers info and store them into database
func FindNode(cfg *Config) {

	var allPeerInfos []rpc.PeerInfo
	for i := 0; i < len(cfg.RPCNodes); i++ {
		rpcURL := cfg.RPCNodes[i]
		rpc.RPCURL = rpcURL
		rpcSeeleRPC, err := rpc.GetSeeleRPC()
		if err != nil {
			rpc.ReleaseSeeleRPC()
			log.Error(err)
			continue
		}

		peerInfos, err := rpcSeeleRPC.GetPeersInfo()
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
		ipAndPort := strings.Split(peer.RemoteAddress, ":")
		host := ipAndPort[0]
		if v, ok := nodeMap[host]; ok {
			v.LastSeen = time.Now().Unix()
			cnum <- 1
		} else {
			go ProcessSinglePeer(&peer, cnum)
		}
	}

	for i := 0; i < len(allPeerInfos); i++ {
		<-cnum
	}
}

//RestoreNodeFromDB restore data from database into nodemap
func RestoreNodeFromDB() {
	nodes, err := database.GetNodeInfos()
	if err != nil {
		return
	}

	now := time.Now().Unix()
	for i := 0; i < len(nodes); i++ {
		nodes[i].LastSeen = now
		nodeMap[nodes[i].Host] = *nodes[i]
	}
}

//StartFindNodeService start the node map service
func StartFindNodeService(cfg *Config) {
	RestoreNodeFromDB()
	FindNode(cfg)

	ticks := time.NewTicker(cfg.Interval * time.Second)
	tick := ticks.C
	go func() {
		for range tick {
			FindNode(cfg)
			_, ok := <-tick
			if !ok {
				break
			}
		}
	}()
}

func init() {
	nodeMap = make(map[string]database.DBNodeInfo)
}
