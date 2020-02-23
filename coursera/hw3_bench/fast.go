package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// более быстрая и оптимальная версия функции SlowSearch
// заменил regexp.MatchString на более быструю strings.Contains
// и r.ReplaceAllString на strings.ReplaceAll
// следующее слабое место - bytes makeSlice (22% памяти)
// добавил структуру для хранения данных пользователя,
// чтобы не накапливать эти данные в слайсе + получилось убрать
// ненужные преобразования типов из-за этого
// Перебор браузеров в одном цикле вместо двух
// убрал накапливание найденных пользователей:
// foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, email)

type User struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Browsers []string `json:"browsers"`
}

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	seenBrowsers := []string{}
	uniqueBrowsers := 0

	lines := strings.Split(string(fileContents), "\n")

	fmt.Fprintln(out, "found users:")

	user := User{}
	for i, line := range lines {
		err := json.Unmarshal([]byte(line), &user)
		if err != nil {
			panic(err)
		}
		isAndroid := false
		isMSIE := false
		for _, browser := range user.Browsers {
			if ok := strings.Contains(browser, "Android"); ok {
				isAndroid = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}
			if ok := strings.Contains(browser, "MSIE"); ok {
				isMSIE = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}
		}
		if !(isAndroid && isMSIE) {
			continue
		}
		email := strings.ReplaceAll(user.Email, "@", " [at] ")
		fmt.Fprintf(out, "[%d] %s <%s>\n", i, user.Name, email)
	}
	fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))
}
