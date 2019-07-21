package ctrl

import (
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"io/ioutil"
	"model"
	"net/http"
	"regexp"
	"strconv"
	"tools"

	"github.com/martini-contrib/render"
)

type ByDrama []model.Drama

func (s ByDrama) Len() int {
	return len(s)
}
func (s ByDrama) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByDrama) Less(i, j int) bool {
	reg := regexp.MustCompile(`.*([0-9]+).*`)
	si := reg.FindAllStringSubmatch(s[i].Name, 1)
	sj := reg.FindAllStringSubmatch(s[j].Name, 1)
	siInt, _ := strconv.Atoi(si[0][1])
	sjInt, _ := strconv.Atoi(sj[0][1])
	return siInt < sjInt
}
func getAnimeByAreaAndYear(area string, year string) (rs []model.AnimeInfo) {
	dbConn := tools.GetDefDb()

	anime := []model.Anime{}
	if year > "" {
		_, err := dbConn.DbMap.Select(&anime, "select * from anime where area=? and year=? limit 0,11", area, year)
		tools.CheckErr(err)
	} else {
		_, err := dbConn.DbMap.Select(&anime, "select * from anime where  (year=? and area=? ) or area=?  limit 0,11", year, area, area)
		tools.CheckErr(err)
	}

	for _, item := range anime {
		obj := model.AnimeInfo{}

		director := []model.Director{}
		_, err := dbConn.DbMap.Select(&director, "select * from director where anime_id=?", item.Id)
		tools.CheckErr(err)
		obj.Anime = item
		if len(director) > 0 {
			obj.Director = director[0].Name
		}
		star := []model.Star{}
		_, err = dbConn.DbMap.Select(&star, "select * from star where anime_id=?", item.Id)
		tools.CheckErr(err)

		starStr := ""
		if len(star) > 0 {
			for _, o := range star {
				starStr = starStr + " " + o.Name
			}
		}
		obj.Star = starStr

		rs = append(rs, obj)
	}
	return
}

//默认Home页
func DefaultGetHome(req *http.Request, r render.Render) {

	//items := getAllItems()

	jp := getAnimeByAreaAndYear("日本", "2019")
	gc := getAnimeByAreaAndYear("大陆", "2019")
	om := getAnimeByAreaAndYear("欧美", "")

	r.HTML(200, "default/home", map[string]interface{}{
		"title": "I am title",
		"jp":    jp,
		"gc":    gc,
		"om":    om,
		//"items": items,
	})
}

//增加一个记录
func DefaultPostHome(req *http.Request, r render.Render) {
	type Input struct {
		Name string `json:"name"` //需要指定json的名称, 解析的时候根据这个匹配
	}

	//接收客户端post的json, 需要先读取body的内容, 然后把json解析成golang的类型
	b, err := ioutil.ReadAll(req.Body)
	tools.CheckErr(err)

	var data Input
	err = json.Unmarshal(b, &data)
	tools.CheckErr(err)

	fmt.Printf("Input item:%+v\n", data)

	dbConn := tools.GetDefDb()
	res, err := dbConn.DbMap.Exec("insert into `table1` (`name`) values(?)", data.Name)
	tools.CheckErr(err)

	lastId, err := res.LastInsertId()
	tools.CheckErr(err)

	fmt.Println("lastId:", lastId)

	//r.JSON会把第二个参数自动转为json
	r.JSON(200, map[string]interface{}{
		"insert_id": lastId,
	})
}

//修改一个记录
func DefaultPutHome(req *http.Request, r render.Render) {
	type Input struct {
		Id   int
		Name string
	}

	b, err := ioutil.ReadAll(req.Body)
	tools.CheckErr(err)

	var data Input
	err = json.Unmarshal(b, &data)
	tools.CheckErr(err)

	fmt.Printf("Input item:%+v\n", data)

	dbConn := tools.GetDefDb()
	res, err := dbConn.DbMap.Exec("update `table1` set `name`=? where `id`=?", data.Name, data.Id)
	tools.CheckErr(err)

	affectedCount, err := res.RowsAffected()
	tools.CheckErr(err)

	fmt.Println("affectedCount:", affectedCount)

	r.JSON(200, nil)
}

//删除
func DefaultDeleteHome(req *http.Request, r render.Render) {

	req.ParseForm()

	id, _ := strconv.ParseInt(req.Form.Get("id"), 10, 64)

	dbConn := tools.GetDefDb()
	res, err := dbConn.DbMap.Exec("delete from table1 where id=?", id)
	tools.CheckErr(err)

	affectedCount, err := res.RowsAffected()
	tools.CheckErr(err)

	fmt.Println("affectedCount:", affectedCount)

	r.JSON(200, nil)

}
func getAnimeById(id string) (rs []model.AnimeInfo) {
	dbConn := tools.GetDefDb()

	anime := []model.Anime{}
	_, err := dbConn.DbMap.Select(&anime, "select * from anime where id=?", id)
	tools.CheckErr(err)

	for _, item := range anime {
		obj := model.AnimeInfo{}

		director := []model.Director{}
		_, err := dbConn.DbMap.Select(&director, "select * from director where anime_id=?", item.Id)
		tools.CheckErr(err)
		obj.Anime = item
		if len(director) > 0 {
			obj.Director = director[0].Name
		}
		star := []model.Star{}
		_, err = dbConn.DbMap.Select(&star, "select * from star where anime_id=?", item.Id)
		tools.CheckErr(err)

		starStr := ""
		if len(star) > 0 {
			for _, o := range star {
				starStr = starStr + " " + o.Name
			}
		}
		obj.Star = starStr

		drama := []model.Drama{}
		_, err = dbConn.DbMap.Select(&drama, "select * from drama where anime_id=?", item.Id)
		tools.CheckErr(err)
		//fmt.Println(drama)
		dramaMap := map[string][]model.Drama{}

		for _, item := range drama {
			fmt.Println(item)
			dramaMap[item.Source] = append(dramaMap[item.Source], item)
		}
		//for _, item := range dramaMap {
		//	sort.Sort(ByDrama(item))
		//}
		fmt.Println(dramaMap)
		obj.Drama = dramaMap
		rs = append(rs, obj)
	}
	return
}
func DefaultGetDetail(params martini.Params, req *http.Request, r render.Render) {
	id := params["id"]
	fmt.Println(id)
	anime := getAnimeById(id)

	fmt.Println(anime)
	if len(anime) == 0 {
		r.Redirect("/404", 200)
	} else {
		r.HTML(200, "default/detail", map[string]interface{}{
			"title": "I am title",
			"anime": anime[0],
		})
	}

}
func DefaultGetPlay(req *http.Request, r render.Render) {
	r.HTML(200, "default/play", map[string]interface{}{
		"title": "I am title",
		//"items": items,
	})
}
