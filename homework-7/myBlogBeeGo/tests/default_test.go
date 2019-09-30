/*
 * HomeWork-6: Mongo in BeeGo
 * Created on 28.09.19 22:22
 * Copyright (c) 2019 - Eugene Klimov
 */

package test

import (
	"context"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"myBlog/models"
	_ "myBlog/routers"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func init() {
	// connect to Mongo
	mdb, err := mongo.NewClient(options.Client().ApplyURI("mongodb://Klim.Go:27017"))
	if err != nil {
		log.Fatalln("Can't create MongoDB client:", err)
	}
	models.MDB = mdb
	if err = models.MDB.Connect(context.TODO()); err != nil {
		log.Fatalln("Can't connect to MongoDB server:", err)
	}
	if err = models.MDB.Ping(context.TODO(), nil); err != nil {
		log.Fatalln("Can't ping MongoDB server:", err)
	}

	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

// TestBeego is a sample to run an endpoint test
func TestBeego(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	lg := logs.NewLogger(10)
	lg.Trace("Testing myBlog, code[%d]\n", w.Code)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
	})
}

func TestCreatePost(t *testing.T) {
	posts := models.NewPosts()
	post := models.Post{
		Title:   "Test",
		Date:    time.Now(),
		Summary: "Test from tests",
		Body:    "### Testik",
		Created: time.Now(),
		Deleted: time.Unix(0, 0),
	}
	posts.Posts = append(posts.Posts, post)
	if err := posts.CreatePost(); err != nil {
		t.Errorf("Eror create new post: %s\n", err)
	}
	t.Logf("PASS: Create post")
}