function log(target: any, propertyKey: string, descriptor: PropertyDescriptor) {
  const originalMethod = descriptor.value;

  descriptor.value = function (...args: any[]) {
    console.log(`Chamando ${propertyKey} com args:`, args);
    const result = originalMethod.apply(this, args);
    console.log(`Resultado:`, result);
    return result;
  };

  return descriptor;
}

class Calculator {
  @log()
  add(a: number, b: number) {
    return a + b;
  }
}

const calc = new Calculator();
calc.add(2, 3);
