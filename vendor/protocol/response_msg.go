package protocol

/*
{
	"code": 200,        // http状态码
	"enmsg": "ok",      // 报错string常量
	"cnmsg": "导入成功", // 报错信息
	"data": null        // 数据，本接口没有数据
}
*/

// ResponseMsg res msg
type ResponseMsg struct {
	Code  int         `json:"code"`
	EnMsg string      `json:"enmsg"`
	CnMsg string      `json:"cnmsg"`
	Data  interface{} `json:"data"`
}
