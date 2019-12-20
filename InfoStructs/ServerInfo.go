package InfoStructs

import "encoding/json"

//存储于客户端和其他注册中心的的注册中心信息。注册中心发回客户端的心跳数据以及注册中心间的互通数据也要是这个格式
type ServerInfo struct {
	ServiceType    string   `json:"service_type"`    //注册中心接受的服务类型
	Status         int8     `json:"status"`          //注册中心发送信息的标记状态，有正常、停机和服务类型不兼容三种
	ID             string   `json:"id"`              //注册中心的ID
	Addr           string   `json:"addr"`            //注册中心的地址
	RelatedServers []string `json:"related_servers"` //与此注册中心相连的其他注册中心地址
}

//从一个JSON字符串中解析出ClientInfo
func ParseServer(jsonData []byte) (*ServerInfo, error) {
	info := ServerInfo{"", 0, "", "", []string{}}
	err := json.Unmarshal(jsonData, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

//将一个ClientInfo转化为JSON字符串
func (info ServerInfo) String() ([]byte, error) {
	return json.Marshal(info)
}
