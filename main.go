package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
)

type Email struct {
	Address     string `json:"address"`
	Username    string `json:"username"`
	Domain      string `json:"domain"`
	HostExists  bool   `json:"hostExists"`
	Deliverable bool   `json:"deliverable"`
	FullInbox   bool   `json:"fullInbox"`
	CatchAll    bool   `json:"catchAll"`
	Disposable  bool   `json:"disposable"`
	Gravatar    bool   `json:"gravatar"`
}

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

	excelFileName := "emails.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Printf(err.Error())
	}
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			for _, cell := range row.Cells {
				text := cell.String()
				url := "https://trumail.io/json/" + strings.TrimSpace(strings.Replace(text, " ", "", -1))
				response, err := http.Get(url)
				if err != nil {
					fmt.Printf("%s", err)
					os.Exit(1)
				} else {
					defer response.Body.Close()
					contents, err := ioutil.ReadAll(response.Body)
					if err != nil {
						fmt.Printf("%s", err)
						os.Exit(1)
					}
					var email Email
					bodyString := string(contents)
					fmt.Println(bodyString)
					json.Unmarshal(contents, &email)
					if !email.Deliverable {
						fmt.Println(email.Address)
					}
				}
			}
		}
	}
}
