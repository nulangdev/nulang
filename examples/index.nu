import dns from "dns";

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
