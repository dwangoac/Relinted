import Foundation

class User {
    var name: String?
    var age: Int?

    init(name: String, age: Int?) {
        self.name = name
        self.age = age
    }

    func greet() -> String {
        let displayName = name ?? "Anonymous"
        return "Hello, \(displayName)!"
    }
}

func main() {
    let users = [User(name: "Alice", age: 30), User(name: "Bob", age: nil)]

    for user in users {
        let message = user.greet()
        print(message)
    }

    let optionalValue: Int? = nil
    let result = optionalValue ?? 42
    print("Result: \(result)")
}

main()
