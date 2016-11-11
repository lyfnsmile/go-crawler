// main package

package main

import (
    "fmt"
    "crawler/crawlerData"
)

func main () {
    // 使用crawldata包里面的Crawl()抓取需要的数据存到数据库
    crawlerData.Crawl()
    fmt.Println("主函数")
}