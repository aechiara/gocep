package gocep

import (
	//"fmt"

	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/bjarneh/latinx"
)

// CEP representa os dados do Logradouro
type CEP struct {
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	Cep        string `json:"cep"`
}

// ToJSON return CEP struct as Json
func (c *CEP) ToJSON() (string, error) {
	jsonRet, err := json.Marshal(c)
	return string(jsonRet), err
}

// BuscaCep consulta o CEP informado no site dos correios
func BuscaCep(cep string) (CEP, error) {

	if len(cep) != 8 {
		return CEP{}, errors.New("O CEP DEVE ter 8 digitos")
	}

	_, errorAtoi := strconv.Atoi(cep)
	if errorAtoi != nil {
		return CEP{}, errors.New("O CEP DEVE ser apenas digitos")
	}

	const cepURL = "http://www.buscacep.correios.com.br/sistemas/buscacep/resultadoBuscaCepEndereco.cfm"

	v := url.Values{}
	v.Set("relaxation", cep)
	v.Set("tipoCEP", "ALL")
	v.Set("semelhante", "N")

	s := v.Encode()
	//fmt.Println("Posting data: " + s)

	req, _ := http.NewRequest("POST", cepURL, strings.NewReader(s))
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
	c, err := converter.Decode(b)
	body := string(c)

	//fmt.Println(body)

	/* capture the part of HTML with the data */
	re := regexp.MustCompile("(?s)(?m)<table class=\"tmptabela\">(.*?)</table>")
	output := re.FindString(body)

	if len(output) == 0 {
		return CEP{}, errors.New("DADOS NAO ENCONTRADOS")
	}

	/* strip some special chars */
	reg, err := regexp.Compile("&nbsp;|\\t|\\r")
	cleanString := reg.ReplaceAllString(output, "")

	fieldValues := getFieldsValue(cleanString)

	cepRet := CEP{
		Logradouro: fieldValues[0],
		Bairro:     fieldValues[1],
		Localidade: fieldValues[2],
		Cep:        fieldValues[3]}

	// jsonRet, err := json.Marshal(cepRet)
	//fmt.Println(string(json_ret))

	// return string(jsonRet), err
	return cepRet, nil
}

/* grab field Names */
func getFieldsName(s string) []string {

	retorno := make([]string, 0)

	results := regexp.MustCompile(`<th.*?>(.*?):</th>`).FindAllStringSubmatch(s, -1)
	//for i, match := range results {
	for _, match := range results {
		//full := match[0]
		submatches := match[1:len(match)]
		//fmt.Printf("%v => \"%v\" from \"%v\"\n", i, submatches[0], full)
		retorno = append(retorno, submatches[0])
	}

	return retorno
}

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
