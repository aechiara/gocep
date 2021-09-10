 # gocep
 Bibiloteca para buscar endereço a partir de um CEP no site dos Correios em Golang

 ### Install
```sh
 $ go get github.com/aechiara/gocep
```

### Use:
 ```go
package something 

import (
    "github.com/aechiara/gocep"
 )

func Something(cep string) {
    structCep, err := gocep.BuscaCep("01310000")
    jsonCep, err = stringCep.ToJSON()
}
```