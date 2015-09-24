package controllers

import (
	// "fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	// "log"
	"encoding/json"
	"net/http"
	// "net/url"
	"github.com/astaxie/beego/config"
	"qax580go/models"
	"strings"
)

type WxAuthoController struct {
	beego.Controller
}

func (c *WxAuthoController) Get() {

	code := c.Input().Get("code")
	state := c.Input().Get("state")
	beego.Debug("-------")
	if len(code) != 0 && len(state) != 0 {
		beego.Debug(code)
		beego.Debug(state)
		getAccessToken(code, c)
	}
	c.TplNames = "wxautho.html"
}

func (c *WxAuthoController) Post() {
	c.TplNames = "wxautho.html"
}

func getAccessToken(code string, c *WxAuthoController) {
	isdebug := "true"
	iniconf, err := config.NewConfig("json", "conf/myconfig.json")
	if err != nil {
		beego.Debug(err)
	} else {
		isdebug = iniconf.String("qax580::isdebug")
	}
	wx_url := "[REALM]?appid=[APPID]&secret=[SECRET]&code=[CODE]&&grant_type=authorization_code"
	// if beego.AppConfig.Bool("qax580::isdebug") {
	realm_name := ""
	if isdebug == "true" {
		realm_name = "http://localhost:9090"
	} else {
		realm_name = "https://api.weixin.qq.com/sns/oauth2/access_token"
	}
	appid := "wx570bbcc8cf9fdd80"
	secret := "c4b26e95739bc7defcc42e556cc7ae42"
	wx_url = strings.Replace(wx_url, "[REALM]", realm_name, -1)
	wx_url = strings.Replace(wx_url, "[APPID]", appid, -1)
	wx_url = strings.Replace(wx_url, "[SECRET]", secret, -1)
	wx_url = strings.Replace(wx_url, "[CODE]", code, -1)
	beego.Debug("----------------get Token --------------------")
	beego.Debug(wx_url)
	resp, err := http.Get(wx_url)
	if err != nil {
		beego.Debug(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	beego.Debug("----------------get Token body--------------------")
	if err != nil {
		beego.Debug(err)
	} else {
		beego.Debug(string(body))
	}

	var atj models.AccessTokenJson
	if err := json.Unmarshal(body, &atj); err == nil {
		beego.Debug("----------------get Token json--------------------")
		beego.Debug(atj)
		if atj.ErrCode == 0 {
			getUserInfo(atj.AccessToken, atj.OpenID, c)
		}
	} else {
		beego.Debug("----------------get Token json error--------------------")
		beego.Debug(err)
	}
}
func getUserInfo(access_toke, openid string, c *WxAuthoController) {
	isdebug := "true"
	iniconf, err := config.NewConfig("json", "conf/myconfig.json")
	if err != nil {
		beego.Debug(err)
	} else {
		isdebug = iniconf.String("qax580::isdebug")
	}
	wx_url := "[REALM]?access_token=[ACCESS_TOKEN]&openid=[OPENID]&lang=zh_CN"
	// if beego.AppConfig.Bool("qax580::isdebug") {
	realm_name := ""
	if isdebug == "true" {
		realm_name = "http://localhost:9091"
	} else {
		realm_name = "https://api.weixin.qq.com/sns/userinfo"
	}
	wx_url = strings.Replace(wx_url, "[REALM]", realm_name, -1)
	wx_url = strings.Replace(wx_url, "[ACCESS_TOKEN]", access_toke, -1)
	wx_url = strings.Replace(wx_url, "[OPENID]", openid, -1)
	beego.Debug("----------------get UserInfo --------------------")
	beego.Debug(wx_url)

	resp, err := http.Get(wx_url)
	if err != nil {
		beego.Debug(err)
	}

	defer resp.Body.Close()
	beego.Debug("----------------get UserInfo body--------------------")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		beego.Debug(err)
	} else {
		beego.Debug(string(body))
	}
	var uij models.UserInfoJson
	if err := json.Unmarshal(body, &uij); err == nil {
		beego.Debug("----------------get UserInfo json--------------------")
		beego.Debug(uij)
		if uij.ErrCode == 0 {
			c.Data["Uij"] = uij
		}
	} else {
		beego.Debug("----------------get UserInfo json error--------------------")
		beego.Debug(err)
	}
}
