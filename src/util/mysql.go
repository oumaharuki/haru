package util

import (
	"model"
	"tools"
)

func Save2Mysql(AnimeData model.AnimeData, Picture string) {
	dbConn := tools.GetDefDb()

	AnimeDatas := []model.Anime{}
	_, err := dbConn.DbMap.Select(&AnimeDatas, "select * from anime where name=?", AnimeData.Title)
	tools.CheckErr(err)

	if len(AnimeDatas) == 0 {
		AnimeDatasRs, err := dbConn.DbMap.Exec("insert into `anime` (`name`,`em_num`,`year`,`area`,`picture`,`introduction`,`form`"+
			") values(?,?,?,?,?,?,?)", AnimeData.Title, AnimeData.EmNum, AnimeData.Year, AnimeData.Area, Picture, AnimeData.Introduction,
			"pilipili")
		tools.CheckErr(err)
		lastId, err := AnimeDatasRs.LastInsertId()
		tools.CheckErr(err)

		for _, item := range AnimeData.Star {
			stars := []model.Star{}
			_, err := dbConn.DbMap.Select(&stars, "select * from star where name=?", item)
			tools.CheckErr(err)

			if len(stars) == 0 {
				_, err := dbConn.DbMap.Exec("insert into `star` (`name`,`anime_id`"+
					") values(?,?)", item, lastId)
				tools.CheckErr(err)
			}
		}

		for _, item := range AnimeData.Director {
			directors := []model.Director{}
			_, err = dbConn.DbMap.Select(&directors, "select * from director where name=?", item)
			tools.CheckErr(err)
			if len(directors) == 0 {
				_, err := dbConn.DbMap.Exec("insert into `director` (`name`,`anime_id`"+
					") values(?,?)", item, lastId)
				tools.CheckErr(err)
			}
		}

		if len(AnimeData.Drama) > 0 {
			for _, item := range AnimeData.Drama {
				dramas := []model.Drama{}
				_, err = dbConn.DbMap.Select(&dramas, "select * from drama where name=? and play_url=? and source=? and play_method=?",
					item.Name, item.Url, item.From, item.PlayMethod)
				tools.CheckErr(err)
				if len(dramas) == 0 {
					_, err := dbConn.DbMap.Exec("insert into `drama` (`name`,`anime_id`,`play_url`,`source`,play_method"+
						") values(?,?,?,?,?)", item.Name, lastId, item.Url, item.From, item.PlayMethod)
					tools.CheckErr(err)
				}
			}
		}
	} else {
		if len(AnimeData.Drama) > 0 {
			for _, item := range AnimeData.Drama {
				dramasById := []model.Drama{}
				_, err = dbConn.DbMap.Select(&dramasById, "select * from drama where anime_id=?",
					AnimeDatas[0].Id)
				tools.CheckErr(err)
				if len(dramasById) == 0 {
					_, err := dbConn.DbMap.Exec("insert into `drama` (`name`,`anime_id`,`play_url`,`source`"+
						") values(?,?,?,?)", item.Name, AnimeDatas[0].Id, item.Url, item.From)
					tools.CheckErr(err)
				} else {
					dramas := []model.Drama{}
					_, err = dbConn.DbMap.Select(&dramas, "select * from drama where name=? and play_url=? and source=? and anime_id=? and play_method=?",
						item.Name, item.Url, item.From, AnimeDatas[0].Id, item.PlayMethod)
					tools.CheckErr(err)
					if len(dramas) == 0 {
						_, err := dbConn.DbMap.Exec("update `drama` set `name`=?,`play_url`=?,`source`=?,play_method=?"+
							" where anime_id=?", item.Name, item.Url, item.From, item.PlayMethod, AnimeDatas[0].Id)
						tools.CheckErr(err)
					}
				}
			}
		}
	}
}
