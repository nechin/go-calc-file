package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что метод запроса - POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Парсим multipart форму
	err := r.ParseMultipartForm(10 << 20) // Лимит на размер файла - 10 МБ
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Получаем файл из формы
	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
		http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fmt.Printf("Received file: %s", header.Filename)

	// Возвращаем указатель временного файла в начало
	_, err = file.Seek(0, 0)
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusInternalServerError)
		return
	}

	// Считываем содержимое файла и суммируем числа
	sum, err := sumNumbersInFile(file)
	if err != nil {
		http.Error(w, "Error processing file", http.StatusInternalServerError)
		return
	}

	// Отправляем сумму в ответ
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Sum of numbers: %d", sum)
	fmt.Fprintln(w, "")
}

func sumNumbersInFile(file io.Reader) (int, error) {
	var sum int
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		// Разделяем строку на слова (числа)
		numbers := strings.Fields(line)
		for _, numStr := range numbers {
			num, err := strconv.Atoi(numStr)
			if err != nil {
				// Если не удалось преобразовать в число, пропускаем
				continue
			}
			sum += num
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return sum, nil
}

func main() {
	// Регистрируем обработчик для POST-запросов
	http.HandleFunc("/upload", uploadHandler)

	// Запускаем сервер на порту 8080
	fmt.Println("Server started at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

// curl -X POST -F "file=@file.txt" http://localhost:8080/upload
