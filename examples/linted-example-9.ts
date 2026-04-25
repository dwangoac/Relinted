class Greeter {
    private greeting: string;

    constructor(message: string) {
        this.greeting = message;
    }

    greet(): string {
        return "Hello, " + this.greeting;
    }
}

function main() {
    const greeter = new Greeter("world");
    console.log(greeter.greet());

    const numbers: number[] = [1, 2, 3];
    for (const n of numbers) {
        console.log(n);
    }
}

main();
