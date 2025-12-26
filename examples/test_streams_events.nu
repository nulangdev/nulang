const stream = require("stream");

console.log("Defining MinhaStream class...");

class MinhaStream extends stream.Readable {
  constructor() {
    super();
  }
  _read(size) {
    console.log("MinhaStream._read called");
    this.push("dados dinâmicos");
    this.push(null); // fim
  }
}

console.log("Instantiating MinhaStream...");
const myStream = new MinhaStream();

console.log("Attaching data listener...");
myStream.on("data", (chunk) => {
  console.log("Received data: " + chunk.toString());
});

myStream.on("end", () => {
  console.log("Stream ended");
});

console.log("Manually calling read(0) to kickstart...");
myStream.read(0);

console.log("Done setup.");
