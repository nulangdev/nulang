# Módulo DNS

O módulo `dns` fornece funções para realizar consultas DNS (Domain Name System), permitindo resolver nomes de domínio em endereços IP e vice-versa.

## Importação

```javascript
const dns = require("dns");
```

---

## Funções Principais

### `dns.lookup(hostname, options?, callback)`

Resolve um hostname para o primeiro registro A (IPv4) ou AAAA (IPv6) encontrado.

```javascript
// Uso básico
dns.lookup("google.com", (err, address, family) => {
  if (err) {
    console.error("Erro:", err.message);
    return;
  }
  console.log(`Endereço: ${address}`);
  console.log(`Família: IPv${family}`);
});

// Com opções
dns.lookup("google.com", { family: 4 }, (err, address, family) => {
  console.log(`IPv4: ${address}`);
});

dns.lookup("google.com", { family: 6 }, (err, address, family) => {
  console.log(`IPv6: ${address}`);
});
```

#### Opções

| Opção    | Tipo   | Descrição                              |
| -------- | ------ | -------------------------------------- |
| `family` | Number | 4 para IPv4, 6 para IPv6, 0 para ambos |

---

### `dns.resolve(hostname, rrtype?, callback)`

Resolve registros DNS de um tipo específico.

```javascript
// Registros A (IPv4)
dns.resolve("google.com", "A", (err, addresses) => {
  console.log("IPv4:", addresses);
});

// Registros AAAA (IPv6)
dns.resolve("google.com", "AAAA", (err, addresses) => {
  console.log("IPv6:", addresses);
});

// Registros MX (Mail Exchange)
dns.resolve("google.com", "MX", (err, records) => {
  records.forEach((record) => {
    console.log(`Exchange: ${record.exchange}, Priority: ${record.priority}`);
  });
});

// Registros TXT
dns.resolve("google.com", "TXT", (err, records) => {
  console.log("TXT:", records);
});

// Registros NS (Name Server)
dns.resolve("google.com", "NS", (err, servers) => {
  console.log("Name Servers:", servers);
});

// Registros CNAME
dns.resolve("www.google.com", "CNAME", (err, cnames) => {
  console.log("CNAMEs:", cnames);
});
```

#### Tipos de Registro (rrtype)

| Tipo      | Descrição                      |
| --------- | ------------------------------ |
| `'A'`     | Registros IPv4 (padrão)        |
| `'AAAA'`  | Registros IPv6                 |
| `'CNAME'` | Registros de nome canônico     |
| `'MX'`    | Registros de servidor de email |
| `'NS'`    | Registros de servidor de nomes |
| `'TXT'`   | Registros de texto             |

---

## Métodos Específicos por Tipo

### `dns.resolve4(hostname, options?, callback)`

Resolve registros A (IPv4).

```javascript
dns.resolve4("google.com", (err, addresses) => {
  console.log("IPv4 addresses:", addresses);
  // ['142.250.189.206', '142.250.189.238']
});
```

### `dns.resolve6(hostname, options?, callback)`

Resolve registros AAAA (IPv6).

```javascript
dns.resolve6("google.com", (err, addresses) => {
  console.log("IPv6 addresses:", addresses);
  // ['2607:f8b0:4004:800::200e']
});
```

### `dns.resolveMx(hostname, callback)`

Resolve registros MX (Mail Exchange).

```javascript
dns.resolveMx("google.com", (err, records) => {
  records.forEach((record) => {
    console.log(`${record.exchange} (prioridade: ${record.priority})`);
  });
});
```

### `dns.resolveNs(hostname, callback)`

Resolve registros NS (Name Server).

```javascript
dns.resolveNs("google.com", (err, servers) => {
  console.log("Servidores DNS:", servers);
  // ['ns1.google.com', 'ns2.google.com', ...]
});
```

### `dns.resolveTxt(hostname, callback)`

Resolve registros TXT.

```javascript
dns.resolveTxt("google.com", (err, records) => {
  console.log("Registros TXT:", records);
});
```

### `dns.resolveCname(hostname, callback)`

Resolve registros CNAME.

```javascript
dns.resolveCname("www.google.com", (err, cnames) => {
  console.log("CNAMEs:", cnames);
});
```

---

## DNS Reverso

### `dns.reverse(ip, callback)`

Realiza uma consulta DNS reversa, convertendo um IP para um hostname.

```javascript
dns.reverse("8.8.8.8", (err, hostnames) => {
  if (err) {
    console.error("Erro:", err.message);
    return;
  }
  console.log("Hostnames:", hostnames);
  // ['dns.google']
});

dns.reverse("1.1.1.1", (err, hostnames) => {
  console.log("Cloudflare DNS:", hostnames);
});
```

---

## Exemplos Práticos

### Verificador de Domínios

```javascript
const dns = require("dns");

function checkDomain(domain) {
  console.log(`\n=== Verificando: ${domain} ===\n`);

  // IPv4
  dns.resolve4(domain, (err, addresses) => {
    if (err) {
      console.log("IPv4: Não encontrado");
    } else {
      console.log("IPv4:", addresses.join(", "));
    }
  });

  // IPv6
  dns.resolve6(domain, (err, addresses) => {
    if (err) {
      console.log("IPv6: Não encontrado");
    } else {
      console.log("IPv6:", addresses.join(", "));
    }
  });

  // MX Records
  dns.resolveMx(domain, (err, records) => {
    if (err) {
      console.log("MX: Não encontrado");
    } else {
      console.log("MX Records:");
      records.forEach((r) => {
        console.log(`  ${r.exchange} (${r.priority})`);
      });
    }
  });

  // NS Records
  dns.resolveNs(domain, (err, servers) => {
    if (err) {
      console.log("NS: Não encontrado");
    } else {
      console.log("Name Servers:", servers.join(", "));
    }
  });
}

checkDomain("google.com");
```

### Validador de Email (via MX)

```javascript
const dns = require("dns");

function validateEmailDomain(email) {
  const domain = email.split("@")[1];

  if (!domain) {
    console.log("Email inválido: domínio não encontrado");
    return;
  }

  dns.resolveMx(domain, (err, records) => {
    if (err || records.length === 0) {
      console.log(`Domínio ${domain} não aceita emails`);
    } else {
      console.log(`Domínio ${domain} é válido para emails`);
      console.log("Servidores de email:");
      records.forEach((r) => {
        console.log(`  ${r.exchange}`);
      });
    }
  });
}

validateEmailDomain("user@gmail.com");
```

### Verificador de Latência DNS

```javascript
const dns = require("dns");

function measureDNS(hostname) {
  const start = Date.now();

  dns.lookup(hostname, (err, address) => {
    const duration = Date.now() - start;

    if (err) {
      console.log(`${hostname}: Falha - ${err.message}`);
    } else {
      console.log(`${hostname}: ${address} (${duration}ms)`);
    }
  });
}

measureDNS("google.com");
measureDNS("github.com");
measureDNS("facebook.com");
```

---

## lookup vs resolve

| Característica    | `lookup`                     | `resolve`                 |
| ----------------- | ---------------------------- | ------------------------- |
| Fonte             | SO (cache local, hosts file) | Servidor DNS diretamente  |
| Performance       | Pode ser mais rápido         | Sempre consulta rede      |
| Uso               | Uso geral, imita navegadores | Consultas DNS específicas |
| Tipos de registro | Apenas A/AAAA                | Todos os tipos            |

```javascript
// lookup usa o resolver do SO
dns.lookup("localhost", (err, address) => {
  console.log(address); // 127.0.0.1 (do /etc/hosts)
});

// resolve sempre consulta DNS
dns.resolve4("localhost", (err, addresses) => {
  if (err) console.log("Não resolvido via DNS");
});
```

---

## Tratamento de Erros

```javascript
dns.lookup("dominio-inexistente.xyz", (err, address) => {
  if (err) {
    switch (err.code) {
      case "ENOTFOUND":
        console.log("Domínio não encontrado");
        break;
      case "ENODATA":
        console.log("Nenhum dado encontrado");
        break;
      case "ESERVFAIL":
        console.log("Falha no servidor DNS");
        break;
      case "ETIMEOUT":
        console.log("Timeout na consulta");
        break;
      default:
        console.log("Erro:", err.message);
    }
    return;
  }
  console.log("Endereço:", address);
});
```

---

## Notas de Compatibilidade

| Funcionalidade | Nulang | Node.js |
| -------------- | ------ | ------- |
| `lookup`       | ✅     | ✅      |
| `resolve`      | ✅     | ✅      |
| `resolve4`     | ✅     | ✅      |
| `resolve6`     | ✅     | ✅      |
| `resolveMx`    | ✅     | ✅      |
| `resolveNs`    | ✅     | ✅      |
| `resolveTxt`   | ✅     | ✅      |
| `resolveCname` | ✅     | ✅      |
| `reverse`      | ✅     | ✅      |
| `setServers`   | ❌     | ✅      |
| `getServers`   | ❌     | ✅      |
| `promises` API | ❌     | ✅      |
