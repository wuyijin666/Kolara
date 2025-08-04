package utils

import (
	"Kolara/kiface"
	"encoding/json"
	"os"
)

/**
 * 存储有关kolara框架的全局参数
 * 该参数可由用户在 Kolara.json中进行配置
 */

 type GlobalObj struct {
	/*
	* Server
	*/
	TcpServer kiface.IServer // 当前Kolara框架的全局Server对象
	Host  string               // 当前服务器监听的主机IP
	TcpPort int           // 当前服务器监听的主机端口号
	Name string              // 当前服务器名称
	MaxConn int              // 服务器端支持的最大连接数
	/*
	* Kolara
	*/
	Version string           // 当前Kolara框架的版本号
	MaxPackageSize uint32    // 当前框架允许的最大数据包大小
}




// 先从conf/kolara.json中加载数据 若无，则选择默认数据
func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("conf/kolara.json")
	if err != nil {
		panic(err)
	}
	// 将json数据解析到struct中
	if err := json.Unmarshal(data, &GlobalObject); err != nil {
		panic(err)
	}

}

// 定义一个全局对外的GlobalObj对象
var GlobalObject *GlobalObj
func init() {
	GlobalObject = &GlobalObj {
       // TcpServer : nil,
		Host :  "0.0.0.0",
		TcpPort : 8999,
		Name : "KolaraServerApp",
		MaxConn : 1000,
		Version : "V0.3",
		MaxPackageSize : 4096,
	}

	GlobalObject.Reload()

}

