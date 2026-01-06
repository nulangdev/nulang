import http from "http";

const PORT = 3000;

const server = http.createServer((req, res) => {
  console.log(`${req.method} ${req.url}`);

  res.setHeader("Content-Type", "application/json");
  res.statusCode = 200;

  res.end(
    JSON.stringify({
      message: "Hello from Nulang HTTP Server!",
      timestamp: Date.now(),
      path: req.url,
    })
  );
});

server.listen(PORT, () => {
  console.log(`🚀 Server running at http://localhost:${PORT}`);
});
