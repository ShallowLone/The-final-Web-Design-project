package main

import (
	"fmt"
	"strconv"
	"html/template"
	"net/http"
    "net/url"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

/*---------------------------------------------------------------------------
由于初期设计失误，没有使用父类子类的方式写Database相关操作和模型
有空记得补一下

---------------------------------------------------------------------------*/

//-------------------------------------------//
//数学函数重写
func min(a int, b int) int {
	if a < b {
		return b
	} else {
		return a
	}
}
func max(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

//string转化int
func StrToInt(str string) int {
	result, _ := strconv.Atoi(str)
	return result
}

//int转化为string
func IntToString(num int) string {
	result := strconv.Itoa(num)
	return result
}

//-------------------------------------------//
//初始化数据库
func Init() {
	dsn := "root:Kagarino_Kirie@tcp(127.0.0.1:3306)/renga?charset=utf8mb3&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("DB ERROR~!!")
	}
}
func GetNextId() int {
	ans := 1
	for true {
		if ans > 1000000 {
			return -1
		}
		temp := []Game{}
		DB.Where("id = ?", ans).Find(&temp)
		if len(temp) == 0 {
			return ans
		}
		ans++
	}
	return ans
}
func GetNextCid() int {
	ans := 1
	for true {
		if ans > 1000000 {
			return -1
		}
		temp := []Character{}
		DB.Where("cid = ?", ans).Find(&temp)
		if len(temp) == 0 {
			return ans
		}
		ans++
	}
	return ans
}
func GameWithItsChara(id int)[]Character{
    temp := []Character{}
    DB.Where("id = ?",id).Find(&temp)
    return temp
}
//-------------------------------------------//
//接受图片
func GetCharaPictrue(c *gin.Context) {
	file, err := c.FormFile("imgfile")
	if err != nil {
		return
	}
	if err != nil {
		println("error")
	}
	println(file.Filename)
	//dst := path.Join("/img/chara/", c.PostForm("cid")+".png")
	dst := "/www/wwwroot/myproject/mu_html/picture/chara/" + c.PostForm("cid") + ".png"
	println(dst)
	c.SaveUploadedFile(file, dst)
}
func GetGamePictrue(c *gin.Context) {
    file, _ := c.FormFile("imgfile")
    dst :="/www/wwwroot/myproject/mu_html/picture/game/"+file.Filename
    c.SaveUploadedFile(file, dst)
}
//-------------------------------------------//
//验证权限
func IdConfirm() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("IdConfirm Works")
		username, err1 := c.Cookie("youAss")
		temp := Admin{}
		println("AdminAccontIs" + username)
		DB.Where("usrname = ?", username).Find(&temp)
		println("AdminAccontIs" + temp.Usrname + "  " + temp.Password)
		if temp.Usrname == username && err1 == nil && temp.Usrname != "" {
			fmt.Println("Login Passed")
			println(err1)
			c.Next()
		} else {
			fmt.Println("Login Failed")
			c.Redirect(http.StatusMovedPermanently, "/admin/login")
		}
	}
}

//-------------------------------------------//
func CharaAllSearch(text string) []Character {
	result := []Character{}
	text = "%" + text + "%"
	DB.Where("name like ? or tag like ?", text, text).Find(&result)
	return result
}
func GameAllSearch(text string) []Game {
	result := []Game{}
	text = "%" + text + "%"
	DB.Where("title like ? or tag like ? or brand like ?", text, text, text).Find(&result)
	return result
}

//-------------------------------------------//
type Brand struct {
	Name string
}

//自动加入社团信息
func AutoAddBrand(str string) {
	result := []Brand{}
	println("加入的TAG是：" + str)
	DB.Where("name = ?", str).Find(&result)
	if len(result) == 0 {
		objct := Brand{Name: str}
		DB.Create(&objct)
	}
}
func GetBrandList() []Brand {
	result := []Brand{}
	DB.Where("name > ?", "").Find(&result)
	return result
}
func BrandSlice(start int, end int, lst []Brand) map[int]Brand {
	result := map[int]Brand{}
	sum := 0
	for k, v := range lst {
		if k >= start && k < end {
			result[sum] = v
			sum++
		}
	}
	return result
}

//-------------------------------------------//
type Alltag struct {
	Name string
}

//自动加入tag
func AutoAddTag(str string) {
	tags := ""
	for _, i := range str {
		char := string(i)
		if char == "," {
			lst := []Alltag{}
			DB.Where("name = ?", tags).Find(&lst)
			if len(lst) == 0 && tags != "" {
				objct := Alltag{Name: tags}
				DB.Create(&objct)
			}
			tags = ""
			println("Tags Add Okay")
		} else {
			tags += char
		}
	}
}

//列举tag
func GetTagList() []Alltag {
	result := []Alltag{}
	DB.Where("name > ?", "").Find(&result)
	return result
}
func TagSlice(start int, end int, lst []Alltag) map[int]Alltag {
	result := map[int]Alltag{}
	sum := 0
	for k, v := range lst {
		if k >= start && k < end {
			result[sum] = v
			sum++
		}
	}
	return result
}

//-------------------------------------------//
//character操作
type Character struct {
	Cid     int
	Name    string
	Birth   int
	Tag     string
	Img     string
	Id      int
	Content string
}

//CharaTag搜索
func QueryCharaWithTag(query string) []Character {
	result := []Character{}
	query = "," + "%" + query + "%" + ","
	DB.Where("tag like ?", query).Find(&result)
	return result
}

//Post Data into Chara Object
func PackageChara(c *gin.Context) Character {
	cid := StrToInt(c.PostForm("cid"))
	name := c.PostForm("rwname")
	birth := StrToInt(c.PostForm("rwbirth"))
	tag := c.PostForm("rwtag")
	img := "/img/chara/" + c.PostForm("cid") + ".png"
	id := StrToInt(c.PostForm("pid"))
	content := c.PostForm("rwcontent")
	temp := Character{
		Cid:     cid,
		Name:    name,
		Birth:   birth,
		Tag:     tag,
		Img:     img,
		Id:      id,
		Content: content,
	}
	AutoAddTag(tag)
	println("New Characters Okay")
	return temp
}

//根据cid，返回角色的Object
func ShowChara(cid int) Character {
	result := Character{}
	DB.Where("cid = ?", cid).Find(&result)
	return result
}

//创建新角色
func AddChara(chara Character) {
	result := DB.Create(&chara)
	fmt.Println(result.RowsAffected)
}

//用名称搜索角色
func SearchCharaWithName(text string) []Character {
	result := []Character{}
	text = "%" + text + "%"
	DB.Where("name like ?", text).Find(&result)
	return result
}
func SearchCharaWithTag(text string) []Character {
	result := []Character{}
	text = "%" + text + "%"
	DB.Where("tag like ?", text).Find(&result)
	return result
}
func SearchCharaWithTime(time string) []Character {
	result := []Character{}
	time = time + "%"
	DB.Where("birth like ?", time).Find(&result)
	return result
}

//用给出的角色object覆写原本的角色
func RewriteChara(chara Character) {
	DB.Model(&Character{}).Where("cid = ?", chara.Cid).Save(&chara)
}

//根据Id删除对应的角色
func DeleteChara(cid int) {
	DB.Delete(&Character{}, cid)
}

//得到全部角色object组成的列表
func GetCharaList() []Character {
	println("Function 'GetCharaList' is Woriking")
	result := []Character{}
	DB.Where("name <> ?", "").Find(&result)
	println("CharaList's length:" + IntToString(len(result)))
	return result
}

//切片并返回角色object
func CharaSlice(start int, finish int, list []Character) map[int]Character {
	println("Function 'CharaSlice' is Woriking")
	temp := map[int]Character{}
	sum := 0
	for k, v := range list {
		if k >= start && k < finish {
			temp[sum] = v
			sum = sum + 1
		}
	}
	return temp
}

//-------------------------------------------//
//Game相关数据库操作
type Game struct {
	Id      int
	Title   string
	Content string
	Brand   string
	Tag     string
	Date    int
}

//Gametag搜索
func QueryGameWithTag(query string) []Game {
	result := []Game{}
	query = "," + query + ","
	DB.Where("tag like ?", query).Find(&result)
	return result
}

//给出PosT数据封装游戏
func PackageGame(c *gin.Context) Game {
	brand := c.PostForm("rwbrand")
	title := c.PostForm("rwtitle")
	content := c.PostForm("rwcontent")
	tag := c.PostForm("rwtag")
	id := StrToInt(c.PostForm("pid"))
	date := StrToInt(c.PostForm("rwdate"))
	newGame := Game{
		Date:    date,
		Brand:   brand,
		Content: content,
		Tag:     tag,
		Id:      id,
		Title:   title,
	}
	AutoAddTag(tag)
	AutoAddBrand(brand)
	return newGame
}

//根据id，返回游戏的Object
func ShowGame(id int) Game {
	result := Game{}
	DB.Where("id = ?", id).Find(&result)
	return result
}

//输入一个游戏object，并存入数据库
func AddGame(game Game) {
	result := DB.Create(&game)
	fmt.Println(result.RowsAffected)
}

//根据输入的文本，搜索包含文本的标题
func SearchGame(text string) []Game {
	result := []Game{}
	text = "%" + text + "%"
	DB.Where("title like ?", text).Find(&result)
	return result
}
func SearchGameWithTag(text string) []Game {
	result := []Game{}
	text = "%" + text + "%"
	DB.Where("tag like ?", text).Find(&result)
	return result
}
func SearchGameWithBrand(text string) []Game {
	result := []Game{}
	text = "%" + text + "%"
	DB.Where("brand like ?", text).Find(&result)
	return result
}
func SearchGameWithTime(time string) []Game {
	result := []Game{}
	time = time + "%"
	DB.Where("date like ?", time).Find(&result)
	return result
}

//用给出的游戏object覆写原本的游戏
func RewriteGame(game Game) {
	DB.Model(&Game{}).Where("id = ?", game.Id).Save(&game)
}

//根据Id删除对应的游戏
func DeleteGame(id int) {
	DB.Delete(&Game{}, id)
}

//得到全部游戏object组成的列表
func GetGameList() []Game {
	result := []Game{}
	DB.Where("title >= ?", "").Find(&result)

	return result
}
func GameSlice(start int, finish int, list []Game) map[int]Game {
	println("Function 'CharaSlice' is Woriking")
	temp := map[int]Game{}
	sum := 0
	for k, v := range list {
		if k >= start && k < finish {
			temp[sum] = v
			sum = sum + 1
		}
	}
	return temp
}

//-------------------------------------------//
//usr相关数据库操作
type Admin struct {
	Usrname  string
	Password string
}

//用得到的用户object与数据库中对比，返回能否登录
func Login(user Admin) bool {
	temp := Admin{}
	println("------------------------------------")
	println("Name=" + user.Usrname)
	println("PassWord=" + user.Password)
	println("-------------------------------------")
	DB.Where("usrname = ?", user.Usrname).Find(&temp)
	return (temp.Password == user.Password)
}

//-------------------------------------------//
func main() {
	gin.SetMode(gin.ReleaseMode)
	//初始化Database
	Init()
	//实例化一个引擎
	r := gin.Default()

	//路径的映射，当你访问“/css”时实际在访问“mu_html/css”
	r.Static("/css", "mu_html/css")
	r.Static("/img", "mu_html/picture")
	r.Static("/wb", "mu_html")
	r.Static("/js", "mu_html/javascript")

	//HTML输出加，其他时候可以不用加.预加载html，此为放置Html文件的地址
	r.LoadHTMLGlob("mu_html/**/*.*")

	//配置路由

	//默认路径
	r.GET("/", func(c *gin.Context) {
		c.Request.URL.Path = "/home/"
		r.HandleContext(c)
	})
	//首页
	indexRouters := r.Group("/home")
	{
		indexRouters.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/home/index")
		})
		indexRouters.GET("/index", func(c *gin.Context) {
			c.HTML(200, "index_index.html", gin.H{})
		})
		indexRouters.GET("/about", func(c *gin.Context) {
			c.HTML(200, "index_AboutUs.html", gin.H{})
		})
	}
	//管理页面入口，无需cookie
	adminEnterrouters := r.Group("/admin")
	{
		//登录页面
		adminEnterrouters.GET("/login", func(c *gin.Context) {
			c.HTML(200, "admin_login.html", gin.H{})
		})

		//认证账号密码，生成cookie
		adminEnterrouters.POST("/login_done", func(c *gin.Context) {
			//从POST中获得输入的账号密码
			username := c.PostForm("username")
			password := c.PostForm("password")
			fmt.Println("Post Respond working")
			//根据数据生成用户Object，与数据库对比
			user := Admin{Usrname: username, Password: password}
			if Login(user) {
				//建立cookie，跳到管理首页
				fmt.Println("Cookie Set Works")
				c.SetCookie("youAss", username, 1000000, "/", "www.fullyarmedbunny19c.xyz", false, true)
				c.SetCookie("youAss", username, 1000000, "/", "fullyarmedbunny19c.xyz", false, true)
				c.SetCookie("youAss", username, 1000000, "/", "159.138.155.59", false, true)
				c.Redirect(http.StatusMovedPermanently, "/admin/index")
			} else {
				//认证失败，跳回登陆页面
				fmt.Println("Cookie Set Failed")
				c.Redirect(http.StatusMovedPermanently, "/admin/login")
			}
		})
	}

	//管理界面，需要cookie
	adminrouters := r.Group("/admin").Use(IdConfirm())
	{
		adminrouters.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/admin/index")
		})
		adminrouters.GET("/index", func(c *gin.Context) {
			c.HTML(200, "admin_index.html", gin.H{
				"total": len(GetGameList()) + len(GetCharaList()),
			})
		})
		/*adminrouters.GET("/rewrite", func(c *gin.Context) {
			c.HTML(200, "admin_rewrite.html", gin.H{
				"total": len(GetGameList()),
			})
		})*/
		//重写游戏文章

		//可更改的游戏的列表
		adminrouters.GET("/rewrite/games/:id", func(c *gin.Context) {
			id := StrToInt(c.Param("id"))
			list := GameSlice(id*10, (id+1)*10, GetGameList())
			c.HTML(200, "admin_rewritegame.html", gin.H{
				"list": list,
				"npage": id+1,
				"lpage":id-1,
			})
		})
		//动态路由，游戏覆写页面
		adminrouters.GET("/rewrite/game/:pid", func(c *gin.Context) {
			c.Header("Content-Type", "text/html; charset=utf-8")
			id := StrToInt(c.Param("pid"))
			temp := ShowGame(id)
			c.HTML(200, "admin_rewritegame_model.html", gin.H{
				"id":      temp.Id,
				"title":   temp.Title,
				"date":    temp.Date,
				"brand":   temp.Brand,
				"tag":     temp.Tag,
				"content": temp.Content,
			})
		})
		//覆写信息接受页面
		adminrouters.POST("/rw/game", func(c *gin.Context) {
			println("rw Page Working")
			newGame := PackageGame(c)
			println("Load Pictrue")
			RewriteGame(newGame)
			c.Redirect(http.StatusMovedPermanently, "/admin/rewrite/games/0")
		})
		//移除信息接受页面
		adminrouters.POST("/rm/game", func(c *gin.Context) {
			id := StrToInt(c.PostForm("pid"))
			DeleteGame(id)
			c.Redirect(http.StatusMovedPermanently, "/admin/rewrite/games/0")
		})
		adminrouters.GET("/add/game", func(c *gin.Context) {
			c.HTML(200, "admin_addgame.html", gin.H{
				"nextId": GetNextId(),
			})
		})
		adminrouters.POST("/addgame", func(c *gin.Context) {
			println("add Page Working")
			newGame := PackageGame(c)
			AddGame(newGame)
			c.Redirect(http.StatusMovedPermanently, "/admin/rewrite/games/0")
		})

		//重写角色文章界面

		//可更改的角色的列表
		adminrouters.GET("/rewrite/charas/:pg", func(c *gin.Context) {
			id := StrToInt(c.Param("pg"))
			list := CharaSlice(id*10, (id+1)*10, GetCharaList())
			c.HTML(200, "admin_rewritechara.html", gin.H{
				"list": list,
				"npage": id+1,
				"lpage":id-1,
			})
		})
		//动态路由，角色覆写页面
		adminrouters.GET("/rewrite/chara/:cid", func(c *gin.Context) {
			c.Header("Content-Type", "text/html; charset=utf-8")
			cid := StrToInt(c.Param("cid"))
			println("You are connecting to " + c.Param("cid"))
			temp := ShowChara(cid)
			c.HTML(200, "admin_rewritechara_model.html", gin.H{
				"id":      temp.Id,
				"tag":     temp.Tag,
				"cid":     temp.Cid,
				"name":    temp.Name,
				"content": temp.Content,
				"birth":   temp.Birth,
				"img":     temp.Img,
			})
		})
		//覆写信息接受页面
		adminrouters.POST("/rw/chara", func(c *gin.Context) {
			println("rw Page Working")
			newChara := PackageChara(c)
			RewriteChara(newChara)
			GetCharaPictrue(c)
			c.Redirect(http.StatusMovedPermanently, "/admin/rewrite/charas/0")
		})
		//移除信息接受页面
		adminrouters.POST("/rm/chara", func(c *gin.Context) {
			id := StrToInt(c.PostForm("cid"))
			DeleteChara(id)
			c.Redirect(http.StatusMovedPermanently, "/admin/rewrite/charas/0")
		})
		adminrouters.GET("/add/chara", func(c *gin.Context) {
			c.HTML(200, "admin_addchara.html", gin.H{
				"nextCid": GetNextCid(),
			})
		})
		adminrouters.POST("/addchara", func(c *gin.Context) {
			println("add Page Working")
			newChara := PackageChara(c)
			GetCharaPictrue(c)
			AddChara(newChara)
			c.Redirect(http.StatusMovedPermanently, "/admin/rewrite/charas/0")
		})

		//用户信息界面
		adminrouters.GET("/user", func(c *gin.Context) {
			usrname, _ := c.Cookie("youAss")
			c.HTML(200, "admin_user.html", gin.H{
				"usrname": usrname,
				"roots":   "Yes",
			})
		})
		adminrouters.GET("/readme", func(c *gin.Context) {
			c.HTML(200, "admin_readme.html", gin.H{})
		})
	}
	SearchRouster := r.Group("/search")
	{
		SearchRouster.GET("/:mode/:str/:pg", func(c *gin.Context) {
			println("SearchEngine Working")
			lst1 := map[int]Game{}
			lst2 := map[int]Character{}
			mode := c.Param("mode")
			txt := c.Param("str")
			println(mode + "," + txt)
			pg := StrToInt(c.Param("pg"))
			if mode == "time" {
				lst1 = GameSlice(pg*5, pg*5+5, SearchGameWithTime(txt))
				lst2 = CharaSlice(pg*5, pg*5+5, SearchCharaWithTime(txt))
			}
			if mode == "brand" {
				lst1 = GameSlice(pg*5, pg*5+5, SearchGameWithBrand(txt))
			}
			if mode == "tag" {
				lst1 = GameSlice(pg*5, pg*5+5, SearchGameWithTag(txt))
				lst2 = CharaSlice(pg*5, pg*5+5, SearchCharaWithTag(txt))
			}
			if mode == "all" {
				lst1 = GameSlice(pg*5, pg*5+5, GameAllSearch(txt))
				lst2 = CharaSlice(pg*5, pg*5+5, CharaAllSearch(txt))
			}
			c.HTML(200, "class_model.html", gin.H{
				"list1": lst1,
				"list2": lst2,
				"npage": pg+1,
				"lpage":pg-1,
				"mod":mode,
				"str":txt,
			})
		})
	}
	ClassRouster := r.Group("/class")
	{
		ClassRouster.GET("/index", func(c *gin.Context) {
			c.HTML(200, "class_index.html", gin.H{})
		})
		ClassRouster.GET("/time", func(c *gin.Context) {
			c.HTML(200, "class_time.html", gin.H{})
		})
		ClassRouster.GET("/tag", func(c *gin.Context) {
			lst := TagSlice(0, 1000, GetTagList())
			c.HTML(200, "class_tag.html", gin.H{
				"list": lst,
			})
		})
		ClassRouster.GET("/brand", func(c *gin.Context) {
			lst := BrandSlice(0, 1000, GetBrandList())
			c.HTML(200, "class_brand.html", gin.H{
				"list": lst,
			})
		})
		ClassRouster.GET("/ours", func(c *gin.Context) {
			c.HTML(200, "class_ours.html", gin.H{})
		})
	}
	//浏览游戏
	GameRouster := r.Group("/game")
	{
		GameRouster.GET("/:pid", func(c *gin.Context) {
			c.Header("Content-Type", "text/html; charset=utf-8")
			id := StrToInt(c.Param("pid"))
			temp := ShowGame(id)
			fmt.Println("Game Page Working")
			c.HTML(200, "Game_model.html", gin.H{
				"title":   template.HTML(temp.Title),
				"content": template.HTML(temp.Content),
				"brand":   temp.Brand,
				"Date":    temp.Date,
			})
		})
	}

	CharaRouster := r.Group("/chara")
	{
		CharaRouster.GET("/:cid", func(c *gin.Context) {
			c.Header("Content-Type", "text/html; charset=utf-8")
			cid := StrToInt(c.Param("cid"))
			temp := ShowChara(cid)
			fmt.Println("Chara Page Working")
			c.HTML(200, "Chara_model.html", gin.H{
				"name":    template.HTML(temp.Name),
				"content": template.HTML(temp.Content),
				"birth":   template.HTML(temp.Birth),
				"game":    template.HTML(ShowGame(temp.Id).Title),
				"imgpath": template.HTML(temp.Img),
			})
		})
	}
	r.POST("/pictures",func(c *gin.Context) {
	    GetGamePictrue(c)
	    c.Redirect(http.StatusMovedPermanently, "/admin/readme")
	})
	r.POST("/advice",func(c *gin.Context) {
	    title:=url.QueryEscape("意见反馈")
	    content:=url.QueryEscape(c.PostForm("content"))
	    println("This is Content:"+content)
	    http.Get("https://api.wer.plus/api/qqmail?name=Crusherbunny&me=2366817971@qq.com&to=2366817971@qq.com&title="+title+"&text="+content+"&key=ppigbpyhrqxfdidd")
	    c.Redirect(http.StatusMovedPermanently,"/home/about")
	})
	r.Run() // 监听并在 localhost:8080 上启动服务
}
