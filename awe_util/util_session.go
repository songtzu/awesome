package awe_util

import (
	"sync"
)

const RedisSessionHashMap  = "userSession"
var sessionMap sync.Map //session的map,string(userid)--->*model.UserSession




