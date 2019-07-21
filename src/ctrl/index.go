package ctrl

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"util"
)

func IndexGetDoReptiles(params martini.Params, req *http.Request, r render.Render) {
	start := params["start"]
	end := params["end"]
	util.Anime(start, end)

	r.JSON(200, map[string]interface{}{
		"insert_id": "",
	})
	return

}
