/*
 * HomeWork-5: Start BeeGo
 * Created on 25.09.19 23:21
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"fmt"
	"github.com/go-ini/ini"
	"os"
	"path"
)

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
