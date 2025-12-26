# Crypto Module

O módulo `crypto` fornece funcionalidades criptográficas básicas.

## Importação

```javascript
const crypto = require("crypto");
// ou
import crypto from "crypto";
```

## Funções de Hash

### crypto.createHash(algorithm)

Cria um objeto Hash para gerar digests.

```javascript
const hash = crypto.createHash("sha256");
```

#### Algoritmos Suportados

| Algoritmo | Descrição                                         |
| --------- | ------------------------------------------------- |
| `md5`     | MD5 (128 bits) - não recomendado para segurança   |
| `sha1`    | SHA-1 (160 bits) - não recomendado para segurança |
| `sha256`  | SHA-256 (256 bits)                                |
| `sha512`  | SHA-512 (512 bits)                                |

### Objeto Hash

#### hash.update(data)

Adiciona dados ao hash.

```javascript
const hash = crypto.createHash("sha256");
hash.update("Hello");
hash.update(", World!");
```

**Retorno**: O próprio objeto Hash (permite encadeamento)

#### hash.digest([encoding])

Calcula e retorna o digest.

```javascript
const hash = crypto.createHash("sha256");
hash.update("Hello, World!");
const digest = hash.digest("hex");
console.log(digest);
// "dffd6021bb2bd5b0af676290809ec3a53191dd81c7f70a4b28688a362182986f"
```

**Encodings suportados**:
| Encoding | Descrição |
|----------|-----------|
| `hex` | String hexadecimal |
| `base64` | String base64 |
| (nenhum) | Buffer |

### Exemplos de Hash

```javascript
// SHA-256
function sha256(text) {
  return crypto.createHash("sha256").update(text).digest("hex");
}

console.log(sha256("senha123"));
// "ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f"

// MD5
function md5(text) {
  return crypto.createHash("md5").update(text).digest("hex");
}

console.log(md5("Hello"));
// "8b1a9953c4611296a827abf8c47804d7"
```

## HMAC

### crypto.createHmac(algorithm, key)

Cria um objeto HMAC (Hash-based Message Authentication Code).

```javascript
const hmac = crypto.createHmac("sha256", "chave-secreta");
```

### Objeto HMAC

#### hmac.update(data)

Adiciona dados ao HMAC.

#### hmac.digest([encoding])

Calcula e retorna o HMAC.

### Exemplo de HMAC

```javascript
function createSignature(data, secret) {
  return crypto.createHmac("sha256", secret).update(data).digest("hex");
}

const signature = createSignature("mensagem", "minha-chave");
console.log(signature);

// Verificar assinatura
function verifySignature(data, secret, signature) {
  const expected = createSignature(data, secret);
  return expected === signature;
}
```

## Números Aleatórios

### crypto.randomBytes(size)

Gera bytes aleatórios criptograficamente seguros.

```javascript
const bytes = crypto.randomBytes(16);
console.log(bytes.toString("hex"));
// "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4"
```

### crypto.randomUUID()

Gera um UUID v4 aleatório.

```javascript
const uuid = crypto.randomUUID();
console.log(uuid);
// "550e8400-e29b-41d4-a716-446655440000"
```

## Exemplos Práticos

### Hash de Senha

```javascript
function hashPassword(password, salt) {
  return crypto
    .createHash("sha256")
    .update(password + salt)
    .digest("hex");
}

function generateSalt() {
  return crypto.randomBytes(16).toString("hex");
}

// Criar hash
const salt = generateSalt();
const hash = hashPassword("minhaSenha123", salt);
console.log(`Salt: ${salt}`);
console.log(`Hash: ${hash}`);

// Verificar senha
function verifyPassword(password, salt, hash) {
  return hashPassword(password, salt) === hash;
}

console.log(verifyPassword("minhaSenha123", salt, hash)); // true
console.log(verifyPassword("senhaErrada", salt, hash)); // false
```

### Token de Sessão

```javascript
function generateSessionToken() {
  return crypto.randomBytes(32).toString("hex");
}

const token = generateSessionToken();
console.log(token);
// "a1b2c3d4e5f6..."  (64 caracteres hex)
```

### API Key

```javascript
function generateApiKey() {
  const bytes = crypto.randomBytes(24);
  return bytes.toString("base64");
}

const apiKey = generateApiKey();
console.log(apiKey);
// "Base64EncodedString=="
```

### Checksum de Arquivo

```javascript
const fs = require("fs");

function fileChecksum(filepath) {
  const data = fs.readFileSync(filepath);
  return crypto.createHash("md5").update(data).digest("hex");
}

const checksum = fileChecksum("arquivo.txt");
console.log(`MD5: ${checksum}`);
```

### Assinatura de Webhook

```javascript
function signWebhook(payload, secret) {
  const hmac = crypto.createHmac("sha256", secret);
  hmac.update(JSON.stringify(payload));
  return hmac.digest("hex");
}

function verifyWebhook(payload, signature, secret) {
  const expected = signWebhook(payload, secret);
  return expected === signature;
}

// Enviar webhook
const payload = { event: "user.created", data: { id: 123 } };
const secret = "webhook-secret";
const signature = signWebhook(payload, secret);

// Verificar no receptor
const isValid = verifyWebhook(payload, signature, secret);
console.log("Webhook válido:", isValid);
```

### Token JWT Simplificado

```javascript
function base64UrlEncode(str) {
  return Buffer.from(str)
    .toString("base64")
    .replaceAll("+", "-")
    .replaceAll("/", "_")
    .replaceAll("=", "");
}

function createToken(payload, secret) {
  const header = { alg: "HS256", typ: "JWT" };

  const headerB64 = base64UrlEncode(JSON.stringify(header));
  const payloadB64 = base64UrlEncode(JSON.stringify(payload));

  const signature = crypto
    .createHmac("sha256", secret)
    .update(`${headerB64}.${payloadB64}`)
    .digest("base64")
    .replaceAll("+", "-")
    .replaceAll("/", "_")
    .replaceAll("=", "");

  return `${headerB64}.${payloadB64}.${signature}`;
}

const token = createToken(
  { userId: 123, exp: Date.now() + 3600000 },
  "jwt-secret"
);
console.log(token);
```

### Hash Consistente

```javascript
function consistentHash(key, buckets) {
  const hash = crypto.createHash("md5").update(key).digest("hex");

  // Converte primeiros 8 chars hex para número
  const num = parseInt(hash.substring(0, 8), 16);
  return num % buckets;
}

// Distribuir chaves em 10 buckets
console.log(consistentHash("user:123", 10)); // 7
console.log(consistentHash("user:456", 10)); // 2
console.log(consistentHash("user:789", 10)); // 4
```

## Comparação Segura de Tempo Constante

```javascript
function secureCompare(a, b) {
  if (a.length !== b.length) {
    return false;
  }

  let result = 0;
  for (let i = 0; i < a.length; i++) {
    result |= a.charCodeAt(i) ^ b.charCodeAt(i);
  }
  return result === 0;
}

// Uso para comparar hashes
const hash1 = crypto.createHash("sha256").update("test").digest("hex");
const hash2 = crypto.createHash("sha256").update("test").digest("hex");
console.log(secureCompare(hash1, hash2)); // true
```

## Boas Práticas

### ✅ Use SHA-256 ou superior para hashes

```javascript
// ✅ Bom
crypto.createHash("sha256");

// ❌ Evite para segurança
crypto.createHash("md5");
crypto.createHash("sha1");
```

### ✅ Use HMAC para autenticação de mensagens

```javascript
// ✅ Bom para assinar dados
crypto.createHmac("sha256", secret).update(data);
```

### ✅ Use randomBytes para segredos

```javascript
// ✅ Criptograficamente seguro
crypto.randomBytes(32);

// ❌ Não use para segurança
Math.random();
```

## Veja Também

- [Buffer](./buffer.md) - Manipulação de dados binários
- [HTTP](./http.md) - Requisições HTTP
