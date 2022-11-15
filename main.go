package main

import (
	"fmt"
	//"path"
	"strconv"

	"html/template"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

//初始化数据库
func Init() {
	dsn := "root:Kagarino_Kirie@tcp(127.0.0.1:3306)/renga?charset=utf8mb3&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("DB ERROR~!!")
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

//接受图片
func GetCharaPictrue(c *gin.Context) {
	file, err := c.FormFile("imgfile")
	if file.Filename == "0.txt" {
		return
	}
	if err != nil {
		println("error")
	}
	println(file.Filename)
	//dst := path.Join("/", c.PostForm("cid")+".png")
	dst := "F:/Go_Works/src/myproject/mu_html/picture/chara/" + c.PostForm("cid") + ".png"
	println(dst)
	c.SaveUploadedFile(file, dst)
}
func GetGamePictrue(c *gin.Context) {
}

//验证权限
func IdConfirm() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("IdConfirm Works")
		username, err1 := c.Cookie("username")
		temp := Admin{}
		DB.Where("usrname = ?", username).Find(&temp)
		if temp.Usrname == username && err1 == nil {
			fmt.Println("Login Passed")
			c.Next()
		} else {
			fmt.Println("Login Failed")
			c.Redirect(http.StatusMovedPermanently, "/admin/login")
		}
	}
}

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

//Post Data To Chara Object
func PackageChara(c *gin.Context) Character {
	cid := StrToInt(c.PostForm("cid"))
	name := c.PostForm("rwname")
	birth := StrToInt(c.PostForm("rwbirth"))
	tag := c.PostForm("rwtag")
	img := "F:/Go_Works/src/myproject/mu_html/picture/chara/" + c.PostForm("cid") + ".png"
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
	println("New Object Okay")
	return temp
}

//根据cid，返回角色的Object
func ShowChara(cid int) Character {
	result := Character{}
	DB.Where("id = ?", cid).Find(&result)
	return result
}

//创建新角色
func AddChara(chara Character) {
	result := DB.Create(&chara)
	fmt.Println(result.RowsAffected)
}

//用名称搜索角色
func SearchChara(text string) []Character {
	result := []Character{}
	DB.Where("name <= ?", text).Find(&result)
	return result
}

//用给出的角色object覆写原本的角色
func RewriteChara(chara Character) {
	temp := Character{Id: chara.Id}
	DB.Find(&temp)
	temp = chara
	DB.Save(&temp)
}

//根据Id删除对应的角色
func DeleteChara(cid int) {
	DB.Delete(&Character{}, cid)
}

//得到全部角色object组成的列表
func GetCharaList() []Character {
	result := []Character{}
	DB.Where("name >= ?", "").Find(&result)
	return result
}

//将角色内容转化为map输出得到列表
func CharaToList(list []Character) map[int]Character {
	temp := map[int]Character{}
	for i, art := range list {
		temp[i] = art
	}
	return temp
}

//Game相关数据库操作
type Game struct {
	Id      int
	Title   string
	Content string
	Brand   string
	Tag     string
	Date    int
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
	DB.Where("title <= ?", text).Find(&result)
	return result
}

//用给出的游戏object覆写原本的游戏
func RewriteGame(game Game) {
	temp := Game{Id: game.Id}
	DB.Find(&temp)
	temp = game
	DB.Save(&temp)
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

//将游戏内容转化为map输出得到列表
func GameToList(list []Game) map[int]Game {
	temp := map[int]Game{}
	for i, art := range list {
		temp[i] = art
	}
	return temp
}

//usr相关数据库操作
type Admin struct {
	Usrname  string
	Password string
}

//用得到的用户object与数据库中对比，返回能否登录
func Login(user Admin) bool {
	temp := Admin{}
	DB.Where("usrname = ?", user.Usrname).Find(&temp)
	return (temp.Password == user.Password)
}
func main() {
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
				c.SetCookie("username", username, 100000000, "/", "127.0.0.1/", false, false)
				c.Redirect(http.StatusMovedPermanently, "/admin/index")
			} else {
				//认证失败，跳回登陆页面
				fmt.Println("Cookie Set Failed")
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
		adminrouters.GET("/rewrite", func(c *gin.Context) {
			c.HTML(200, "admin_rewrite.html", gin.H{
				"total": len(GetGameList()),
			})
		})
		//重写游戏文章

		//可更改的游戏的列表
		adminrouters.GET("/rewrite/game", func(c *gin.Context) {
			list := GameToList(GetGameList())
			c.HTML(200, "admin_rewritegame.html", list)
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
			c.Redirect(http.StatusMovedPermanently, "/admin/rewrite/game")
		})
		//移除信息接受页面
		adminrouters.POST("/rm/game", func(c *gin.Context) {
			id := StrToInt(c.PostForm("pid"))
			DeleteGame(id)
			c.Redirect(http.StatusMovedPermanently, "/admin/rewrite/game")
		})
		adminrouters.GET("/add/game", func(c *gin.Context) {
			c.HTML(200, "admin_addgame.html", gin.H{})
		})
		adminrouters.POST("/addgame", func(c *gin.Context) {
			println("add Page Working")
			newGame := PackageGame(c)
			AddGame(newGame)
			c.Redirect(http.StatusMovedPermanently, "/admin/rewrite/game")
		})

		//重写角色文章界面

		//可更改的角色的列表
		adminrouters.GET("/rewrite/chara", func(c *gin.Context) {
			list := CharaToList(GetCharaList())
			c.HTML(200, "admin_rewritechara.html", list)
		})
		//动态路由，角色覆写页面
		adminrouters.GET("/rewrite/chara/:cid", func(c *gin.Context) {
			c.Header("Content-Type", "text/html; charset=utf-8")
			id := StrToInt(c.Param("cid"))
			temp := ShowChara(id)
			c.HTML(200, "admin_rewritechara_model.html", gin.H{
				"id":      temp.Id,
				"tag":     temp.Tag,
				"cid":     temp.Cid,
				"name":    temp.Name,
				"content": temp.Content,
				"birth":   temp.Birth,
				"img":     template.HTML(temp.Img),
			})
		})
		//覆写信息接受页面
		adminrouters.POST("/rw/chara", func(c *gin.Context) {
			println("rw Page Working")
			newChara := PackageChara(c)
			RewriteChara(newChara)
			GetCharaPictrue(c)
			c.Redirect(http.StatusMovedPermanently, "/admin/rewrite/chara")
		})
		//移除信息接受页面
		adminrouters.POST("/rm/chara", func(c *gin.Context) {
			id := StrToInt(c.PostForm("cid"))
			DeleteChara(id)
			c.Redirect(http.StatusMovedPermanently, "/admin/rewrite/chara")
		})
		adminrouters.GET("/add/chara", func(c *gin.Context) {
			c.HTML(200, "admin_addchara.html", gin.H{})
		})
		adminrouters.POST("/addchara", func(c *gin.Context) {
			println("add Page Working")
			newChara := PackageChara(c)
			GetCharaPictrue(c)
			AddChara(newChara)
			c.Redirect(http.StatusMovedPermanently, "/admin/rewrite/chara")
		})

		//用户信息界面
		adminrouters.GET("/user", func(c *gin.Context) {
			usrname, _ := c.Cookie("username")
			c.HTML(200, "admin_user.html", gin.H{
				"usrname": usrname,
				"roots":   "Yes",
			})
		})
		adminrouters.GET("/readme", func(c *gin.Context) {
			c.HTML(200, "admin_readme.html", gin.H{})
		})
	}
	//浏览游戏
	GameRouster := r.Group("/Game")
	{
		GameRouster.GET("/:pid", func(c *gin.Context) {
			c.Header("Content-Type", "text/html; charset=utf-8")
			id := StrToInt(c.Param("pid"))
			temp := ShowGame(id)
			fmt.Println("Game Page Working")
			c.HTML(200, "Game_model.html", gin.H{
				"title":   template.HTML(temp.Title),
				"context": template.HTML(temp.Content),
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
	r.Run() // 监听并在 localhost:8080 上启动服务
}
