package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"net/http"
	"net/url"
	"encoding/json"
	"strings"
	"strconv"
	"os"
)
var ratesDir string = "./rates"
var pathList string = ratesDir+"/list.txt"
var pathDate string = ratesDir+"/date.txt"
var pathNumber string = ratesDir+"/number.txt"
var pathDenniKurz string = "./denni_kurz.txt"
var coinCode string = "czk"

type userBaseResponse struct {
	Code []string `json:"code"`
	Volume []string `json:"volume"`
	Value []string `json:"value"`
	CoinCode string `json:"coin_code"`
	Date string `json:"date"`
}

type indexUrlParams struct {
	Code []string `json:"code"`
	Amount []string `json:"amount"`
}

func main() {
	runShellScript("rates.sh")

	http.HandleFunc("/index.html", index_handler)
	http.HandleFunc("/list", list_handler)
	http.HandleFunc("/date", date_handler)
	http.HandleFunc("/json", json_handler)
	http.HandleFunc("/number", number_handler)
	http.HandleFunc("/denni_kurz.txt", dailyrates_handler)
	http.HandleFunc("/", index_handler)

	http.ListenAndServe(":8902", nil)
}

func index_handler(w http.ResponseWriter, r *http.Request) {

	q, _ := url.PathUnescape(r.URL.RawQuery)
	if len(q) != 0 {
		var code string = ""
		var amount float64 = 1.0
		var answer string = ""
		m, err := url.ParseQuery(q)
		if err != nil {
			fmt.Println(err)
		}
		
		js, err := json.Marshal(m)
		if err != nil {
			fmt.Println(err)
		}
		var param *indexUrlParams
		json.Unmarshal(js, &param)

		if param.Code != nil && len(param.Code[0]) > 0 {
			code = strings.ToUpper(param.Code[0])
			
			if ! isExistCode(code) {
				if code == "LIST" {
					http.Redirect(w, r, "/list", http.StatusMovedPermanently)
				}
				if code == "DATE" {
					http.Redirect(w, r, "/date", http.StatusMovedPermanently)
				}
				return	
			}
			info := getRate(param.Code[0])
			value, _ := strconv.ParseFloat(info[0], 64)
			amountBase, _ := strconv.ParseFloat(info[1], 64)
			if param.Amount != nil && len(param.Amount[0]) > 0 {
				amount, _ = strconv.ParseFloat(param.Amount[0], 64)
			} else {
				amount = amountBase
			}
			value = value * (amount/amountBase)
			answer = fmt.Sprintf("%.3f\n%.f", value, amount)
		}
		w.Write([]byte(answer))
		return
	}

	http.ServeFile(w, r, "web/index.html")

}


func list_handler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, pathList)
}

func date_handler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, pathDate)
}

func dailyrates_handler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, pathDenniKurz)
}

func number_handler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, pathNumber)

}

func json_handler(w http.ResponseWriter, r *http.Request) {
	var baseResp = userBaseResponse{}

	q, _ := url.PathUnescape(r.URL.RawQuery)
	if len(q) != 0 {

	}
	for i, value := range getList() {
		if i == 0 {
			continue
		}
		infos := getRate(value)
		baseResp.Code = append(baseResp.Code, value)
		baseResp.Value = append(baseResp.Value, infos[0])
		baseResp.Volume = append(baseResp.Volume, infos[1])
	}	
	baseResp.CoinCode = coinCode 
	baseResp.Date = getDate()

	raw, err := json.Marshal(&baseResp)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(raw)
}

func isExistCode(code string) bool {
	for _, value := range getList() {
		if code == value {
			return true
		}
	} 
	return false
}

func getList() []string {
	var list []string
		
	file, err := os.Open(pathList)
	if err != nil {
		e := fmt.Sprintf("Failed to open %s", pathList)
		return []string{e}
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		list = append(list, scanner.Text())
	}
	file.Close()
	return list 
}

func getRate(currency string) []string {
	curr := strings.ToUpper(currency)
	pathFile := fmt.Sprintf("%s/%s.txt", ratesDir, curr)
	
	file, err := os.Open(pathFile)
	if err != nil {
		e := fmt.Sprintf("Failed to open %s", pathFile)
		return []string{e}
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string   // 0 - value, 1 - volume

	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	file.Close()
	
	return text
}

func getDate() string {

	file, err := os.Open(pathDate)
	if err != nil {
		e := fmt.Sprintf("Failed to open %s", pathDate)
		return e 
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text string

	scanner.Scan()
	text = scanner.Text()
	file.Close()
	
	return text

}

func runShellScript(name string) {
	_, err := exec.Command("/bin/sh", name).Output()
	if err != nil {
		fmt.Printf("error %s", err)
	}
}
