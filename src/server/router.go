package main

import (
	"ctrl"
	"github.com/go-martini/martini"
)

//路由设置
//Controller函数命名:<Controller><Method><Action>
func InitRoute(m *martini.ClassicMartini) {
	m.Get("/", ctrl.DefaultGetHome)
	m.Get("/detail/:id", ctrl.DefaultGetDetail)
	m.Get("/play", ctrl.DefaultGetPlay)
	//m.Post("/", ctrl.DefaultPostHome) //Post一般负责添加的
	//m.Put("/", ctrl.DefaultPutHome)   //Put一般负责修改的
	//m.Delete("/", ctrl.DefaultDeleteHome)
	m.Get("/doReptiles/:start/:end", ctrl.IndexGetDoReptiles)
}
