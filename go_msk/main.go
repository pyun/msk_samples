/**
 * @Author: pyun
 * @Date: 2021/2/25 23:01
 * @Desc:
 */

package main

import (
	"flag"
	"go_demo/service"
)

func main() {
	iam := flag.Bool("iam", false, "是否采用iam角色认证，默认false")
	flag.Parse()
	if iam != nil {
	  service.Start(*iam)
	}
}
