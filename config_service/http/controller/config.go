package controller

import (
	. "digicon/config_service/log"
	"digicon/config_service/model"
	. "digicon/proto/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ConfigGroup struct {
}

func (this *ConfigGroup) Router(r *gin.Engine) {
	Config := r.Group("/config")
	{
		Config.POST("/config_put", this.ConfigPut)
		Config.GET("/config_list", this.ConfigList)

		//Config.POST("/config_put_key", this.ConfigPutOne)
		//Config.DELETE("/config_delete", this.ConfigDelete)

	}
}

func (this *ConfigGroup) ConfigPut(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	var err error

	req := struct {
		PutType int32 `json:"put_type"`
	}{}

	if err := c.ShouldBind(&req); err != nil {
		Log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}

	err = model.PutToConsul(req.PutType)
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}
	ret.SetErrCode(ERRCODE_SUCCESS, GetErrorMessage(ERRCODE_SUCCESS))
	return
}

func (this *ConfigGroup) ConfigList(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()
	req := struct {
		PutType int32 `form:"put_type"    json:"put_type"  binding:"required"`
	}{}

	if err := c.ShouldBindQuery(&req); err != nil {
		Log.Errorln(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, GetErrorMessage(ERRCODE_PARAM))
		return
	}
	result, err := model.GetFromConsul(req.PutType)
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, GetErrorMessage(ERRCODE_UNKNOWN))
		return
	}
	ret.SetDataSection("list", result)
	//ret.SetDataValue(result)
	ret.SetErrCode(ERRCODE_SUCCESS, GetErrorMessage(ERRCODE_SUCCESS))
	return
}

//
//func (this *ConfigGroup) ConfigDelete(c *gin.Context){
//	ret := NewPublciError()
//	defer func(){
//		c.JSON(http.StatusOK, ret.GetResult())
//	}()
//	consulClient := dao.DB.GetConsulCli()
//	fmt.Println(consulClient)
//
//	kv := consulClient.KV()
//	dw, err := kv.Delete("k", nil)
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(dw)
//
//	ret.SetErrCode(ERRCODE_SUCCESS, GetErrorMessage(ERRCODE_SUCCESS))
//	return
//}
//
//func (this *ConfigGroup) ConfigPutOne(c *gin.Context) {
//	ret := NewPublciError()
//	defer func() {
//		c.JSON(http.StatusOK, ret.GetResult())
//	}()
//
//	mconfig := new(model.ConfigQuenesModel)
//	mconfig.GetAllQuenes()
//
//	consulClient := dao.DB.GetConsulCli()
//	fmt.Println(consulClient)
//
//	kv := consulClient.KV()
//	p := &api.KVPair{Key: "k", Value: []byte("vvvvvvvvv")}
//	wm, err := kv.Put(p, nil)
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(wm)
//
//	p = &api.KVPair{Key: "ads", Value: []byte("vvvvvvvasdfffffffffffffljlkfdjlaksjflkasjdflkkkkkkkkkkkkkkkkkkkkkkkkkkkkkksvv")}
//	wm, err = kv.Put(p, nil)
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(wm)
//
//	ret.SetErrCode(ERRCODE_SUCCESS, GetErrorMessage(ERRCODE_SUCCESS))
//	return
//}
