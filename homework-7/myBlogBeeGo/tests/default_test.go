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

var testDBCases = []struct {
	action      string
	description string
}{
	{action: "CreatePost", description: "create new post"},
	{action: "GetAllPosts", description: "get all posts"},
	{action: "GetOnePost", description: "get one post"},
	{action: "UpdatePost", description: "update post"},
	{action: "DeletePost", description: "delete post"},
}

var testMethodsCases = []struct {
	method string
	api    string
	action string
	code   int
}{
	{
		action: "PostCreatePost",
		method: "POST",
		api:    "/api/v1/posts/create",
		code:   http.StatusCreated},
	{
		action: "GetPost",
		method: "GET",
		api:    "/posts",
		code:   http.StatusOK,
	},
	{
		action: "PutUpdatePost",
		method: "PUT",
		api:    "/api/v1/posts",
		code:   http.StatusOK,
	},
	{
		action: "GetEditPost",
		method: "GET",
		api:    "/posts/edit",
		code:   http.StatusOK,
	},
	{
		action: "GetCreatePost",
		method: "GET",
		api:    "/posts/create",
		code:   http.StatusOK,
	},
	{
		action: "DeleteDeletePost",
		method: "DELETE",
		api:    "/api/v1/posts",
		code:   http.StatusOK,
	},
}

//func TestDB(t *testing.T) {
//	for _, test := range testDBCases {
//		var err error
//		posts := models.NewPosts()
//
//		switch test.action {
//		case "CreatePost":
//			posts.Posts = append(posts.Posts, post)
//			err = posts.CreatePost()
//			post.ID = post.OID.Hex()
//		case "GetAllPosts":
//			err = posts.GetPosts("")
//		case "GetOnePost":
//			err = posts.GetPosts(post.ID)
//		case "UpdatePost":
//			post.Title = post.Title + "_Updated"
//			post.Summary = post.Summary + "_Updated"
//			post.Body = post.Body + "_Updated"
//			posts.Posts = append(posts.Posts, post)
//			err = posts.UpdatePost(post.ID, false)
//		case "DeletePost":
//			err = posts.DeletePost(post.ID)
//		}
//
//		if err != nil {
//			posts.Lg.Error("Error %s in DB: %s", test.description, err)
//			return
//		}
//		posts.Lg.Informational("PASS: %s in DB", test.description)
//	}
//}

func TestMethods(t *testing.T) {
	for _, test := range testMethodsCases {
		var r *http.Request
		posts := models.NewPosts()

		switch test.action {
		case "PostCreatePost":
			post.ID = randomHex(12)
			body := `
			{
			   "id":"` + post.ID + `",
				"title":"qqq",
				"summary":"www",
				"body":"**eee**"
			}`
			r, _ = http.NewRequest(test.method, test.api, strings.NewReader(body))
		case "GetPost", "GetEditPost":
			r, _ = http.NewRequest(test.method, test.api+"/?id="+post.ID, nil)
		case "PutUpdatePost":
			body := `
			{
				"title":"111",
				"summary":"222",
				"body":"*333*"
			}`
			r, _ = http.NewRequest(test.method, test.api+"/"+post.ID, strings.NewReader(body))
		case "DeleteDeletePost":
			r, _ = http.NewRequest(test.method, test.api+"/"+post.ID, nil)
		case "GetCreatePost":
			r, _ = http.NewRequest(test.method, test.api, nil)
		}

		w := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		if w.Code != test.code {
			posts.Lg.Error("Error %s %s, code[%d]", test.method, test.api, w.Code)
			return
		}
		posts.Lg.Informational("PASS: %s %s", test.method, test.api)
	}
}

func randomHex(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}
