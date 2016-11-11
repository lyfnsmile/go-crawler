// store crawler Data

package crawlerData

import (
    "fmt"
    "github.com/PuerkitoBio/goquery"
    "strconv"
    s "strings"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

// 定义一个存储一条数据的结构体
type ImageData struct {
    Src    string
    Tp     string
    Title  string
    Width  int
    Height int
}

// 定义切片用于存储抓取的全部数据
type ImageDatas []ImageData

func Crawl() {
    fmt.Println("包crawlerdata中的Crawl函数")

    // 定义一个切片存储所有数据
    var datas ImageDatas
    // 抓取数据
    imageDatas := CrawlData(&datas)

    for _, imageData := range imageDatas {
        fmt.Println(imageData.Src, imageData.Title, imageData.Tp, imageData.Height, imageData.Width)
    }

}

func OpenDatabase() (*sql.DB, error) {
    // 连接数据库
    db, err := sql.Open("mysql", "root:mysql@tcp(127.0.0.1:3306)/lizx?charset=utf8")
    if err != nil {
        return nil, err
    }
    return db, nil
}


/*
   该函数用来抓取数据，并将存储的值返回到主函数
*/
func CrawlData(datas *ImageDatas) (imageDatas ImageDatas) {
    imageDatas = *datas
    // 规定抓取时匹配的元素
    var types = [...]string{
        "people",
        "objects",
        "whimsical",
        "nature",
        "urban",
        "animals"}

    doc, err := goquery.NewDocument("http://www.gratisography.com/")
    if err != nil {
        fmt.Printf(err.Error())
    }

    for _, tp := range types {
        doc.Find("#container ul").Find(s.Join([]string{".", tp}, "")).Each(func(i int, s *goquery.Selection) {
            img := s.Find("img.lazy")
            src, _ := img.Attr("data-original")
            title, _ := img.Attr("alt")
            width, _ := img.Attr("width")
            height, _ := img.Attr("height")

            // 将宽度和高度的字符串类型转为数值型
            wd, error := strconv.Atoi(width)
            if error != nil {
                fmt.Println("字符串转换成整数失败")
            }
            hg, error := strconv.Atoi(height)
            if error != nil {
                fmt.Println("字符串转换成整数失败")
            }
            // fmt.Printf("Review %d: %s - %s - %s - %d - %d\n", i, src, tp, title, wd, hg)
            imageData := ImageData{src, tp, title, wd, hg}
            imageDatas = append(imageDatas, imageData)
        })
    }

    InsertData(&imageDatas)
    return
}

/*
   该函数将获取的数据存储到数据库
*/

func GetAllImages() (imageDatas ImageDatas, err error) {
    // 连接数据库
    db, err := OpenDatabase()
    if err != nil {
        fmt.Printf(s.Join([]string{"连接数据库失败", err.Error()}, "-->"))
        return nil, err
    }
    defer db.Close()

    // Prepare statement for inserting data
    imgOut, err := db.Query("SELECT * FROM gratisography")
    if err != nil {
        fmt.Println(s.Join([]string{"获取数据失败", err.Error()}, "-->"))
        return nil, err
    }
    defer imgOut.Close()

    // 定义扫描select到的数据库字段的变量
    var (
        id          int
        img_url     string
        type_name   string
        title       string
        width       int
        height      int
        create_time string
    )
    for imgOut.Next() {
        // db.Query()中select几个字段就需要Scan多少个字段
        err := imgOut.Scan(&id, &img_url, &type_name, &title, &width, &height, &create_time)
        if err != nil {
            fmt.Println(s.Join([]string{"查询数据失败", err.Error()}, "-->"))
            return nil, err
        } else {
            imageData := ImageData{img_url, type_name, title, width, height}
            imageDatas = append(imageDatas, imageData)
        }
    }

    return imageDatas, nil
}

/*
   该函数将获取的数据存储到数据库
*/
func InsertData(datas *ImageDatas) {
    imageDatas := *datas
    // 连接数据库
    db, err := OpenDatabase()
    if err != nil {
        fmt.Printf(s.Join([]string{"连接数据库失败", err.Error()}, "-->"))
    }
    defer db.Close()

    for i := 0; i < len(imageDatas); i++ {
        imageData := imageDatas[i]
        // Prepare statement for inserting data
        imgIns, err := db.Prepare("INSERT INTO gratisography (img_url, type_name, title, width, height) VALUES( ?, ?, ?, ?, ? )") // ? = placeholder
        if err != nil {
            fmt.Println(s.Join([]string{"拼装数据格式", err.Error()}, "-->"))
        }
        defer imgIns.Close() // Close the statement when we leave main()

        img, err := imgIns.Exec(s.Join([]string{"http://www.gratisography.com", imageData.Src}, "/"), imageData.Tp, imageData.Title, imageData.Width, imageData.Height)
        if err != nil {
            fmt.Println(s.Join([]string{"插入数据失败", err.Error()}, "-->"))
        } else {
            success, _ := img.LastInsertId()
            // 数字变成字符串,success是int64型的值，需要转为int，网上说的Itoa64()在strconv包里不存在
            insertId := strconv.Itoa(int(success))
            fmt.Println(s.Join([]string{"成功插入数据：", insertId}, "\t-->\t"))
        }
    }
}

func GetTpImages(tp string) (imageDatas ImageDatas, err error) {
    // 连接数据库
    db, err := OpenDatabase()
    if err != nil {
        fmt.Printf(s.Join([]string{"连接数据库失败", err.Error()}, "-->"))
        return nil, err
    }
    defer db.Close()

    // Prepare statement for inserting data
    fmt.Printf("SELECT * FROM gratisography where type_name="+tp)
    imgOut, err := db.Query("SELECT * FROM gratisography where type_name=?",tp)
    if err != nil {
        fmt.Println(s.Join([]string{"获取数据失败", err.Error()}, "-->"))
        return nil, err
    }
    defer imgOut.Close()

    // 定义扫描select到的数据库字段的变量
    var (
        id          int
        img_url     string
        type_name   string
        title       string
        width       int
        height      int
        create_time string
    )
    for imgOut.Next() {
        // db.Query()中select几个字段就需要Scan多少个字段
        err := imgOut.Scan(&id, &img_url, &type_name, &title, &width, &height, &create_time)
        if err != nil {
            fmt.Println(s.Join([]string{"查询数据失败", err.Error()}, "-->"))
            return nil, err
        } else {
            imageData := ImageData{img_url, type_name, title, width, height}
            imageDatas = append(imageDatas, imageData)
        }
    }

    return imageDatas, nil
}