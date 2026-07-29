# API Validator - Validação Tributária Automática

Projeto pessoal para validação tributária completa de produtos do varejo alimentício. Como já trabalhei na área, sempre quis tentar automatizar este desafio cruzando dados fiscais com inteligência artificial para checar cadastros de produtos.

> **Sobre a continuidade:** Eu não levei o projeto para frente por conta das constantes alterações na legislação, o que demandaria um tempo muito grande ou um outro serviço de monitoramento e atualização desses dados tributários.

---

## 🔍 Como Funciona o Validador

A API recebe a descrição do produto (ex: *"Biscoito Recheado Chocolate 130g"*) junto com a NCM e executa uma esteira de validações encadeadas:

```
[Entrada: Produto + NCM] 
       │
       ▼
 1. Validação NCM (IA) ──(Incompatível?)──► [Para a Validação]
       │ (Compatível)
       ├─► 2. Validação IPI (Consulta + IA para Enquadramento Ex)
       ├─► 3. Validação CEST (Consulta + IA para Desambiguação)
       ├─► 4. Validação PIS/COFINS (Consulta)
       ├─► 5. Validação Isenção (Consulta + IA) ──(É Isento?)──► [Retorna Resposta]
       └─► 6. Validação Diferimento (Consulta + IA) ──(É Diferido?)──► [Retorna Resposta]
             │
             ├─► 7. Validação Alíquota ICMS (Consulta + IA)
             ├─► 8. Validação Redução de BC (Consulta + IA)
             └─► 9. Validação ICMS ST (Valida com base no CEST)
```

### 📌 Mapeamento das Validações

1. **NCM (Nomenclatura Comum do Mercosul):** 
   - *Onde entra a IA:* A IA (Google Gemini) analisa semanticamente se a descrição comercial do produto condiz com a descrição oficial da NCM no banco de dados.
   - *Gatilho de parada:* Se a NCM for considerada **incompatível**, o processo é interrompido imediatamente para evitar cálculos tributários sobre um NCM errado.

2. **IPI (Imposto sobre Produtos Industrializados):**
   - Verifica a alíquota da tabela de IPI por NCM.
   - *Onde entra a IA:* Quando existem exceções tarifárias (Ex-IPI), a IA analisa o detalhamento do produto para verificar se ele se enquadra na exceção.

3. **CEST (Código Especificador da Substituição Tributária):**
   - Relaciona o NCM com a tabela de CEST no banco.
   - *Onde entra a IA:* Se houverem múltiplos CESTs possíveis para a mesma NCM, a IA analisa a descrição do item para selecionar o CEST exato.

4. **PIS / COFINS:**
   - Consulta a tributação (Alíquota zero, Tributado, Monofásico) e a fundamentação legal correspondente na base MongoDB.

5. **Isenção de ICMS:**
   - Consulta regras e convênios fiscais de isenção vinculados à NCM.
   - *Onde entra a IA:* Avalia se o produto atende às condições e especificidades exigidas pelo texto legal.
   - *Prevalência:* Se for confirmado que o produto é **Isento**, o fluxo é finalizado (a isenção prevalece sobre alíquotas, reduções e ICMS ST).

6. **Diferimento de ICMS:**
   - Verifica regras de adiamento do imposto.
   - *Onde entra a IA:* Analisa se a aplicação do diferimento é válida para a descrição e categoria do produto.
   - *Prevalência:* Se for **Diferido**, o fluxo é finalizado.

7. **Alíquota Interna de ICMS:**
   - Consulta no banco de dados a alíquota padrão ou específica cadastrada para o NCM no estado.
   - *Onde entra a IA:* Auxilia na identificação de alíquotas diferenciadas (ex: produtos de cesta básica ou alíquotas específicas).

8. **Redução da Base de Cálculo (ICMS):**
   - Verifica se há benefício de redução de BC cadastrado para a combinação NCM + Alíquota.
   - *Onde entra a IA:* Analisa se as características do produto atendem ao regulamento do benefício.

9. **ICMS ST (Substituição Tributária):**
   - Valida a incidência de Substituição Tributária cruzando o enquadramento do NCM e do CEST validado anteriormente.

---

## 🛠️ Tecnologias

- **Go (Golang)** - Linguagem principal e esteira de validação
- **Gin Web Framework** - Roteamento HTTP e API REST
- **MongoDB** - Armazenamento de tabelas e regras fiscais
- **Google Gemini API** - Análise semântica e enquadramento de regras fiscais

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
