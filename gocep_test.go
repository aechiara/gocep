package gocep

import "testing"

func TestBuscaCep(t *testing.T) {
	expected := `{"logradouro":"Avenida Paulista - até 610 - lado par","bairro":"Bela Vista","localidade":"São Paulo/SP","cep":"01310-000"}`
	expected_cep := CEP{Logradouro: "Avenida Paulista - até 610 - lado par", Bairro: "Bela Vista", Localidade: "São Paulo/SP", Cep: "01310-000"}
	actual, actual_cep := BuscaCep("01310000")
	if actual != expected {
		t.Errorf("Error getting CEP, got: %v, want: %v.", actual, expected)
	}
	if actual_cep != expected_cep {
		t.Errorf("Error getting CEP, got: %v, want: %v.", actual_cep, expected_cep)
	}
}
