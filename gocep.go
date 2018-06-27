package gocep

import (
	//"fmt"
    "log"
	"net/http"
	"net/url"
	"io/ioutil"
	"strings"
	"regexp"
	"encoding/json"
	"github.com/bjarneh/latinx"
)

type CEP struct {
	Logradouro	string `json:"logradouro"`
	Bairro		string `json:"bairro"`
	Localidade	string `json:"localidade"`
	Cep			string `json:"cep"`
}


func BuscaCep(cep string) string {

	cepUrl := "http://www.buscacep.correios.com.br/sistemas/buscacep/resultadoBuscaCepEndereco.cfm"

	v := url.Values{}
	v.Set("relaxation", cep)
	v.Set("tipoCEP", "ALL")
	v.Set("semelhante", "N")

	s := v.Encode()
	//fmt.Println("Posting data: " + s)

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

	/* convert text to ISO */
	converter := latinx.Get(latinx.ISO_8859_1)
	c,err := converter.Decode(b)
	body := string(c)

	//fmt.Println(body)

	/* capture the part of HTML with the data */
	re := regexp.MustCompile("(?s)(?m)<table class=\"tmptabela\">(.*?)</table>")
	output := re.FindString(body)

	/* strip some special chars */
	reg, err := regexp.Compile("&nbsp;|\\t|\\r")
	cleanString := reg.ReplaceAllString(output, "")

	/* uncoment to grap the field names
	fieldNames := getFieldsName(cleanString)
	for _, item := range fieldNames {
		fmt.Println("nome: [", item, "]")
	}
	*/

	fieldValues := getFieldsValue(cleanString)

	cep_ret := CEP{fieldValues[0], fieldValues[1], fieldValues[2], fieldValues[3]}
	json_ret, _ := json.Marshal(cep_ret)
	//fmt.Println(string(json_ret))

	return string(json_ret)
}

/* uncoment to grab field Names
func getFieldsName(s string) [] string {

	retorno := make([]string, 0)

	results := regexp.MustCompile(`<th.*?>(.*?):</th>`).FindAllStringSubmatch(s, -1)
	for i, match := range results {
		full := match[0]
		submatches := match[1:len(match)]
		fmt.Printf("%v => \"%v\" from \"%v\"\n", i, submatches[0], full)
		retorno = append(retorno, submatches[0])
	}

	return retorno
}
*/

/* private function to deal with the string and grab the data */
func getFieldsValue(s string) []string {

	retorno := make([]string, 0)

	results := regexp.MustCompile(`<td.*?>(.*?)</td>`).FindAllStringSubmatch(s, -1)
	//for i, match := range results {
	for _, match := range results {
		//full := match[0]
		submatches := match[1:len(match)]
		//fmt.Printf("%v => \"%v\" from \"%v\"\n", i, submatches[0], full)
		retorno = append(retorno, submatches[0])
	}

	return retorno
}
