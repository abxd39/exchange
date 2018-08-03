package controller

import (
	"crypto/sha1"
	"digicon/proto/common"
	"digicon/ws_service/model"
	"digicon/ws_service/rpc/client"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gopkg.in/olahol/melody.v1"
	"time"
)

type WebChatGroup struct {
	m *melody.Melody
}

func (this *WebChatGroup) Router(r *gin.Engine) {
	this.m = melody.New()
	WebChat := r.Group("/ws")
	{
		//WebChat.GET("/web_chat", this.WebChat)
		WebChat.GET("/channel/:channelid", this.WSChannel)
	}
}

type Message struct {
	InfoType int32  `form:"info_type"   json:"info_type"  binding:"required"`  // 消息类型   ,1 认证消息，2，内容消息
	Token    string `form:"token"       json:"token"       `                   // token验证
	OrderId  string `form:"order_id"    json:"order_id"    binding:"required"` // 订单ID
	SellerId uint64 `form:"seller_id"    json:"seller_id"  `                   // 卖家id
	Buyer_id uint64 `form:"buyer_id"     json:"buyer_id"`                      // 买家id
	Uid      uint64 `form:"uid"         json:"uid"`                            // 当前聊天id
	UserName string `form:"username"    json:"username" `                      // 当前聊天用户名
	Content  string `form:"content"     json:"content"`                        // 聊天内容
}

type ErrorDT struct{}

type ErrorRspMessage struct {
	Code     int32   `json:"code"`
	Data     ErrorDT `json:"data"`
	Msg      string  `json:"msg"`
	RespType int32   `json:"resp_type"` // 消息类型(1: 系统消息, 2: 用户消息)
}

type RespMessage struct {
	InfoType    int32  `form:"info_type"   json:"info_type"  binding:"required"`  // 消息类型   ,1 认证消息，2，内容消息, 3 关闭连接, 4, 订单已完成, 5 订单已取消
	OrderId     string `form:"order_id"    json:"order_id"    binding:"required"` // 订单ID
	Uid         uint64 `form:"uid"         json:"uid"  `                          // 用户ID
	UserName    string `form:"username"   json:"username" `
	Content     string `form:"content"     json:"content"` // neirong
	CreatedTime string `xorm:"DATETIME" json:"created_time"`
}

type ResponseMessage struct {
	Code     int32       `json:"code"`      // 0: 成功,  1: 未知错误, 2 : 参数错误,   201: 登陆失效
	RespType int32       `json:"resp_type"` // 消息类型(0: 系统消息, 1: 用户消息)
	Data     RespMessage `json:"data"`
	Msg      string      `json:"msg"`
}

/*

 */
func (this *WebChatGroup) WebChat(c *gin.Context) {
	this.m.HandleRequest(c.Writer, c.Request)
	this.m.HandleMessage(func(s *melody.Session, msg []byte) {
		fmt.Println(string(msg))
		this.m.Broadcast(msg)
	})
}

/*
	func:
*/
func (this *WebChatGroup) WSChannel(c *gin.Context) {
	this.m.HandleRequest(c.Writer, c.Request)
	this.m.HandleMessage(func(s *melody.Session, msg []byte) {
		// todo msg
		var mesg Message
		if err := json.Unmarshal(msg, &mesg); err == nil {
			log.Errorln("mesg:", mesg)
			switch mesg.InfoType {
			// 认证
			case 1:
				if this.CheckAuth(mesg.Token, mesg.Uid) {
					hashChannelId := this.GenerateHashChannelId(mesg)
					s.Set("channelId", hashChannelId)
					message := &ErrorRspMessage{
						Code: 0,
						Msg:  "成功",
					}
					data, _ := json.Marshal(message)
					s.Write(data)
				}
			// 发送消息
			case 2:
				channelid, _ := s.Get("channelId")
				hashChannelId := this.GenerateHashChannelId(mesg)
				if channelid == hashChannelId {
					go this.ChatBroadCast(s, mesg, msg)
				} else {
					this.CloseSession(s, 201, "auth error!")
				}
			// 关闭聊天
			case 3:
				channelid, _ := s.Get("channelid")
				hashChannelId := this.GenerateHashChannelId(mesg)
				if channelid == hashChannelId {
					this.CloseSession(s, 0, "close connect!")
				}
			// 订单完成
			case 4:
				channelid, _ := s.Get("channelId")
				hashChannelId := this.GenerateHashChannelId(mesg)
				if channelid == hashChannelId {
					go this.OrderBroadCast(s, mesg, msg, 4)
				} else {
					this.CloseSession(s, 201, "auth error!")
				}
			// 订单取消
			case 5:
				channelid, _ := s.Get("channelId")
				hashChannelId := this.GenerateHashChannelId(mesg)
				if channelid == hashChannelId {
					go this.OrderBroadCast(s, mesg, msg, 5)
				} else {
					this.CloseSession(s, 201, "auth error!")
				}
			default:
				this.CloseSession(s, 1, "not found message type!")
			}
		} else {
			this.CloseSession(s, 2, "data struct error!")
		}

	})
}

/*
	generate hash channelid
*/
func (this *WebChatGroup) GenerateHashChannelId(mesg Message) (hashChannelId string) {
	channelid := fmt.Sprintf("%v-%v-%v", mesg.OrderId, mesg.SellerId, mesg.Buyer_id)
	Sha1Inst := sha1.New()
	Sha1Inst.Write([]byte(channelid))
	hashChannelId = fmt.Sprintf("%x", Sha1Inst.Sum([]byte("")))
	//fmt.Println("hashChannelId:", hashChannelId)
	return
}

/*
	校验token
*/
func (this *WebChatGroup) CheckAuth(token string, uid uint64) bool {
	log.Println("token:", token, uid)
	rsp, err := client.InnerService.UserService.CallTokenVerify(uid, []byte(token))
	//fmt.Println("rsp:", rsp, err)
	if err != nil {
		fmt.Println(err.Error())
		log.Errorln(err.Error())
		return false
	}
	if rsp.Err == errdefine.ERRCODE_SUCCESS {
		fmt.Println("success")
		return true
	} else {
		log.Println("auth token error!")
		log.Println(rsp.Err, rsp.Message)
		return false
	}
}

/*
func : close
*/
func (this *WebChatGroup) CloseSession(s *melody.Session, code int32, msg string) {
	closesMsg := &ErrorRspMessage{
		Code:     code,
		RespType: 0,
		Msg:      msg,
	}
	data, err := json.Marshal(closesMsg)
	err = this.m.BroadcastFilter(data, func(q *melody.Session) bool {
		qv, _ := q.Get("channelId")
		sv, _ := s.Get("channelId")
		if qv == sv {
			return true
		} else {
			return false
		}
	})
	if err != nil {
		fmt.Println(err.Error())
		log.Errorln(err.Error())
	}
	time.Sleep(5 * time.Second)
	s.Close()

}

/*
	订单相关消息
*/
func (this *WebChatGroup) OrderBroadCast( s *melody.Session, mesg Message, msg []byte, RespType int32){
	var content string
	if RespType == 4 {
		content = "订单已完成!"
	}else if RespType == 5 {
		content = "订单已取消!"
	}else{
		content = ""
	}

	rmsg := RespMessage{
		Content:     content,
		UserName:    mesg.UserName,
		Uid:         mesg.Uid,
		OrderId:     mesg.OrderId,
		InfoType:    RespType,
		CreatedTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	message := &ResponseMessage{
		Code:     0,
		Data:     rmsg,
		RespType: 0,
		Msg:      "成功",
	}
	data, _ := json.Marshal(message)
	this.m.BroadcastFilter(data, func(q *melody.Session) bool {
		qv, _ := q.Get("channelId")
		sv, _ := s.Get("channelId")
		if qv == sv {
			return true
		} else {
			return false
		}
	})
}




/*
	send message
*/
func (this *WebChatGroup) ChatBroadCast(s *melody.Session, mesg Message, msg []byte) {
	rmsg := RespMessage{
		Content:     mesg.Content,
		UserName:    mesg.UserName,
		Uid:         mesg.Uid,
		OrderId:     mesg.OrderId,
		InfoType:    mesg.InfoType,
		CreatedTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	message := &ResponseMessage{
		Code:     0,
		Data:     rmsg,
		RespType: 1,
		Msg:      "成功",
	}
	data, _ := json.Marshal(message)
	go SaveChatMsg(mesg)
	this.m.BroadcastFilter(data, func(q *melody.Session) bool {
		qv, _ := q.Get("channelId")
		sv, _ := s.Get("channelId")
		if qv == sv {
			return true
		} else {
			return false
		}
	})
}

/*

 */
func SaveChatMsg(mesg Message) {
	log.Errorln("go run to save msg :", mesg)
	chat := new(model.Chats)
	chat.OrderId = mesg.OrderId
	chat.Uid = mesg.Uid
	chat.Uname = mesg.UserName
	chat.Content = mesg.Content
	chat.States = 1
	chat.CreatedTime = time.Now().Format("2006-01-02 15:04:05")
	code := chat.Add()
	fmt.Println(code)
	log.Println("write to mysql code:", code)

}
