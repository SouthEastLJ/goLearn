package main

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
)

/**
在dao层遇到sql.ErrNoRows的时候，应该Wrap这个error。
原因：sql包返回了根错误值，在业务中可以wrap这个error，然后往上抛;最终捕获这个错误，记录堆栈信息，并作出处理。
 */

func query() error {
	// 省略db.QueryRow的业务代码，直接模拟ErrNoRows错误
	return errors.Wrap(sql.ErrNoRows, "query failed")
}


func main() {
	err := query()
	if errors.Cause(err) == sql.ErrNoRows {
		fmt.Printf("data not found, %v\n", err)
		fmt.Printf("%+v\n", err)
		return
	}
	if err != nil {
		// 未知其他错误，做其他处理
	}
}
