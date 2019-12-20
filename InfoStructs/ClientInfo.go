package InfoStructs

import "encoding/json"

//存储于服务器端的用户信息。客户端发来的心跳数据也要是这个格式
type ClientInfo struct {
	ServiceType string `json:"service_type"` //客户端的服务类型，和Server类的服务类型含义相同
	Status      int8   `json:"status"`       //从注册中心一侧看到的客户端状态，有正常和停机
	ID          string `json:"id"`           //客户端的ID
	Addr        string `json:"addr"`         //客户端的地址
	MaxRegister uint64 `json:"max_register"` //客户端最多可以连多少个服务端
}

//从一个JSON字符串中解析出ClientInfo
func ParseClient(jsonData []byte) (*ClientInfo, error) {
	info := ClientInfo{"", 0, "", "", 0}
	err := json.Unmarshal(jsonData, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

//将一个ClientInfo转化为JSON字符串
func (info ClientInfo) String() ([]byte, error) {
	return json.Marshal(info)
}
