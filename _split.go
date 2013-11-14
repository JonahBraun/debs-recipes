package main

import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "io/ioutil"
import "os"

func main(){
	db, err := sql.Open("mysql", "root:q@tcp(127.0.0.1:3306)/recipe?charset=utf8")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	var(
		acreated, calias, ctitle, aalias, atitle, aintro, atext string
		aid int
		cats string
	)

	crows, err := db.Query("select IFNULL(c.alias, \"alt\"), IFNULL(c.title, \"NULL\") from jos_categories as c order by c.title;")
	if err != nil {
		panic(err)
	}
	defer crows.Close()

	for crows.Next(){
		err := crows.Scan(&calias, &ctitle)
		if err != nil {
			panic(err)
		}

		println(calias, ctitle)
		err = os.Mkdir("public/"+calias, 0777)
		err = os.Mkdir("menu/"+calias, 0777)
		if err != nil {
			//panic(err)
		}
		cats += calias+" "+ctitle+"\n"
	}

	/*err = ioutil.WriteFile("menu/categories.yaml", []byte(cats), 0666)
	if err != nil {
		panic(err)
	}*/


	rows, err := db.Query(`select date_format(a.created, '%Y-%m-%d'),
		IFNULL(c.alias, "alt"), IFNULL(c.title, "NULL"),
		a.id, a.alias, a.title, a.introtext, a.fulltext
		from jos_content as a
		left join jos_categories as c
		on a.catid = c.id order by a.id;`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next(){
		err := rows.Scan(&acreated, &calias, &ctitle, &aid, &aalias, &atitle, &aintro, &atext)
		if err != nil {
			panic(err)
		}

		txt := []byte("---\nlayout: recipe\ntitle: "+atitle+"\ncategory: "+calias+"\n---\n"+aintro+"\n"+atext)

		//ioutil.WriteFile("menu/"+calias+"/"+acreated+"-"+aalias+".html", txt, 0666)
		ioutil.WriteFile("_posts/"+acreated+"-"+aalias+".html", txt, 0666)

	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

}
