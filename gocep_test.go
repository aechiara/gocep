package gocep

import (
	"testing"
)

func TestBuscaCep(t *testing.T) {
	// expected := `{
	// 	"logradouro":"Avenida Paulista - até 610 - lado par",
	// 	"bairro":"Bela Vista",
	// 	"localidade":"São Paulo/SP",
	// 	"cep":"01310-000"
	// 	}`
	expectedCep := &CEP{
		Logradouro: "Avenida Paulista - até 610 - lado par",
		Bairro:     "Bela Vista",
		Localidade: "São Paulo/SP",
		Cep:        "01310-000",
	}
	actualCep, _ := BuscaCep("01310000")
	if *actualCep != *expectedCep {
		t.Errorf("Error getting CEP, got: %v, want: %v.", actualCep, expectedCep)
	}
}

func TestCepAsJson(t *testing.T) {
	actualCep, err := BuscaCep("01310000")
	if err != nil {
		t.Errorf("Got an Error trying to Marshal CEP to JSON: %v", err.Error())
	}

	_, err = actualCep.ToJSON()

	if err != nil {
		t.Errorf("Got an Error trying to Marshal CEP to JSON: %v", err.Error())
	}
}

func TestCepNaoExiste(t *testing.T) {
	actualCep, err := BuscaCep("00000000")

	if err == nil {
		t.Errorf("Should return nil, got: %v", actualCep)
	}
}

func TestCepDeveTerApenasDigitos(t *testing.T) {
	_, err := BuscaCep("123asd45")
	if err == nil {
		t.Errorf("Shoud return error, got nil")
	}
}

func TestCepDeveTerTamonho8(t *testing.T) {
	_, err := BuscaCep("234")
	if err == nil {
		t.Errorf("Shoud return error, got nil")
	}
}
