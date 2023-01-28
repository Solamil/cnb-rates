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
var pathHolytrinity string = ratesDir+"/svata_trojice.txt"
var coinCode string = "czk"

type userBaseResponse struct {
	Code []string `json:"code"`
	Volume []string `json:"volume"`
	Value []string `json:"value"`
	CoinCode string `json:"coin_code"`
	Number int `json:"number"`
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
	http.HandleFunc("/svata_trojice.txt", holytrinity_handler)
	http.HandleFunc("/svata_trojice", holytrinity_handler)
	http.HandleFunc("/holy_trinity", holytrinity_handler)
	http.HandleFunc("/denni_kurz.txt", dailyrates_handler)
	http.HandleFunc("/", index_handler)

	http.ListenAndServe(":8902", nil)
}

func index_handler(w http.ResponseWriter, r *http.Request) {

	q, _ := url.PathUnescape(r.URL.RawQuery)
	if len(q) != 0 {
		var code string = ""
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
			var amount float64 = 0.0
			if param.Amount != nil && len(param.Amount[0]) > 0 {
				amount, _ = strconv.ParseFloat(param.Amount[0], 64)
				value = value * (amount/amountBase)
				answer = fmt.Sprintf("%.3f", value)
			} else {
				answer = fmt.Sprintf("%.3f\n%.f", value, amountBase)
			}
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

func holytrinity_handler(w http.ResponseWriter, r *http.Request) {
	var result string = ""
	output := readFile(pathHolytrinity)
	q, _ := url.PathUnescape(r.URL.RawQuery)
	switch q {
		case "p": // p as pretty
			result = fmt.Sprintf("💵%s",output[0])
		default:
			http.ServeFile(w, r, pathHolytrinity)
			return
	}
	fmt.Println(q)
	w.Write([]byte(result))
}

func json_handler(w http.ResponseWriter, r *http.Request) {
	var baseResp = userBaseResponse{}

//	q, _ := url.PathUnescape(r.URL.RawQuery)
//	if len(q) != 0 {
//
//	}
	list := readFile(pathList) 
	for i := 1; i < len(list); i++ {
		infos := getRate(list[i])
		baseResp.Code = append(baseResp.Code, list[i])
		baseResp.Value = append(baseResp.Value, infos[0])
		baseResp.Volume = append(baseResp.Volume, infos[1])
	}

	baseResp.CoinCode = coinCode 

	baseResp.Date = readFile(pathDate)[0]
	baseResp.Number, _ = strconv.Atoi(readFile(pathNumber)[0])

	raw, err := json.Marshal(&baseResp)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(raw)
}

func isExistCode(code string) bool {
	for _, value := range readFile(pathList) {
		if code == value {
			return true
		}
	} 
	return false
}

func getRate(currency string) []string {
	curr := strings.ToUpper(currency)
	pathFile := fmt.Sprintf("%s/%s.txt", ratesDir, curr)

	output := readFile(pathFile) // 0 - value, 1 - volume
	return output
}

func readFile(pathFile string) []string {
	var output []string
	file, err := os.Open(pathFile)
	if err != nil {
		e := fmt.Sprintf("Failed to open %s", pathFile)
		return []string{e}
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		output = append(output, scanner.Text())
	}
	file.Close()
	return output
}

func runShellScript(name string) {
	_, err := exec.Command("/bin/sh", name).Output()
	if err != nil {
		fmt.Printf("error %s", err)
	}
}
