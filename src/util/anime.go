package util

import (
	"fmt"
	"model"
	"regexp"
	"strconv"
	"sync"
	"time"
)

var (
	wgGetPage sync.WaitGroup
	wgSave    sync.WaitGroup
	dataChan  = make(chan model.AnimeData, 2000)
	doChan    = make(chan int, 10)

	pathUrl = "http://www.aiqidm.com"
	Area    = "日本"
)

func Anime(start, end string) {
	startInt, _ := strconv.Atoi(start)
	endInt, _ := strconv.Atoi(end)
	for i := startInt; i <= endInt; i++ {
		fmt.Printf("第%d页开始\n", i)
		fmt.Printf("到第%d页\n", endInt)
		url := pathUrl + "/list/ribendonghua/index-" + strconv.Itoa(i) + ".html"
		if i == 1 {
			url = pathUrl + "/list/ribendonghua/"
		}

		wgGetPage.Add(1)
		go func(pageUrl string) {
			GetPageData(pageUrl)
			wgGetPage.Done()
		}(url)

	}
	go func() {
		wgGetPage.Wait()
		close(dataChan)
		fmt.Println("chan close")
	}()

	for item := range dataChan {

		wgSave.Add(1)
		go func(AnimeData model.AnimeData) {
			doChan <- 1
			str := extractHandle(AnimeData.Picture, `/([0-9a-z]+\.[a-z]+)`, 1)
			nowTime := int(time.Now().Unix())
			timestr := strconv.Itoa(nowTime)
			path := "./public/upload/anime"
			imgPath := path + "/" + timestr + ".jpg"
			img := "/upload/anime/" + timestr + ".jpg"
			if len(str) > 0 {
				imgPath = path + "/" + str[0]
				img = "/upload/anime/" + str[0]
			}
			SaveImg(AnimeData.Picture, imgPath, path, str)
			Save2Mysql(AnimeData, img)
			<-doChan
			wgSave.Done()
		}(item)

	}
	wgSave.Wait()
}
func GetPageData(url string) {
	//得到每一页动漫数据
	rs, err := HttpGet(url)
	if err != nil {
		return
	}

	//iconv.Convert(input, out, "gb2312", "utf-8")
	//utf8 := mahonia.NewDecoder("gbk").ConvertString(string(rs))

	reg := regexp.MustCompile(`.*<li><a href="(.*)"><img alt=".*`)
	allUrl := reg.FindAllStringSubmatch(string(rs), -1)

	for i, item := range allUrl {
		fmt.Println("第" + strconv.Itoa(i+1) + "个：" + item[1])
		//处理单个动漫
		//if item[1] == "/vod/detail/id/2337.html" {
		//	continue
		//}
		//if item[1] == `/rss.xml" title="RSS订阅` {
		//	continue
		//}
		AnimeItemHandle(i+1, item[1])
	}

}

//处理单个动漫
func AnimeItemHandle(i int, url string) {
	url = pathUrl + url
	rs, err := HttpGet(url)
	if err != nil {
		return
	}
	AnimeData := model.AnimeData{}
	//utf8 := mahonia.NewDecoder("gbk").ConvertString(string(rs))
	AnimeData = AnimeItemExtractHandle(rs)
	dataChan <- AnimeData

}

//单个动漫数据提取
func AnimeItemExtractHandle(rs string) model.AnimeData {
	//title
	title := extractHandle(rs, `<p class="names">(?s:(.*?))</p>`, 1)
	//em_num
	em_num := extractHandle(rs, `<p class="red">最近更新：更新到(?s:(.*?))</p>`, 1)
	//year
	year := extractHandle(rs, `<p>年代：(?s:(.*?))</p>`, 1)
	//area
	//area := extractHandle(rs, `<span>地区：.">(?s:(.*?))</a>`, 1)
	area := Area
	//star 多个
	stars := extractHandle(rs, `<p>配音：(?s:(.*?))</p>`, 1)
	//director
	director := extractHandle(rs, `<p>原著作者:(?s:(.*?))</p>`, 1)
	//picture
	picture := extractHandle(rs, `<img src="(.*)" alt="`, 1)
	//drama 剧集，多个，不同源
	dramas := extractHandle(rs, `<li><a href="(.*)</li>.*`, -1)
	//introduction 简介
	introduction := extractHandle(rs, `<div class="xxjs_con">(?s:(.*?))</div>`, 1)

	//starStr := ""
	//if len(stars) > 0 {
	//	starStr = stars[0]
	//}
	//star := extractHandle(starStr, `target="_blank">(?s:(.*?))</a>`, -1)
	fmt.Println("dramas:", dramas)
	fmt.Println("picture:", picture)
	drama := ExtractDramaHandle(dramas)
	arr := LinkHandle(drama)

	AnimeData := model.AnimeData{}
	if len(title) > 0 {
		AnimeData.Title = title[0]
	}
	if len(em_num) > 0 {
		AnimeData.EmNum = em_num[0]
	}
	if len(year) > 0 {
		AnimeData.Year = year[0]
	}
	AnimeData.Area = area
	if len(stars) > 0 {
		AnimeData.Star = stars[0]
	}
	if len(director) > 0 {
		AnimeData.Director = director[0]
	}

	if len(picture) > 0 {
		AnimeData.Picture = picture[0]
	}
	if len(introduction) > 0 {
		AnimeData.Introduction = introduction[0]
	}
	fmt.Println(arr)
	AnimeData.Drama = arr

	return AnimeData
}

//提取数据函数
func extractHandle(rs, regStr string, num int) (content []string) {
	reg := regexp.MustCompile(regStr)
	allUrl := reg.FindAllStringSubmatch(rs, num)
	for _, item := range allUrl {
		//content=item[1]
		content = append(content, item[1])
	}
	return
}

//提取每一集
func ExtractDramaHandle(rs []string) (content []model.Dramas) {
	for _, item := range rs {
		url := extractHandle(item, `(.*)" title="`, 1)
		name := extractHandle(item, `.*">(.*)</a>`, 1)

		fmt.Println("name:", name)
		var drama model.Dramas
		drama.Url = url[0]
		drama.Name = name[0]
		content = append(content, drama)
	}
	return
}

//获取视频链接
func LinkHandle(content []model.Dramas) (DramaData []model.DramaData) {
	for i, item := range content {
		obj := model.DramaData{}
		//处理单个动漫
		obj = ExtractLinkHandle(i, item.Url, item.Name)
		DramaData = append(DramaData, obj)
	}
	return DramaData
}
func ExtractLinkHandle(index int, url, name string) (arr model.DramaData) {
	//url = pathUrl + url
	rs, err := HttpGet(url)
	if err != nil {
		return model.DramaData{}
	}
	//rs = mahonia.NewDecoder("gbk").ConvertString(string(rs))
	//drama 剧集，多个，不同源
	dramas := extractHandle(rs, `Play.init\((.*)\)`, 1)
	obj := model.DramaData{}

	if len(dramas) > 0 {

		//url = "http://www.imomoe.io" + dramas[0]
		//rsjs, err := HttpGet(url)
		//if err != nil {
		//	return model.DramaData{}
		//}
		//playUrl = mahonia.NewDecoder("gbk").ConvertString(string(rsjs))

		playUrl := extractHandle(dramas[0], `'(.*)',`, 1)
		fmt.Println("dramas[0]:", dramas[0])
		fmt.Println("playUrl:", playUrl)
		if len(playUrl) > 0 {
			obj.Url = playUrl[0]
		}

		From := extractHandle(dramas[0], `','(.*)'`, 1)
		if len(From) > 0 {
			obj.From = From[0]
			if From[0] == "youku" || From[0] == "tudou" {
				obj.PlayMethod = "https://v.youku.com/v_show/id_().html"
			} else if From[0] == "qiyi" {
				obj.PlayMethod = "http://dispatcher.video.qiyi.com/common/shareplayer.html?vid=()&tvId=239459700&coop=coop_171_58dm&autoplay=1&bd=1&fullsreen=1"
			}
		}

		obj.Name = name

	}
	return obj
}
