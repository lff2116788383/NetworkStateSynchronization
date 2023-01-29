package server

import (
	"sync"
)

type ServerInfo struct {
	Maps sync.Map // userId *User
	OnlineMap sync.Map
}

type GameServerTest struct {
	Info ServerInfo

}
