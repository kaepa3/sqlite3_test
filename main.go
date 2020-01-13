package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

type Record struct {
	ID    int
	Title string
	Value string
}

func main() {

	list := createRecords(100)
	log.Println("Create record:" + strconv.Itoa(len(list)))

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "graf"
	p.X.Label.Text = "x"
	p.Y.Label.Text = "y"

	//addDBtime := calc(addDB, list, "addDB", 10)
	//addDBpertime := calc(addDBper, list, "addDBper", 10)
	//csvTime := calc(outCsv, list, "outCsv", 10)
	//jsonTime := calc(outJson, list, "outJson", 10)

	//if err := plotutil.AddLinePoints(p,
	//	"db commit", point(addDBtime),
	//	"csv", point(csvTime),
	//	"json", point(jsonTime),
	//	"db exec", point(addDBpertime)); err != nil {
	//	panic(err)
	//}

	//if err := p.Save(5*vg.Inch, 5*vg.Inch, "graf.png"); err != nil {
	//	panic(err)
	//}
	list1000 := createRecords(1000)
	list10000 := createRecords(10000)

	result := make([]float64, 3)
	result[0] = calc(addDB, list, "addDB", 1)[0]
	result[1] = calc(addDB, list1000, "addDB", 1)[0]
	result[2] = calc(addDB, list10000, "addDB", 1)[0]

	resultjson := make([]float64, 3)
	resultjson[0] = calc(outJson, list, "outJson", 1)[0]
	resultjson[1] = calc(outJson, list1000, "outJson", 1)[0]
	resultjson[2] = calc(outJson, list10000, "outJson", 1)[0]

	resultcsv := make([]float64, 3)
	resultcsv[0] = calc(outCsv, list, "outCsv", 1)[0]
	resultcsv[1] = calc(outCsv, list1000, "outCsv", 1)[0]
	resultcsv[2] = calc(outCsv, list10000, "outCsv", 1)[0]

	if err := plotutil.AddLinePoints(p,
		"db ", point(result),
		"csv", point(resultcsv),
		"json", point(resultjson)); err != nil {
		panic(err)
	}
	if err := p.Save(5*vg.Inch, 5*vg.Inch, "plot.png"); err != nil {
		panic(err)
	}
}

func point(l []float64) plotter.XYs {

	pts := make(plotter.XYs, len(l))

	for i, v := range l {
		pts[i].X = float64(i)
		pts[i].Y = v
	}

	return pts
}

func calc(f func([]Record), r []Record, comment string, count int) []float64 {
	list := make([]float64, count)
	log.Println(comment)
	for idx := range list {
		start := time.Now()
		f(r)
		end := time.Now()
		list[idx] = (end.Sub(start)).Seconds()
		log.Println(list[idx])
	}
	return list
}

func createRecords(count int) []Record {
	list := make([]Record, count)
	for i := range list {
		list[i].ID = i
		list[i].Title = RandString(5)
		list[i].Value = RandString(7)
	}
	return list
}

var rs1Letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = rs1Letters[rand.Intn(len(rs1Letters))]
	}
	return string(b)
}

func initDB() *sql.DB {

	// データベースのコネクションを開く
	db, err := sql.Open("sqlite3", "./test.db")
	if err == nil {
		_, err = db.Exec(
			`CREATE TABLE IF NOT EXISTS "BOOKS" ("ID" INTEGER PRIMARY KEY, "TITLE" TEXT, "VALUE" TEXT )`,
		)
		if err != nil {
			log.Println("err")
			log.Println(err)
		}
		db.Exec("DELETE FROM BOOKS")

	} else {
		log.Println("err")
		log.Println(err)
	}

	return db
}
func addDB(r []Record) {
	db := initDB()
	if db != nil {
		defer db.Close()
	}
	if tx, err := db.Begin(); err == nil {
		if stmt, err := tx.Prepare(`INSERT INTO BOOKS (ID,TITLE,VALUE) VALUES (?, ?, ?);`); err == nil {
			defer stmt.Close()
			for _, v := range r {
				_, err := stmt.Exec(
					v.ID,
					v.Title,
					v.Value,
				)
				if err != nil {
					log.Println(err)
				}
			}
			tx.Commit()
		} else {
			log.Println(err)
		}
	} else {
		log.Println(err)
	}
}
func addDBper(r []Record) {
	db := initDB()
	if db != nil {
		defer db.Close()
	}
	for _, v := range r {
		_, err := db.Exec(
			`INSERT INTO BOOKS (ID, TITLE, VALUE) VALUES (?, ?, ?)`,
			v.ID+1000,
			v.Title,
			v.Value,
		)
		if err != nil {
			log.Println(err)
		}
	}

}
func deleteFile() {
	list := []string{"test.csv", "test.json"}
	for _, v := range list {

		if err := os.Remove(v); err != nil {
			fmt.Println(err)
		}
	}
}

func outCsv(r []Record) {
	file, err := os.Create("test.csv")
	if err == nil {
		defer file.Close()
		for _, v := range r {
			text := fmt.Sprintf("%d, %s, %d\n", v.ID, v.Title, v.Value)
			file.WriteString(text)
		}
	} else {
		fmt.Println("error")

	}
}

func outJson(r []Record) {
	file, err := os.Create("test.json")
	if err == nil {
		outputJson, err := json.Marshal(&r)
		if err == nil {
			// jsonデータを出力
			fmt.Fprint(file, string(outputJson))
		}
	} else {
		fmt.Println("error")
	}
}
