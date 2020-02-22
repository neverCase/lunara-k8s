package service

import (
	"github.com/gin-gonic/gin"
)

type auth struct {
	Status    int `json:"status"`
	LoginTime int `json:"login_time"`
}

func (h *httpService) checkAuth(login *loginRequest) bool {
	groups := map[string]string{
		"admin": "123456",
		"user2": "123456",
	}
	if login.Account == "" || login.Password == "" {
		return false
	}
	if t, ok := groups[login.Account]; ok {
		if t == login.Password {
			return true
		}
	}
	return false
}

//func (h *httpService) getCookie(c *gin.Context) (signed bool, res httpResponse) {
//	session := sessions.Default(c)
//	if v := session.Get(_authKey); v != nil {
//		var (
//			authCache auth
//			str       []byte
//			err       error
//		)
//		if err := json.Unmarshal(v.([]byte), &authCache); err != nil {
//			log.Printf("handleSignOut Unmarshal data:%v err:%v\n", v.(string), err)
//			return false, h.getResponse(errUnknown, msgFailed)
//		}
//		log.Printf("handleSignOut -------- cookie data:%v \n", v)
//		authCache.Status = authInActive
//		if str, err = json.Marshal(authCache); err != nil {
//			log.Printf("handleSignIn Marshal data:%v err:%v\n", string(str), err)
//			return h.getResponse(errUnknown, msgFailed)
//		}
//		session.Delete(_authKey)
//		session.Options(sessions.Options{MaxAge: _cookieTimeout})
//		session.Save()
//		return h.getResponse(errNone, msgSuccess)
//	} else {
//		log.Println("no session")
//		return h.getResponse(errUnknown, msgFailed)
//	}
//}

func (h *httpService) setCookie(c *gin.Context) {

}