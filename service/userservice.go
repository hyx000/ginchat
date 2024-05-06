package service

import (
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// GetUserList
// @summary 所有用户
// @Tags 用户模块
// @Success 200 {string} json{"code","message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()
	c.JSON(200, gin.H{
		"message": data,
	})
}

// FindUserByNameAndPwd
// @summary 登录
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/findUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) {
	data := models.UserBasic{}
	//name := c.Query("name")
	//pssword := c.Query("password")
	name := c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	fmt.Println(name, password)
	fmt.Println(name)
	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(200, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "该用户不存在",
			"data":    data,
		})
		return
	}
	flag := utils.ValidPassword(password, user.Salt, user.PassWord)
	if !flag {
		c.JSON(200, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "密码不正确",
			"data":    data,
		})
		return
	}
	pwd := utils.MakePassword(password, user.Salt)
	data = models.FindUserByNameAndPwd(name, pwd)
	c.JSON(200, gin.H{
		"code":    0, //0成功 -1 失败
		"message": "登录成功",
		"data":    data,
	})
}

// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [get]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}

	user.Name = c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	repassword := c.Request.FormValue("repassword")

	data := models.FindUserByName(user.Name)
	salt := fmt.Sprintf("%06d", rand.Int31())

	if user.Name == "" || password == "" {
		c.JSON(-1, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "用户名或者密码不能为空！",
			"data":    user,
		})
		return
	}
	if data.Name != "" {
		c.JSON(-1, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "用户名已被注册！",
			"data":    user,
		})
		return
	}
	if password != repassword {
		c.JSON(-1, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "两次密码不一致！",
			"data":    user,
		})
		return
	}
	user.PassWord = utils.MakePassword(password, salt)
	user.Salt = salt
	models.CreateUser(user)
	c.JSON(200, gin.H{
		"code":    0, //  0成功   -1失败
		"message": "新增用户成功！",
		"data":    user,
	})
}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id query string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /user/deleteUser [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(200, gin.H{
		"message": "删除成功",
	})
}

// UpdateUser
// @Summary 修改用户
// @Tags 用户模块
// @param id formData string false "id"
// @param name formData string false "name"
// @param password formData string false "password"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @Success 200 {string} json{"code","message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(-1, gin.H{
			"code":    -1,
			"message": "修改的参数不匹配!",
			"data":    user,
		})
		return
	}
	models.UpdateUser(user)
	c.JSON(200, gin.H{
		"code":    0,
		"message": "修改成功",
		"data":    user,
	})
}

// 防止跨域
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 发送
func SendMsg(c *gin.Context) {
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	MsgHandler(c, ws)
}

func MsgHandler(c *gin.Context, ws *websocket.Conn) {
	for {
		msg, err := utils.Subscribe(c, utils.PublishKey)
		if err != nil {
			fmt.Println(" MsgHandler 发送失败", err)
		}
		tm := time.Now().Format("2006-01-02 15:04:05")
		m := fmt.Sprintf("[ws][%s]:%s", tm, msg)

		err = ws.WriteMessage(1, []byte(m))
		if err != nil {
			log.Fatalln(err)
		}
	}

}

func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}

func SearchFriends(c *gin.Context) {
	I, _ := strconv.Atoi(c.Request.FormValue("userId"))
	id := uint(I)
	users := models.SearchFriends(id)
	utils.RespOkList(c.Writer, users, len(users))

}
func RedisMsg(c *gin.Context) {
	userIdA, _ := strconv.Atoi(c.PostForm("userIdA"))
	userIdB, _ := strconv.Atoi(c.PostForm("userIdB"))
	start, _ := strconv.Atoi(c.PostForm("start"))
	end, _ := strconv.Atoi(c.PostForm("end"))
	isRev, _ := strconv.ParseBool(c.PostForm("isRev"))
	res := models.RedisMsg(int64(userIdA), int64(userIdB), int64(start), int64(end), isRev)
	utils.RespOkList(c.Writer, "ok", res)
}
