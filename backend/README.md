# Projeto Final de Gerenciamento de Restaurantes

Este projeto é uma API RESTful para um sistema de gerenciamento de restaurantes, construído usando Go. Ele permite que os usuários interajam com os dados do restaurante, incluindo a criação, leitura, atualização e exclusão de informações do restaurante. O projeto também inclui recursos de gerenciamento de usuários.

https://documenter.getpostman.com/view/36585246/2sA3dskt8T

## Características

- **Criação de Recursos**: Crie novos restaurantes ou contas de usuário.
- **Leitura de Recursos**: Obtenha informações sobre todos os restaurantes, tipos de comida e recomendações.
- **Atualização de Recursos**: Atualize as informações de um restaurante ou de um usuário.
- **Exclusão de Recursos**: Exclua um restaurante ou um usuário.
- **Tratamento de Erros**: A API responde com códigos e mensagens de status HTTP apropriados em caso de erros.

# Endpoints da API

Aqui estão os endpoints disponíveis na API:

## Usuário Normal
- `POST /conta`: Cria uma nova conta de usuário.
- `GET /login/user/:username`: Encontra um usuário existente no banco de dados.
- `POST /login`: Faz login com uma conta de usuário.
- `GET /auth/protected`: Endpoint protegido que requer autenticação. Retorna o ID e o nome de usuário do usuário autenticado.

## Usuário Restaurante
- `POST /login/restaurante`: Faz login com uma conta de restaurante.
- `POST /conta/restaurante`: Cria uma nova conta de restaurante.
- `GET /restaurante`: Obtém todos os restaurantes.
- `GET /login/rest/:restname`: Encontra um restaurante existente no banco de dados.
- `POST /saveLists/:userName`: Salva listas para um usuário.
- `DELETE /restaurante/:userName/lista/:nome`: Exclui uma lista de um usuário.
- `GET /restaurante/:userName/lista/get`: Obtém todas as listas de um usuário.
- `POST /restaurante/:userName/imagem`: Atualiza a personalização de um restaurante.
- `GET /restaurante/:userName/pegarimagem`: Obtém a imagem de personalização de um restaurante.

## Checagem de Permissões
- `GET /auth/checkPermissions`: Checa as permissões do usuário autenticado.

## Tratamento de Erros

O servidor responde com códigos e mensagens de status HTTP apropriados em caso de erros. Por exemplo:

- `400 Bad Request`: A solicitação estava malformada ou contém dados inválidos.
- `401 Unauthorized`: Credenciais inválidas.
- `404 Not Found`: O recurso solicitado não foi encontrado.
- `500 Internal Server Error`: Ocorreu um erro no lado do servidor.

Cada resposta de erro é um objeto JSON contendo `code`, `message` e `info`.

## Como executar o projeto

Para executar este projeto, você precisará ter Go instalado em seu ambiente de desenvolvimento. Siga estas etapas:

1. Clone o repositório para o seu ambiente local.
2. Navegue até a pasta do projeto.
3. Execute `go mod download` para baixar as dependências necessárias.
4. Execute `go run main.go` para iniciar o servidor.

Agora, você deve ser capaz de acessar a API em `http://localhost:8080`.

## Banco de Dados

O banco de dados está rodando localmente, utilizando o MongoDB. Você pode acessá-lo em `mongodb://localhost:27017`.
https://www.mongodb.com

## Frontend Vue

O frontend do projeto, construído com Vue.js, também está rodando localmente. Você pode acessá-lo em `http://localhost:5173`.
https://vuejs.org/guide/quick-start