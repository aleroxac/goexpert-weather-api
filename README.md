# goexpert-weather-api
Projeto do Laboratório "Deploy com Cloud Run" do treinamento GoExpert(FullCycle).



## O desafio
Desenvolver um sistema em Go que receba um CEP, identifica a cidade e retorna o clima atual (temperatura em graus celsius, fahrenheit e kelvin). Esse sistema deverá ser publicado no Google Cloud Run.


## Como rodar o projeto
``` shell
# build the container image
make build

# push the container image and deploy to GCP Cloud Run
make deploy

# run locally
make run
```



## Funcionalidades da Linguagem Utilizadas
- context
- net/http
- encoding/json
- testing
- testify



## Requisitos: sistema
- [x] O sistema deve receber um CEP válido de 8 digitos
- [x] O sistema deve realizar a pesquisa do CEP e encontrar o nome da localização, a partir disso, deverá retornar as temperaturas e formata-lás em: Celsius, Fahrenheit, Kelvin.
- [x] O sistema deve responder adequadamente nos seguintes cenários:
    - Em caso de sucesso:
        - [x] Código HTTP: 200
        - [x] Response Body: { "temp_C": 28.5, "temp_F": 28.5, "temp_K": 28.5 }
    - Em caso de falha, caso o CEP não seja válido (com formato correto):
        - [x] Código HTTP: 422
        - [x] Mensagem: invalid zipcode
    - ​​​Em caso de falha, caso o CEP não seja encontrado:
        - [x] Código HTTP: 404
        - [x] Mensagem: can not find zipcode
- [x] Deverá ser realizado o deploy no Google Cloud Run.

## Requisitos: entrega
- [x] O código-fonte completo da implementação.
- [x] Testes automatizados demonstrando o funcionamento.
- [x] Utilize docker/docker-compose para que possamos realizar os testes de sua aplicação.
- [x] Deploy realizado no Google Cloud Run (free tier) e endereço ativo para ser acessado.