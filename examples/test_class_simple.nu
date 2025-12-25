// Simple class test
class Animal {
  constructor(name) {
    this.name = name;
  }

  speak() {
    console.log(this.name + " makes a sound");
  }
}

console.log("Animal type:", typeof Animal);
console.log("Animal:", Animal);

let animal = new Animal("Dog");
console.log("Created animal:", animal);
console.log("Animal name:", animal.name);
animal.speak();
