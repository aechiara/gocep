package gocep

import (
	"testing"
)

func TestCEPRunner(t *testing.T) {
	t.Run("busca cep", func(t *testing.T) {
		expectedCep := &CEP{
			Logradouro: "Avenida Paulista - até 610 - lado par",
			Bairro:     "Bela Vista",
			Localidade: "São Paulo",
			Cep:        "01310-000",
		}
		actualCep, _ := BuscaCep("01310000")
		if expectedCep.Logradouro != actualCep.Logradouro &&
			expectedCep.Localidade != actualCep.Localidade &&
			expectedCep.Cep != actualCep.Cep &&
			expectedCep.Bairro != actualCep.Bairro {
			t.Errorf("Error getting CEP, got: %v, want: %v.", actualCep, expectedCep)
		}
	})

	t.Run("test cep as json", func(t *testing.T) {
		actualCep, err := BuscaCep("01310000")
		if err != nil {
			t.Errorf("Got an Error trying to Marshal CEP to JSON: %v", err.Error())
		}

		_, err = actualCep.ToJSON()

		if err != nil {
			t.Errorf("Got an Error trying to Marshal CEP to JSON: %v", err.Error())
		}
	})

	t.Run("cep nao exisite", func(t *testing.T) {
		actualCep, err := BuscaCep("00000000")

		if err == nil {
			t.Errorf("Should return nil, got: %v", actualCep)
		}
	})

	t.Run("cep deve ter apenas digitos", func(t *testing.T) {
		_, err := BuscaCep("123asd45")
		if err == nil {
			t.Errorf("Shoud return error, got nil")
		}
	})

	t.Run("cep deve ter tamanho 8", func(t *testing.T) {
		_, err := BuscaCep("234")
		if err == nil {
			t.Errorf("Shoud return error, got nil")
		}
	})
}