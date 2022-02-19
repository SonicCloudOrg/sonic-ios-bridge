package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

func main() {
	var url = "https://www.theiphonewiki.com/wiki/Models"
	p, _ := goquery.NewDocument(url)

	pTitle := p.Find(".wikitable").Eq(0) //直接提取title的内容

	th := pTitle.Find("th")
	var gIndex = 0
	var Iindex = 0
	th.Map(func(i int, s *goquery.Selection) string {
		if s.Text() == "Generation\n" {
			gIndex = i
		}
		if s.Text() == "Identifier\n" {
			Iindex = i
		}
		return s.Text()
	})
	tr:=pTitle.Find("td")
	fmt.Println(tr.Eq(gIndex).Text())
	fmt.Println(tr.Eq(Iindex).Text())
	//	if t.Data=="Generation" {
	//		fmt.Println(i)
	//	}
	//	if t.Data=="Identifier"{
	//		fmt.Println(i)
	//	}
	//}
	//o:=i.Find("td").Eq(0).Text()
	//fmt.Println(o)

	//cmd.Execute()
}
