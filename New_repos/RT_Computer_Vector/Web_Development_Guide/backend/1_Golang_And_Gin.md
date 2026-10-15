### What is Go

**Go (Golang)** is a statically typed, compiled language created by Google, designed for simplicity, fast compilation, and built-in concurrency support. It's widely used for backend services, CLIs, and cloud infrastructure tools (Docker and Kubernetes are both written in Go).

**Why people choose Go:**
- Compiles to a single static binary — no runtime/dependencies needed to deploy
- Built-in concurrency via goroutines (lightweight threads)
- Fast compilation and execution
- Simple syntax with a small, opinionated feature set
- Strong standard library, especially for networking/web

---

## PART 1: The Go Language

---

### Installation & Setup

```bash
# Linux/macOS
curl -OL https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

**What's happening here:** you're downloading Go's official binary release, extracting it into `/usr/local` (the conventional place for manually-installed software on Linux), and then adding its `bin` folder to your `PATH` so the `go` command works from any terminal. The `export PATH` line only lasts for your current terminal session — add it to `~/.bashrc` or `~/.zshrc` to make it permanent.

**Verify:**
```bash
go version
```
This just confirms the install worked and shows which version you have — useful since Go's behavior/features can vary between versions.

**Windows:** download the installer from `https://go.dev/dl/` and run it — it sets up PATH automatically, so there's no manual export step like on Linux/macOS.

---

### Your First Program

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

**Run it:**
```bash
go run main.go
```
This compiles and runs the program in one step — good for quick testing, but it doesn't leave behind a reusable binary.

**Build a binary:**
```bash
go build main.go
./main
```
This compiles your code into a standalone executable file (`main`) that you can run directly, copy to another machine, or ship — this is the actual deployment artifact for a Go program, unlike `go run` which is just for development.

#### Key Concepts Here
- **`package main`** — every executable Go program needs a `main` package. This tells the compiler "this isn't a library, this is something that should produce a runnable program."
- **`import`** — brings in other packages (here, `fmt` for formatted I/O). Go's standard library is organized into small, focused packages, and you only pay for what you import.
- **`func main()`** — the entry point of the program. When you run the compiled binary, execution starts here — every Go executable needs exactly one of these.

---

### Variables & Constants

```go
// Explicit type
var name string = "Alice"

// Type inferred
var age = 25

// Short declaration (only inside functions)
city := "New York"

// Multiple variables at once
var x, y int = 10, 20

// Constants
const Pi = 3.14159
```

**Why three different ways to declare a variable exist:** `var name string = "Alice"` is the most explicit — useful at the package level or when you want to be clear about the type regardless of the initial value. `var age = 25` lets Go infer the type from the value (here, `int`) — less typing, same safety. `city := "New York"` is Go's shorthand and by far the most common style you'll see inside functions, but notice it **only works inside functions** — you can't use `:=` at the package/global level.

**Constants** (`const`) are values that can never change after being set — the compiler enforces this, catching accidental reassignment as an error rather than a runtime bug.

**Zero values** — variables declared without a value get a default:
```go
var i int     // 0
var s string  // ""
var b bool    // false
```
This is a deliberate Go design choice: unlike some languages where an uninitialized variable holds garbage/undefined memory, Go guarantees every variable starts in a predictable, usable state — so `var i int` is always safely usable as `0`, never a crash waiting to happen.

---

### Basic Data Types

| Type | Examples |
|---|---|
| `int`, `int8`, `int16`, `int32`, `int64` | Whole numbers |
| `uint`, `uint8` (`byte`) | Unsigned whole numbers |
| `float32`, `float64` | Decimal numbers |
| `string` | Text |
| `bool` | `true` / `false` |
| `rune` | A single Unicode character (alias for `int32`) |

**Why so many integer sizes exist:** the number after `int` (8/16/32/64) tells you how many bits of memory it uses, which caps the range of values it can hold and affects performance/memory usage. Plain `int` automatically picks the most efficient size for your machine (usually 64-bit on modern systems), so you'd only reach for `int8`/`int32`/etc. when memory layout matters — e.g. reading binary file formats or working with hardware-level data.

**`rune` vs `byte`:** a `byte` is one raw 8-bit value (good for ASCII text or binary data). A `rune` represents a full Unicode character, which can take more than one byte — this distinction matters the moment you deal with non-English text.

```go
var temperature float64 = 98.6
var isActive bool = true
var initial rune = 'A'
```

---

### Operators

```go
// Arithmetic
a + b, a - b, a * b, a / b, a % b

// Comparison
a == b, a != b, a < b, a > b, a <= b, a >= b

// Logical
a && b, a || b, !a
```

These behave the same as in most C-family languages: arithmetic operators do math (`%` is remainder/modulo), comparison operators return a `bool`, and logical operators (`&&` AND, `||` OR, `!` NOT) combine boolean expressions. Nothing Go-specific here — worth knowing mainly because `if`/`for` conditions in Go rely entirely on these.

---

### Control Structures

#### If / Else
```go
if age >= 18 {
    fmt.Println("Adult")
} else if age >= 13 {
    fmt.Println("Teenager")
} else {
    fmt.Println("Child")
}
```
**One Go quirk worth noting:** conditions are never wrapped in parentheses (`if age >= 18`, not `if (age >= 18)`) — Go's syntax intentionally drops them since they're unnecessary noise.

**If with a short statement:**
```go
if score := getScore(); score > 90 {
    fmt.Println("Excellent")
}
```
This pattern lets you declare a variable (`score`) that only exists for the duration of the `if`/`else` block — useful for keeping temporary values from leaking into the rest of your function.

#### Switch
```go
switch day {
case "Mon", "Tue", "Wed", "Thu", "Fri":
    fmt.Println("Weekday")
case "Sat", "Sun":
    fmt.Println("Weekend")
default:
    fmt.Println("Invalid day")
}
```
**Difference from C/Java's `switch`:** Go's `case` doesn't "fall through" to the next case automatically — each case implicitly breaks after running, which avoids one of the most common bugs in C-style switch statements. You can also match multiple values in one `case` line (like `"Mon", "Tue", ...` above), which most other languages require separate stacked `case` lines for.

#### For Loops (Go's only loop keyword)
```go
// Classic
for i := 0; i < 5; i++ {
    fmt.Println(i)
}

// While-style
i := 0
for i < 5 {
    i++
}

// Infinite loop
for {
    break // stop manually
}

// Range loop (arrays, slices, maps, strings)
for index, value := range mySlice {
    fmt.Println(index, value)
}
```
**Go deliberately has no `while` or `do-while` keyword** — `for` covers every looping need by simply omitting parts of its syntax (no init/condition/post = infinite loop; only a condition = behaves like `while`). The `range` form is the one you'll use constantly for iterating collections — it automatically gives you both the index/key and the value on each iteration.

---

### Functions

```go
func add(a int, b int) int {
    return a + b
}

// Shorter parameter syntax (same type)
func add(a, b int) int {
    return a + b
}
```
When consecutive parameters share the same type, Go lets you skip repeating the type name — both versions above behave identically.

#### Multiple Return Values (very common in Go)
```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("cannot divide by zero")
    }
    return a / b, nil
}

result, err := divide(10, 2)
if err != nil {
    fmt.Println("Error:", err)
}
```
**Why this matters:** Go has no exceptions (see the Error Handling section below), so returning a result *and* an error side-by-side is the language's standard way of saying "here's your value, and here's whether something went wrong." Checking `if err != nil` immediately after almost every function call becomes a very familiar Go pattern.

#### Named Return Values
```go
func split(sum int) (x, y int) {
    x = sum * 4 / 9
    y = sum - x
    return // "naked" return
}
```
Naming the return values up front (`x, y int`) both documents what the function returns and lets you use a bare `return` at the end — Go automatically returns whatever `x` and `y` currently hold. Useful for short functions, but can hurt readability in longer ones since it's less obvious what's being returned at a glance.

#### Variadic Functions
```go
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

sum(1, 2, 3, 4) // 10
```
The `...int` syntax means "accept any number of `int` arguments," which Go packages into a slice (`nums`) inside the function. This is how functions like `fmt.Println` can accept any number of arguments.

#### Closures
```go
func counter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

next := counter()
next() // 1
next() // 2
```
A closure is a function that "remembers" variables from where it was created, even after that outer function has finished running. Here, the inner function keeps a private reference to `count` — each call to `next()` increments the *same* `count`, rather than resetting it, because the closure captured that variable rather than copying it.

---

### Arrays & Slices

#### Arrays (fixed size)
```go
var numbers [5]int
numbers[0] = 10
fruits := [3]string{"apple", "banana", "cherry"}
```
An array's size is part of its type — a `[5]int` and a `[10]int` are considered different types entirely, and the size can never change after creation. This rigidity is why arrays are rarely used directly in everyday Go code.

#### Slices (dynamic, more commonly used)
```go
fruits := []string{"apple", "banana", "cherry"}
fruits = append(fruits, "date")

// Slicing
sub := fruits[1:3] // ["banana", "cherry"]

// Make a slice with initial length/capacity
nums := make([]int, 5)      // length 5, all zero
nums := make([]int, 0, 10)  // length 0, capacity 10
```
**A slice is a flexible, resizable view over an underlying array**, which is why it's used almost everywhere in real Go code instead of a plain array. `append()` adds an element, growing the underlying storage automatically when needed. The `make([]int, 0, 10)` form pre-allocates room for 10 elements without actually populating them — a performance optimization when you know roughly how much you'll be adding, avoiding repeated reallocation as the slice grows.

---

### Maps

```go
ages := map[string]int{
    "Alice": 30,
    "Bob":   25,
}

ages["Charlie"] = 35        // add
value, exists := ages["Bob"] // check existence
delete(ages, "Alice")        // remove

for key, val := range ages {
    fmt.Println(key, val)
}
```
A map is Go's key-value lookup structure (similar to a dictionary/hashmap in other languages). The two-value form `value, exists := ages["Bob"]` matters because looking up a missing key doesn't error — it silently returns the zero value for that type, so `exists` (a `bool`) is how you tell "key wasn't found" apart from "key exists and its value happens to be zero." **Note:** map iteration order in Go is intentionally randomized each run — never rely on it being consistent.

---

### Structs

```go
type Person struct {
    Name string
    Age  int
}

p := Person{Name: "Alice", Age: 30}
fmt.Println(p.Name)

// Anonymous struct
point := struct {
    X, Y int
}{X: 1, Y: 2}
```
A struct groups related fields together into one custom type — Go's equivalent of a class's data (without inheritance). Anonymous structs (no `type` name) are handy for quick, throwaway groupings you don't need to reuse elsewhere.

#### Methods on Structs
```go
func (p Person) Greet() string {
    return "Hello, " + p.Name
}

p.Greet()
```
The `(p Person)` before the function name is called a **receiver** — it's what attaches this function to the `Person` type as a method, letting you call it as `p.Greet()` instead of `Greet(p)`. This is Go's alternative to classes with methods, without needing actual class syntax.

**Pointer receiver** (needed to modify the struct):
```go
func (p *Person) Birthday() {
    p.Age++
}

p.Birthday() // p.Age is now incremented
```
Using `Person` (value) as a receiver gives the method its own *copy* of the struct — changes inside don't affect the original. Using `*Person` (pointer) gives the method direct access to the original struct's memory, so modifications actually stick. As a rule of thumb: use a pointer receiver whenever the method needs to change the struct's fields.

---

### Pointers

```go
x := 10
p := &x        // p holds the address of x
fmt.Println(*p) // dereference: prints 10

*p = 20        // modifies x through the pointer
fmt.Println(x)  // 20
```
`&x` gets the memory address where `x` lives, and `*p` "follows" that address to read or write the actual value stored there. This is exactly why the pointer-receiver distinction above matters — a pointer lets you reach back and change the original data instead of working on a disposable copy.

**Why pointers matter in Go:** passing a struct by value copies it; passing a pointer lets a function modify the original. For large structs, pointers are also more efficient — you're passing a small memory address around instead of duplicating a big chunk of data every time.

---

### Interfaces

**An interface defines behavior**, not data — any type that implements its methods automatically satisfies it (no explicit "implements" keyword needed).

```go
type Shape interface {
    Area() float64
}

type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return 3.14159 * c.Radius * c.Radius
}

func printArea(s Shape) {
    fmt.Println(s.Area())
}

printArea(Circle{Radius: 5}) // works because Circle implements Area()
```
This is called **"structural typing"** or "duck typing" at compile time — Go doesn't require you to explicitly declare "Circle implements Shape" anywhere. As long as `Circle` has an `Area() float64` method, it automatically satisfies the `Shape` interface and can be passed anywhere a `Shape` is expected. This makes interfaces in Go lightweight and encourages small, single-purpose interfaces (often just one or two methods).

**The empty interface** (`interface{}` or `any` in modern Go) can hold any type:
```go
func printAnything(val interface{}) {
    fmt.Println(val)
}
```
Since every type in Go satisfies "no methods required," `interface{}` (or its modern alias `any`) acts as an escape hatch for accepting literally any value — useful for generic-ish utility functions, though modern Go often prefers actual generics for type safety instead.

---

### Error Handling

**Go doesn't use exceptions** — errors are just regular values, usually returned as the last value from a function.

```go
func readFile(path string) ([]byte, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    return data, nil
}

data, err := readFile("file.txt")
if err != nil {
    log.Fatal(err)
}
```
This is one of Go's most opinionated design decisions: instead of `try`/`catch` blocks that can silently skip over error-prone code, Go forces you to explicitly check `err` after nearly every operation that could fail. It's more verbose, but the tradeoff is that error handling is impossible to accidentally ignore — it's right there in the code path, not hidden in a stack unwind.

**Creating custom errors:**
```go
import "errors"

err := errors.New("something went wrong")

// Or with formatting
err := fmt.Errorf("failed to process %s: %w", filename, originalErr)
```
`errors.New` creates a simple static error message. `fmt.Errorf` with `%w` does the same but also "wraps" an existing error inside it — preserving the original error so calling code can later unwrap it and check what specifically went wrong, while still adding your own context on top.

---

### Defer, Panic, Recover

```go
func readFile() {
    file, _ := os.Open("data.txt")
    defer file.Close() // runs when the function returns, no matter what
    // ... use file ...
}
```
`defer` schedules a function call to run right before the enclosing function returns — regardless of *how* it returns (normal completion, early `return`, or even a panic). This is Go's clean way of guaranteeing cleanup code (closing files, unlocking mutexes, etc.) always runs, without needing a `finally` block like other languages.

**Panic/Recover** (Go's rough equivalent of try/catch, used sparingly):
```go
func safeDivide(a, b int) (result int) {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from panic:", r)
            result = 0
        }
    }()
    return a / b // panics if b == 0
}
```
A `panic` is Go's mechanism for genuinely unrecoverable, unexpected situations (like dividing by zero above) — it immediately stops normal execution and starts unwinding the call stack. `recover()`, only useful inside a `defer`, can catch that panic and let the program continue instead of crashing entirely. **Convention matters here:** Go code is expected to use regular `error` returns for anything expected/handleable, and reserve `panic`/`recover` for truly exceptional, programmer-error-level situations — not as a general substitute for error handling.

---

### Goroutines (Concurrency)

**A goroutine is a lightweight thread** managed by the Go runtime — you can spawn thousands of them cheaply.

```go
func sayHello() {
    fmt.Println("Hello")
}

go sayHello() // runs concurrently, doesn't block
```
Adding the `go` keyword before a function call runs it *concurrently* — the calling code immediately continues to the next line without waiting for `sayHello()` to finish. Goroutines are much cheaper than OS threads (a few KB of memory each vs megabytes for a real thread), which is why Go programs can casually spin up thousands of them where other languages would need careful thread-pool management.

#### Channels — communication between goroutines

```go
ch := make(chan string)

go func() {
    ch <- "message from goroutine" // send
}()

msg := <-ch // receive (blocks until a value arrives)
fmt.Println(msg)
```
Since goroutines run independently, you need a safe way for them to pass data to each other without race conditions — that's what a **channel** is for. `ch <- value` sends data into the channel, `<-ch` receives it. By default, both operations **block** — the sender waits until someone's ready to receive, and vice versa — which is what makes channels useful for synchronizing goroutines, not just passing data.

#### Buffered Channels
```go
ch := make(chan int, 3) // can hold 3 values without a receiver ready
ch <- 1
ch <- 2
ch <- 3
```
A buffered channel can hold a fixed number of values *without* immediately needing a receiver — sends only block once the buffer is full. This is useful when you want to decouple the speed of a sender from a receiver slightly, rather than forcing them to stay in lockstep.

#### WaitGroups (waiting for multiple goroutines to finish)
```go
var wg sync.WaitGroup

for i := 0; i < 5; i++ {
    wg.Add(1)
    go func(n int) {
        defer wg.Done()
        fmt.Println("Worker", n)
    }(i)
}

wg.Wait() // blocks until all 5 goroutines call Done()
```
`WaitGroup` solves a common problem: how does your main function know when a bunch of goroutines it launched have all actually finished? `wg.Add(1)` increments a counter before launching each goroutine, `wg.Done()` (called via `defer`, so it always runs) decrements it when that goroutine finishes, and `wg.Wait()` blocks the calling code until that counter hits zero.

#### Select (waiting on multiple channels)
```go
select {
case msg1 := <-ch1:
    fmt.Println(msg1)
case msg2 := <-ch2:
    fmt.Println(msg2)
default:
    fmt.Println("no message ready")
}
```
`select` lets a goroutine wait on *multiple* channels at once, and proceeds with whichever one becomes ready first — similar in spirit to a `switch`, but for channel operations instead of values. The optional `default` case makes the whole thing non-blocking: if no channel has a value ready immediately, it runs `default` instead of waiting.

---

### Packages & Modules

**Every Go project is organized into packages.** A module is a collection of packages with dependency versioning.

```bash
go mod init github.com/yourname/myproject
```
This creates a `go.mod` file tracking your module name and dependencies.

**What `go.mod` actually does:** it's Go's equivalent of `package.json` (Node) or `requirements.txt` (Python) — it records your project's module path (used for imports) and pins the exact versions of every external package you depend on, so builds are reproducible across machines.

**Installing a dependency:**
```bash
go get github.com/gin-gonic/gin
```
This downloads the package, adds it to `go.mod`, and records its exact resolved version (plus a checksum) in a companion `go.sum` file for integrity verification.

**Project structure example:**
```
myproject/
├── go.mod
├── main.go
└── handlers/
    └── user.go
```

**Importing your own package:**
```go
import "github.com/yourname/myproject/handlers"
```
Notice you import your own subfolders using the *full module path* (from `go.mod`) plus the folder name — not a relative path like `./handlers`. This is a common early trip-up for people coming from other languages.

---

## PART 2: The Gin Framework (Backend Web Development)

---

### What is Gin

**Gin is a fast, minimalist HTTP web framework for Go** — the most popular choice for building REST APIs in Go, known for its performance (built on a fast HTTP router) and simple, expressive syntax.

---

### Installation

```bash
go mod init myapi
go get github.com/gin-gonic/gin
```
Same two commands as any Go dependency: initialize the module, then fetch the package. Gin isn't part of Go's standard library — it's a third-party package, hence needing `go get`.

---

### Basic Server

```go
package main

import "github.com/gin-gonic/gin"

func main() {
    router := gin.Default() // includes logging + recovery middleware

    router.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })

    router.Run(":8080") // listens on localhost:8080
}
```
`gin.Default()` gives you a router pre-configured with two useful pieces of middleware: request logging (prints each incoming request to your console) and panic recovery (so a bug in one handler doesn't crash your entire server). `c *gin.Context` is the object every handler receives — it carries the incoming request data and gives you methods to write the response. `router.Run(":8080")` starts the actual HTTP server and blocks, listening for connections on port 8080.

```bash
go run main.go
```

---

### Routing

```go
router.GET("/users", getUsers)
router.POST("/users", createUser)
router.PUT("/users/:id", updateUser)
router.DELETE("/users/:id", deleteUser)
router.PATCH("/users/:id", partialUpdateUser)
```
Each line maps an HTTP method + URL pattern to a handler function — this is the core of how a REST API is structured: `GET` for reading, `POST` for creating, `PUT`/`PATCH` for updating, `DELETE` for removing.

#### Route Parameters
```go
router.GET("/users/:id", func(c *gin.Context) {
    id := c.Param("id")
    c.JSON(200, gin.H{"user_id": id})
})
```
The `:id` in the route pattern is a placeholder that matches any value in that URL segment (e.g. `/users/42` → `id = "42"`). `c.Param("id")` retrieves whatever was actually in that position of the incoming request's URL.

#### Query Parameters
```go
// GET /search?q=golang&page=2
router.GET("/search", func(c *gin.Context) {
    query := c.Query("q")                     // "golang"
    page := c.DefaultQuery("page", "1")       // "2", or "1" if missing
    c.JSON(200, gin.H{"query": query, "page": page})
})
```
Query parameters are the `?key=value` pairs after a URL. `c.Query()` reads one, returning an empty string if it's missing; `c.DefaultQuery()` does the same but lets you specify a fallback value instead of getting an empty string.

---

### Handling Request Bodies (JSON Binding)

```go
type CreateUserInput struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
    Age   int    `json:"age" binding:"gte=0,lte=130"`
}

func createUser(c *gin.Context) {
    var input CreateUserInput

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // input.Name, input.Email, input.Age are now populated and validated
    c.JSON(201, gin.H{"message": "user created", "user": input})
}
```
The text in backticks after each struct field (called **struct tags**) tells Gin two things: how to map incoming JSON fields to this struct (`json:"name"`), and what validation rules to enforce (`binding:"required"`). `c.ShouldBindJSON(&input)` does all the heavy lifting in one call — it parses the incoming request body as JSON, populates your struct, and runs every validation rule, returning an error if anything fails (missing required field, invalid email format, age out of range, etc.), so you don't have to write that validation logic by hand.

**Common validation tags:** `required`, `email`, `min=`, `max=`, `gte=`, `lte=`, `len=`.

---

### Response Methods

```go
c.JSON(200, gin.H{"key": "value"})       // JSON response
c.String(200, "Plain text response")      // plain text
c.XML(200, gin.H{"key": "value"})         // XML response
c.Data(200, "image/png", imageBytes)      // raw bytes
c.File("./file.pdf")                      // serve a file
c.Redirect(302, "/new-location")          // redirect
```
Each of these sets the correct HTTP status code and `Content-Type` header automatically, so you don't have to manage those manually — `c.JSON` sends `application/json`, `c.String` sends plain text, and so on. `gin.H` is just Gin's shorthand for `map[string]interface{}`, a convenient way to build ad-hoc JSON objects without defining a struct for every response.

---

### Route Groups

**Group related routes** under a shared prefix and middleware:

```go
api := router.Group("/api/v1")
{
    api.GET("/users", getUsers)
    api.POST("/users", createUser)

    admin := api.Group("/admin")
    admin.Use(AuthMiddleware()) // only applies within this group
    {
        admin.DELETE("/users/:id", deleteUser)
    }
}
```
Grouping avoids repeating the same URL prefix (`/api/v1`) on every single route, and lets you attach middleware (like authentication) to only a *subset* of your routes rather than the whole app. Nesting groups (like `admin` inside `api`) is how you build up more specific rules layer by layer — here, only routes under `/api/v1/admin` require the `AuthMiddleware`, while the rest of `/api/v1` doesn't.

---

### Middleware

**Middleware runs before/after your handler** — used for logging, auth, CORS, rate limiting, etc.

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"error": "unauthorized"})
            c.Abort() // stops the request from reaching the handler
            return
        }
        c.Next() // continue to the next middleware/handler
    }
}
```
Middleware is just a function that runs *before* your actual route handler gets a chance to — this example checks for an `Authorization` header before letting the request through. `c.Abort()` stops the request pipeline right there (the actual handler never runs), while `c.Next()` tells Gin "this middleware is done, move on to whatever's next" — either another middleware, or the final route handler.

**Applying middleware:**
```go
router.Use(AuthMiddleware())              // global, applies to all routes
router.GET("/admin", AuthMiddleware(), adminHandler) // single route
```
`router.Use()` attaches middleware globally, so it runs on every incoming request. Passing it directly as an extra argument to a single route (like the second example) scopes it to just that one route instead.

---

### CORS (Cross-Origin Requests)

```bash
go get github.com/gin-contrib/cors
```

```go
import "github.com/gin-contrib/cors"

router.Use(cors.Default()) // allows all origins (dev only)

// Production-safe config
router.Use(cors.New(cors.Config{
    AllowOrigins: []string{"https://yourfrontend.com"},
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
}))
```
**Why CORS matters:** browsers block a webpage on one domain from calling an API on a different domain unless that API explicitly allows it — this is a browser security feature, not something your Go server does by default. `cors.Default()` opens this up to any origin, which is fine for local testing but a real security risk in production (any website could call your API from a user's browser). The `cors.New(...)` config restricts this to only the specific frontend domain(s) and HTTP methods you actually want to allow.

---

### File Uploads

```go
router.POST("/upload", func(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    c.SaveUploadedFile(file, "./uploads/"+file.Filename)
    c.JSON(200, gin.H{"message": "uploaded", "filename": file.Filename})
})
```
`c.FormFile("file")` pulls an uploaded file out of a `multipart/form-data` request (the standard way browsers submit file uploads), matching whatever field name the client used (here, `"file"`). `c.SaveUploadedFile` then writes it to disk at the path you specify — you're responsible for choosing where and under what filename it gets stored.
)
---

### Serving Static Files

```go
router.Static("/static", "./public")     // serve a whole folder
router.StaticFile("/favicon.ico", "./favicon.ico")
```
`router.Static` maps a URL prefix directly to a folder on disk, so a request to `/static/logo.png` serves `./public/logo.png` automatically — handy for serving a built frontend, images, or other assets without writing a handler for each file. `router.StaticFile` does the same for a single specific file.

---

### Error Handling Pattern

```go
func getUser(c *gin.Context) {
    id := c.Param("id")
    user, err := findUserByID(id)

    if err != nil {
        c.JSON(404, gin.H{"error": "user not found"})
        return
    }

    c.JSON(200, user)
}
```
This follows the same Go convention as regular error handling — check `err`, and if it's not `nil`, respond accordingly (here, a 404) and `return` immediately so the rest of the function doesn't run with bad/missing data.

**Centralized error handling with middleware:**
```go
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        if len(c.Errors) > 0 {
            c.JSON(500, gin.H{"error": c.Errors[0].Error()})
        }
    }
}
```
Instead of repeating error-response logic in every handler, this middleware runs `c.Next()` first (letting the actual handler execute), and then checks afterward if the handler added anything to `c.Errors` (via `c.Error(err)` inside a handler) — centralizing your error-response formatting in one place instead of scattering `c.JSON(500, ...)` calls throughout your codebase.

---

### Connecting to a Database (Example: PostgreSQL with GORM)

```bash
go get gorm.io/gorm
go get gorm.io/driver/postgres
```
**GORM** is Go's most popular ORM (Object-Relational Mapper) — it lets you interact with a database using Go structs and method calls instead of writing raw SQL by hand.

```go
import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type User struct {
    ID    uint   `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

var db *gorm.DB

func initDB() {
    dsn := "host=localhost user=postgres password=secret dbname=mydb port=5432"
    db, _ = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    db.AutoMigrate(&User{})
}

func getUsers(c *gin.Context) {
    var users []User
    db.Find(&users)
    c.JSON(200, users)
}
```
The `dsn` (Data Source Name) string tells GORM how to connect to your Postgres instance. `db.AutoMigrate(&User{})` inspects your `User` struct and automatically creates/updates the corresponding database table to match its fields — handy in development, though many teams switch to explicit migration files for production to keep schema changes reviewable. `db.Find(&users)` fetches every row from the `users` table and populates the `users` slice with the results, converting database rows into Go structs automatically.

---

### Typical Project Structure

```
myapi/
├── go.mod
├── main.go
├── handlers/
│   ├── user.go
│   └── auth.go
├── middleware/
│   └── auth.go
├── models/
│   └── user.go
├── routes/
│   └── routes.go
└── config/
    └── database.go
```
This layout separates concerns by responsibility: `handlers/` holds the actual request-handling logic, `middleware/` holds reusable request-pipeline logic (auth, logging), `models/` defines your data structures, `routes/` wires URLs to handlers, and `config/` holds setup code like database connections. This isn't enforced by Go itself — it's just a common convention that keeps larger projects navigable.

**`routes/routes.go` example:**
```go
package routes

import (
    "github.com/gin-gonic/gin"
    "myapi/handlers"
)

func SetupRoutes(router *gin.Engine) {
    api := router.Group("/api")
    {
        api.GET("/users", handlers.GetUsers)
        api.POST("/users", handlers.CreateUser)
    }
}
```
Keeping route definitions in their own file/function (`SetupRoutes`) means `main.go` stays short — it just calls `routes.SetupRoutes(router)` — while the actual URL-to-handler mapping lives somewhere dedicated and easy to scan.

---

### Quick Reference — Common Gin Methods

| Method | Purpose |
|---|---|
| `c.Param("id")` | Get URL path parameter |
| `c.Query("key")` | Get query string parameter |
| `c.ShouldBindJSON(&obj)` | Parse + validate JSON request body |
| `c.JSON(code, data)` | Send JSON response |
| `c.Abort()` | Stop request from continuing to next handler |
| `c.Next()` | Continue to next middleware/handler |
| `c.Set("key", val)` | Store data in the request context |
| `c.Get("key")` | Retrieve data from the request context |
| `router.Group("/prefix")` | Group routes under a shared path |
| `router.Use(middleware)` | Apply middleware globally |

**On `c.Set` / `c.Get`:** these let one middleware pass data along to the next middleware or the final handler within the same request — for example, an auth middleware could `c.Set("userID", 42)` after validating a token, and the handler further down the chain could `c.Get("userID")` to retrieve it, without needing a shared global variable.

---

### Full Mental Model

```
Go language          → syntax, types, structs, goroutines, error handling
     ↓
go mod / go get       → manage dependencies (like Gin itself)
     ↓
Gin framework          → routes + handlers + middleware
     ↓
c.ShouldBindJSON        → parse & validate incoming requests
     ↓
Business logic          → your own Go functions/structs
     ↓
Database (GORM, etc.)  → persist data
     ↓
c.JSON(...)             → send the response back
```

---
 
## PART 3: Building a Complete API in Go with Gin
 
This section ties everything above together into one small, runnable REST API — a simple "Todo" API with in-memory storage (no database needed to try it out).
 
### 1. Project Setup
 
```bash
mkdir todo-api && cd todo-api
go mod init todo-api
go get github.com/gin-gonic/gin
```
 
### 2. Project Structure
 
```
todo-api/
├── go.mod
├── main.go
├── models/
│   └── todo.go
└── handlers/
    └── todo.go
```
 
### 3. The Model
 
`models/todo.go`
```go
package models
 
type Todo struct {
    ID        int    `json:"id"`
    Title     string `json:"title" binding:"required"`
    Completed bool   `json:"completed"`
}
```
This struct is both the data shape stored in memory *and* the JSON shape sent to/from clients — the `json:` tags control serialization, and `binding:"required"` means `Title` must be present when a client creates a new todo.
 
### 4. The Handlers
 
`handlers/todo.go`
```go
package handlers
 
import (
    "net/http"
    "strconv"
 
    "github.com/gin-gonic/gin"
    "todo-api/models"
)
 
// In-memory "database" — a slice, plus a counter for generating IDs.
var todos = []models.Todo{}
var nextID = 1
 
// GET /todos
func GetTodos(c *gin.Context) {
    c.JSON(http.StatusOK, todos)
}
 
// GET /todos/:id
func GetTodo(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }
 
    for _, t := range todos {
        if t.ID == id {
            c.JSON(http.StatusOK, t)
            return
        }
    }
    c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
}
 
// POST /todos
func CreateTodo(c *gin.Context) {
    var input models.Todo
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
 
    input.ID = nextID
    nextID++
    todos = append(todos, input)
 
    c.JSON(http.StatusCreated, input)
}
 
// PUT /todos/:id
func UpdateTodo(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }
 
    var input models.Todo
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
 
    for i, t := range todos {
        if t.ID == id {
            input.ID = id
            todos[i] = input
            c.JSON(http.StatusOK, input)
            return
        }
    }
    c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
}
 
// DELETE /todos/:id
func DeleteTodo(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }
 
    for i, t := range todos {
        if t.ID == id {
            todos = append(todos[:i], todos[i+1:]...)
            c.JSON(http.StatusOK, gin.H{"message": "deleted"})
            return
        }
    }
    c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
}
```
Each handler follows the same shape: pull data out of the request (`c.Param`, `c.ShouldBindJSON`), do something with the in-memory `todos` slice, and send a JSON response with an appropriate status code. `strconv.Atoi` converts the `:id` path parameter (always a string) into an `int` so it can be compared against `t.ID`. Swapping this in-memory slice for a real database later (GORM, as shown earlier in this doc) wouldn't change any handler's signature — only what happens inside the function body.
 
### 5. Wiring It Up
 
`main.go`
```go
package main
 
import (
    "github.com/gin-gonic/gin"
    "todo-api/handlers"
)
 
func main() {
    router := gin.Default()
 
    api := router.Group("/api/v1")
    {
        api.GET("/todos", handlers.GetTodos)
        api.GET("/todos/:id", handlers.GetTodo)
        api.POST("/todos", handlers.CreateTodo)
        api.PUT("/todos/:id", handlers.UpdateTodo)
        api.DELETE("/todos/:id", handlers.DeleteTodo)
    }
 
    router.Run(":8080")
}
```
`main.go` stays intentionally thin — it just creates the router, groups the routes under `/api/v1`, and starts listening. All the actual logic lives in `handlers/`, and the data shape lives in `models/`, matching the project-structure convention from earlier in this doc.
 
### 6. Run It
 
```bash
go run main.go
```
 
### 7. Try It (with curl)
 
```bash
# Create a todo
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Learn Gin", "completed": false}'
 
# List all todos
curl http://localhost:8080/api/v1/todos
 
# Get one todo
curl http://localhost:8080/api/v1/todos/1
 
# Update a todo
curl -X PUT http://localhost:8080/api/v1/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Learn Gin", "completed": true}'
 
# Delete a todo
curl -X DELETE http://localhost:8080/api/v1/todos/1
```
 
### 8. Natural Next Steps
 
- Swap the in-memory `todos` slice for a real database using GORM (see the **Connecting to a Database** section above) — only `handlers/todo.go` needs to change, since `main.go` and the routes stay the same.
- Add the `AuthMiddleware()` from earlier to protect write routes (`POST`/`PUT`/`DELETE`) so only authenticated requests can modify data.
- Add the CORS middleware if this API will be called from a frontend running on a different domain/port.
- Add table-driven tests using Go's built-in `testing` package and `net/http/httptest` to simulate requests against your handlers without needing a real server running.
 

