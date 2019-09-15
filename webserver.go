// webserver.go

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type Value struct {
	V int `json:"value"`
}

func main() {
	http.HandleFunc("/api/v1/figures", Sum)
	http.ListenAndServe("localhost:3000", nil)
}

func Sum(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		MethodPost(w, r)
	} else if r.Method == "GET" {
		MethodGet(w, r)
	} else {
		fmt.Fprintln(w, "{ error: несоответствие метода }")
	}
}

func File(w http.ResponseWriter, value int) { //запись в файл, суммирование

	file, err := os.Open("sum.log")
	if err != nil {

		file, err = os.Create("sum.log")
		if err != nil {
			fmt.Fprintln(w, "{ error: ошибка создания файла }")
			return
		}
		_, err = file.WriteString(strconv.Itoa(value))
		if err != nil {
			fmt.Fprintln(w, "{ error: ошибка при записи в файл }")
			return
		}

		fmt.Fprintln(w, "{ status: OK }")

		file.Close()
		return
	}

	reader := bufio.NewReader(file)
	line, _ := reader.ReadString('\n')
	old_value, err := strconv.Atoi(line)

	if err != nil {
		fmt.Fprintln(w, "{ error : ошибка при работе с данными из файла }")
		return
	}

	res := old_value + value

	file.Close()

	file, err = os.Create("sum.log")
	if err != nil {
		fmt.Fprintln(w, "{ error: ошибка при записи в файл }")
		return
	}

	_, err = file.WriteString(strconv.Itoa(res)) // запись строки

	if err != nil {
		fmt.Fprintln(w, "{ error: ошибка при записи в файл }")
		return
	}
	fmt.Fprintln(w, "{ status: OK }")
	file.Close()

}

func MethodPost(w http.ResponseWriter, r *http.Request) {

	var val Value
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintln(w, "{ error: ошибка при работе с json}")
		return
	} else {
		err = json.Unmarshal(body, &val)
		if err != nil {
			//http.Error(w, http.StatusText(400), 400)
			fmt.Fprintln(w, "{ error: ошибка при работе с json}")
			return
		}
	}
	File(w, val.V) //обрабатываем запись в файл

}
func MethodGet(w http.ResponseWriter, r *http.Request) {

	br, err := ioutil.ReadFile("sum.log")

	if err != nil {
		fmt.Fprintln(w, "{ error: ошибка при чтении из файла }")
		return
	}
	line := string(br)
	_, err = strconv.Atoi(line)
	if err != nil {
		fmt.Fprintln(w, "{ error: ошибка при чтении из файла }")
		return
	}
	fmt.Fprintln(w, "{ value: ", line, " }")

}
