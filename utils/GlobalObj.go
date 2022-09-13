package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"test/iface"
)

type GlobalObj struct {
	TcpServer iface.IServer
	Host      string
	TcpPort   int
	Name      string

	Version        string
	MaxConn        int
	MaxPackageSize uint32

	WorkPoolSize   uint32
	MaxWorkTaskLen uint32
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {

	path, _ := filepath.Abs("src/zinx/testDemo/conf/zinx.json")
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}

}

func init() {

	GlobalObject = &GlobalObj{
		Name:           "testServer",
		Host:           "192.168.1.106",
		TcpPort:        8999,
		Version:        "v0.4",
		MaxConn:        1000,
		MaxPackageSize: 4096,
		WorkPoolSize:   10,
		MaxWorkTaskLen: 1024,
	}

	//GlobalObject.Reload()

}
