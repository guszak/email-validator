package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	trumail "github.com/sdwolfe32/trumail/verifier"
)

func main() {
	port := "3001"
	g := gin.Default()
	//g.Use(gin.Logger())
	//g.Use(gin.Recovery())
	g.GET("/email", readSheet)
	http.Handle("/", g)

	http.ListenAndServe(":"+port, nil)
}

// CreateProduct add product
func readSheet(c *gin.Context) {
	file, error := os.Open("emails.csv")
	if error != nil {
		fmt.Println("Error:", error)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	lineCount := 0
	for {
		record, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			fmt.Println("Error:", error)
			return
		}
		// record is an array of string so is directly printable
		//fmt.Println("Record", lineCount, "is", record, "and has", len(record), "fields")
		// and we can iterate on top of that
		v := trumail.NewVerifier(20, "gmail.com", "engenheiroecp@gmail.com")
		for i := 0; i < len(record); i++ {
			fmt.Println(" ", record[i])
			res := v.Verify(record[i])
			log.Println(*res[0])
		}
		fmt.Println()
		lineCount++
	}
}
