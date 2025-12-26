const http = require("http");

console.log("Starting HTTP Server Test");

const PORT = 8085;
const state = { done: false };

const server = http.createServer((req, res) => {
  console.log("Server received request: " + req.method + " " + req.url);

  let body = "";
  req.on("data", (chunk) => {
    body += chunk;
  });

  req.on("end", () => {
    console.log("Server Body received: " + body);
    res.writeHead(200, { "Content-Type": "application/json" });
    res.end(
      JSON.stringify({
        message: "Hello from Nulang HTTP Server",
        received: body,
      })
    );
  });
});

server.listen(PORT, () => {
  console.log("Server listening on port " + PORT);

  const options = {
    method: "POST",
    headers: {
      "Content-Type": "text/plain",
    },
  };

  const req = http.request(
    "http://localhost:" + PORT + "/test",
    options,
    (res) => {
      console.log("Client received response: " + res.statusCode);

      let resBody = "";
      res.on("data", (chunk) => {
        resBody += chunk;
      });

      res.on("end", () => {
        console.log("Client Response body: " + resBody);
        server.close();
        console.log("Test Passed");
        state.done = true;
      });
    }
  );

  req.on("error", (e) => {
    console.error("Request error: " + e.message);
    server.close();
    state.done = true;
  });

  req.write("Initial Data");
  req.end();
});

console.log("Waiting for test completion...");
