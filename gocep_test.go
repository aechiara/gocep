package gocep_test

import (
	"github.com/aechiara/gocep"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func mockServer(responseBody string, statusCode int) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		encoder := charmap.ISO8859_1.NewEncoder()
		encodedBody, _, _ := transform.String(encoder, responseBody)
		_, _ = w.Write([]byte(encodedBody))
	})
	return httptest.NewServer(handler)
}

func TestCEPRunner(t *testing.T) {
	t.Run("cep deve ter apenas digitos", func(t *testing.T) {
		_, err := gocep.Buscar("123asd45")
		if err == nil {
			t.Errorf("Deveria retornar um erro, mas retornou nulo")
		}
		expectedError := "O CEP DEVE conter apenas dígitos"
		if err.Error() != expectedError {
			t.Errorf("Esperava a mensagem de erro '%s', mas obteve '%s'", expectedError, err.Error())
		}
	})

	t.Run("cep deve ter tamanho 8", func(t *testing.T) {
		_, err := gocep.Buscar("234")
		if err == nil {
			t.Errorf("Deveria retornar um erro, mas retornou nulo")
		}
		expectedError := "O CEP DEVE ter 8 digitos"
		if err.Error() != expectedError {
			t.Errorf("Esperava a mensagem de erro '%s', mas obteve '%s'", expectedError, err.Error())
		}
	})

	t.Run("cep nao existe", func(t *testing.T) {
		server := mockServer(`{"erro": false, "total": 0}`, http.StatusOK)
		defer server.Close()

		originalURL := gocep.GetCepURL()
		gocep.SetCepURL(server.URL)
		defer gocep.SetCepURL(originalURL)

		_, err := gocep.Buscar("00000000")
		if err == nil {
			t.Errorf("Deveria retornar um erro, mas retornou nulo")
		}

		expectedError := "CEP não encontrado"
		if err.Error() != expectedError {
			t.Errorf("Esperava a mensagem de erro '%s', mas obteve '%s'", expectedError, err.Error())
		}
	})

	t.Run("busca cep", func(t *testing.T) {
		responseBody := `{
			"erro": false,
			"mensagem": "Dados encontrados com sucesso.",
			"total": 1,
			"dados": [
				{
					"uf": "SP",
					"localidade": "São Paulo",
					"logradouroDNEC": "Avenida Paulista - até 610 - lado par",
					"bairro": "Bela Vista",
					"cep": "01310000"
				}
			]
		}`
		server := mockServer(responseBody, http.StatusOK)
		defer server.Close()

		originalURL := gocep.GetCepURL()
		gocep.SetCepURL(server.URL)
		defer gocep.SetCepURL(originalURL)

		expectedCep := &gocep.CEP{
			Logradouro: "Avenida Paulista - até 610 - lado par",
			Bairro:     "Bela Vista",
			Localidade: "São Paulo",
			Cep:        "01310000",
			UF:         "SP",
		}
		actualCep, err := gocep.Buscar("01310000")
		if err != nil {
			t.Fatalf("Não esperava erro, mas obteve %v", err)
		}

		if !reflect.DeepEqual(actualCep, expectedCep) {
			t.Errorf("Erro ao buscar CEP, obteve: %v, esperava: %v.", actualCep, expectedCep)
		}
	})

	t.Run("test cep as json", func(t *testing.T) {
		responseBody := `{
			"erro": false,
			"mensagem": "Dados encontrados com sucesso.",
			"total": 1,
			"dados": [
				{
					"uf": "SP",
					"localidade": "São Paulo",
					"logradouroDNEC": "Avenida Paulista - até 610 - lado par",
					"bairro": "Bela Vista",
					"cep": "01310000"
				}
			]
		}`
		server := mockServer(responseBody, http.StatusOK)
		defer server.Close()

		originalURL := gocep.GetCepURL()
		gocep.SetCepURL(server.URL)
		defer gocep.SetCepURL(originalURL)

		actualCep, err := gocep.Buscar("01310000")
		if err != nil {
			t.Fatalf("Ocorreu um erro ao tentar buscar o CEP: %v", err.Error())
		}

		jsonString, err := actualCep.ToJSON()
		if err != nil {
			t.Errorf("Ocorreu um erro ao tentar converter o CEP para JSON: %v", err.Error())
		}

		expectedJSON := `{"logradouro":"Avenida Paulista - até 610 - lado par","bairro":"Bela Vista","localidade":"São Paulo","cep":"01310000","uf":"SP"}`
		if jsonString != expectedJSON {
			t.Errorf("JSON esperado '%s', obteve '%s'", expectedJSON, jsonString)
		}
	})

	t.Run("http request fails", func(t *testing.T) {
		server := mockServer("", http.StatusInternalServerError)
		server.Close() // Close server to simulate a connection error

		originalURL := gocep.GetCepURL()
		gocep.SetCepURL(server.URL)
		defer gocep.SetCepURL(originalURL)

		_, err := gocep.Buscar("01310000")
		if err == nil {
			t.Fatalf("Esperava um erro, mas retornou nulo")
		}
	})

	t.Run("malformed json response", func(t *testing.T) {
		server := mockServer(`{"erro": false, "total": 1, "dados": [{}`, http.StatusOK)
		defer server.Close()

		originalURL := gocep.GetCepURL()
		gocep.SetCepURL(server.URL)
		defer gocep.SetCepURL(originalURL)

		_, err := gocep.Buscar("01310000")
		if err == nil {
			t.Fatalf("Esperava um erro, mas retornou nulo")
		}
		expectedError := "Erro ao decodificar a resposta:"
		if !strings.HasPrefix(err.Error(), expectedError) {
			t.Errorf("Esperava que o erro começasse com '%s', mas obteve '%s'", expectedError, err.Error())
		}
	})
}
