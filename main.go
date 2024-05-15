package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Quiz struct {
	Data struct {
		Quiz struct {
			Info struct {
				Questions []struct {
					Structure struct {
						Query struct {
							Text string `json:"text"`
						} `json:"query"`
						Answer  int `json:"answer"`
						Options []struct {
							Text string `json:"text"`
						} `json:"options"`
					} `json:"structure"`
				} `json:"questions"`
			} `json:"info"`
		} `json:"quiz"`
	} `json:"data"`
}

func removeHTMLTags(text string) string {
	text = strings.ReplaceAll(text, "<br>", "\n")
	re := regexp.MustCompile(`<[^>]*>`)
	text = re.ReplaceAllString(text, "")
	text = strings.ReplaceAll(text, "&nbsp;", " ")
	return text
}

func getAnswer(id string) {
	resp, err := http.Get("https://quizizz.com/quiz/" + id)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var data Quiz
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	wow := ""
	for i, question := range data.Data.Quiz.Info.Questions {
		nomor := i + 1
		questionText := removeHTMLTags(question.Structure.Query.Text)
		answer := removeHTMLTags(question.Structure.Options[question.Structure.Answer].Text)
		wow += fmt.Sprintf("%d. %s\nAnswer: %s\n\n", nomor, questionText, answer)
	}
	fmt.Println(wow)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide the summary id of the Quizizz quiz.")
		return
	}
	getAnswer(os.Args[1])
}
