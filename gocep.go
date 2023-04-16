package gocep

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"unicode"

	"github.com/bjarneh/latinx"
)

// CEP representa os dados do Logradouro
type CEP struct {
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	Cep        string `json:"cep"`
}

type cepResponse struct {
	Erro     bool   `json:"erro"`
	Mensagem string `json:"mensagem"`
	Total    int    `json:"total"`
	Dados    []struct {
		Uf                       string        `json:"uf"`
		Localidade               string        `json:"localidade"`
		LocNoSem                 string        `json:"locNoSem"`
		LocNu                    string        `json:"locNu"`
		LocalidadeSubordinada    string        `json:"localidadeSubordinada"`
		LogradouroDNEC           string        `json:"logradouroDNEC"`
		LogradouroTextoAdicional string        `json:"logradouroTextoAdicional"`
		LogradouroTexto          string        `json:"logradouroTexto"`
		Bairro                   string        `json:"bairro"`
		BaiNu                    string        `json:"baiNu"`
		NomeUnidade              string        `json:"nomeUnidade"`
		Cep                      string        `json:"cep"`
		TipoCep                  string        `json:"tipoCep"`
		NumeroLocalidade         string        `json:"numeroLocalidade"`
		Situacao                 string        `json:"situacao"`
		FaixasCaixaPostal        []interface{} `json:"faixasCaixaPostal"`
		FaixasCep                []interface{} `json:"faixasCep"`
	} `json:"dados"`
}

const (
	cepURL = "https://buscacepinter.correios.com.br/app/endereco/carrega-cep-endereco.php"
)

// ToJSON return CEP struct as Json
func (c *CEP) ToJSON() (string, error) {
	jsonRet, err := json.Marshal(c)
	return string(jsonRet), err
}

// Buscar consulta o CEP informado no site dos correios
func Buscar(cep string) (*CEP, error) {

	if len(cep) != 8 {
		return nil, errors.New("O CEP DEVE ter 8 digitos")
	}

	if !isStringOnlyDigits(cep) {
		return nil, errors.New("O CEP DEVE conter apenas dígitos")
	}

	v := url.Values{}
	v.Set("endereco", cep)
	v.Set("tipoCEP", "ALL")
	v.Set("cepaux", "")
	v.Set("mensagem_alerta", "")
	v.Set("pagina", "/app/endereco/index.php")

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

	//log.Println("response Status:", resp.Status)
	//log.Println("response Headers:", resp.Header)

	//b, _ := ioutil.ReadAll(resp.Body)
	b, _ := io.ReadAll(resp.Body)

	/* convert text to ISO */
	converter := latinx.Get(latinx.ISO_8859_1)
	c, err := converter.Decode(b)
	if err != nil {
		c = b
		return nil, errors.New("Erro durante o Decode da resposta: " + err.Error())
	}
	//body := string(c)

	//log.Println(body)

	// bind to JSON
	jsonResponse := new(cepResponse)
	err = json.Unmarshal(c, jsonResponse)
	if err != nil {
		return nil, errors.New("Erro unmarshalling response: " + err.Error())
	}

	if jsonResponse.Total == 0 {
		return nil, errors.New("CEP não encontrado")
	}

	r := CEP{
		Logradouro: jsonResponse.Dados[0].LogradouroDNEC,
		Cep:        cep,
		Bairro:     jsonResponse.Dados[0].Bairro,
		Localidade: jsonResponse.Dados[0].Localidade,
	}
	return &r, nil
}

func isStringOnlyDigits(str string) bool {
	for _, r := range str {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
