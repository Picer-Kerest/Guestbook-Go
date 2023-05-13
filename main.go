package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
)

type GuestBook struct {
	SignatureCount int
	Signatures     []string
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func GetStrings(fileName string) ([]string, error) {
	var lines []string

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}
	err = file.Close()
	if err != nil {
		return nil, err
	}

	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return lines, nil
}

func viewHandler(writer http.ResponseWriter, request *http.Request) {
	signature, err := GetStrings("signatures.txt")
	check(err)
	html, err := template.ParseFiles("main.html")
	check(err)

	guestbook := GuestBook{
		SignatureCount: len(signature),
		Signatures:     signature,
	}
	err = html.Execute(writer, guestbook)
	// Выполняем шаблон HTML, передавая данные из guestbook в writer
	check(err)
}

func addHandler(writer http.ResponseWriter, request *http.Request) {
	html, err := template.ParseFiles("add.html")
	check(err)
	err = html.Execute(writer, nil)
	check(err)
}

func createHandler(writer http.ResponseWriter, request *http.Request) {
	signature := request.FormValue("signature")
	// Читаем значение формы с именем signature

	options := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	file, err := os.OpenFile("signatures.txt", options, os.FileMode(0600))
	check(err)
	//signature = "\n" + signature
	_, err = fmt.Fprintln(file, signature)
	check(err)
	err = file.Close()
	check(err)
	//_, err := writer.Write([]byte(signature))
	// Записываем значение в http ответ
	// [] byte, потому что информация, которую мы передаём, хранится в byte
	http.Redirect(writer, request, "/guestbook", http.StatusFound)
}

func main() {
	http.HandleFunc("/guestbook", viewHandler)
	http.HandleFunc("/guestbook/add", addHandler)
	http.HandleFunc("/guestbook/create", createHandler)
	// Обработчик по адресу
	err := http.ListenAndServe("localhost:8080", nil)
	// Запуск сервера. Запуск бесконечный
	log.Fatal(err)
	// Попадаем сюда если есть ошибка
}
