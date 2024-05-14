package main

/*
ширина таблицы 6
высота 161
1:03
*/

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Data struct {
	Name         string
	Average      string
	Median       string
	Inflation    string
	Costofliving string
	Healthcare   string
}
type Test struct {
	Info []Data
}

var data = parse()

func homepage(w http.ResponseWriter, r *http.Request) {

	tmpl, _ := template.ParseFiles("front/index.html")
	q := Test{data}
	tmpl.Execute(w, q)
}

func handleRequest() {
	http.HandleFunc("/", homepage)
	http.ListenAndServe(":2000", nil)
}

func main() {

	handleRequest()
}

// парсер
func parse() []Data {
	countries, median, average := parsesalary()
	countries2, inflation := inflation()
	countries3, costofliving := costofliving()
	countries4, healthcare := healthcare()

	// вся информация будет содержаться в этом срезе
	data := make([]Data, 0)

	for i := range countries {
		data = append(data, Data{Name: countries[i], Average: average[i],
			Median: median[i]})
	}

	// добавляем инфляцию
	for i := range countries2 {
		for j := range countries {
			if countries2[i] == countries[j] {
				data[j].Inflation = inflation[i]
			}
		}
	}
	// добавляем costofliving
	for i := range countries3 {
		for j := range countries {
			if countries3[i] == countries[j] {
				data[j].Costofliving = costofliving[i]
			}
		}
	}

	// добавляем healthcare
	for i := range countries4 {
		for j := range countries {
			if countries4[i] == countries[j] {
				data[j].Healthcare = healthcare[i]
			}
		}
	}
	// заменяем отсутствующие данные на -

	// changing inflation
	for i := range data {
		if data[i].Inflation == "" {
			data[i].Inflation = "-"
		}
	}

	// changing costofliving
	for i := range data {
		if data[i].Costofliving == "" {
			data[i].Costofliving = "-"
		}
	}

	// changing healthcare
	for i := range data {
		if data[i].Healthcare == "" {
			data[i].Healthcare = "-"
		}
	}

	return data
}

// информация о зарплатах
func parsesalary() ([]string, []string, []string) {
	url := "https://worldpopulationreview.com/country-rankings/median-income-by-country"
	countries := make([]string, 0) // список стран
	medium := make([]string, 0)
	average := make([]string, 0)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	// зарплаты
	c := 0
	inner_counter := 0
	doc.Find("tr>td").Each(func(i int, s *goquery.Selection) {
		// пропускаем первые 40 значений, тк они включют лишнюю информацию
		c += 1

		if c > 40 {
			inner_counter += 1
			switch inner_counter {
			case 1:
				medium = append(medium, s.Text())
			case 2:
				average = append(average, s.Text())
			}
		}
		if c%4 == 0 && c > 40 {
			//fmt.Println(" ")
			inner_counter = 0
		}
	})

	// страны
	c = 0
	doc.Find("tr>th").Each(func(i int, s *goquery.Selection) {

		c += 1
		if c > 9 {
			// fmt.Println(s.Text())
			countries = append(countries, s.Text())
		}
	})

	return countries, medium, average
}

// уровень инфляции
func inflation() ([]string, []string) {
	url := "https://wisevoter.com/country-rankings/inflation-by-country/"
	countries3 := make([]string, 0)
	inflation := make([]string, 0)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(resp.Body)

	// добавляем страну в срез
	doc.Find("tbody>tr>td.shdb-on-page-table-body-Geo").Each(func(i int, s *goquery.Selection) {
		countries3 = append(countries3, s.Text())
	})

	// добавляем уровень инфляции
	doc.Find("tbody>tr>td.shdb-on-page-table-body-Data").Each(func(i int, s *goquery.Selection) {
		inflation = append(inflation, s.Text())
	})

	return countries3, inflation
}

// costofliving
func costofliving() ([]string, []string) {
	url := "https://wisevoter.com/country-rankings/cost-of-living-by-country/"
	countries4 := make([]string, 0)
	cost := make([]string, 0)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(resp.Body)

	// добавляем страну в срез
	doc.Find("tbody>tr>td.shdb-on-page-table-body-Geo").Each(func(i int, s *goquery.Selection) {
		countries4 = append(countries4, s.Text())
	})

	// добавляем уровень инфляции
	doc.Find("tbody>tr>td.shdb-on-page-table-body-Data").Each(func(i int, s *goquery.Selection) {
		cost = append(cost, s.Text())
	})

	return countries4, cost
}

// healthcare (стоимость ежегодного здравоохранения)
func healthcare() ([]string, []string) {
	url := "https://en.wikipedia.org/wiki/List_of_countries_by_total_health_expenditure_per_capita"
	countries5 := make([]string, 0)
	health := make([]string, 0)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(resp.Body)

	// добавляем страну в срез
	c := 0
	inner_counter := 0
	second_counter := -6
	doc.Find("table.wikitable>tbody>tr>td").Each(func(i int, s *goquery.Selection) {

		c += 1
		inner_counter++
		second_counter++
		if c < 229 {

			if (inner_counter-1)%6 == 0 {
				test := strings.Split(s.Text(), "*")[0]
				b := string(test[2:])
				countries5 = append(countries5, b)
			}

			if second_counter%6 == 0 {
				health = append(health, s.Text())

			}
		}

	})
	return countries5, health
}
