package controllers

import (
	"fmt"
	"github.com/guonaihong/gout"
	"github.com/tidwall/gjson"
)

type GjsonController struct {
	BaseController
}

//type res struct{
//	ID int `json:"id"`
//	NanShengxiao string `json:"nan_shengxiao"`
//	NvShengxiao string `json:"nv_shengxiao"`
//	ZhiShu string `json:"zhishu"`
//	Jieguo string `json:"jieguo"`
//	Pingshu string `json:"pingshu"`
//}

func (c GjsonController) Index() {

	var response = ""

	_ = gout.GET("http://aseests.quhuitu.com/shengxiaoshuju.json").BindBody(&response).Do()

	// 取数组中第一个集合中的 pingshu 字段
	//c.Output(gjson.Get(response, "1.pingshu").String())

	/***************** 结果相同 *************************/
	// 取数组中所有集合的 pingshu 字段的值
	//c.Output(gjson.Get(response, "#.pingshu").Value())

	// 数组
	item := gjson.Get(response, "#.pingshu").Array()
	for key, val := range item {
		fmt.Println(key, val.String())
	}
	/***************** 结果相同 *************************/

	c.Output(item)
}
