/*
 * HomeWork-3: Simple blog
 * Created on 19.09.2019 19:23
 * Copyright (c) 2019 - Eugene Klimov
 */
package main

import (
	"gopkg.in/russross/blackfriday.v2"
	"html/template"
)

// Posts storage todo - need xSQL storage instead map
var Posts map[string]Post

var tIndex, tPost *template.Template

func init() {
	//tIndex = template.Must(template.ParseFiles(path.Join(templatePath, indexTemplate)))
	//tPost = template.Must(template.ParseFiles(path.Join(templatePath, postTemplate)))

	Posts = map[string]Post{
		"1": {
			"Мой первый пост!",
			"18-е Сентября 2019 года",
			"Lorem ipsum dolor sit amet, consectetur adipisicing elit. Odio praesentium, quos. Aspernatur assumenda cupiditate deserunt ducimus, eveniet, expedita inventore laboriosam magni modi non odio, officia qui sequi similique unde voluptatem.",
			template.HTML(blackfriday.Run([]byte(`
Здесь основной текст.
# Markdown!
*Это* **круто**!
`))),
		},
		"2": {
			"Это уже второй пост!",
			"19-е Сентября 2019 года",
			`
Lorem ipsum dolor sit amet, consectetur adipisicing elit. Odio praesentium, quos. Aspernatur assumenda cupiditate deserunt ducimus, eveniet, expedita inventore laboriosam magni modi non odio, officia qui sequi similique unde voluptatem.
Lorem ipsum dolor sit amet, consectetur adipisicing elit. Odio praesentium, quos. Aspernatur assumenda cupiditate deserunt ducimus, eveniet, expedita inventore laboriosam magni modi non odio, officia qui sequi similique unde voluptatem.
`,
			template.HTML(blackfriday.Run([]byte(`
Разобрался в шаблонах и маркдаунах, как их совместить.

Теперь понять, как переходить на отдельные посты.
# Anybody!
*Hz* **cool**!
`))),
		},
	}
}
