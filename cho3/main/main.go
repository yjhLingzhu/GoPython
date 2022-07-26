package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	gorm.Model
	ID   uint
	Name string
	Age  int
}

func test(arr []int) {
	arr[0] = 999
	fmt.Println(arr)
}

func main() {
	// 链接数据库
	db, err := gorm.Open("mysql", "root:123456@(127.0.0.1:3306)/article_spider?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	// 自动迁移
	// db.AutoMigrate(&User{})

	// u1 := User{Name: "yjh", Age: 20} // 这里如果不写上ID 的话，下面的NewRecord一直是true
	// u2 := User{Name: "fzj", Age: 19}

	// fmt.Println(db.NewRecord(&u1)) // 这个NewRecord是对比这个记录的所有字段的，完全相同才算已存在
	// db.Create(&u1)
	// db.Create(&u2)

	// 查询
	var u User
	// db.First(&u) // 将查询出来的这条数据和u做关联
	// fmt.Printf("%#v\n", u)

	// 更新
	db.Model(&u).Update("name", "hahah") // 对查出来的数据进行更新
	m1 := map[string]interface{}{
		"Name": "fhj",
		"Age":  19,
	}
	db.Debug().Model(&User{}).Where("name = ?", "yjh11").Updates(m1)

	// 删除
	// db.Delete(&u) // 删除数据，如果有delete_at的话就将时间写进去, 否则真正删除
	// db.Debug().Where("name=?", "yjh").Delete(&User{})
	// db.Debug().Delete(&User{}, "name=? and age=?", "fhj", 20)

	// 真正删除还有一个方式，即使它包含了delete_at字段，就是使用Unscoped()这个方法
	// db.Debug().Unscoped().Where("name=? and age=?", "daha", 10).Delete(&User{})

	// var arr = []int{1, 2} // 指针、接口、map、切片这四个是引用类型
	// test(arr)
	// fmt.Println(arr)

	// FirstOrInit/Attrs
	// 未找到
	// db.FirstOrInit(&u, User{Name: "non_existing"})
	// db.Where(&User{Name: "non_existing"}).Attrs(&User{Age: 30}).FirstOrInit(&u)
	// fmt.Printf("%#v", u)
	// db.Create(&u)

	//Scan
	// type Info struct {
	// 	Name string "json:name"
	// 	Age  int    "json:age"
	// }
	// var info []Info
	// db.Debug().Model(&User{}).Select("name, age").Where("name = ?", "yjh").Scan(&info)
	// b, _ := json.Marshal(info)
	// fmt.Printf("%#v", string(b))

	// Rows
	// rows, err := db.Debug().Model(&User{}).Select("name, age").Where("name = ?", "yjh").Rows()
	// defer rows.Close()
	// fmt.Println(rows)
	// for rows.Next() {
	// 	var info Info
	// 	// rows.Scan(&info.Age, &info.Name)		// 这个赋值不了
	// 	db.ScanRows(rows, &info)
	// 	fmt.Printf("%#v\n", info)
	// }
}
