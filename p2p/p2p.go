package p2p

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/usgeeus/geecoin.git/blockchain"
	"github.com/usgeeus/geecoin.git/utils"
)

var upgrader = websocket.Upgrader{}

func Upgrade(w http.ResponseWriter, r *http.Request) {
	//Port:3000 will upgrade the reqeust from :4000
	openPort := r.URL.Query().Get("openPort")
	ip := utils.Splitter(r.RemoteAddr, ":", 0)
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return openPort != "" && ip != ""
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	utils.HandlerErr(err)
	initPeer(conn, ip, openPort)
}

func AddPeer(address, port, openPort string, broadcast bool) {
	//Port :4000 is requesting un upgrade from the port:3000 (:4000이 :3000으로  upgrade를 요청해요)
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort), nil)
	utils.HandlerErr(err)
	p := initPeer(conn, address, port)
	if broadcast {
		broadcastNewPeer(p)
		return
	}
	sendNewestBlock(p)
}

func BroadcastNewBlock(b *blockchain.Block) {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	for _, p := range Peers.v {
		notifyNewBlock(b, p)
	}
}

func BoradcastNewTx(tx *blockchain.Tx) {
	for _, p := range Peers.v {
		notifyNewTx(tx, p)
	}
}

func broadcastNewPeer(newPeer *peer) {
	for key, p := range Peers.v {
		if key != newPeer.key {
			payload := fmt.Sprintf("%s:%s", newPeer.key, p.port) //p.port  내포트
			notifyNewPeer(payload, p)
		}
	}
}
