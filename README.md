# API Validator - Automated Tax Validation

A personal project that validates the full tax classification of retail food products in Brazil. I used to work in tax compliance, and I always wanted to automate this task by combining curated tax data with AI to check product registrations.

> **Project status:** I did not take the project further because Brazilian tax legislation changes constantly. Keeping the data current would require a large ongoing effort or a separate service to monitor and update the tax rules.

---

## 🔍 How the Validator Works

The API receives a product description (e.g. *"Biscoito Recheado Chocolate 130g"*) together with its NCM code and runs a chain of validations:

```
[Input: Product + NCM]
       │
       ▼
 1. NCM validation (AI) ──(Incompatible?)──► [Stop validation]
       │ (Compatible)
       ├─► 2. IPI validation (Lookup + AI for Ex-tariff exceptions)
       ├─► 3. CEST validation (Lookup + AI for disambiguation)
       ├─► 4. PIS/COFINS validation (Lookup)
       ├─► 5. Exemption validation (Lookup + AI) ──(Exempt?)──► [Return response]
       └─► 6. Deferral validation (Lookup + AI) ──(Deferred?)──► [Return response]
             │
             ├─► 7. ICMS rate validation (Lookup + AI)
             ├─► 8. Tax base reduction validation (Lookup + AI)
             └─► 9. ICMS ST validation (based on the validated CEST)
```

### 📌 Validation Steps

1. **NCM (Mercosur Common Nomenclature, the product classification code):**
   - *Where AI comes in:* Google Gemini checks semantically whether the product's commercial description matches the official NCM description stored in the database.
   - *Stop condition:* If the NCM is considered **incompatible**, the process stops immediately to avoid calculating taxes on a wrong NCM.

2. **IPI (federal tax on manufactured products):**
   - Looks up the IPI rate for the NCM.
   - *Where AI comes in:* When tariff exceptions exist (Ex-IPI), the AI analyzes the product details to check whether it fits the exception.

3. **CEST (tax substitution code):**
   - Maps the NCM to the CEST table in the database.
   - *Where AI comes in:* When several CEST codes are possible for the same NCM, the AI analyzes the product description to select the exact one.

4. **PIS / COFINS (federal social contributions):**
   - Looks up the tax treatment (zero rate, taxed, single-phase) and the corresponding legal basis in MongoDB.

5. **ICMS exemption (state VAT):**
   - Looks up exemption rules and tax agreements linked to the NCM.
   - *Where AI comes in:* Evaluates whether the product meets the conditions required by the legal text.
   - *Precedence:* If the product is confirmed as **exempt**, the flow ends (exemption takes precedence over rates, reductions and ICMS ST).

6. **ICMS deferral:**
   - Checks rules that defer the tax payment.
   - *Where AI comes in:* Analyzes whether the deferral applies to the product's description and category.
   - *Precedence:* If the product is **deferred**, the flow ends.

7. **Internal ICMS rate:**
   - Looks up the standard or specific rate registered for the NCM in the state.
   - *Where AI comes in:* Helps identify differentiated rates (e.g. staple food products or specific rates).

8. **ICMS tax base reduction:**
   - Checks whether a tax base reduction benefit exists for the NCM + rate combination.
   - *Where AI comes in:* Analyzes whether the product's characteristics meet the benefit's regulation.

9. **ICMS ST (tax substitution):**
   - Validates whether tax substitution applies by cross-checking the NCM classification with the CEST validated earlier.

> Scope: São Paulo state, intrastate operations.

---

## 🛠️ Tech Stack

- **Go (Golang)** - Main language and validation pipeline
- **Gin Web Framework** - HTTP routing and REST API
- **MongoDB** - Storage for tax tables and rules
- **Google Gemini API** - Semantic analysis and tax rule classification

---

## 🚀 Running the Project

1. Clone the repository and open the folder:
   ```bash
   git clone https://github.com/raphaeldertinatti/api-validator.git
   cd api-validator
   ```

2. Create a `.env` file based on `.env.example`:
   ```env
   MONGO_URI=mongodb://localhost:27017
   DB_NAME=validator_db
   GEMINI_API_KEY=your_key_here
   PORT=8080
   ```

3. Run the API:
   ```bash
   go run cmd/validator/main.go
   ```

The API runs on `http://localhost:8080` by default.
