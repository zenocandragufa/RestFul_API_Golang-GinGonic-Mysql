package main

import (
	"bytes"
	"fmt"
	"net/http"
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/uts")
	err = db.Ping()
	if err != nil {
		panic("Gagal Menghubungkan ke Database...")
	}
	defer db.Close()

	router := gin.Default()

	type Progdi struct {
		Id			int		`json: "id"`
		Jenjang		string	`json: "jenjang"`
		NmProgdi	string	`json: "nmprodi"`
	}
	
	// Menampilkan Detail Data Berdasarkan ID
	router.GET("/:id", func(c *gin.Context) {
		var (
			progdi Progdi
			result gin.H
		)
		id := c.Param("id")
		row := db.QueryRow("select id, jenjang, nmprogdi from progdi where id = ?;", id)
		err = row.Scan(&progdi.Id, &progdi.Jenjang, &progdi.NmProgdi)
		if err != nil {
			// If no results send null
			result = gin.H{
				"Hasile" : "Yahh Tidak ada data yang ditemukan",
				"jumlahe" : 0,
			}
		} else {
			result = gin.H{
				"Hasile": progdi,
				"Jumlahe":  1,
			}
		}
		c.JSON(http.StatusOK, result)
	})




	// Menampilkan Semua Data 
	router.GET("/", func(c *gin.Context) {
		var (
			progdi  Progdi
			progdis []Progdi
		)
		rows, err := db.Query("select id, jenjang, nmprogdi from progdi;")
		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&progdi.Id, &progdi.Jenjang, &progdi.NmProgdi)
			progdis = append(progdis, progdi)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		c.JSON(http.StatusOK, gin.H{
			"Hasile": progdis,
			"Jumlahe":  len(progdis),
		})
	})


	// Menambahkan Data Program Studi
	router.POST("/", func(c *gin.Context) {
		var buffer bytes.Buffer
		id := c.PostForm("id")
		jenjang := c.PostForm("jenjang")
		nmprogdi := c.PostForm("nmprogdi")
		stmt, err := db.Prepare("insert into progdi (id, jenjang, nmprogdi) values(?,?,?);")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(id, jenjang, nmprogdi)

		if err != nil {
			fmt.Print(err.Error())
		}

		// Fastest way to append strings
		buffer.WriteString(jenjang)
		buffer.WriteString(" ")
		buffer.WriteString(nmprogdi)
		defer stmt.Close()
		datane := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"Pesane": fmt.Sprintf(" Yeyyy Berhasil menambahkan Progdi %s ", datane),
		})
	})


	// PUT Merubah Data
	router.PUT("/", func(c *gin.Context) {
		var buffer bytes.Buffer
		id := c.PostForm("id")
		jenjang := c.PostForm("jenjang")
		nmprogdi := c.PostForm("nmprogdi")
		stmt, err := db.Prepare("update progdi set jenjang= ?, nmprogdi= ? where id= ?;")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(jenjang, nmprogdi, id)
		if err != nil {
			fmt.Print(err.Error())
		}

		// Fastest way to append strings
		buffer.WriteString(jenjang)
		buffer.WriteString(" ")
		buffer.WriteString(nmprogdi)
		defer stmt.Close()
		datane := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"Pesane": fmt.Sprintf("Berhasil Merubah Menjadi %s", datane),
		})
	})


	// Delete resources
	router.DELETE("/", func(c *gin.Context) {
		id := c.PostForm("id")
		stmt, err := db.Prepare("delete from progdi where id= ?;")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(id)
		if err != nil {
			fmt.Print(err.Error())
		}
		c.JSON(http.StatusOK, gin.H{
			"Pesane": fmt.Sprintf("Berhasil Menghapus %s", id),
		})
	})
	router.Run(":8080")
}