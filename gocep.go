package main

import "fmt"
import "log"
import "net/http"
import "net/url"
import "io/ioutil"
import "strings"
import "regexp"

type SitemapIndex struct {
	Locations [] string `xml:"sitemap>loc"`
}


func main() {
	cepUrl := "http://www.buscacep.correios.com.br/sistemas/buscacep/resultadoBuscaCepEndereco.cfm"
	var cep = "02242005"

	v := url.Values{}
	v.Set("relaxation", cep)
	v.Set("tipoCEP", "ALL")
	v.Set("semelhante", "N")

	s := v.Encode()
	fmt.Println("Posting data: " + s)

	req, _ := http.NewRequest("POST", cepUrl, strings.NewReader(s))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)

	b, _ := ioutil.ReadAll(resp.Body)
	body := string(b)

	//fmt.Println(body)

	re := regexp.MustCompile("(?s)(?m)<table class=\"tmptabela\">(.*?)</table>")
	output := re.FindString(body)
	fmt.Println("--------")
	fmt.Printf("[%q]\n", output)
	fmt.Println("--------")

	reg, err := regexp.Compile("[^a-zA-Z0-9<>= /]+")
	processedString := reg.ReplaceAllString(output, "")
	fmt.Printf("[%q]\n", processedString)


}
