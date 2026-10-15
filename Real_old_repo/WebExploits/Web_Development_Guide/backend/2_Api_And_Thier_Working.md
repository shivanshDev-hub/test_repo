# Guide to API (for dummies) & Golang nerds 

This is a plain-English map of the major API "styles" you'll run into — what each one actually is, when it's used, and how they compare. Think of this as the "zoomed out" view; your Go/Gin doc is the "zoomed in" view of building one specific style (REST) in one specific tool (Gin).
 
---
 
## What is an API, really?
 
**API = a way for one piece of software to ask another piece of software for something, or tell it to do something.**
 
A restaurant analogy: you (the client) don't walk into the kitchen (the server) yourself. You give your order to a waiter (the API), the waiter takes it to the kitchen, and brings back your food (the response). The API is just the agreed-upon *way* you communicate that request — the menu format, the language you use, whether you can ask follow-up questions, etc.
 
Different API "styles" are just different rules for how that conversation happens.
 
---
 
## 1. REST (Representational State Transfer)
 
**The default. Most common style on the web.**
 
- Data is organized around **resources** — a user, a post, an order — each with its own URL.
- You act on resources using standard HTTP methods: `GET` (read), `POST` (create), `PUT`/`PATCH` (update), `DELETE` (remove).
- Responses are usually JSON.
```
GET    /users/42        → get user 42
POST   /users            → create a new user
PUT    /users/42        → update user 42
DELETE /users/42        → delete user 42
```
 
**Good for:** most typical web/mobile apps, public APIs, CRUD-heavy systems.
**Downside:** you often over-fetch (get more data than you need) or under-fetch (need multiple requests to get everything), especially with deeply nested data.
 
*(This is what the Gin doc we built teaches — Gin is a REST-style framework by convention.)*
 
---
 
## 2. GraphQL
 
**The "ask for exactly what you want" style.**
 
- One single endpoint (usually `/graphql`), not many URLs.
- The client sends a **query** describing exactly which fields it wants, even across related data, and gets back exactly that — nothing more, nothing less.
```graphql
{
  user(id: "42") {
    name
    posts {
      title
    }
  }
}
```
 
That one request replaces what might be 2–3 separate REST calls (get user, then get their posts).
 
**Good for:** apps with complex, nested data needs — especially mobile apps where minimizing requests/bandwidth matters, or frontends that need different data shapes from the same backend.
**Downside:** more complex to set up on the server, caching is trickier than REST's built-in HTTP caching, and it has a steeper learning curve.
 
---
 
## 3. WebSockets
 
**The "keep the phone line open" style — for real-time, two-way communication.**
 
- Unlike REST/GraphQL (client asks, server answers, connection closes), a WebSocket opens **one persistent connection** that both sides can send messages over, anytime, in either direction.
- No need to keep "asking" for updates — the server can just push new data the instant it happens.
```javascript
const socket = new WebSocket("wss://example.com/chat");
socket.onmessage = (msg) => console.log(msg.data);
socket.send("hello");
```
 
**Good for:** live chat, multiplayer games, live dashboards/stock tickers, collaborative editing (like Google Docs).
**Downside:** the connection has to stay open, which is heavier on server resources and harder to scale than stateless REST calls.
 
---
 
## 4. gRPC
 
**The "fast, strict, machine-to-machine" style — built by Google.**
 
- Uses HTTP/2 and a binary format called **Protocol Buffers (protobuf)** instead of human-readable JSON — much faster and smaller over the wire.
- You define your API's shape upfront in a `.proto` file (a strict contract), and tooling auto-generates client/server code in whatever language you're using.
- Supports streaming (client streams data to server, server streams to client, or both at once).
**Good for:** service-to-service communication inside a company's backend (microservices talking to each other) — very popular in the Go ecosystem specifically.
**Downside:** not browser-friendly out of the box (needs a gateway to expose it to web/mobile clients), and you lose the human-readability of JSON — harder to debug by just looking at raw traffic.
 
---
 
## 5. SOAP (Simple Object Access Protocol)
 
**The older, strict, enterprise-y style.**
 
- Uses XML for everything (requests, responses), with a rigid, formally-defined contract (WSDL) describing every operation.
- Has built-in support for things like transactions and strong security standards, which is why it's stuck around in banking, healthcare, and government systems.
**Good for:** legacy enterprise systems, and industries with heavy compliance/security requirements.
**Downside:** verbose, heavy, and slow to work with compared to REST/JSON — rarely chosen for new projects today.
 
---
 
## 6. Webhooks
 
**Not really its own "style" — more like "REST, but the server calls you."**
 
- Normally *you* call an API. With a webhook, *you* register a URL with someone else's service, and *they* call **your** URL whenever something happens.
- Example: Stripe calls your `/webhooks/payment-success` endpoint the instant a payment completes, instead of you having to constantly ask "did it complete yet?"
**Good for:** event notifications from third-party services (payments, GitHub pushes, form submissions) without constantly polling.
 
---
 
## 7. WebRTC (bonus — peer-to-peer)
 
**Browser-to-browser, no server in the middle (mostly).**
 
- Lets two devices talk directly to each other — used for video calls, screen sharing, file transfer.
- Powers things like Zoom (web version), Google Meet, and peer-to-peer file sharing tools.
**Good for:** video/audio calls, real-time peer-to-peer data.
**Downside:** dealing with firewalls/NAT (getting two random devices on the internet to find each other) is genuinely tricky.
 
---
 
## Quick Comparison Table
 
| Style | Data format | Connection type | Best for |
|---|---|---|---|
| REST | JSON (usually) | Request → response, then closes | General-purpose CRUD APIs |
| GraphQL | JSON | Request → response, then closes | Complex/nested data, flexible clients |
| WebSockets | Text or binary | Persistent, two-way | Real-time updates (chat, live data) |
| gRPC | Protobuf (binary) | HTTP/2, can stream | Fast internal microservice communication |
| SOAP | XML | Request → response, then closes | Legacy enterprise systems |
| Webhooks | JSON (usually) | Server calls you, one-off | Event notifications from third parties |
| WebRTC | Binary/media streams | Peer-to-peer, persistent | Video/audio calls, P2P data |
 
---
 
## Which one should *you* actually learn first?
 
If you're learning backend development (which is what the Go/Gin doc is for): **REST first, always.** It's the most common, most transferable, and every other style is easier to understand once REST clicks. GraphQL and gRPC are natural "next steps" once REST feels comfortable — WebSockets is worth learning whenever you want to build something real-time.
 
---
 
## Videos to Learn These Thoroughly
 
**Best single overview (all styles compared, beginner-friendly):**
[You Can't Become a Backend Developer Without Understanding These (REST vs GraphQL vs gRPC)](https://www.youtube.com/watch?v=QojX9Nt4G1o) — explains REST, GraphQL, and gRPC clearly with practical framing (when/why you'd pick each).
 
**More in-depth, methodical comparison (part of a larger course):**
[REST API vs GraphQL vs gRPC - The Complete Guide](https://www.youtube.com/watch?v=n1mQBUv7o80) — a more structured, course-style breakdown of the three, good if you want more depth than a single overview video.
 
**Performance-focused comparison (benchmarks):**
[gRPC vs REST vs GraphQL: Comparison & Performance](https://www.youtube.com/watch?v=uH0SxYdsjv4) — useful once you understand the basics and want to see actual speed/overhead differences.
 
**For WebSockets specifically** (not covered in the above three):
Search "WebSockets crash course" on YouTube — Traversy Media and Web Dev Simplified both have solid, focused explainer videos on WebSockets specifically, since it's a fundamentally different mental model (persistent connection) from the request/response styles above.
 
None of these single videos cover *literally everything* (that's what this doc is for) — but between this guide and the videos above, you'll have a working understanding of every major API style used in real-world backend development today.
 
---

## Now the part for golang nerds 

This is meant to be the **last resource you ever need** for deciding which API style to use and how to actually build it in Go — whether that's a web project, an internal tool, or something in your cybersecurity/tooling career (scanners, agents, dashboards, integrations with threat-intel APIs, etc.).
 
Pairs with `Guide_to_API_for_dummies.md` (the conceptual overview). This doc is the **implementation + decision-making** companion — less theory, more "here's the code, here's when to use it."
 
---
 
# PART 0: The Decision Framework
 
Before writing any code, answer these questions. They'll point you to the right style almost every time.
 
### Question 1 — Who's on the other end?
 
| Consumer | Likely style |
|---|---|
| A web frontend / mobile app | REST or GraphQL |
| Another internal service (Go-to-Go, microservices) | gRPC |
| A third-party developer / public API | REST (most universally understood) |
| A CLI tool, agent, or script you control on both ends | REST, gRPC, or raw TCP/UDP |
| A browser needing live updates | WebSockets |
| Two peers directly (calls, file transfer) | WebRTC |
 
### Question 2 — Request/response, or ongoing stream?
 
- **One-off ask, one-off answer** → REST, GraphQL, gRPC (unary), SOAP
- **Server needs to keep pushing data without being asked again** → WebSockets, gRPC (streaming)
- **Third party needs to tell *you* when something happens, unprompted** → Webhooks
### Question 3 — How rigid vs. flexible does the data need to be?
 
- **Fixed, well-known shape, simple CRUD** → REST
- **Clients need different slices of deeply related data** → GraphQL
- **Strict typed contract between services, performance matters** → gRPC (protobuf enforces the contract)
### Question 4 — Performance & bandwidth constraints?
 
- **Human-facing, moderate traffic** → REST/JSON is fine, don't overthink it
- **High-volume, machine-to-machine, telemetry/agents reporting constantly** → gRPC (binary, smaller payloads, HTTP/2 multiplexing) or a custom protocol over raw TCP
- **Embedded/constrained devices** → consider raw TCP/UDP with a minimal custom protocol, or gRPC if the platform supports it
### Question 5 — Does it need to work in a browser?
 
- Browsers can natively speak: HTTP (REST, GraphQL), WebSockets, WebRTC
- Browsers **cannot** natively speak raw gRPC (needs grpc-web + a proxy) or raw TCP/UDP sockets
- If "someone's frontend needs to call this directly," cross gRPC and raw sockets off the list unless you're adding a gateway
### Question 6 — Public API, or internal/legacy enterprise?
 
- **New public or internal API** → REST (default) or GraphQL (if data-fetching flexibility matters)
- **Integrating with an old bank/gov/enterprise system** → you'll likely be *consuming* SOAP, not building it. Rarely a first choice for new projects.
### The One-Page Cheat Sheet
 
| Scenario | Use this |
|---|---|
| Standard web/mobile backend, CRUD operations | **REST** |
| Frontend needs flexible, nested data from one endpoint | **GraphQL** |
| Live dashboard, chat, real-time scan/log streaming | **WebSockets** |
| Internal microservice-to-microservice, Go-to-Go, high performance | **gRPC** |
| Agent/telemetry reporting at scale (many clients, frequent small messages) | **gRPC**, or raw TCP if you need full control |
| Notify your system when a third party has an event | **Webhooks** (you receive) |
| Notify other systems when *your* tool has an event | **Webhooks** (you send) |
| Peer-to-peer audio/video/data, no relay server | **WebRTC** |
| Legacy enterprise/bank/gov integration | **SOAP** (consume, don't build new) |
| Custom lightweight protocol for a scanner, agent, or internal tool you fully control | **Raw TCP/UDP** |
| Calling a third-party threat-intel/security API from a tool (VirusTotal, Shodan, etc.) | **REST client** |
 
**Real systems almost always combine styles** — e.g., REST for CRUD + WebSockets for live updates + Webhooks for third-party event notifications, all in the same project. Don't feel like you have to pick just one.
 
---
 
# PART 1: REST — Implementation in Go
 
You already have a full Gin-based walkthrough in your other doc (`golang_And_Gin.md`). Here's the **minimal, dependency-free version** using only the standard library — useful for small tools where you don't want to pull in a framework — plus the **client side**, which matters just as much for tool-building (calling *other* people's REST APIs).
 
### 1a. A REST server with zero dependencies (`net/http`)
 
```go
package main
 
import (
    "encoding/json"
    "net/http"
)
 
type Target struct {
    ID   string `json:"id"`
    Host string `json:"host"`
}
 
var targets = map[string]Target{}
 
func targetsHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        list := []Target{}
        for _, t := range targets {
            list = append(list, t)
        }
        json.NewEncoder(w).Encode(list)
 
    case http.MethodPost:
        var t Target
        if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        targets[t.ID] = t
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(t)
 
    default:
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    }
}
 
func main() {
    http.HandleFunc("/targets", targetsHandler)
    http.ListenAndServe(":8080", nil)
}
```
No router, no middleware, no third-party code — just `net/http`. Good for tiny internal tools or when you want zero dependencies (e.g., a single-binary CLI tool that also exposes a local API).
 
**When to reach for Gin instead:** the moment you need route parameters (`/targets/:id`), route groups, built-in validation, or middleware chains — `net/http` makes you hand-roll all of that yourself.
 
### 1b. A REST client in Go (consuming someone else's API)
 
This is the part most tutorials skip, but it's what you'll actually use constantly — calling VirusTotal, Shodan, AbuseIPDB, internal services, etc.
 
```go
package main
 
import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)
 
func getRequest(url, apiKey string) ([]byte, error) {
    client := &http.Client{Timeout: 10 * time.Second}
 
    req, err := http.NewRequest(http.MethodGet, url, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Authorization", "Bearer "+apiKey)
    req.Header.Set("Accept", "application/json")
 
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
 
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
    }
    return io.ReadAll(resp.Body)
}
 
func postRequest(url, apiKey string, payload any) ([]byte, error) {
    body, err := json.Marshal(payload)
    if err != nil {
        return nil, err
    }
 
    client := &http.Client{Timeout: 10 * time.Second}
    req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Authorization", "Bearer "+apiKey)
    req.Header.Set("Content-Type", "application/json")
 
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
 
    return io.ReadAll(resp.Body)
}
```
**Always set a `Timeout`** on the client — the default `http.Client{}` has no timeout, which means a hung server on the other end can hang your tool forever. This is a very common bug in hand-rolled Go tools. Always `defer resp.Body.Close()` too, or you'll leak connections over a long-running tool's lifetime.
 
### 1c. REST Authentication Patterns (Go code for each)
 
**API Key (simplest — good for internal tools/agents):**
```go
func APIKeyMiddleware(validKey string) gin.HandlerFunc {
    return func(c *gin.Context) {
        key := c.GetHeader("X-API-Key")
        if key != validKey {
            c.JSON(401, gin.H{"error": "invalid api key"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```
 
**JWT (stateless tokens — good for user sessions or authenticated agents):**
```bash
go get github.com/golang-jwt/jwt/v5
```
```go
import "github.com/golang-jwt/jwt/v5"
 
func GenerateToken(userID string, secret []byte) (string, error) {
    claims := jwt.MapClaims{
        "sub": userID,
        "exp": time.Now().Add(time.Hour * 24).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(secret)
}
 
func ValidateToken(tokenStr string, secret []byte) (jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
        return secret, nil
    })
    if err != nil || !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }
    return token.Claims.(jwt.MapClaims), nil
}
```
 
**mTLS (mutual TLS — the gold standard for machine-to-machine/agent auth, common in security tooling):**
```go
func startMTLSServer() {
    caCert, _ := os.ReadFile("ca.crt")
    caPool := x509.NewCertPool()
    caPool.AppendCertsFromPEM(caCert)
 
    server := &http.Server{
        Addr: ":8443",
        TLSConfig: &tls.Config{
            ClientCAs:  caPool,
            ClientAuth: tls.RequireAndVerifyClientCert, // reject any client without a valid cert
        },
    }
    server.ListenAndServeTLS("server.crt", "server.key")
}
```
This is the pattern you want for something like an agent-to-management-server setup: every client must present a certificate signed by your CA, or the connection is rejected outright — no password/token to leak, no risk of a stolen API key working from anywhere.
 
---
 
# PART 2: GraphQL — Implementation in Go
 
The standard library for Go GraphQL is **`gqlgen`** — it generates type-safe Go code from a schema you write, rather than you hand-writing resolvers from scratch.
 
```bash
go get github.com/99designs/gqlgen
```
 
### 2a. Define your schema (`schema.graphqls`)
 
```graphql
type Host {
  id: ID!
  address: String!
  openPorts: [Int!]!
}
 
type Query {
  host(id: ID!): Host
  hosts: [Host!]!
}
 
type Mutation {
  addHost(address: String!): Host!
}
```
This schema **is the contract** — gqlgen reads this file and generates matching Go interfaces/structs for you to implement. You never hand-write the parsing/validation logic; the schema defines what's allowed.
 
### 2b. Implement the resolver (the generated interface)
 
```go
type Resolver struct{}
 
var hosts = map[string]*Host{}
 
func (r *queryResolver) Host(ctx context.Context, id string) (*Host, error) {
    h, ok := hosts[id]
    if !ok {
        return nil, fmt.Errorf("host not found")
    }
    return h, nil
}
 
func (r *queryResolver) Hosts(ctx context.Context) ([]*Host, error) {
    list := []*Host{}
    for _, h := range hosts {
        list = append(list, h)
    }
    return list, nil
}
 
func (r *mutationResolver) AddHost(ctx context.Context, address string) (*Host, error) {
    h := &Host{ID: uuid.NewString(), Address: address}
    hosts[h.ID] = h
    return h, nil
}
```
Each field in your schema (`Query.host`, `Mutation.addHost`) maps to exactly one Go function — gqlgen wires the routing for you based on the schema, so you only write the "what happens when this field is requested" logic.
 
### 2c. Querying GraphQL from Go (client side)
 
GraphQL doesn't need a special client library — it's just a POST request with a JSON body containing your query string:
 
```go
func queryGraphQL(url, query string) ([]byte, error) {
    body, _ := json.Marshal(map[string]string{"query": query})
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    return io.ReadAll(resp.Body)
}
 
// Usage:
result, _ := queryGraphQL("https://api.example.com/graphql", `
    query {
        hosts {
            id
            address
            openPorts
        }
    }
`)
```
 
**Reality check for tool-building:** GraphQL is far more common on the *consuming* side (calling GitHub's or Shopify's GraphQL API) than the *building* side for internal tools. Most of your own tools/agents are better served by REST or gRPC — reach for GraphQL when you're building a platform with a real frontend team that needs flexible queries.
 
---
 
# PART 3: WebSockets — Implementation in Go
 
The standard library for this is **`gorilla/websocket`** — the de facto choice in the Go ecosystem.
 
```bash
go get github.com/gorilla/websocket
```
 
### 3a. WebSocket server
 
```go
package main
 
import (
    "log"
    "net/http"
 
    "github.com/gorilla/websocket"
)
 
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true }, // tighten this in production
}
 
func wsHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("upgrade failed:", err)
        return
    }
    defer conn.Close()
 
    for {
        msgType, msg, err := conn.ReadMessage()
        if err != nil {
            log.Println("client disconnected:", err)
            return
        }
        log.Printf("received: %s", msg)
 
        // Echo it back, or push your own data (e.g. live scan results)
        if err := conn.WriteMessage(msgType, msg); err != nil {
            log.Println("write failed:", err)
            return
        }
    }
}
 
func main() {
    http.HandleFunc("/ws", wsHandler)
    http.ListenAndServe(":8080", nil)
}
```
`upgrader.Upgrade` converts a normal HTTP request into a persistent WebSocket connection — this only happens once, at the start. After that, `ReadMessage`/`WriteMessage` run in a loop for the entire lifetime of the connection. **`CheckOrigin` returning `true` unconditionally is a security hole in production** — it means any website can open a WebSocket to your server from a user's browser. Lock this down to your actual known origins.
 
**Pushing data without waiting for a client message** (e.g., streaming live scan progress to a dashboard):
```go
func pushUpdates(conn *websocket.Conn, updates <-chan string) {
    for update := range updates {
        if err := conn.WriteMessage(websocket.TextMessage, []byte(update)); err != nil {
            return
        }
    }
}
```
This is the real power of WebSockets over REST polling — a goroutine can push to the connection the instant new data exists (e.g., a port scan finding a new open port), instead of the client having to repeatedly ask "anything new yet?"
 
### 3b. WebSocket client
 
```go
func connectWS(url string) {
    conn, _, err := websocket.DefaultDialer.Dial(url, nil)
    if err != nil {
        log.Fatal("dial failed:", err)
    }
    defer conn.Close()
 
    conn.WriteMessage(websocket.TextMessage, []byte("hello server"))
 
    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            log.Println("connection closed:", err)
            return
        }
        log.Printf("server says: %s", msg)
    }
}
```
 
---
 
# PART 4: gRPC — Implementation in Go
 
gRPC is a first-class citizen in Go (it was co-developed with Go in mind), and it's the go-to choice for Go-to-Go service communication — very common in agent/telemetry systems for security tooling.
 
```bash
go get google.golang.org/grpc
go get google.golang.org/protobuf/cmd/protoc-gen-go
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
```
 
### 4a. Define the contract (`agent.proto`)
 
```protobuf
syntax = "proto3";
package agent;
option go_package = "myproject/agentpb";
 
service AgentService {
  rpc ReportStatus (StatusRequest) returns (StatusResponse);
  rpc StreamTelemetry (TelemetryRequest) returns (stream TelemetryEvent); // server streaming
}
 
message StatusRequest {
  string agent_id = 1;
}
 
message StatusResponse {
  bool healthy = 1;
  string message = 2;
}
 
message TelemetryRequest {
  string agent_id = 1;
}
 
message TelemetryEvent {
  string event_type = 1;
  string details = 2;
  int64 timestamp = 3;
}
```
This `.proto` file **is the strict, language-agnostic contract** — you generate matching client/server code for Go (or any other language) from this single source of truth. Unlike REST/JSON, there's no ambiguity about field types or optionality; the compiler enforces it.
 
**Generate the Go code:**
```bash
protoc --go_out=. --go-grpc_out=. agent.proto
```
 
### 4b. Implement the server
 
```go
type server struct {
    agentpb.UnimplementedAgentServiceServer
}
 
func (s *server) ReportStatus(ctx context.Context, req *agentpb.StatusRequest) (*agentpb.StatusResponse, error) {
    return &agentpb.StatusResponse{Healthy: true, Message: "ok"}, nil
}
 
func (s *server) StreamTelemetry(req *agentpb.TelemetryRequest, stream agentpb.AgentService_StreamTelemetryServer) error {
    for i := 0; i < 5; i++ {
        event := &agentpb.TelemetryEvent{
            EventType: "heartbeat",
            Details:   fmt.Sprintf("tick %d", i),
            Timestamp: time.Now().Unix(),
        }
        if err := stream.Send(event); err != nil {
            return err
        }
        time.Sleep(time.Second)
    }
    return nil
}
 
func main() {
    lis, _ := net.Listen("tcp", ":50051")
    grpcServer := grpc.NewServer()
    agentpb.RegisterAgentServiceServer(grpcServer, &server{})
    grpcServer.Serve(lis)
}
```
`StreamTelemetry` is a **server-streaming RPC** — the client makes one request, and the server keeps sending events over the same connection until it's done. This is exactly the pattern you'd use for an agent continuously reporting telemetry back to a central collector, without the overhead of a new connection per message.
 
### 4c. Implement the client
 
```go
func main() {
    conn, _ := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
    defer conn.Close()
 
    client := agentpb.NewAgentServiceClient(conn)
 
    // Unary call
    resp, _ := client.ReportStatus(context.Background(), &agentpb.StatusRequest{AgentId: "agent-01"})
    log.Println(resp.Healthy, resp.Message)
 
    // Streaming call
    stream, _ := client.StreamTelemetry(context.Background(), &agentpb.TelemetryRequest{AgentId: "agent-01"})
    for {
        event, err := stream.Recv()
        if err == io.EOF {
            break
        }
        log.Println(event.EventType, event.Details)
    }
}
```
**Note:** `insecure.NewCredentials()` is for local testing only — in any real deployment, use TLS (`credentials.NewTLS(...)`) or mTLS (same idea as the REST mTLS example above) so agents can't be impersonated.
 
**Why this matters for security tooling specifically:** if you're building something like a fleet-wide agent (EDR-style, asset inventory, log shipper), gRPC's binary format + streaming + strict typing makes it far more efficient than a REST agent that polls every N seconds — at scale (thousands of agents), this is a real difference in bandwidth and server load.
 
---
 
# PART 5: SOAP — Implementation in Go
 
Go doesn't have first-class SOAP support (no equivalent of Python's `zeep`) because SOAP is rarely chosen for new Go projects — you'll mostly be **consuming** an existing legacy SOAP API, not building one.
 
### 5a. Consuming a SOAP API manually (works with any SOAP endpoint)
 
```go
func callSOAP(url, soapAction string, envelope []byte) ([]byte, error) {
    req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(envelope))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "text/xml; charset=utf-8")
    req.Header.Set("SOAPAction", soapAction)
 
    client := &http.Client{Timeout: 15 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    return io.ReadAll(resp.Body)
}
 
// Example envelope you'd build for a specific WSDL operation:
const envelope = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetUserInfo xmlns="http://example.com/">
      <UserID>12345</UserID>
    </GetUserInfo>
  </soap:Body>
</soap:Envelope>`
```
SOAP is fundamentally just **XML over HTTP with a rigid envelope structure** — there's no magic, just a specific XML shape the server expects. You parse the response the same way you'd parse any XML in Go, using `encoding/xml`.
 
### 5b. Using a helper library instead of hand-rolling XML
 
```bash
go get github.com/tiaguinho/gosoap
```
```go
soapClient, _ := gosoap.SoapClient("https://example.com/service.wsdl", nil)
 
params := gosoap.Params{"UserID": "12345"}
res, _ := soapClient.Call("GetUserInfo", params)
 
var result MyResponseStruct
res.Unmarshal(&result)
```
This handles WSDL parsing and envelope construction for you — worth it if you're integrating with more than one or two SOAP operations.
 
---
 
# PART 6: Webhooks — Implementation in Go
 
Webhooks have two sides: **receiving** (you expose an endpoint, a third party calls it) and **sending** (you call someone else's endpoint when your tool has an event).
 
### 6a. Receiving a webhook (with signature verification)
 
Third parties (GitHub, Stripe, etc.) sign their webhook payloads with a shared secret so you can verify the request actually came from them and wasn't forged. **Always verify this** — an unauthenticated webhook receiver is an open door for anyone who finds the URL.
 
```go
func webhookHandler(secret []byte) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        body, _ := io.ReadAll(r.Body)
 
        signature := r.Header.Get("X-Signature-256")
        mac := hmac.New(sha256.New, secret)
        mac.Write(body)
        expected := hex.EncodeToString(mac.Sum(nil))
 
        if !hmac.Equal([]byte(signature), []byte(expected)) {
            http.Error(w, "invalid signature", http.StatusUnauthorized)
            return
        }
 
        var event map[string]interface{}
        json.Unmarshal(body, &event)
 
        // handle the event (e.g., "scan_complete", "payment_success", etc.)
        log.Println("verified webhook event:", event)
        w.WriteHeader(http.StatusOK)
    }
}
```
`hmac.Equal` (not `==`) matters here — it does a constant-time comparison, which prevents a timing attack where an attacker could otherwise guess the correct signature one byte at a time by measuring response times.
 
### 6b. Sending a webhook (notifying other systems of your tool's events)
 
```go
func sendWebhook(url string, secret []byte, payload any) error {
    body, _ := json.Marshal(payload)
 
    mac := hmac.New(sha256.New, secret)
    mac.Write(body)
    signature := hex.EncodeToString(mac.Sum(nil))
 
    req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-Signature-256", signature)
 
    client := &http.Client{Timeout: 5 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    return nil
}
```
**Production tip:** wrap this with retry logic (exponential backoff) — the receiving endpoint might be briefly down, and a real webhook system (Stripe's, for example) retries failed deliveries for hours/days rather than giving up after one attempt.
 
---
 
# PART 7: WebRTC — Implementation in Go
 
WebRTC is the most complex style here, and it's rarely relevant to backend/tool-building work (it's mostly for browser-to-browser audio/video/data). Including it for completeness, but you'll reach for this far less often than the others.
 
The Go library is **Pion** (`pion/webrtc`) — the most complete WebRTC implementation outside a browser.
 
```bash
go get github.com/pion/webrtc/v4
```
 
```go
import "github.com/pion/webrtc/v4"
 
func createPeerConnection() (*webrtc.PeerConnection, error) {
    config := webrtc.Configuration{
        ICEServers: []webrtc.ICEServer{
            {URLs: []string{"stun:stun.l.google.com:19302"}}, // helps peers find each other through NAT
        },
    }
    return webrtc.NewPeerConnection(config)
}
```
**Important reality check:** WebRTC needs a **separate signaling channel** (usually a plain WebSocket or REST endpoint you build yourself) for the two peers to exchange connection info *before* they can talk directly. It also typically needs STUN/TURN servers to work reliably across real-world networks (NAT traversal). This is genuinely one of the more involved things to get right — if you ever need it, budget real time for it and look at Pion's official example repo rather than trying to build it from a snippet.
 
---
 
# PART 8: Raw TCP/UDP — Custom Protocols in Go
 
This is the one that matters most for a lot of security/cyber tooling — scanners, agents, network utilities, and anything where you're building **both ends** of a lightweight custom protocol and don't want the overhead of HTTP framing at all.
 
### 8a. TCP server
 
```go
func startTCPServer(addr string) error {
    listener, err := net.Listen("tcp", addr)
    if err != nil {
        return err
    }
    defer listener.Close()
 
    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        go handleConn(conn) // handle each connection concurrently
    }
}
 
func handleConn(conn net.Conn) {
    defer conn.Close()
    reader := bufio.NewReader(conn)
 
    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            return // client disconnected
        }
        log.Println("received:", strings.TrimSpace(line))
        conn.Write([]byte("ack\n"))
    }
}
```
 
### 8b. TCP client
 
```go
func connectTCP(addr string) error {
    conn, err := net.Dial("tcp", addr)
    if err != nil {
        return err
    }
    defer conn.Close()
 
    conn.Write([]byte("hello\n"))
 
    reader := bufio.NewReader(conn)
    response, _ := reader.ReadString('\n')
    log.Println("server said:", response)
    return nil
}
```
 
### 8c. UDP (connectionless — good for lightweight, high-volume, loss-tolerant messages)
 
```go
// Server
func startUDPServer(addr string) error {
    udpAddr, _ := net.ResolveUDPAddr("udp", addr)
    conn, err := net.ListenUDP("udp", udpAddr)
    if err != nil {
        return err
    }
    defer conn.Close()
 
    buf := make([]byte, 1024)
    for {
        n, remoteAddr, _ := conn.ReadFromUDP(buf)
        log.Printf("from %s: %s", remoteAddr, buf[:n])
        conn.WriteToUDP([]byte("ack"), remoteAddr)
    }
}
 
// Client
func sendUDP(addr, message string) error {
    conn, err := net.Dial("udp", addr)
    if err != nil {
        return err
    }
    defer conn.Close()
    _, err = conn.Write([]byte(message))
    return err
}
```
**TCP vs UDP, quick rule of thumb:** use TCP when you need every message delivered, in order (agent reporting, control commands). Use UDP when occasional message loss is fine and speed/low overhead matters more (metrics pings, some scanning techniques) — but you're responsible for handling loss/ordering yourself if it matters.
 
### 8d. Designing your own protocol — a minimal length-prefixed framing pattern
 
Raw TCP gives you a stream of bytes, not "messages" — you need to define where one message ends and the next begins. The simplest robust way is **length-prefixing**:
 
```go
func sendMessage(conn net.Conn, payload []byte) error {
    length := uint32(len(payload))
    header := make([]byte, 4)
    binary.BigEndian.PutUint32(header, length)
 
    if _, err := conn.Write(header); err != nil {
        return err
    }
    _, err := conn.Write(payload)
    return err
}
 
func readMessage(conn net.Conn) ([]byte, error) {
    header := make([]byte, 4)
    if _, err := io.ReadFull(conn, header); err != nil {
        return nil, err
    }
    length := binary.BigEndian.Uint32(header)
 
    payload := make([]byte, length)
    if _, err := io.ReadFull(conn, payload); err != nil {
        return nil, err
    }
    return payload, nil
}
```
This pattern (4-byte length header, then that many bytes of payload) is the backbone of most custom binary protocols — including, under the hood, things like gRPC's own framing. Once you have this, you can send anything as the payload: JSON, protobuf, or your own binary format.
 
**Security note that applies to everything in this section:** the moment this carries anything sensitive (credentials, commands, scan results with real target data), wrap the connection in TLS (`crypto/tls`) rather than sending plaintext over raw sockets. Never invent your own encryption — use the standard library's `crypto/tls` and call it done.
 
---
 
# PART 9: Master Library Cheat Sheet
 
| Style | Go library | Install |
|---|---|---|
| REST (framework) | Gin | `go get github.com/gin-gonic/gin` |
| REST (zero-dep) | `net/http` (standard library) | built-in |
| GraphQL | gqlgen | `go get github.com/99designs/gqlgen` |
| WebSockets | gorilla/websocket | `go get github.com/gorilla/websocket` |
| gRPC | google.golang.org/grpc | `go get google.golang.org/grpc` |
| SOAP (client) | gosoap | `go get github.com/tiaguinho/gosoap` |
| WebRTC | Pion | `go get github.com/pion/webrtc/v4` |
| Raw TCP/UDP | `net` (standard library) | built-in |
| JWT auth | golang-jwt | `go get github.com/golang-jwt/jwt/v5` |
| OAuth2 (client) | golang.org/x/oauth2 | `go get golang.org/x/oauth2` |
| Rate limiting | golang.org/x/time/rate | `go get golang.org/x/time/rate` |
| Packet capture / raw network analysis (adjacent to "API" work, common in cyber tooling) | gopacket | `go get github.com/google/gopacket` |
 
---
 
# PART 10: Scenario Playbook for Cyber/Security Tooling
 
Concrete "I'm building X, use Y" scenarios, since this is where most of your future project decisions will actually come from:
 
| You're building... | Use | Why |
|---|---|---|
| A vulnerability scanner with a web dashboard | REST (Gin) backend + WebSockets for live scan progress | Dashboard needs standard CRUD (scan history, targets) plus real-time feedback while a scan runs |
| A fleet of monitoring/EDR-style agents reporting to a central server | gRPC with mTLS | Efficient binary format, strict contract, strong mutual auth so agents can't be spoofed |
| A CLI tool that enriches IPs/hashes via VirusTotal, Shodan, AbuseIPDB, etc. | REST client with API key auth | You're consuming their existing REST APIs — no choice needed, just implement the client correctly with timeouts |
| An authorized red-team C2/management console | mTLS + gRPC or a custom length-prefixed TCP protocol | Full control over both ends, strong auth is non-negotiable, minimal overhead matters |
| A Slack/Discord bot that posts alerts when your tool finds something | REST (to send messages) + Webhooks (to receive triggers) | Most chat platforms expose REST APIs for sending and webhook systems for receiving |
| A packet sniffer / network traffic analyzer | Not really an "API" at all — use `gopacket` for raw packet capture | This operates below the API layer, directly on network interfaces |
| Integrating with an old enterprise SIEM or ticketing system | Whatever they expose — often SOAP for older ones, REST for modern ones | You don't get to choose here; check their docs first |
| A tool that needs to notify you the instant a scan finishes, without polling | WebSockets (if it's your own UI) or Webhooks (if it's a separate system) | Push beats poll whenever "instant" matters |
| An internal automation tool gluing together several of your other internal Go services | gRPC | Best performance/type-safety for Go-to-Go, especially as the number of services grows |
 
---
 
# PART 11: Universal Best Practices (apply to every style above)
 
- **Always set timeouts** on both clients and servers — a hung connection with no timeout is one of the most common bugs in hand-rolled network code.
- **Always use TLS** the moment real credentials, targets, or sensitive data cross the wire — `crypto/tls` in the standard library, never a hand-rolled scheme.
- **Prefer mTLS over API keys** for machine-to-machine/agent communication where you control both ends — a stolen API key works from anywhere; a stolen private key without the matching cert doesn't.
- **Log and rate-limit everything** you expose publicly — this matters even more for security tooling, where your own API surface can itself become a target.
- **Version your API early** (`/api/v1/...` or a version field in gRPC) — retrofitting versioning after you have real clients is painful.
- **Validate all input server-side**, regardless of style — client-side validation (or a strict GraphQL/protobuf schema) helps, but never assume it's enough.
---
 
This doc + `Guide_to_API_for_dummies.md` together should cover: what each API style is, when to choose it, how to implement both server and client in Go, how to authenticate properly, and how it maps onto real projects — including the kind of tooling you'll build in a security career. Come back to Part 0 and Part 10 any time you're starting a new project and aren't sure which style fits.
 

