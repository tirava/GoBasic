package test

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/go-ini/ini"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"myBlog/models"
	_ "myBlog/routers"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"
)

func init() {
	var err error
	dbName := beego.AppConfig.String("DBNAME")
	if dbName == "" {
		dbName = "blog"
	}
	dsn := beego.AppConfig.String("DSN")
	if dsn == "" {
		dsn = "?charset=utf8&interpolateParams=true"
	}
	// BeeGo ORM
	err = orm.RegisterDataBase("default", "mysql", myCnf("client")+"/"+dbName+dsn)
	if err != nil {
		log.Fatalln("Can't open BeeGo DB:", err)
	}
	orm.RegisterModel(new(models.Post))
	models.ORM = orm.NewOrm()

	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

// TestBeego is a sample to run an endpoint test
func TestBeego(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("testing", "TestBeego", "Code[%d]\n%s", w.Code, w.Body.String())

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
	})
}

// myCnf reads MySQL parameters from .my.cnf
func myCnf(profile string) string {
	cnf := path.Join(os.Getenv("HOME"), ".my.cnf")
	cfg, err := ini.LoadSources(ini.LoadOptions{AllowBooleanKeys: true}, cnf)
	if err != nil {
		return ""
	}
	for _, s := range cfg.Sections() {
		if s.Name() != profile {
			continue
		}
		user := s.Key("user")
		password := s.Key("password")
		host := s.Key("host")
		port := s.Key("port")
		return fmt.Sprintf("%s:%s@tcp(%s:%s)", user, password, host, port)
	}
	return ""
}
