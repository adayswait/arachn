package handler

import "time"
import "crypto/md5"
import "encoding/hex"
import "encoding/json"
import "github.com/gofiber/fiber"
import "github.com/gofiber/session"
import "github.com/adayswait/mojo/db"
import "github.com/adayswait/mojo/global"
import "github.com/adayswait/mojo/utils"

import (
	"fmt"
	"github.com/google/goexpect"
	"github.com/google/goterm/term"
	"golang.org/x/crypto/ssh"
	// "google.golang.org/grpc/codes"
	"regexp"
)

// import "github.com/adayswait/mojo/cmd"

var sessions *session.Session

func init() {
	sessions = session.New()
}

type SecInfo struct {
	User   string `json:"user"`
	Passwd string `json:"passwd"`
}

func Login(c *fiber.Ctx) {
	var secInfo SecInfo
	var accessToken string
	if err := c.QueryParser(&secInfo); err == nil {
		if len(secInfo.User) != 0 && len(secInfo.Passwd) != 0 {
			//优先使用用户名和密码进行验证
			tempMD5 := md5.New()
			tempMD5.Write([]byte(secInfo.Passwd))
			tempMD5.Write([]byte(global.MD5_SALT))
			passwdMD5 := hex.EncodeToString(tempMD5.Sum(nil))
			eau := db.Auth(secInfo.User, passwdMD5)
			if eau != nil {
				c.JSON(fiber.Map{"code": global.RET_ERR_USER_PASSWD,
					"data": eau.Error()})
				return
			}
			token, etk := db.GetAccessToken(secInfo.User)
			if etk != nil {
				c.JSON(fiber.Map{"code": global.RET_ERR_DB,
					"data": etk.Error()})
			}
			accessToken = token
		}
		//取不到user和passwd, 后面尝试access_token登录
	}
	if len(accessToken) == 0 {
		//用户名密码验证未成功, 再尝试使用access_token进行验证
		accessToken = c.Cookies(global.ACCESS_TOKEN)
	} else {
		cookie := new(fiber.Cookie)
		cookie.Name = global.ACCESS_TOKEN
		cookie.Expires = time.Now().Add(24 * time.Hour)
		cookie.Value = accessToken
		c.Cookie(cookie)
	}

	if len(accessToken) == 0 {
		c.JSON(fiber.Map{"code": global.RET_ERR_ACCESS_TOKEN,
			"data": "no access_token found"})
		return
	}
	userInfo, elg := db.Login(accessToken)
	if elg != nil {
		c.JSON(fiber.Map{"code": global.RET_ERR_ACCESS_TOKEN,
			"data": elg.Error()})
		return
	}
	data, _ := json.Marshal(userInfo)
	c.JSON(fiber.Map{"code": global.RET_OK, "data": string(data)})

	//create session
	(func() {
		store := sessions.Get(c)
		store.Set(global.SESSION_KEY_USER, userInfo.User)
		store.Set(global.SESSION_KEY_GROUP, userInfo.Group)
		defer store.Save()
	})()
}

func Logout(c *fiber.Ctx) {
	store := sessions.Get(c)
	store.Destroy()

	//todo
	cookie := new(fiber.Cookie)
	cookie.Name = global.ACCESS_TOKEN
	cookie.Value = ""
	c.Cookie(cookie)

	defer store.Save()
	c.JSON(fiber.Map{"code": global.RET_OK, "data": nil})
}

func Register(c *fiber.Ctx) {
	var secInfo SecInfo
	err := c.QueryParser(&secInfo)
	if err != nil {
		c.JSON(fiber.Map{"code": global.RET_ERR_HTTP_QUERY,
			"data": err.Error()})
		return
	}
	if len(secInfo.User) == 0 || len(secInfo.Passwd) == 0 {
		c.JSON(fiber.Map{"code": global.RET_ERR_INPUT,
			"data": "user or password is empty"})
		return
	}
	tempMD5 := md5.New()
	tempMD5.Write([]byte(secInfo.Passwd))
	tempMD5.Write([]byte(global.MD5_SALT))
	passwdMD5 := hex.EncodeToString(tempMD5.Sum(nil))

	edb := db.Register(secInfo.User, passwdMD5)
	if edb != nil {
		c.JSON(fiber.Map{"code": global.RET_ERR_DB, "data": edb.Error()})
		return
	}
	accessToken, etk := db.GetAccessToken(secInfo.User)
	if etk != nil {
		c.JSON(fiber.Map{"code": global.RET_ERR_DB, "data": etk.Error()})
		return
	}
	cookie := new(fiber.Cookie)
	cookie.Name = global.ACCESS_TOKEN
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Value = accessToken
	c.Cookie(cookie)
	c.JSON(fiber.Map{"code": global.RET_OK, "data": nil})
}

func ViewDB(c *fiber.Ctx) {
	store := sessions.Get(c)
	// user := store.Get(global.SESSION_KEY_USER)
	group := store.Get(global.SESSION_KEY_GROUP)
	if group == nil {
		c.JSON(fiber.Map{"code": global.RET_ERR_SESSION_INVALID,
			"data": "session invalid"})
		return
	}
	if int(group.(int64)) > int(global.GROUP_UNDEF) {
		c.JSON(fiber.Map{"code": global.RET_ERR_NO_RIGHT,
			"data": "no right to do this"})
		return
	}
	table := c.Params("table")
	if len(table) == 0 {
		ret, err := db.Buckets()
		if err != nil {
			c.JSON(fiber.Map{"code": global.RET_ERR_DB, "data": err.Error()})
			return
		}
		c.JSON(fiber.Map{"code": global.RET_OK, "data": ret})
		return
	}
	key := c.Params("key")
	if len(key) == 0 {
		ret, err := db.Keys(table)
		if err != nil {
			c.JSON(fiber.Map{"code": global.RET_ERR_DB, "data": err.Error()})
			return
		}
		c.JSON(fiber.Map{"code": global.RET_OK, "data": ret})
		return
	} else {
		ret, err := db.Get(table, key)
		if err != nil {
			c.JSON(fiber.Map{"code": global.RET_ERR_DB, "data": err.Error()})
			return
		} else {
			c.JSON(fiber.Map{"code": global.RET_OK, "data": string(ret)})
			return
		}
	}
}

func UpdateDB(c *fiber.Ctx) {
	store := sessions.Get(c)
	// user := store.Get(global.SESSION_KEY_USER)
	group := store.Get(global.SESSION_KEY_GROUP)
	if group == nil {
		c.JSON(fiber.Map{"code": global.RET_ERR_SESSION_INVALID,
			"data": "session invalid"})
		return
	}
	if int(group.(int64)) > int(global.GROUP_UNDEF) {
		c.JSON(fiber.Map{"code": global.RET_ERR_NO_RIGHT,
			"data": "no right to do this"})
		return
	}
	table := c.Params("table")
	key := c.Params("key")
	if len(table) == 0 {
		c.JSON(fiber.Map{"code": global.RET_ERR_URL_PARAM,
			"data": "request param err"})
		return
	}

	body := struct {
		Value string `json:"value"`
	}{}

	if errBp := c.BodyParser(&body); errBp != nil {
		c.JSON(fiber.Map{"code": global.RET_ERR_BODY_PARAM,
			"data": errBp.Error()})
		return
	}
	err := db.Set(table, key, body.Value)
	if err != nil {
		c.JSON(fiber.Map{"code": global.RET_ERR_DB,
			"data": err.Error()})
		return
	}
	c.JSON(fiber.Map{"code": global.RET_OK, "data": nil})
}

func DeleteDB(c *fiber.Ctx) {
	store := sessions.Get(c)
	// user := store.Get(global.SESSION_KEY_USER)
	group := store.Get(global.SESSION_KEY_GROUP)
	if group == nil {
		c.JSON(fiber.Map{"code": global.RET_ERR_SESSION_INVALID,
			"data": "session invalid"})
		return
	}
	if int(group.(int64)) > int(global.GROUP_UNDEF) {
		c.JSON(fiber.Map{"code": global.RET_ERR_NO_RIGHT,
			"data": "no right to do this"})
		return
	}
	table := c.Params("table")
	key := c.Params("key")
	if len(table) == 0 || len(key) == 0 {
		c.JSON(fiber.Map{"code": global.RET_ERR_URL_PARAM,
			"data": "request param err"})
		return
	}
	err := db.Delete(table, key)
	if err != nil {
		c.JSON(fiber.Map{"code": global.RET_ERR_DB,
			"data": err.Error()})
		return
	}
	c.JSON(fiber.Map{"code": global.RET_OK, "data": nil})
}

type ParamSvn struct {
	RepoUrl string `json:"repourl"`
	Limit   string `json:"limit"`
	Version string `json:"version"`
}

func SvnLog(c *fiber.Ctx) {
	const timeout = time.Minute
	var paramSvn ParamSvn
	if err := c.QueryParser(&paramSvn); err == nil {
		if len(paramSvn.Version) != 0 {
			cmd := fmt.Sprintf("svn log -%s %s",
				paramSvn.Version, paramSvn.RepoUrl)
			e, _, err := expect.Spawn(cmd, -1)
			defer e.Close()
			if err != nil {
				c.JSON(fiber.Map{"code": global.RET_ERR_SPAWN,
					"data": cmd})
				return
			}
			// allRE := regexp.MustCompile(`[\s\S]`)
			ret, _, _ := e.Expect(nil, timeout)
			c.JSON(fiber.Map{"code": global.RET_OK,
				"data": ret})
			return
		}
		if len(paramSvn.RepoUrl) != 0 {
			if len(paramSvn.Limit) == 0 {
				paramSvn.Limit = "10"
			}
			cmd := fmt.Sprintf("svn log -l %s -q %s",
				paramSvn.Limit, paramSvn.RepoUrl)
			e, _, err := expect.Spawn(cmd, -1)
			defer e.Close()
			if err != nil {
				c.JSON(fiber.Map{"code": global.RET_ERR_SPAWN,
					"data": cmd})
				return
			}
			// allRE := regexp.MustCompile(`[\s\S]`)
			ret, _, _ := e.Expect(nil, timeout)
			c.JSON(fiber.Map{"code": global.RET_OK,
				"data": ret})
			return
		} else {
			c.JSON(fiber.Map{"code": global.RET_ERR_URL_PARAM,
				"data": `can't find param repourl`})
			return
		}
	} else {
		c.JSON(fiber.Map{"code": global.RET_ERR_HTTP_QUERY,
			"data": err.Error()})
		return
	}
}

type SubmitDepParam struct {
	DepId string `json:"depid"`
}
type DepInfo struct {
	Type     string `json:"type"`
	RepoUrl  string `json:"repourl"`
	Rversion string `json:"rversion"`
}

func SubmitDep(c *fiber.Ctx) {
	const timeout = time.Minute
	var subDepParam SubmitDepParam
	if err := c.QueryParser(&subDepParam); err == nil {
		if len(subDepParam.DepId) != 0 {
			depInfoInDB, _ := db.Get(global.BUCKET_OPS_DEPBIL, subDepParam.DepId)
			var depInfo DepInfo
			json.Unmarshal(depInfoInDB, &depInfo)
			path := utils.GetRepoPath()
			if len(path) == 0 {
				path = "."
			}
			cmd := fmt.Sprintf("svn checkout -%s %s %s/%s",
				depInfo.Rversion, depInfo.RepoUrl, path, depInfo.Type)
			e, _, err := expect.Spawn(cmd, -1)
			defer e.Close()
			if err != nil {
				c.JSON(fiber.Map{"code": global.RET_ERR_SPAWN,
					"data": cmd})
				return
			}
			checkedOutRE := regexp.MustCompile(fmt.Sprintf("Checked out revision %s", depInfo.Rversion))
			ret, _, _ := e.Expect(checkedOutRE, timeout)
			c.JSON(fiber.Map{"code": global.RET_OK,
				"data": ret})
			return
		} else {
			c.JSON(fiber.Map{"code": global.RET_ERR_URL_PARAM,
				"data": `can't find param depid`})
			return
		}
	} else {
		c.JSON(fiber.Map{"code": global.RET_ERR_HTTP_QUERY,
			"data": err.Error()})
		return
	}
}

func Rsync(c *fiber.Ctx) {
	const timeout = time.Minute
	fmt.Println(term.Bluef("SSH Example"))

	sshClt, err := ssh.Dial("tcp", "10.1.1.43:22000", &ssh.ClientConfig{
		User:            "jesse",
		Auth:            []ssh.AuthMethod{ssh.Password("m@@n1i9ht")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		fmt.Println("ssh.Dial failed")
	}
	defer sshClt.Close()

	e, _, err := expect.SpawnSSH(sshClt, timeout)
	if err != nil {
		fmt.Println(err)
	}
	defer e.Close()

	// e.ExpectBatch([]expect.Batcher{
	// 	&expect.BExp{R: "password:"},
	// 	&expect.BSnd{S: "m@@n1i9ht" + "\n"},
	// 	// &expect.BExp{R: "Last login:"},
	// 	// &expect.BSnd{S: "node\n"},
	// }, timeout)
	// allRE := regexp.MustCompile(`[\s\S]`)
	// // passRE := regexp.MustCompile("password:")
	logined := regexp.MustCompile("$")
	result, _, _ := e.Expect(logined, timeout)
	fmt.Println(result)
	e.Send("/opt/jesse/server/online/start.sh\n")
	started := regexp.MustCompile("启动")
	result2, _, _ := e.Expect(started, timeout)
	fmt.Println(result2)
	fmt.Println(term.Greenf("All done"))
	// const timeout = 10 * time.Minute
	// // userRE := regexp.MustCompile("username:")
	// // passRE := regexp.MustCompile("password:")
	// // allRE := regexp.MustCompile(`[\s\S]*`)
	// allRE := regexp.MustCompile(">")
	// fmt.Println(term.Bluef("Telnet 1 example"))

	// e, _, err := expect.Spawn("node", -1)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer e.Close()
	// // userRE = regexp.MustCompile("username:")
	// // passRE = regexp.MustCompile("password:")
	// // promptRE = regexp.MustCompile("%")

	// // e.Expect(userRE, timeout)
	// // e.Send("jesse" + "\n")
	// // ret, _, _ := e.Expect(passRE, timeout)
	// // fmt.Println(ret)
	// // e.Send("m@@n1i9ht" + "\n")
	// // fmt.Println("---")
	// ret2, _, _ := e.Expect(allRE, timeout)
	// fmt.Println("2", ret2)
	// // e.Expect(promptRE, timeout)
	// // e.Send("cmd" + "\n")
	// // result, _, _ := e.Expect(promptRE, timeout)
	// e1, _, err := expect.Spawn("cd vendor", -1)
	// e1.Expect(allRE, timeout)
	// fmt.Println("ls")
	// e3, _, err := expect.Spawn("node", -1)
	// ret3, _, _ := e3.Expect(allRE, timeout)
	// fmt.Println("3", ret3)
	// // fmt.Println(term.Greenf("%s: result: %s\n", "cmd", result))
}

var NewDB = UpdateDB
