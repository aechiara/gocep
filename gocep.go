package gocep

import (
	"encoding/json"
	"errors"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"unicode"
)

// CEP representa os dados do Logradouro
type CEP struct {
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	Cep        string `json:"cep"`
	UF         string `json:"uf"`
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
	cepURL = "https://buscacepinter.correios.com.br/app/consulta/html/consulta-detalhes-cep.php"
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
	v.Set("cep", cep)
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
	// header simulating mozilla firefox
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:98.0) Gecko/20100101 Firefox/98.0")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Erro fechando o reader do Body")
		}
	}(resp.Body)

	//log.Println("response Status:", resp.Status)
	//log.Println("response Headers:", resp.Header)

	///* convert text to ISO */
	reader := transform.NewReader(resp.Body, charmap.ISO8859_1.NewDecoder())

	b, err := io.ReadAll(reader)

	if err != nil {
		return nil, errors.New("Erro durante o Decode da resposta: " + err.Error())
	}

	// bind to JSON
	jsonResponse := new(cepResponse)
	err = json.Unmarshal(b, jsonResponse)
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
		UF:         jsonResponse.Dados[0].Uf,
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
