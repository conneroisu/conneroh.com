package main

import (
	"github.com/conneroisu/conneroh.com/internal/data/css"
	"github.com/rotisserie/eris"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {

	for k, v := range css.ClassMapStr {
		for k2, v2 := range css.ClassMapStr {
			if v == v2 {
				return eris.Errorf(`
					'%s' : %s
					'%s' : %s
					Duplicate class name
				`, k, k2, v, v2)
			}
		}
	}
	return nil
}
