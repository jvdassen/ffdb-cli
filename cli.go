package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

type macroNutrientProfile struct {
	protein       string
	carbohydrates string
	fats          string
}

type hyperLink struct {
	text string
	url  string
}

func searchForProducts(searchTerm string) []hyperLink {
	res, err := http.Get("https://fddb.info/db/de/suche/?udd=0&cat=site-de&search=" + searchTerm)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s\n", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	fmt.Printf("Retrieved")
	if err != nil {
		log.Fatal(err)
	}

	hyperLinks := parseSearchResults(doc)

	return hyperLinks
}

func getNutritionInfoForProduct(product string) macroNutrientProfile {
	res, err := http.Get("https://fddb.info/db/de/lebensmittel/" + product + "/index.html")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s\n", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	nutritionInfo := parseNutritionInfo(doc)

	return nutritionInfo
}

func parseSearchResults(doc *goquery.Document) []hyperLink {
	hyperLinks := []hyperLink{}
	linkCSSSelector := ".standardcontent > table > tbody > tr > td > div > a"

	elements := doc.Find(linkCSSSelector)
	elements.Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			item := hyperLink{
				text: s.Text(),
				url:  href,
			}
			hyperLinks = append(hyperLinks, item)
		}
	})

	return hyperLinks
}

func parseNutritionInfo(doc *goquery.Document) macroNutrientProfile {
	nutritionInfo := macroNutrientProfile{
		protein:       "0g",
		carbohydrates: "0g",
		fats:          "0g",
	}

	doc.Find(".standardcontent > div:nth-child(2) > div").Each(func(i int, s *goquery.Selection) {
		macros := s.Find("div")

		var macro = "none"
		var value = "none"
		macros.Each(func(i int, s *goquery.Selection) {
			text := s.Text()
			if i == 0 {
				macro = text
			} else {
				value = text
				if macro == "Protein" {
					nutritionInfo.protein = value
				}
				if macro == "Kohlenhydrate" {
					nutritionInfo.carbohydrates = value
				}
				if macro == "Fett" {
					nutritionInfo.fats = value
				}
			}
		})
	})
	return nutritionInfo
}

func main() {
	res := searchForProducts("Kiwi")
	fmt.Printf("%+v\n", res)
	nutri1 := getNutritionInfoForProduct("naturprodukt_apfel_frisch")
	fmt.Printf("Getting nutrition info for apple: %+v\n", nutri1)
	nutri2 := getNutritionInfoForProduct("imbiss_doener_kebap")
	fmt.Printf("Getting nutrition info for doner: %+v", nutri2)
}
