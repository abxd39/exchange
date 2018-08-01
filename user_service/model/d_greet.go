package model

import (
	"github.com/GeeTeam/GtGoSdk"
	 cf "digicon/user_service/conf"
)
var Gt *GtGoSdk.Geetest

func init() {
	Gt = GtGoSdk.GeetestLib(cf.GtPrivateKey, cf.GtCaptchaID)
}
