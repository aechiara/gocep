//go:build integration

package gocep_test

import (
	"github.com/aechiara/gocep"
	"reflect"
	"testing"
)

func TestBuscarIntegration(t *testing.T) {
	t.Run("busca cep real", func(t *testing.T) {
		// Precisamos garantir que estamos usando a URL real para este teste.
		// A URL padrão é a real, mas outro teste pode tê-la alterado.
		// Então, salvamos a URL atual, a definimos como a padrão e a restauramos depois.
		// Uma abordagem melhor para o futuro seria não ter um estado global para a URL.
		originalURL := gocep.GetCepURL()
		gocep.SetCepURL("https://buscacepinter.correios.com.br/app/consulta/html/consulta-detalhes-cep.php")
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
}
