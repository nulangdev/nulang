const http = require("http");

console.log("Starting CRUD API Server...");

const PORT = 8086;

// Mock Database
let users = [];
let nextId = 1;

// Helper to parse body
const parseBody = (req, callback) => {
  let body = "";
  req.on("data", (chunk) => {
    body += chunk;
  });
  req.on("end", () => {
    callback(body);
  });
};

// Helper to send JSON
const sendJSON = (res, status, data) => {
  res.writeHead(status, { "Content-Type": "application/json" });
  res.end(JSON.stringify(data));
};

const server = http.createServer((req, res) => {
  console.log("Request: " + req.method + " " + req.url);

  // Simple CORS
  res.setHeader("Access-Control-Allow-Origin", "*");
  res.setHeader(
    "Access-Control-Allow-Methods",
    "GET, POST, PUT, DELETE, OPTIONS"
  );

  if (req.method === "OPTIONS") {
    res.writeHead(204);
    res.end();
    return;
  }

  const urlParts = req.url.split("/"); // e.g., ["", "users", "1"]
  const resource = urlParts[1];
  const id = urlParts[2];

  if (resource === "users") {
    if (req.method === "GET") {
      if (id) {
        // Get one
        const user = users.find((u) => u.id == id);
        if (user) {
          sendJSON(res, 200, user);
        } else {
          sendJSON(res, 404, { error: "User not found" });
        }
      } else {
        // List all
        sendJSON(res, 200, { users: users });
      }
    } else if (req.method === "POST") {
      // Create
      parseBody(req, (body) => {
        let data = {};
        try {
          data = JSON.parse(body);
        } catch (e) {}
        const newUser = { id: nextId, name: data.name, email: data.email };
        nextId = nextId + 1;
        users.push(newUser);
        sendJSON(res, 201, newUser);
      });
    } else if (req.method === "PUT" && id) {
      // Update
      parseBody(req, (body) => {
        let data = {};
        try {
          data = JSON.parse(body);
        } catch (e) {}
        const userIndex = users.findIndex((u) => u.id == id);
        if (userIndex !== -1) {
          users[userIndex].name = data.name || users[userIndex].name;
          users[userIndex].email = data.email || users[userIndex].email;
          sendJSON(res, 200, users[userIndex]);
        } else {
          sendJSON(res, 404, { error: "User not found" });
        }
      });
    } else if (req.method === "DELETE" && id) {
      // Delete
      const initialLength = users.length;
      users = users.filter((u) => u.id != id);
      if (users.length < initialLength) {
        sendJSON(res, 204, null);
      } else {
        sendJSON(res, 404, { error: "User not found" });
      }
    } else {
      sendJSON(res, 405, { error: "Method Not Allowed" });
    }
  } else {
    if (req.url === "/") {
      sendJSON(res, 200, { message: "Welcome to Nulang CRUD API" });
    } else {
      sendJSON(res, 404, { error: "Not Found" });
    }
  }
});

server.listen(PORT, () => {
  console.log("CRUD Server listening on port " + PORT);
  console.log("Ready for manual testing.");
});
