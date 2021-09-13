package gocep

import (
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
func BuscaCep(cep string) (*CEP, error) {

	if len(cep) != 8 {
		return nil, errors.New("O CEP DEVE ter 8 digitos")
	}

	_, err := strconv.Atoi(cep)
	if err != nil {
		return nil, errors.New("O CEP DEVE conter apenas dígitos")
	}

	const cepURL = "http://www.buscacep.correios.com.br/sistemas/buscacep/resultadoBuscaCepEndereco.cfm"

	v := url.Values{}
	v.Set("relaxation", cep)
	v.Set("tipoCEP", "ALL")
	v.Set("semelhante", "N")

	s := v.Encode()
	//log.Println("Posting data: " + s)

	req, err := http.NewRequest("POST", cepURL, strings.NewReader(s))
	if err != nil {
		return nil, errors.New("Erro criando Requisição: " + err.Error())
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)

	b, _ := ioutil.ReadAll(resp.Body)

	/* convert text to ISO */
	converter := latinx.Get(latinx.ISO_8859_1)
	c, err := converter.Decode(b)
	if err != nil {
		return nil, errors.New("Erro durante o Decode da resposta: " + err.Error())
	}
	body := string(c)

	//log.Println(body)

	/* capture the part of HTML with the data */
	re := regexp.MustCompile("(?s)(?m)<table class=\"tmptabela\">(.*?)</table>")
	output := re.FindString(body)

	if len(output) == 0 {
		return nil, errors.New("DADOS NAO ENCONTRADOS")
	}

	/* strip some special chars */
	reg, _ := regexp.Compile("&nbsp;|\\t|\\r")
	cleanString := reg.ReplaceAllString(output, "")

	fieldValues := getFieldsValue(cleanString)

	cepRet := CEP{
		Logradouro: fieldValues[0],
		Bairro:     fieldValues[1],
		Localidade: fieldValues[2],
		Cep:        fieldValues[3]}

	return &cepRet, nil
}

/* grab field Names */
func getFieldsName(s string) []string {

	retorno := make([]string, 0)

	results := regexp.MustCompile(`<th.*?>(.*?):</th>`).FindAllStringSubmatch(s, -1)
	for idx, match := range results {
		subMatches := match[1:len(match)]
		log.Printf("idx %v => \"%v\" from \"%v\"\n", idx, subMatches[0], match[0])
		retorno = append(retorno, subMatches[0])
	}

	return retorno
}

/* private function to deal with the string and grab the data */
func getFieldsValue(s string) []string {

	retorno := make([]string, 0)

	results := regexp.MustCompile(`<td.*?>(.*?)</td>`).FindAllStringSubmatch(s, -1)
	for idx, match := range results {
		subMatches := match[1:len(match)]
		log.Printf("idx %v => \"%v\" from \"%v\"\n", idx, subMatches[0], match[0])
		retorno = append(retorno, subMatches[0])
	}

	return retorno
}
