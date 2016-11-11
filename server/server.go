// main package

package main

import (
    "fmt"
    "github.com/go-martini/martini"
    "github.com/martini-contrib/encoder"
    "crawler/crawlerData"
    "net/http"
    s "strings"
)

type Results struct {
    Err   int                  // 错误码
    Msg   string               // 错误信息
    Datas crawlerData.ImageDatas // 数据，无数据时为nil
}

func main () {
    // 使用crawldata包里面的Crawl()抓取需要的数据存到数据库
    //crawlerData.Crawl()
    //fmt.Println("主函数")
    m := martini.New()
    route := martini.NewRouter()

    var (
        results Results
        err     error
    )

    m.Use(func(c martini.Context, w http.ResponseWriter, r *http.Request) {
        // 将encoder.JsonEncoder{}按照encoder.Encoder接口（注意大小写）类型注入到内部
        c.MapTo(encoder.JsonEncoder{}, (*encoder.Encoder)(nil))
        w.Header().Set("Content-Type", "application/json; charset=utf-8")
    })

    route.Get("/", func(enc encoder.Encoder) (int, []byte) {
        result := Results{10001, "Not Found Data", nil}
        return http.StatusOK, encoder.Must(enc.Encode(result))
    })

    route.Get("/api", func(enc encoder.Encoder) (int, []byte) {
        results.Datas, err = crawlerData.GetAllImages()
        if err != nil {
            fmt.Println(s.Join([]string{"获取数据失败", err.Error()}, "-->"))
            result := Results{10001, "Data Error", nil}
            return http.StatusOK, encoder.Must(enc.Encode(result))
        } else {
            results.Err = 10001
            results.Msg = "获取数据成功"
            return http.StatusOK, encoder.Must(enc.Encode(results))
        }
    })

    route.Get("/api/:tp", func(enc encoder.Encoder,params martini.Params) (int, []byte) {
        tp :=params["tp"]
        fmt.Println("params %d\n",tp)
        results.Datas, err = crawlerData.GetTpImages(tp)
        if err != nil {
            fmt.Println(s.Join([]string{"获取数据失败", err.Error()}, "-->"))
            result := Results{10001, "Data Error", nil}
            return http.StatusOK, encoder.Must(enc.Encode(result))
        } else {
            results.Err = 10001
            results.Msg = "获取数据成功"
            return http.StatusOK, encoder.Must(enc.Encode(results))
        }
    })

    m.Action(route.Handle)
    m.Run()
}