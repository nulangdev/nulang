// Testing new features: fs, crypto, path, Promise

console.log("=== Testing File System (fs) ===");

// Create a test file
fs.writeFileSync("test_file.txt", "Hello, Nulang!");
console.log("File written successfully");

// Check if file exists
console.log("File exists:", fs.existsSync("test_file.txt"));

// Read the file
let content = fs.readFileSync("test_file.txt", "utf8");
console.log("File content:", content);

// Get file stats
let stats = fs.statSync("test_file.txt");
console.log("File size:", stats.size, "bytes");
console.log("Is file:", stats.isFile());
console.log("Is directory:", stats.isDirectory());

// Append to file
fs.appendFileSync("test_file.txt", "\nAppended line!");
console.log("Content after append:", fs.readFileSync("test_file.txt", "utf8"));

// Create a directory
if (!fs.existsSync("test_dir")) {
    fs.mkdirSync("test_dir");
    console.log("Directory created");
}

// Copy file
fs.copyFileSync("test_file.txt", "test_dir/copy.txt");
console.log("File copied");

// List directory
console.log("Directory contents:", fs.readdirSync("test_dir"));

// Clean up
fs.unlinkSync("test_file.txt");
fs.unlinkSync("test_dir/copy.txt");
fs.rmdirSync("test_dir");
console.log("Cleanup complete");

console.log("\n=== Testing Path ===");

console.log("path.join:", path.join("folder", "subfolder", "file.txt"));
console.log("path.dirname:", path.dirname("/path/to/file.txt"));
console.log("path.basename:", path.basename("/path/to/file.txt"));
console.log("path.extname:", path.extname("/path/to/file.txt"));
console.log("path.isAbsolute('/path'):", path.isAbsolute("/path"));
console.log("path.isAbsolute('./path'):", path.isAbsolute("./path"));
console.log("path.parse:", path.parse("/path/to/file.txt"));
console.log("path.sep:", path.sep);

console.log("\n=== Testing Crypto ===");

// Create hash
let hash = crypto.createHash("sha256");
hash.update("Hello, World!");
console.log("SHA256:", hash.digest("hex"));

// Create HMAC
let hmac = crypto.createHmac("sha256", "secret-key");
hmac.update("Hello, World!");
console.log("HMAC-SHA256:", hmac.digest("hex"));

// Random bytes
let randomBytes = crypto.randomBytes(16);
console.log("Random bytes:", randomBytes.toString("hex"));

// Random UUID
console.log("UUID:", crypto.randomUUID());

console.log("\n=== Testing Buffer ===");

// Buffer.from string
let buf1 = Buffer.from("Hello, Nulang!");
console.log("Buffer from string:", buf1.toString());
console.log("Buffer length:", buf1.length);

// Buffer.from array
let buf2 = Buffer.from([72, 101, 108, 108, 111]);
console.log("Buffer from array:", buf2.toString());

// Buffer.alloc
let buf3 = Buffer.alloc(10, 0);
console.log("Allocated buffer:", buf3);
buf3.fill(65);
console.log("Filled buffer:", buf3.toString());

// Buffer operations
console.log("Buffer.isBuffer:", Buffer.isBuffer(buf1));
console.log("Buffer hex:", buf1.toString("hex"));
console.log("Buffer base64:", buf1.toString("base64"));
console.log("Buffer slice:", buf1.slice(0, 5).toString());

// Buffer.concat
let concatenated = Buffer.concat([buf1, buf2]);
console.log("Concatenated:", concatenated.toString());

console.log("\n=== Testing Promise ===");

// Promise.resolve
let p1 = Promise.resolve(42);
console.log("Promise.resolve:", p1);
console.log("Promise state:", p1.state);
console.log("Promise value:", p1.value);

// Promise chain
let result = p1.then(x => x * 2);
console.log("After then:", result.value);

// Promise.reject
let p2 = Promise.reject("Error!");
console.log("Promise.reject state:", p2.state);

// Promise.all
let p3 = Promise.all([
    Promise.resolve(1),
    Promise.resolve(2),
    Promise.resolve(3)
]);
console.log("Promise.all result:", p3.value);

// Promise.allSettled
let p4 = Promise.allSettled([
    Promise.resolve("success"),
    Promise.reject("error")
]);
console.log("Promise.allSettled:", p4.value);

console.log("\n=== Testing require() ===");

// require built-in modules
let fsModule = require("fs");
let pathModule = require("path");
let cryptoModule = require("crypto");

console.log("fs via require:", typeof fsModule);
console.log("path via require:", typeof pathModule);
console.log("crypto via require:", typeof cryptoModule);

console.log("\n✅ All new features working!");
