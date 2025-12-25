// Class inheritance test
class Animal {
  constructor(name) {
    this.name = name;
  }

  speak() {
    console.log(this.name + " makes a sound");
  }
}

console.log("Creating Animal...");
let animal = new Animal("Generic Animal");
animal.speak();
console.log("Animal name:", animal.name);

console.log("\nCreating Dog class...");
class Dog extends Animal {
  constructor(name, breed) {
    // super(name);  // Skip super for now
    this.name = name;
    this.breed = breed;
  }

  speak() {
    console.log(this.name + " barks!");
  }

  getInfo() {
    return this.name + " is a " + this.breed;
  }
}

console.log("Dog type:", typeof Dog);

console.log("\nCreating dog instance...");
let dog = new Dog("Rex", "German Shepherd");
console.log("Dog:", dog);
dog.speak();
console.log("Dog info:", dog.getInfo());
