/*
 * HomeWork-6: Mongo in BeeGo
 * Created on 28.09.19 22:22
 * Copyright (c) 2019 - Eugene Klimov
 */

package test

import (
	"context"
	"encoding/hex"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math/rand"
	"myBlog/models"
	_ "myBlog/routers"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"strings"
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
	rand.Seed(time.Now().UnixNano())

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

func TestDBCreatePost(t *testing.T) {
	posts := models.NewPosts()
	posts.Posts = append(posts.Posts, post)
	if err := posts.CreatePost(); err != nil {
		posts.Lg.Error("Error create new post in DB: %s", err)
		return
	}
	post.ID = post.OID.Hex()
	posts.Lg.Informational("PASS: Create post in DB")
}

func TestDBGetAllPosts(t *testing.T) {
	posts := models.NewPosts()
	if err := posts.GetPosts(""); err != nil {
		posts.Lg.Error("Error get all posts from DB: %s", err)
		return
	}
	posts.Lg.Informational("PASS: Get all posts from DB")
}

func TestDBGetOnePost(t *testing.T) {
	posts := models.NewPosts()
	if err := posts.GetPosts(post.ID); err != nil {
		posts.Lg.Error("Error get one post from DB: %s", err)
		return
	}
	posts.Lg.Informational("PASS: Get one post from DB")
}

func TestDBUpdatePost(t *testing.T) {
	posts := models.NewPosts()
	post.Title = post.Title + "_Updated"
	post.Summary = post.Summary + "_Updated"
	post.Body = post.Body + "_Updated"
	posts.Posts = append(posts.Posts, post)
	if err := posts.UpdatePost(post.ID, false); err != nil {
		posts.Lg.Error("Error update post in DB: %s", err)
		return
	}
	posts.Lg.Informational("PASS: Update post in DB")
}

func TestGetPost(t *testing.T) {
	posts := models.NewPosts()
	r, _ := http.NewRequest("GET", "/posts/?id="+post.ID, nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		posts.Lg.Error("Error GET /posts, code[%d]", w.Code)
		return
	}
	posts.Lg.Informational("PASS: GET /posts")
}

func TestGetEditPost(t *testing.T) {
	posts := models.NewPosts()
	r, _ := http.NewRequest("GET", "/posts/edit/?id="+post.ID, nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		posts.Lg.Error("Error GET /posts/edit, code[%d]", w.Code)
		return
	}
	posts.Lg.Informational("PASS: GET /posts/edit")
}

func TestGetCreatePost(t *testing.T) {
	posts := models.NewPosts()
	r, _ := http.NewRequest("GET", "/posts/create", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		posts.Lg.Error("Error GET /posts/create, code[%d]", w.Code)
		return
	}
	posts.Lg.Informational("PASS: GET /posts/create")
}

func TestPutUpdatePost(t *testing.T) {
	posts := models.NewPosts()
	body := `
{
	"title":"111",
	"summary":"222",
	"body":"*333*"
}`
	r, _ := http.NewRequest("PUT", "/api/v1/posts/"+post.ID, strings.NewReader(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		posts.Lg.Error("Error PUT /api/v1/posts, code[%d]", w.Code)
		return
	}
	posts.Lg.Informational("PASS: PUT //api/v1/posts/")
}

func TestDeletePost(t *testing.T) {
	posts := models.NewPosts()
	if err := posts.DeletePost(post.ID); err != nil {
		posts.Lg.Error("Error Delete post in DB: %s", err)
		return
	}
	posts.Lg.Informational("PASS: Delete post in DB")
}

func TestPostCreatePost(t *testing.T) {
	posts := models.NewPosts()
	post.ID = randomHex(12)
	body := `
{
    "id":"` + post.ID + `",
	"title":"qqq",
	"summary":"www",
	"body":"**eee**"
}`
	r, _ := http.NewRequest("POST", "/api/v1/posts/create", strings.NewReader(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	if w.Code != http.StatusCreated {
		posts.Lg.Error("Error POST /api/v1/posts/create, code[%d]", w.Code)
		return
	}
	posts.Lg.Informational("PASS: POST /api/v1/posts/create")
}

func TestDeleteDeletePost(t *testing.T) {
	posts := models.NewPosts()
	r, _ := http.NewRequest("DELETE", "/api/v1/posts/"+post.ID, nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		posts.Lg.Error("Error DELETE /api/v1/posts, code[%d]", w.Code)
		return
	}
	posts.Lg.Informational("PASS: DELETE /api/v1/posts")
}

func randomHex(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}
