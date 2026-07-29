# API Validator - Validação Tributária Automática

Projeto pessoal para validação tributária completa de produtos do varejo alimentício. Como já trabalhei na área, sempre quis tentar automatizar este desafio cruzando dados fiscais com inteligência artificial para checar cadastros de produtos.

> **Sobre a continuidade:** Eu não levei o projeto para frente por conta das constantes alterações na legislação, o que demandaria um tempo muito grande ou um outro serviço de monitoramento e atualização desses dados tributários.

---

## 🛠️ Tecnologias

- **Go (Golang)** - API e lógica de negócio
- **Gin Web Framework** - Roteamento HTTP
- **MongoDB** - Armazenamento das tabelas fiscais (NCM, CEST, PIS/COFINS, Alíquotas, ICMS ST, etc.)
- **Google Gemini API** - Análise semântica da descrição dos produtos

---

## 🚀 Como rodar o projeto

1. Clone o repositório e acesse a pasta:
   ```bash
   git clone https://github.com/raphaeldertinatti/api-validator.git
   cd api-validator
   ```

2. Crie um arquivo `.env` baseado no `.env.example`:
   ```env
   MONGO_URI=mongodb://localhost:27017
   DB_NAME=validator_db
   GEMINI_API_KEY=sua_chave_aqui
   PORT=8080
   ```

3. Execute a API:
   ```bash
   go run cmd/validator/main.go
   ```

A API rodará por padrão em `http://localhost:8080`.
