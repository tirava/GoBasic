/*
 * HomeWork-5: Start BeeGo
 * Created on 26.09.19 19:50
 * Copyright (c) 2019 - Eugene Klimov
 */

package models

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"net/http"
)

// Error model.
type Error struct {
	ErrCode  int             `json:"code"`
	ErrText  string          `json:"error"`
	ErrDescr string          `json:"descr"`
	Lg       *logs.BeeLogger `json:"-"`
}

// SendError is errors helper
func (e *Error) SendError(w http.ResponseWriter, code int, err error, descr string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	errMsg := Error{
		ErrCode:  code,
		ErrText:  err.Error(),
		ErrDescr: descr,
	}
	data, err := json.Marshal(errMsg)
	if err != nil {
		e.Lg.Error("Can't marshal error data: %s", err)
		return
	}
	if _, err = w.Write(data); err != nil {
		e.Lg.Error("Can't write to ResponseWriter: %s", err)
	}
}
