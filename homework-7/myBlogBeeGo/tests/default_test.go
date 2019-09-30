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
	"go.mongodb.org/mongo-driver/bson/primitive"
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

var post = models.Post{
	ID:      "fake",
	OID:     primitive.NewObjectID(),
	Title:   "Test",
	Date:    time.Now(),
	Summary: "Test from tests",
	Body:    "### Testik",
	Created: time.Now(),
	Deleted: time.Unix(0, 0),
}

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

	models.Lg = logs.NewLogger(10)

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
	lg.Trace("Testing myBlog endpoint, code[%d]", w.Code)

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
	posts.Posts = append(posts.Posts, post)
	if err := posts.CreatePost(); err != nil {
		posts.Lg.Error("Error create new post: %s", err)
		return
	}
	post.ID = post.OID.Hex()
	posts.Lg.Informational("PASS: Create post")
}

func TestGetAllPosts(t *testing.T) {
	posts := models.NewPosts()
	if err := posts.GetPosts(""); err != nil {
		posts.Lg.Error("Error get all posts: %s", err)
		return
	}
	posts.Lg.Informational("PASS: Get all posts")
}

func TestGetOnePost(t *testing.T) {
	posts := models.NewPosts()
	if err := posts.GetPosts(post.ID); err != nil {
		posts.Lg.Error("Error get one post: %s", err)
		return
	}
	posts.Lg.Informational("PASS: Get one post")
}

func TestUpdatePost(t *testing.T) {
	posts := models.NewPosts()
	post.Title = post.Title + "_Updated"
	post.Summary = post.Summary + "_Updated"
	post.Body = post.Body + "_Updated"
	posts.Posts = append(posts.Posts, post)
	if err := posts.UpdatePost(post.ID, false); err != nil {
		posts.Lg.Error("Error update post: %s", err)
		return
	}
	posts.Lg.Informational("PASS: Update post")
}

func TestDeletePost(t *testing.T) {
	posts := models.NewPosts()
	if err := posts.DeletePost(post.ID); err != nil {
		posts.Lg.Error("Error delete post: %s", err)
		return
	}
	posts.Lg.Informational("PASS: Delete post")
}
