 # gocep
 Bibiloteca para buscar endereço a partir de um CEP no site dos Correios em Golang

 ### Install
```sh
 go get github.com/aechiara/gocep
```

### Use:
 ```go
package something 

import (
    "github.com/aechiara/gocep"
 )

func Something(cep string) (string, error) {
    structCep, err := gocep.Buscar(cep)
	if err != nil {
		//...
    }
    jsonCep, err := structCep.ToJSON()
	
	return jsonCep, nil
}
```

### Testes
Para rodar os testes unitários, use o comando padrão `go test`:
```sh
go test -v ./...
```

Para rodar os testes de integração, que fazem requisições reais para a API dos Correios, use a build tag `integration`:
```sh
go test -v -tags=integration ./...
```