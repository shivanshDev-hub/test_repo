# The Complete DBMS Practitioner's Guide
### Theory → PostgreSQL → Redis → MongoDB, in the order you'd actually build a system

Companion videos (short, practical, watch in this order):
1. **DBMS Concepts:** [13 Database Concepts Every Backend Engineer Must Know](https://www.youtube.com/watch?v=pEzrn6PIxRo)
2. **PostgreSQL:** [Learn PostgreSQL in 40 Minutes](https://www.youtube.com/watch?v=chz_g0VfP4U)
3. **Redis:** [Redis Crash Course](https://www.youtube.com/watch?v=jgpVdJB2sKQ)
4. **MongoDB:** [The Ultimate Beginner's Guide to MongoDB in 30 Minutes](https://www.youtube.com/watch?v=9tkDwRAF58Q)

This guide is written so a complete beginner — zero prior DBMS knowledge — can walk it top to bottom and come out the other side able to design, build, query, optimize, secure, and recover a production-grade database system. It's organized in **implementation order**: model your data → build the schema → write queries → optimize them → protect them from concurrent access and crashes → lock them down → then the same journey for Redis and MongoDB → then how to combine all three. Go code is used throughout since that's the stack this pairs with (your Gin/API docs).

**How to read this:** each numbered Part builds on the one before it. If you're already comfortable with a Part, skip to the next — but the ordering matters the first time through, because later sections assume earlier ones.

---

# PART 0 — Orientation

## 0.1 Why a DBMS instead of just files?

Storing data in flat files (CSV, JSON files) breaks down fast:
- **No safe concurrent writes** — two processes writing at once corrupts data.
- **No fast lookup** without scanning everything.
- **No way to enforce rules** — nothing stops "email must be unique" from being violated.
- **No recovery** if the process dies mid-write — you get a half-written file.

A DBMS solves all four: it manages concurrent access safely, indexes data for fast lookup, enforces integrity rules, and guarantees that a crash mid-operation doesn't leave you with corrupted data. Everything in this guide is really just an explanation of *how* it does those four things, and how you use that machinery correctly.

## 0.2 SQL (Relational) vs NoSQL — the real decision

| | SQL (Postgres, MySQL) | NoSQL (MongoDB, Redis) |
|---|---|---|
| Structure | Fixed schema, tables + rows | Flexible — documents (Mongo) or key-value (Redis) |
| Relationships | Native — foreign keys, joins | Usually denormalized/embedded, or handled in app code |
| Best for | Data with clear structure and relationships (users, orders, transactions) | Rapidly changing shapes, huge scale, or simple key-value lookups/caching |
| Consistency | Strong (ACID) by default | Varies — often "eventual consistency" traded for speed/scale |

**Practical rule:** default to Postgres unless you have a specific reason not to. Reach for MongoDB when your data genuinely doesn't fit neat tables (deeply nested, wildly varying shapes). Reach for Redis when you need something in-memory and fast (caching, sessions, rate limits) — not as your primary source of truth. Part 10 turns this into a full, mechanism-level decision framework once you've seen how each engine actually works internally.

## 0.3 Environment Setup — get all three running now

Do this once, up front, so every code sample below is something you can actually run as you read.

```bash
# PostgreSQL
docker run --name pg -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres:16

# Redis
docker run --name redis -p 6379:6379 -d redis:7

# MongoDB
docker run --name mongo -p 27017:27017 -d mongo:7
```

Connect to each via CLI:

```bash
docker exec -it pg psql -U postgres
docker exec -it redis redis-cli
docker exec -it mongo mongosh
```

Keep these three containers running for the rest of this guide — every example assumes they're reachable at `localhost:5432`, `localhost:6379`, and `localhost:27017`.

## 0.4 Roadmap

| Part | What you'll be able to do after it |
|---|---|
| 1 — Data Modeling | Turn a real-world requirement into keys, an ER diagram, and a normalized schema |
| 2 — PostgreSQL | Build that schema for real and write any query against it |
| 3 — Query Processing & Optimization | Explain *why* a query is slow and fix it |
| 4 — Transactions, Concurrency & Recovery | Run correct, concurrent, crash-safe transactions |
| 5 — Security | Lock the database down: who can touch what, and encrypt what matters |
| 6 — Redis | Add a caching/session/real-time layer correctly |
| 7 — MongoDB | Model and query document data correctly |
| 8–9 — NoSQL landscape & distributed systems | Know what else exists and what you trade away at scale |
| 10–12 — Decision framework, architecture, cheat sheets | Choose the right tool per job and combine them like a real system |

---
# PART 1 — Data Modeling: From Requirements to Schema

Before you write a single `CREATE TABLE`, you model the data on paper (or in a diagram). Skipping this step is the single biggest cause of painful schema migrations six months into a project. This Part gives you the full toolkit: keys, ER diagrams, functional dependencies, and normalization — then walks one running example (a small e-commerce system: **customers**, **orders**, **products**) all the way from requirement to real schema, which Part 2 then builds for real in Postgres.

## 1.1 Keys — every kind, and why it exists

A **key** is one or more columns that let you uniquely identify or relate rows. Getting keys right is the foundation everything else (constraints, joins, indexes, normalization) is built on.

| Key type | Definition | Example |
|---|---|---|
| **Super key** | Any set of columns that uniquely identifies a row (may contain extra, unnecessary columns) | `(email)`, `(email, name)`, `(id, email)` all uniquely identify a `users` row |
| **Candidate key** | A *minimal* super key — remove any column and it stops being unique | `(email)` alone, if email is unique |
| **Primary key** | The candidate key you actually choose to identify rows in practice | `id` (usually a surrogate, see below) |
| **Alternate key** | A candidate key that exists but wasn't chosen as primary | `email`, if `id` was chosen as primary |
| **Composite key** | A key made of more than one column, needed together | `(order_id, product_id)` in an order-line-items table |
| **Foreign key** | A column (or set of columns) that references another table's primary key | `orders.customer_id` → `customers.id` |
| **Surrogate key** | An artificial, meaningless identifier the database generates | auto-increment integer, or a UUID |
| **Natural key** | A real-world attribute that already uniquely identifies the entity | a national ID number, an ISBN, an email address |

**Practical implementation in SQL:**

```sql
CREATE TABLE customers (
    id       SERIAL PRIMARY KEY,           -- surrogate primary key
    email    VARCHAR(255) UNIQUE NOT NULL, -- alternate key, enforced
    name     VARCHAR(100) NOT NULL
);

CREATE TABLE order_items (
    order_id    INT REFERENCES orders(id),
    product_id  INT REFERENCES products(id),
    quantity    INT NOT NULL CHECK (quantity > 0),
    PRIMARY KEY (order_id, product_id)     -- composite primary key
);
```

**Practical rule — surrogate vs natural:** default to a surrogate key (`SERIAL`/`BIGSERIAL`/`UUID`). Natural keys look convenient but tend to change (people change emails, get married, correct typos), and a primary key should never need to change since everything else references it. Use `UUID` instead of `SERIAL` specifically when the ID will ever be exposed publicly — sequential IDs leak how many rows you have and are guessable.

```sql
CREATE TABLE customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),  -- requires: CREATE EXTENSION pgcrypto;
    email VARCHAR(255) UNIQUE NOT NULL
);
```

## 1.2 Entities, Relationships, and ER Diagrams

An **Entity-Relationship (ER) diagram** is how you model data on paper before it becomes tables — it forces you to nail down every entity, its attributes, and exactly how entities relate, *before* you're fighting a migration in production.

**Building blocks:**
- **Entity** — a thing you store data about (`Customer`, `Order`, `Product`). Becomes a table.
- **Attribute** — a property of an entity (`name`, `email`). Becomes a column. Attributes can be *simple* (atomic, like `name`), *composite* (splittable, like an `address` made of street/city/zip), *derived* (computable from other data, like `age` from `birth_date` — usually not stored), or *multivalued* (a customer can have several `phone_numbers` — becomes its own table, never a comma-separated column).
- **Relationship** — how two entities connect (`Customer places Order`). Becomes a foreign key, or its own junction table.
- **Cardinality** — how many of one entity relate to how many of another:
  - **1:1** — one row relates to at most one row on the other side (a `User` has one `Profile`).
  - **1:N** — one row relates to many rows on the other side (one `Customer` has many `Orders`).
  - **M:N** — many rows relate to many rows (many `Orders` contain many `Products`).
- **Participation** — *total* (every entity instance must participate, e.g. every `OrderItem` must belong to an `Order`) vs *partial* (optional, e.g. a `Customer` doesn't have to have placed an `Order`).

**The running example, as an ER diagram:**

```mermaid
erDiagram
    CUSTOMER ||--o{ ORDER : places
    ORDER ||--|{ ORDER_ITEM : contains
    PRODUCT ||--o{ ORDER_ITEM : "appears in"

    CUSTOMER {
        int id PK
        varchar name
        varchar email UK
        varchar city
    }
    ORDER {
        int id PK
        int customer_id FK
        date order_date
    }
    ORDER_ITEM {
        int order_id PK, FK
        int product_id PK, FK
        int quantity
    }
    PRODUCT {
        int id PK
        varchar name
        decimal price
    }
```

Reading the notation: `||` means "exactly one", `o{` means "zero or many", `|{` means "one or many". So `CUSTOMER ||--o{ ORDER` reads as "one customer relates to zero-or-many orders" — a 1:N relationship with partial participation on the customer side (a customer doesn't have to have ordered anything yet).

**Translating an ER diagram into tables (the mapping rules):**

1. **Every entity becomes a table.** Its attributes become columns.
2. **A 1:N relationship** becomes a foreign key on the "many" side, pointing at the "one" side's primary key. (`orders.customer_id → customers.id`.)
3. **An M:N relationship** becomes its own junction table (also called a bridge or associative table) holding foreign keys to both sides, usually with a composite primary key. (`order_items(order_id, product_id, quantity)` — this is exactly how `Orders` and `Products` are M:N: one order has many products, one product appears in many orders.)
4. **A 1:1 relationship** becomes a foreign key on either side with a `UNIQUE` constraint, or the two entities get merged into one table if they're always accessed together.
5. **A multivalued attribute** becomes its own table with a foreign key back to the owning entity.

This is exactly the schema Part 2 builds — the ER diagram above *is* the plan for the `customers` / `orders` / `order_items` / `products` tables you'll create in SQL.

## 1.3 Functional Dependencies

A **functional dependency (FD)**, written `X → Y`, means: given a value of `X`, `Y` is completely determined — two rows with the same `X` must have the same `Y`. FDs are the formal reasoning tool behind normalization; every normalization rule in 1.4 is really just "eliminate a specific kind of bad functional dependency."

**Notation and examples**, using an unnormalized `orders` table `(OrderID, Product, CustomerID, CustomerName, CustomerCity)`:

- `OrderID, Product → Quantity` — the pair together determines quantity. This is a **full functional dependency** on the composite key.
- `OrderID → CustomerID` — but `Product` has nothing to do with it. Since `OrderID` alone determines `CustomerID`, while the key is `(OrderID, Product)`, this is a **partial dependency** — `CustomerID` depends on only *part* of the composite key. (This is exactly what 2NF eliminates.)
- `CustomerID → CustomerName`, and separately `OrderID → CustomerID → CustomerName` — `CustomerName` depends on `CustomerID`, which depends on `OrderID`, not directly on `OrderID`. This chain is a **transitive dependency**. (This is exactly what 3NF eliminates.)
- `OrderID → OrderID` is a **trivial** FD (the right side is a subset of the left side) — always true, tells you nothing useful.

**Why this matters practically:** every anomaly normalization fixes — update anomalies (change a customer's name, must update it in 50 rows), insertion anomalies (can't add a product with no orders yet because the table is keyed by order), deletion anomalies (delete the last order and you lose the customer's name entirely) — is a *symptom* of a bad functional dependency. Once you can spot "this column doesn't actually depend on the whole key" or "this column depends on another non-key column," you can normalize any schema without memorizing the rules by rote.

**Finding candidate keys from FDs (the practical technique):** compute the **closure** of an attribute set — repeatedly apply known FDs to see what else you can derive. If the closure of some attribute set reaches *every* column in the table, that set is a superkey; if no smaller subset does, it's a candidate key. Example: given `OrderID → CustomerID` and `CustomerID → CustomerName, CustomerCity`, the closure of `{OrderID}` is `{OrderID, CustomerID, CustomerName, CustomerCity}` — if that's every column, `OrderID` alone is a candidate key for that table.

## 1.4 Normalization

Normalization is the process of organizing tables to reduce data redundancy and avoid the update/insert/delete anomalies described above. Each normal form builds on the one before it — you can't skip a level. You start with an **unnormalized table (UNF)** — often flat, with repeating groups or mixed data — and progressively split it into smaller, related tables until each one models exactly one thing and every functional dependency is "well-behaved" (non-key columns depend on the whole key, and only the key).

---

### Step 1: First Normal Form (1NF)

**Rule:** Every column holds a single (atomic) value, no repeating groups, and each row is uniquely identifiable.

**Before (UNF):**

| OrderID | Customer | Products                 |
|---------|----------|---------------------------|
| 1       | Alice    | Pen (2), Notebook (1)     |
| 2       | Bob      | Pen (5)                   |

**After (1NF)** — split the repeating group into separate rows:

| OrderID | Customer | Product   | Quantity |
|---------|----------|-----------|----------|
| 1       | Alice    | Pen       | 2        |
| 1       | Alice    | Notebook  | 1        |
| 2       | Bob      | Pen       | 5        |

Now the primary key is the composite `(OrderID, Product)`.

---

### Step 2: Second Normal Form (2NF)

**Rule:** Must be in 1NF, and every non-key column must depend on the **whole** primary key — not just part of it. (Only relevant when you have a composite key.)

In the 1NF table above, `Customer` depends only on `OrderID`, not on `Product` — that's a **partial dependency** (see 1.3). Fix by splitting:

**Orders**

| OrderID | Customer |
|---------|----------|
| 1       | Alice    |
| 2       | Bob      |

**OrderItems**

| OrderID | Product   | Quantity |
|---------|-----------|----------|
| 1       | Pen       | 2        |
| 1       | Notebook  | 1        |
| 2       | Pen       | 5        |

---

### Step 3: Third Normal Form (3NF)

**Rule:** Must be in 2NF, and non-key columns must depend **only** on the key — not on other non-key columns (no transitive dependencies).

Suppose `Orders` grows to include customer details:

| OrderID | CustomerID | CustomerName | CustomerCity |
|---------|------------|--------------|--------------|
| 1       | C01        | Alice        | Boston       |
| 2       | C02        | Bob          | Denver       |

Here `CustomerName` and `CustomerCity` depend on `CustomerID`, not directly on `OrderID` — a **transitive dependency** (`OrderID → CustomerID → CustomerName/City`, see 1.3). Fix by splitting:

**Orders**

| OrderID | CustomerID |
|---------|------------|
| 1       | C01        |
| 2       | C02        |

**Customers**

| CustomerID | CustomerName | CustomerCity |
|------------|--------------|--------------|
| C01        | Alice        | Boston       |
| C02        | Bob          | Denver       |

---

### Step 4: Boyce-Codd Normal Form (BCNF)

**Rule:** Must be in 3NF, and **every determinant must be a candidate key** (a stricter version of 3NF that catches edge cases 3NF misses, usually with overlapping composite keys).

Example: each instructor teaches only one course, but a course can have several instructors, and students can have several instructors per course:

| StudentID | Course  | Instructor |
|-----------|---------|------------|
| S1        | Math    | Dr. Lee    |
| S1        | Physics | Dr. Han    |
| S2        | Math    | Dr. Lee    |

Key is `(StudentID, Course)`, but `Instructor → Course` also holds, and `Instructor` isn't a candidate key. That violates BCNF. Fix by splitting:

**StudentInstructor**

| StudentID | Instructor |
|-----------|------------|
| S1        | Dr. Lee    |
| S1        | Dr. Han    |
| S2        | Dr. Lee    |

**InstructorCourse**

| Instructor | Course  |
|------------|---------|
| Dr. Lee    | Math    |
| Dr. Han    | Physics |

---

### Step 5: Fourth Normal Form (4NF)

**Rule:** Must be in BCNF, and the table must have **no multi-valued dependencies** — two independent multi-valued facts about the same entity shouldn't be combined in one table.

An employee can have several skills and several languages, independently of each other:

| Employee | Skill   | Language |
|----------|---------|----------|
| Alice    | Python  | English  |
| Alice    | Python  | French   |
| Alice    | SQL     | English  |
| Alice    | SQL     | French   |

This forces every skill to be paired with every language — pure redundancy. Fix by splitting into two independent tables:

**EmployeeSkill**

| Employee | Skill  |
|----------|--------|
| Alice    | Python |
| Alice    | SQL    |

**EmployeeLanguage**

| Employee | Language |
|----------|----------|
| Alice    | English  |
| Alice    | French   |

---

### Step 6: Fifth Normal Form (5NF)

**Rule:** Must be in 4NF, and the table must have **no join dependency** that isn't already implied by its candidate keys — i.e., it can't be losslessly split into smaller tables without losing information when rejoined.

An agent sells a company's product only if that agent carries that company **and** that company makes that product **and** the agent sells that product — but not every combination is valid, and no two-table split preserves that constraint:

| Agent | Company | Product |
|-------|---------|---------|
| A1    | Acme    | Widget  |
| A1    | Globex  | Gadget  |

If the valid combinations follow a genuine three-way join dependency, splitting into two tables would reintroduce invalid combinations when joined back. Fix by decomposing into three binary tables:

**AgentCompany** `(Agent, Company)` · **CompanyProduct** `(Company, Product)` · **AgentProduct** `(Agent, Product)`

The original table is reconstructed only by joining all three — no redundancy, no false combinations.

---

### Quick Reference

| Level | Fixes | Requirement |
|-------|--------|-------------|
| 1NF   | Repeating groups, non-atomic values     | Atomic columns, unique rows                    |
| 2NF   | Partial dependency                      | Non-key columns depend on the *whole* key      |
| 3NF   | Transitive dependency                   | Non-key columns depend *only* on the key       |
| BCNF  | Anomalies 3NF misses                    | Every determinant is a candidate key           |
| 4NF   | Multi-valued dependency                 | No independent multi-valued facts combined     |
| 5NF   | Join dependency                         | No lossless decomposition remains possible     |

Most real-world schemas stop at **3NF or BCNF** — that's where the biggest anomalies are eliminated. 4NF/5NF are reserved for specific edge cases.

## 1.5 Denormalization — breaking the rules on purpose

**Denormalization** means deliberately reintroducing some redundancy for performance. It's not a mistake — it's a trade-off you make consciously, after you understand normalization, not instead of understanding it.

**Practical example:** storing a `product_name` snapshot directly on an `order_item` row, even though `product_name` technically belongs to `products` and normalization says reference it by ID:

```sql
CREATE TABLE order_items (
    order_id     INT REFERENCES orders(id),
    product_id   INT REFERENCES products(id),
    product_name VARCHAR(255) NOT NULL,  -- denormalized snapshot, intentional
    unit_price   DECIMAL(10,2) NOT NULL, -- price *at time of purchase*
    quantity     INT NOT NULL
);
```

**Why:** if `products.name` or `products.price` changes later, historical orders shouldn't silently change too — a customer's receipt from last year should show what they actually paid, not today's price. Denormalizing here isn't laziness, it's correctness for a different requirement (historical accuracy) than the one normalization optimizes for (avoiding redundant updates).

**Practical rule:** normalize by default (3NF), then denormalize specific, deliberate spots once you've identified a real read-performance problem or a correctness requirement like the one above — never denormalize preemptively "just in case."

## 1.6 From ER Diagram to Real Schema — the worked example

Putting 1.1–1.5 together, the ER diagram in 1.2 translates directly into this schema (this is exactly what Part 2 builds in Postgres):

```sql
CREATE TABLE customers (
    id    SERIAL PRIMARY KEY,                 -- surrogate key (1.1)
    name  VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,        -- alternate key (1.1)
    city  VARCHAR(100)
);

CREATE TABLE products (
    id    SERIAL PRIMARY KEY,
    name  VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0)
);

CREATE TABLE orders (
    id           SERIAL PRIMARY KEY,
    customer_id  INT NOT NULL REFERENCES customers(id),  -- FK for the 1:N relationship (1.2)
    order_date   DATE DEFAULT CURRENT_DATE
);

CREATE TABLE order_items (                                -- junction table for the M:N relationship (1.2)
    order_id   INT NOT NULL REFERENCES orders(id),
    product_id INT NOT NULL REFERENCES products(id),
    quantity   INT NOT NULL CHECK (quantity > 0),
    PRIMARY KEY (order_id, product_id)                    -- composite key (1.1)
);
```

Every design decision here traces back to something in this Part: `id` columns are surrogate primary keys, `email UNIQUE` is an alternate key, `customer_id` is the foreign key the 1:N relationship maps to, `order_items` is the junction table the M:N relationship maps to, and the whole thing is in 3NF (no partial or transitive dependencies remain — `order_items.quantity` depends on the *whole* composite key, and nothing here depends transitively on a non-key column).

---
# PART 2 — PostgreSQL: Building and Querying the Database

Environment is already running (0.3). Schema is already designed (1.6). Now you build it for real and learn to query it — every example below uses the same `customers` / `products` / `orders` / `order_items` schema from Part 1, so run 1.6's `CREATE TABLE` statements first.

```bash
docker exec -it pg psql -U postgres
```

## 2.1 Seed some data

```sql
INSERT INTO customers (name, email, city) VALUES
    ('Alice', 'alice@example.com', 'Boston'),
    ('Bob',   'bob@example.com',   'Denver'),
    ('Carla', 'carla@example.com', 'Boston');

INSERT INTO products (name, price) VALUES
    ('Pen', 2.50), ('Notebook', 8.00);

INSERT INTO orders (customer_id, order_date) VALUES
    (1, '2026-01-05'), (1, '2026-02-10'), (2, '2026-02-11'), (3, '2026-03-01');

INSERT INTO order_items (order_id, product_id, quantity) VALUES
    (1, 1, 2), (1, 2, 1), (2, 1, 5), (3, 2, 3), (4, 1, 1);
```

## 2.2 How to Write Any Query — a repeatable method

Every query — SQL, MongoDB, even a Redis command chain — comes from answering the same five questions in order. Beginners jump straight to typing SQL and get lost; work through this checklist instead:

1. **What am I trying to get back?** (columns/fields — this becomes `SELECT` / a Mongo projection.)
2. **Where does that data live?** (which table(s)/collection(s) — this becomes `FROM` / `db.<collection>`.)
3. **How do those tables relate, if more than one is involved?** (this becomes `JOIN ... ON` / Mongo's `$lookup`, using the foreign keys from Part 1.)
4. **What rows/documents do I want to keep?** (this becomes `WHERE` / a Mongo filter object.)
5. **Do I need to summarize, order, or limit the result?** (this becomes `GROUP BY`/`HAVING`/`ORDER BY`/`LIMIT` or a Mongo aggregation pipeline.)

Walking through "total amount spent by each Boston customer, highest first" with this method: (1) customer name + total spent, (2) `customers` and `order_items`/`orders`/`products`, (3) join them all on their foreign keys, (4) filter `city = 'Boston'`, (5) group by customer, sum the line totals, sort descending. That mapping *is* the query — the SQL syntax below is just how you type that plan out:

```sql
SELECT c.name, SUM(oi.quantity * p.price) AS total_spent
FROM customers c
JOIN orders o       ON o.customer_id = c.id
JOIN order_items oi ON oi.order_id = o.id
JOIN products p     ON p.id = oi.product_id
WHERE c.city = 'Boston'
GROUP BY c.name
ORDER BY total_spent DESC;
```

Every section below is one piece of this puzzle in more depth.

## 2.3 Retrieving Data with SELECT

```sql
SELECT * FROM customers;                       -- all columns, all rows
SELECT name, city FROM customers;               -- specific columns
SELECT name FROM customers WHERE city = 'Boston';  -- filter rows
```

**Common `WHERE` operators:**

```sql
WHERE price > 5
WHERE price BETWEEN 2 AND 10
WHERE city IN ('Boston', 'Denver')
WHERE name LIKE 'A%'          -- starts with A
WHERE customer_id IS NULL
WHERE price > 2 AND city = 'Boston'
```

## 2.4 Sorting and Limiting Results

```sql
SELECT name, city FROM customers
ORDER BY city ASC, name DESC;

SELECT * FROM products
ORDER BY price DESC
LIMIT 2;                       -- top 2 highest-priced products
```

`ORDER BY` accepts multiple columns (sorted left to right). `LIMIT` caps the row count — pair with `OFFSET` for pagination:

```sql
SELECT * FROM orders
ORDER BY order_date
LIMIT 10 OFFSET 20;            -- rows 21–30
```

## 2.5 Joins — the reason relational databases exist

This is Question 3 from 2.2: combining rows from related tables using the foreign keys defined in Part 1.

**INNER JOIN** — only rows with a match on both sides:

```sql
SELECT c.name, o.order_date
FROM customers c
INNER JOIN orders o ON c.id = o.customer_id;
```

**LEFT JOIN** — all rows from the left table, matched or not (unmatched right-side columns come back `NULL`). Useful for "which customers have *no* orders" — add `WHERE o.id IS NULL`:

```sql
SELECT c.name, o.id AS order_id
FROM customers c
LEFT JOIN orders o ON c.id = o.customer_id;
```

**RIGHT JOIN** — mirror of `LEFT JOIN`. **FULL JOIN** — all rows from both sides, matched where possible.

**Self join** — join a table to itself, e.g. finding customers in the same city:

```sql
SELECT a.name AS customer_a, b.name AS customer_b, a.city
FROM customers a
JOIN customers b ON a.city = b.city AND a.id < b.id;
```

`REFERENCES` (the foreign key from Part 1) is what makes this safe — Postgres rejects any `orders` row whose `customer_id` doesn't correspond to a real customer, so a join can never silently produce nonsense from an orphaned row.

## 2.6 Aggregating Data

**Aggregate functions:** `COUNT()`, `SUM()`, `AVG()`, `MIN()`, `MAX()`.

```sql
SELECT COUNT(*) FROM orders;
SELECT SUM(price) FROM products;
SELECT AVG(price) FROM products;
```

**GROUP BY** — split rows into groups before aggregating:

```sql
SELECT customer_id, COUNT(*) AS order_count
FROM orders
GROUP BY customer_id;
```

**HAVING** — filter *after* aggregation (`WHERE` filters rows, `HAVING` filters groups):

```sql
SELECT customer_id, COUNT(*) AS order_count
FROM orders
GROUP BY customer_id
HAVING COUNT(*) > 1;
```

## 2.7 Subqueries and CTEs

**In `WHERE`** (filter using a computed value):

```sql
SELECT name FROM customers
WHERE id IN (SELECT customer_id FROM orders WHERE order_date > '2026-02-01');
```

**In `FROM`** (treat a query result as a temporary table):

```sql
SELECT city, COUNT(*) AS orders_from_city
FROM (
    SELECT c.city, o.id
    FROM customers c JOIN orders o ON c.id = o.customer_id
) AS per_order
GROUP BY city;
```

**Correlated subquery** (re-evaluated per outer row — each customer's most recent order):

```sql
SELECT o1.*
FROM orders o1
WHERE order_date = (
    SELECT MAX(o2.order_date) FROM orders o2 WHERE o2.customer_id = o1.customer_id
);
```

**CTE (Common Table Expression)** — the same idea as a subquery-in-`FROM`, written more readably with `WITH`:

```sql
WITH customer_order_counts AS (
    SELECT customer_id, COUNT(*) AS order_count
    FROM orders
    GROUP BY customer_id
)
SELECT c.name, coc.order_count
FROM customers c
JOIN customer_order_counts coc ON c.id = coc.customer_id;
```

## 2.8 Combining Result Sets

```sql
SELECT city FROM customers UNION SELECT city FROM customers;      -- combine, remove duplicates
SELECT city FROM customers UNION ALL SELECT city FROM customers;  -- combine, keep duplicates (faster)
SELECT city FROM customers INTERSECT SELECT 'Boston';              -- rows present in both
SELECT city FROM customers EXCEPT SELECT 'Boston';                 -- rows in the first, not the second
```

All columns and types must line up between the combined queries.

## 2.9 Modifying Data

```sql
-- INSERT
INSERT INTO customers (name, email, city) VALUES ('Dana', 'dana@example.com', 'Seattle');

-- UPDATE — always include WHERE, or every row changes
UPDATE customers SET city = 'Cambridge' WHERE id = 1;

-- DELETE — always include WHERE, or the table empties
DELETE FROM order_items WHERE quantity < 1;

-- UPSERT (insert, or update if it already exists)
INSERT INTO customers (id, name, email, city)
VALUES (1, 'Alice', 'alice@example.com', 'Cambridge')
ON CONFLICT (id) DO UPDATE SET city = EXCLUDED.city;
```

Wrap multi-statement writes in a transaction (full theory in Part 4) whenever more than one write must succeed or fail together:

```sql
BEGIN;
UPDATE products SET price = price * 1.1 WHERE id = 1;
UPDATE products SET price = price * 1.1 WHERE id = 2;
COMMIT;   -- or ROLLBACK to undo both
```

## 2.10 Altering Schema

```sql
ALTER TABLE customers ADD COLUMN phone VARCHAR(20);
ALTER TABLE orders DROP COLUMN order_date;
ALTER TABLE customers RENAME COLUMN city TO hometown;
```

A **migration** is a versioned, scripted schema change like the ones above, tracked in code and run in order — so every environment (laptop, staging, production) ends up with an identical, reproducible schema history instead of someone hand-editing production. In Go, `golang-migrate` is the standard tool.

## 2.11 Indexes — the short version (full depth in Part 3)

```sql
CREATE INDEX idx_orders_customer_id ON orders(customer_id);
CREATE UNIQUE INDEX idx_customers_email ON customers(email);
```

An index is a lookup structure — like a book's index — that lets the database find rows fast instead of scanning the whole table. **Without one, `WHERE email = ?` on a million-row table scans all million rows every time.** The trade-off: indexes speed up reads but slow down writes slightly (every insert updates the index too) and use extra disk space — index columns you actually filter, join, or sort by, not everything. Part 3 covers index *types* and how to verify one is actually being used.

## 2.12 Views

A view saves a query as a virtual table you can reuse, without duplicating the underlying data:

```sql
CREATE VIEW customer_spending AS
SELECT c.name, SUM(oi.quantity * p.price) AS total_spent
FROM customers c
JOIN orders o ON c.id = o.customer_id
JOIN order_items oi ON oi.order_id = o.id
JOIN products p ON p.id = oi.product_id
GROUP BY c.name;

SELECT * FROM customer_spending WHERE total_spent > 20;
```

A view re-runs its query every time it's selected from. A **materialized view** caches the result and must be manually or automatically refreshed:

```sql
CREATE MATERIALIZED VIEW customer_spending_cached AS
SELECT c.name, SUM(oi.quantity * p.price) AS total_spent
FROM customers c JOIN orders o ON c.id = o.customer_id
JOIN order_items oi ON oi.order_id = o.id JOIN products p ON p.id = oi.product_id
GROUP BY c.name;

REFRESH MATERIALIZED VIEW customer_spending_cached;
```

## 2.13 Window Functions

**Purpose:** compute a value across a set of related rows *without* collapsing them into one row, unlike `GROUP BY`.

```sql
SELECT
    oi.order_id,
    oi.quantity * p.price AS line_total,
    SUM(oi.quantity * p.price) OVER (PARTITION BY oi.order_id) AS order_total,
    RANK() OVER (ORDER BY oi.quantity * p.price DESC) AS line_rank
FROM order_items oi
JOIN products p ON p.id = oi.product_id;
```

`PARTITION BY` splits rows into groups (like `GROUP BY`) but keeps every row visible. Common window functions: `ROW_NUMBER()`, `RANK()`, `DENSE_RANK()`, `LAG()`/`LEAD()` (previous/next row's value), `SUM()`/`AVG()` as running totals:

```sql
-- running total of spend per customer, ordered by date
SELECT c.id, o.order_date, oi.quantity * p.price AS line_total,
    SUM(oi.quantity * p.price) OVER (PARTITION BY c.id ORDER BY o.order_date) AS running_total
FROM customers c
JOIN orders o ON o.customer_id = c.id
JOIN order_items oi ON oi.order_id = o.id
JOIN products p ON p.id = oi.product_id;
```

## 2.14 JSONB — Postgres's flexible-schema escape hatch

```sql
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    payload JSONB
);

INSERT INTO events (payload) VALUES ('{"type": "login", "ip": "1.2.3.4"}');

SELECT * FROM events WHERE payload->>'type' = 'login';
CREATE INDEX idx_events_payload ON events USING GIN (payload);   -- see Part 3 for index types
```

This lets you store semi-structured data (like MongoDB would) *inside* a relational table — genuinely useful when 95% of your data is structured but one field (arbitrary metadata, webhook payloads) varies. `payload->>'type'` extracts a JSON field as text.

## 2.15 Connecting from Go

```bash
go get github.com/jackc/pgx/v5
```

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/jackc/pgx/v5/pgxpool"
)

func main() {
    pool, err := pgxpool.New(context.Background(), "postgres://postgres:secret@localhost:5432/postgres")
    if err != nil {
        log.Fatal(err)
    }
    defer pool.Close()

    _, err = pool.Exec(context.Background(),
        "INSERT INTO customers (name, email, city) VALUES ($1, $2, $3)", "Eve", "eve@example.com", "Austin")
    if err != nil {
        log.Fatal(err)
    }

    rows, err := pool.Query(context.Background(), "SELECT id, name, email FROM customers")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    for rows.Next() {
        var id int
        var name, email string
        rows.Scan(&id, &name, &email)
        fmt.Println(id, name, email)
    }
}
```

`pgxpool` is the modern, high-performance Postgres driver for Go — it manages a **connection pool** for you (see 3.8 for why that matters). **Always use `$1, $2` placeholders, never string-concatenate values into SQL** — that's how SQL injection happens. The driver handles escaping safely when values are passed as separate arguments.

## 2.16 Practical Gotchas

- Prefer `TIMESTAMPTZ` (timestamp with timezone) over plain `TIMESTAMP` — avoids painful timezone bugs later.
- Use `UUID` instead of `SERIAL` for primary keys that will ever be exposed publicly (see 1.1).
- Always wrap multi-step writes in a transaction (Part 4).
- Back up with `pg_dump` regularly — it's not automatic (Part 4.15 covers real recovery/PITR).

---
# PART 3 — Query Processing & Optimization (Beginner → Advanced)

You can write correct queries after Part 2. This Part is about making them *fast*, and understanding *why* they're slow when they aren't — the skill that separates someone who can write SQL from someone who can run a database in production.

## 3.1 How a Query Actually Runs

Every query you send goes through four stages before you see a result:

```mermaid
flowchart LR
    A["SQL text you wrote"] --> B["Parser<br/>(syntax check, builds a query tree)"]
    B --> C["Planner / Optimizer<br/>(considers many possible execution plans,<br/>picks the cheapest one it can estimate)"]
    C --> D["Executor<br/>(actually runs the chosen plan,<br/>fetches/filters/joins/sorts rows)"]
    D --> E["Result rows sent back to you"]
```

The **parser** just checks your SQL is syntactically valid and turns it into an internal tree. The **planner** is the interesting part: for a query with several joins and filters, there are often dozens of valid ways to execute it (which table to scan first, whether to use an index, which join algorithm to use — see 3.7) — the planner estimates the cost of each candidate plan (using table statistics, see 3.4) and picks the cheapest one it can find. This is why the *same* SQL can run at wildly different speeds depending on table size, indexes present, and how recently statistics were updated — the planner is making an estimate, not a guarantee.

## 3.2 Reading EXPLAIN ANALYZE

This is the single most useful command for "why is this query slow":

```sql
EXPLAIN ANALYZE
SELECT c.name, SUM(oi.quantity * p.price) AS total_spent
FROM customers c
JOIN orders o ON o.customer_id = c.id
JOIN order_items oi ON oi.order_id = o.id
JOIN products p ON p.id = oi.product_id
WHERE c.city = 'Boston'
GROUP BY c.name;
```

Reading the output, from the inside out (innermost/most-indented lines run first):

- **`Seq Scan`** — Postgres read every row in the table, checking each against the filter. Fine for small tables, a red flag on large ones with a selective filter.
- **`Index Scan`** / **`Index Only Scan`** — Postgres used an index to jump straight to matching rows (the latter didn't even need to touch the table, because the index alone had every column the query needed).
- **`Bitmap Heap Scan`** / **`Bitmap Index Scan`** — a hybrid: build a bitmap of matching row locations from the index, then fetch them — used when a scan will return too many rows for a plain Index Scan to be efficient but still fewer than a full Seq Scan would touch.
- **`cost=0.29..8.32`** — the planner's *estimated* cost (arbitrary units, startup..total), not real time. Useful for comparing plans, not for absolute timing.
- **`actual time=0.012..0.045 rows=3`** — this only appears with `ANALYZE`, and is what actually happened when Postgres ran it for real: real elapsed milliseconds and real row count.
- **The gap between estimated and actual rows is the single most valuable thing to look at.** If the planner *estimated* 5 rows but the query *actually* returned 50,000, its whole plan was probably wrong (it likely chose a Nested Loop join — fine for 5 rows, terrible for 50,000) — this points straight at stale statistics (3.4).

## 3.3 Index Types — and when to reach for each

`CREATE INDEX` defaults to a **B-tree**, which is correct for the vast majority of cases (equality, ranges, sorting) — but Postgres has specialized index types for specialized query patterns:

| Index type | Good for | Example |
|---|---|---|
| **B-tree** (default) | `=`, `<`, `>`, `BETWEEN`, `ORDER BY`, most everything | `CREATE INDEX ON orders(order_date);` |
| **Hash** | Pure equality only (`=`), slightly smaller/faster than B-tree for that one case | `CREATE INDEX ON orders USING HASH (customer_id);` |
| **GIN** (Generalized Inverted Index) | Values containing multiple "things" — JSONB keys, arrays, full-text search | `CREATE INDEX ON events USING GIN (payload);` |
| **GiST** | Geometric data, ranges, nearest-neighbor searches | Spatial/PostGIS queries |
| **BRIN** (Block Range Index) | Huge tables where the column correlates with physical row order (e.g. an append-only `created_at`) — tiny index size | `CREATE INDEX ON logs USING BRIN (created_at);` |

**Beyond the type, three refinements matter in practice:**

- **Partial index** — index only the rows you actually query, when most rows are irrelevant:
  ```sql
  CREATE INDEX idx_active_orders ON orders(customer_id) WHERE status = 'active';
  ```
- **Expression index** — index the *result* of a function, when you filter on a transformed value:
  ```sql
  CREATE INDEX idx_lower_email ON customers (LOWER(email));
  -- now WHERE LOWER(email) = 'alice@example.com' can use the index
  ```
- **Covering index (`INCLUDE`)** — add extra columns to the index just so a query can be answered from the index alone (an "Index Only Scan"), without touching the table:
  ```sql
  CREATE INDEX idx_orders_covering ON orders(customer_id) INCLUDE (order_date);
  ```

## 3.4 Statistics and the Planner

The planner's cost estimates come from **table statistics** — a sampled snapshot of how many rows exist, how values are distributed, how selective each column is. These go stale after enough writes:

```sql
ANALYZE orders;                 -- refresh statistics for one table
SELECT * FROM pg_stats WHERE tablename = 'orders';  -- inspect what the planner knows
```

Postgres runs this automatically via **autovacuum's analyze** step under normal settings — but after a bulk load or a big schema change, running `ANALYZE` manually before your next query is cheap insurance against a bad plan.

## 3.5 The Optimization Playbook

Concrete rewriting patterns, roughly in the order to try them:

1. **Never `SELECT *` in application code.** Fetching columns you don't need wastes I/O and can defeat an otherwise-perfect covering index.
2. **Avoid wrapping an indexed column in a function** in your `WHERE` clause (`WHERE LOWER(email) = ...` can't use a plain index on `email` — it needs the expression index from 3.3, or better, normalize the data at write time).
3. **Avoid the N+1 query problem** — the classic app-code performance bug:
   ```go
   // BAD: one query to get orders, then N more queries in a loop — N+1 round trips
   orders, _ := pool.Query(ctx, "SELECT id, customer_id FROM orders")
   for orders.Next() {
       var id, custID int
       orders.Scan(&id, &custID)
       pool.QueryRow(ctx, "SELECT name FROM customers WHERE id = $1", custID) // one query per order!
   }
   ```
   ```go
   // GOOD: one JOIN, one round trip
   rows, _ := pool.Query(ctx, `
       SELECT o.id, c.name FROM orders o JOIN customers c ON c.id = o.customer_id`)
   ```
4. **Use `EXPLAIN ANALYZE` before and after every change** — don't guess, measure (3.2).
5. **Batch writes** instead of one `INSERT`/round trip per row — use multi-row `INSERT`, or `COPY` for bulk loads.
6. **Add the index that matches your actual `WHERE`/`JOIN`/`ORDER BY` columns** — go back to 2.11/3.3 once `EXPLAIN` shows a `Seq Scan` where you expected an `Index Scan`.

## 3.6 VACUUM and MVCC Bloat

Postgres uses **MVCC** (Multi-Version Concurrency Control, detailed in Part 4) — an `UPDATE` doesn't overwrite a row in place, it writes a new row version and marks the old one as dead, so concurrent readers can still see a consistent snapshot (this is *why* Postgres reads don't block writes). Those dead row versions pile up ("**bloat**") unless something cleans them up.

```sql
VACUUM orders;             -- reclaim space from dead row versions
VACUUM ANALYZE orders;     -- also refresh statistics (3.4) while you're at it
VACUUM FULL orders;        -- rewrites the table to reclaim space more aggressively — locks the table, use sparingly
```

**`autovacuum`** does this automatically in the background under normal configuration — don't disable it. A table with heavy `UPDATE`/`DELETE` traffic and a growing on-disk size despite a stable row count is the classic symptom of vacuum falling behind.

## 3.7 Join Algorithms — what the planner is actually choosing between

When the planner picks how to execute a `JOIN`, it's choosing between three algorithms:

- **Nested Loop Join** — for each row in the outer table, scan the inner table (or its index) for matches. Cheap when one side is small or there's a good index on the join column; catastrophic (O(n×m)) on two large tables with no index.
- **Hash Join** — build an in-memory hash table from the smaller table's join column, then stream the larger table through it, probing the hash table for matches. Good default for large, unsorted tables with an equality join condition.
- **Merge Join** — if both inputs are already sorted on the join column (or the planner sorts them), walk both in lockstep like a zipper. Good when both sides are already sorted or indexed that way.

You don't choose these directly — the planner picks based on table size, available indexes, and statistics — but `EXPLAIN ANALYZE` (3.2) tells you which one it chose, and a surprising choice (a Nested Loop on two huge unindexed tables) is a direct pointer to "add an index" or "run `ANALYZE`."

## 3.8 Scaling Reads: Pooling, Replicas, Partitioning

- **Connection pooling** — opening a new database connection per request is slow and will exhaust Postgres's max-connections limit under real traffic. A pool keeps a set of open connections ready to reuse. In Go, `pgxpool` (2.15) does this automatically once configured with sensible limits (`MaxConns`).
- **Read replicas** — for read-heavy workloads, replicate the database to one or more secondary servers and send `SELECT` traffic there, keeping writes on the primary. Trades a small replication lag for a lot more read capacity.
- **Table partitioning** — for a single table that's grown huge (e.g. an events table with a billion rows), split it into smaller physical pieces by range or list, while querying it as one logical table:
  ```sql
  CREATE TABLE events (id SERIAL, created_at DATE, payload JSONB) PARTITION BY RANGE (created_at);
  CREATE TABLE events_2026_01 PARTITION OF events FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');
  ```
  Queries filtered on `created_at` only scan the relevant partition(s) — the planner prunes the rest automatically.

## 3.9 Troubleshooting Checklist — Beginner to Advanced Workflow

When a query is slow, work through this in order:

1. Run `EXPLAIN ANALYZE` (3.2). Don't guess.
2. Compare estimated vs. actual row counts. Big gap → run `ANALYZE` (3.4).
3. Is there a `Seq Scan` on a large table with a selective `WHERE`? → add an index (3.3).
4. Is the index there but unused anyway? → check for a function wrapped around the column, or a type mismatch, that's defeating it (3.5).
5. Is it one query called in a loop from application code? → that's N+1 (3.5), fix in app code, not SQL.
6. Is the plan reasonable but the table is just enormous? → consider partitioning (3.8) or a read replica (3.8).
7. Is disk usage growing faster than row count? → check `autovacuum` isn't falling behind (3.6).
8. Still slow with a good plan? → the bottleneck may be concurrency/locking, not the query itself — that's Part 4.

---
# PART 4 — Transaction Management, Concurrency & Recovery

This is the theory that makes a database trustworthy under real conditions: multiple users hitting it at once, and machines that eventually crash. Everything here is "why your data doesn't get corrupted" — read it once carefully, because it's the part most backend engineers half-understand and it's the part that bites hardest in production.

## 4.1 Transaction Concepts and States

A **transaction** is a logical unit of work — one or more read/write operations that must be treated as a single, indivisible action. "Transfer $100 from account A to account B" is two `UPDATE`s that must succeed or fail *together* (see 2.9's `BEGIN`/`COMMIT` pattern) — that pairing is a transaction.

Every transaction moves through a defined set of states:

```mermaid
stateDiagram-v2
    [*] --> Active
    Active --> PartiallyCommitted: final statement executed
    PartiallyCommitted --> Committed: changes durably saved
    Active --> Failed: error detected
    PartiallyCommitted --> Failed: hardware/system failure before commit completes
    Failed --> Aborted: rolled back to pre-transaction state
    Committed --> [*]
    Aborted --> [*]
```

- **Active** — the transaction is executing its statements.
- **Partially committed** — the last statement has run, but changes aren't yet guaranteed durable on disk.
- **Committed** — the transaction finished successfully and its effects are permanent (Durability, 4.2).
- **Failed** — something went wrong (a constraint violation, a deadlock victim, a crash) before commit completed.
- **Aborted** — the system has rolled back all the transaction's changes; the database looks as if it never ran.

A transaction that aborts can typically be **restarted** (retried from scratch) or, if the failure was a logic error rather than a transient one, must be killed — this distinction matters directly for how you write retry logic around deadlocks in 4.8.

## 4.2 Properties of Transactions — ACID

ACID is the four guarantees a transaction should provide. Each property solves a specific failure mode: a crash mid-write, invalid data slipping in, two transactions colliding, or a server losing power right after a commit.

### Atomicity — All or Nothing

**Rule:** a transaction either completes fully or has no effect at all. No partial writes.

**The problem without it:** transferring money as two separate statements —

```sql
UPDATE accounts SET balance = balance - 100 WHERE id = 1;
-- crash happens here
UPDATE accounts SET balance = balance + 100 WHERE id = 2;
```

Money vanishes — deducted from account 1, never credited to account 2.

**Implementation:** wrap related statements in a transaction block:

```sql
BEGIN;
UPDATE accounts SET balance = balance - 100 WHERE id = 1;
UPDATE accounts SET balance = balance + 100 WHERE id = 2;
COMMIT;
```

If anything fails before `COMMIT`, roll back explicitly (or let the engine do it on error):

```sql
BEGIN;
UPDATE accounts SET balance = balance - 100 WHERE id = 1;
-- error detected
ROLLBACK;
```

Most ORMs expose this as a session/transaction context — wrap related writes in a single `db.transaction(() => { ... })` call rather than firing separate queries.

### Consistency — Valid States Only

**Rule:** a transaction can only bring the database from one valid state to another — it must not violate constraints, types, or business rules.

**Implementation:** push rules into the schema itself so the database enforces them, rather than trusting application code:

```sql
CREATE TABLE accounts (
    id INT PRIMARY KEY,
    balance DECIMAL(10,2) NOT NULL CHECK (balance >= 0),
    owner_id INT NOT NULL REFERENCES customers(id)
);
```

`CHECK (balance >= 0)` blocks an overdraft-producing update outright. `REFERENCES` blocks orphaned rows. `NOT NULL` and typed columns block malformed data at write time. With this in place, the transfer transaction above fails and rolls back automatically if it would ever push `balance` negative — you don't need to check that manually in application code.

### Isolation — Concurrent Transactions Don't Interfere

**Rule:** transactions running at the same time shouldn't see each other's uncommitted changes, and the end result should look as if they ran one after another. This property is big enough that it gets its own full treatment in 4.5–4.10 (schedules, serializability, locking, isolation levels) — the short version:

```sql
SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;
```

```sql
BEGIN;
SELECT balance FROM accounts WHERE id = 1 FOR UPDATE;  -- locks the row
UPDATE accounts SET balance = balance - 10 WHERE id = 1;
COMMIT;
```

`FOR UPDATE` forces a second, concurrent transaction to wait until this one commits, instead of both reading the same stale value and silently losing an update.

### Durability — Committed Means Permanent

**Rule:** once a transaction commits, the change survives a crash, power loss, or restart — even if it hasn't been "seen" anywhere else yet. This is mostly handled by the engine (4.11–4.13 cover exactly how), but you configure how aggressively it protects a commit:

- **Write-Ahead Logging (WAL)** — the engine writes the intended change to a durable log *before* touching the actual data files, so it can replay or roll back correctly after a crash. Enabled by default in Postgres — don't disable it for a performance gain unless you fully understand the risk (4.13).
- **fsync on commit** — ensures the OS actually flushes to disk rather than just to a cache: `SET synchronous_commit = on;` (turning this off trades durability for speed — valid for non-critical logging tables, dangerous for financial data).
- **Replication** — commit changes to a replica before acknowledging success (synchronous) for durability beyond a single disk failing, or accept a small window of risk with asynchronous replication.

### Putting It Together

```sql
CREATE TABLE accounts (
    id INT PRIMARY KEY,
    balance DECIMAL(10,2) NOT NULL CHECK (balance >= 0)   -- Consistency
);

SET TRANSACTION ISOLATION LEVEL READ COMMITTED;            -- Isolation level chosen

BEGIN;                                                      -- Atomicity via transaction block
    SELECT balance FROM accounts WHERE id = 1 FOR UPDATE;  -- Isolation via row lock
    UPDATE accounts SET balance = balance - 100 WHERE id = 1;
    UPDATE accounts SET balance = balance + 100 WHERE id = 2;
COMMIT;                                                      -- Durability guaranteed by WAL + fsync from here
```

If the `CHECK` constraint fails, the whole block rolls back (Consistency protecting Atomicity). If the process crashes right after `COMMIT` returns, the change is still there on restart (Durability).

| Property | Guarantees | How you implement it |
|---|---|---|
| Atomicity | All-or-nothing execution | `BEGIN` / `COMMIT` / `ROLLBACK` |
| Consistency | Only valid states are ever saved | Constraints, foreign keys, types, triggers |
| Isolation | Concurrent transactions don't corrupt data | Isolation levels, row locks |
| Durability | Commits survive crashes | WAL, fsync, replication |

**When to relax ACID:** strict `SERIALIZABLE` isolation and synchronous replication cost throughput and latency. Many large-scale systems deliberately relax it — **BASE** (Basically Available, Soft state, Eventual consistency) is common in distributed/NoSQL systems where availability matters more than instant consistency across nodes (see Part 9's CAP theorem); `READ COMMITTED` is a common default over `SERIALIZABLE` for high-traffic apps that can tolerate minor anomalies for speed; async replication trades a small durability window for lower write latency. Full ACID compliance matters most where correctness is non-negotiable — financial transactions, inventory counts, anything involving money or contractual state. Elsewhere, weigh the trade-off deliberately.

## 4.3 Schedules and Serializability

A **schedule** is the actual order in which operations from one or more transactions execute. Running transactions strictly one after another — a **serial schedule** — is always safe, because there's no interleaving to cause problems, but it throws away all the performance benefit of concurrency (every other transaction waits its turn). A **concurrent (non-serial) schedule** interleaves operations from multiple transactions to get real parallelism — and the entire point of this Part is: how do we interleave safely?

A schedule is **serializable** if its result is guaranteed *equivalent* to some serial execution of the same transactions — meaning you get the performance of concurrency with the safety of running one-at-a-time. There are two common definitions of "equivalent":

- **Conflict serializability** — two operations *conflict* if they come from different transactions, act on the same data item, and at least one is a write. A schedule is conflict-serializable if you can transform it into a serial schedule purely by swapping *adjacent, non-conflicting* operations. This is the practical, checkable definition — 4.4 shows exactly how to test for it.
- **View serializability** — a broader definition: two schedules are view-equivalent if every transaction reads the same values (the same "reads-from" relationships) and produces the same final writes in both. Every conflict-serializable schedule is view-serializable, but not the reverse — view serializability additionally allows harmless "blind writes" (writing a value without reading it first) that conflict serializability would flag unnecessarily. View serializability is more general but much more expensive to test, so real database engines check for conflict serializability in practice.

## 4.4 Test of Serializability — the Precedence Graph

You test whether a schedule is conflict-serializable by building a **precedence graph** (also called a serialization graph):

1. Create one node per transaction.
2. For every pair of conflicting operations where `Ti`'s operation happens before `Tj`'s operation on the same data item, draw an edge `Ti → Tj`.
3. **If the graph has a cycle, the schedule is *not* conflict-serializable.** If it's acyclic, it *is* — and a topological sort of the graph gives you an equivalent serial order.

**Worked example** — two transactions both touching data item `A`:

```
T1: Read(A)  Write(A)             
T2:           Read(A)  Write(A)  Read(A)... T1: Write(A) again
```

Concretely, schedule `S`:  `R1(A), R2(A), W1(A), W2(A)` — T1 reads A, T2 reads A, T1 writes A, T2 writes A.

- `R1(A)` before `W2(A)` → edge `T1 → T2` (T1 read, then T2 wrote the same item).
- `R2(A)` before `W1(A)` → edge `T2 → T1` (T2 read, then T1 wrote the same item).

```mermaid
graph LR
    T1 -->|R1 before W2| T2
    T2 -->|R2 before W1| T1
```

Two edges in opposite directions is a **cycle** (`T1 → T2 → T1`) — this schedule is **not** serializable: there's no serial order (T1-then-T2, or T2-then-T1) that produces the same result, because each transaction's write stomped on a value the other had already read. This exact pattern is the **lost update** problem from 4.6 in formal terms. Compare a *safe* interleaving, `R1(A), W1(A), R2(A), W2(A)`: only one edge (`T1 → T2`), no cycle, conflict-serializable, equivalent to running T1 then T2 serially — which is exactly what row locking (4.7) forces to happen.

## 4.5 Problems of Concurrent Execution

These are the concrete anomalies that motivate everything from here through 4.10 — each is a specific way uncontrolled interleaving corrupts results:

- **Lost Update** — two transactions read the same value, both compute a new value based on it, and the second write silently overwrites the first:
  ```sql
  -- Transaction A (reads 100)         -- Transaction B (reads 100)
  SELECT balance FROM accounts WHERE id = 1;   SELECT balance FROM accounts WHERE id = 1;
  UPDATE accounts SET balance = 90 WHERE id=1; UPDATE accounts SET balance = 80 WHERE id=1;
  COMMIT;                                       COMMIT;
  -- final balance: 80 — A's deduction silently lost
  ```
- **Dirty Read** — a transaction reads a value another transaction wrote but hasn't committed yet; if that other transaction rolls back, the first transaction acted on data that never really existed.
- **Unrepeatable (Non-Repeatable) Read** — a transaction reads the same row twice and gets two different values, because another transaction committed a change to it in between.
- **Phantom Read** — a transaction re-runs the same range query (`WHERE amount > 50`) twice and gets a *different set of rows*, because another transaction inserted or deleted a matching row in between.

Each of these is prevented by a specific isolation level (4.10) and, underneath that, by specific locking behavior (4.6–4.9).

## 4.6 Locking Mechanisms

A **lock** is how a database physically stops the anomalies in 4.5: a transaction that wants to touch a data item must acquire a lock on it first.

**The two basic lock types:**

| Lock | Held for | Compatible with |
|---|---|---|
| **Shared (S)** | Reading a data item | Other Shared locks — many readers can hold an S lock on the same item simultaneously |
| **Exclusive (X)** | Writing a data item | Nothing — an X lock blocks every other lock, shared or exclusive, on that item |

**Compatibility matrix:**

| | S requested | X requested |
|---|---|---|
| **S held** | ✅ Granted | ❌ Must wait |
| **X held** | ❌ Must wait | ❌ Must wait |

**Lock granularity** is the size of what you lock — a single row, a page, a whole table, or the whole database. Coarser granularity (table-level) means less bookkeeping overhead but more contention (everyone blocks everyone); finer granularity (row-level) means more overhead but far less contention. Postgres locks at the row level for normal writes by default:

```sql
SELECT * FROM accounts WHERE id = 1 FOR UPDATE;   -- row-level exclusive lock
SELECT pid, mode, relation::regclass, granted FROM pg_locks WHERE relation = 'accounts'::regclass;
```

## 4.7 Two-Phase Locking Protocol (2PL)

2PL is the rule that *guarantees* conflict serializability (4.3) when transactions acquire and release locks correctly:

- **Growing phase** — a transaction may *acquire* locks, but may not release any.
- **Shrinking phase** — once a transaction releases its *first* lock, it may not acquire any more.

```mermaid
graph LR
    A["Growing Phase<br/>acquire locks only"] --> B["Lock point<br/>(max locks held)"] --> C["Shrinking Phase<br/>release locks only"]
```

Basic 2PL guarantees serializability but has a real problem: if a transaction releases a lock and then *aborts*, any other transaction that already read the now-invalid value has to be rolled back too — a **cascading rollback**.

- **Strict 2PL** fixes this: hold all **exclusive** locks until the transaction actually commits or aborts (never release early). This is what Postgres and most production databases actually do — it avoids cascading rollbacks because nobody can read your uncommitted write in the first place.
- **Rigorous 2PL** goes further: hold *all* locks — shared and exclusive — until commit/abort. Simpler to reason about, slightly more conservative.

**Practical note on Postgres specifically:** Postgres doesn't use classic 2PL for reads — it uses **MVCC** (3.6): readers see a consistent snapshot without taking read locks at all, so reads never block writes and writes never block reads. Writes still take row-level exclusive locks in a strict-2PL style (held until commit) to prevent two transactions from stomping on the same row. This is *why* Postgres reads are so unblockable, and it's what makes `FOR UPDATE` necessary when you specifically *do* need a read to block concurrent writers (as in 4.2's transfer example).

## 4.8 Deadlock

A **deadlock** is a circular wait: Transaction A holds a lock Transaction B needs, while Transaction B holds a lock Transaction A needs — neither can proceed. Four conditions must *all* hold for deadlock to be possible: **mutual exclusion** (a resource can only be held by one transaction at a time), **hold and wait** (a transaction holds a lock while waiting for another), **no preemption** (a lock can't be forcibly taken away), and **circular wait** (a cycle of transactions each waiting on the next).

**Detection** — build a **wait-for graph**: an edge `Ti → Tj` if `Ti` is waiting for a lock held by `Tj`. A cycle in this graph means deadlock. Postgres runs this check automatically (`deadlock_timeout`, default 1 second) and kills one transaction (the "victim," usually the one that would be cheapest to abort) to break the cycle, returning an error the application must catch and retry.

**Practical demo** — two sessions locking rows in *opposite* order will deadlock:

```sql
-- Session 1                              -- Session 2
BEGIN;                                    BEGIN;
UPDATE accounts SET balance=90 WHERE id=1; UPDATE accounts SET balance=80 WHERE id=2;
UPDATE accounts SET balance=90 WHERE id=2; -- waits on Session 2's lock on row 2
                                           UPDATE accounts SET balance=80 WHERE id=1; -- waits on Session 1's lock on row 1
-- deadlock detected; Postgres aborts one transaction with:
-- ERROR: deadlock detected
```

**Prevention** (avoid the *possibility* of deadlock, usually by ordering) vs. **avoidance** (allow it, decide who wins based on runtime information):

- **Wait-Die** and **Wound-Wait** — timestamp-based prevention schemes: assign every transaction a timestamp at start. In Wait-Die, an older transaction is allowed to wait for a younger one, but a younger transaction requesting a lock held by an older one is forced to abort ("die") and restart. Wound-Wait is the mirror image: an older transaction "wounds" (forces the abort of) a younger one holding a lock it needs, while a younger transaction simply waits for an older one.
- **The practical technique for application developers: lock ordering.** Always acquire locks on multiple rows in the *same, consistent global order* (e.g., always lock the lower `id` first) across every code path in your application. This eliminates the circular-wait condition entirely — no deadlock detector needed:
  ```sql
  -- Always update the lower id first, everywhere in the codebase
  BEGIN;
  UPDATE accounts SET balance = balance - 100 WHERE id = LEAST(1, 2);
  UPDATE accounts SET balance = balance + 100 WHERE id = GREATEST(1, 2);
  COMMIT;
  ```
- **App-level fallback:** since deadlocks can still happen from timing you don't control, always wrap transactions prone to lock contention in retry logic that catches the deadlock error and retries the whole transaction from `BEGIN`.

## 4.9 Intent Locking and Multiple Granularity

4.6 mentioned granularity as a trade-off — but it creates a practical problem: if Transaction A wants to lock the *entire* `accounts` table, how does it quickly check that no other transaction already holds a lock on some individual *row* inside it, without scanning every row to check?

**Intent locks** solve this by requiring a transaction to announce its intent at every coarser level *above* the row it actually wants to lock, so a conflict check only has to look at the top level:

| Lock | Meaning |
|---|---|
| **IS** (Intent Shared) | "I intend to place S locks on some finer-grained items below this level" |
| **IX** (Intent Exclusive) | "I intend to place X locks on some finer-grained items below this level" |
| **SIX** (Shared + Intent Exclusive) | "I hold this whole level with S, and intend to place X locks on some finer items within it" |

**The protocol:** to lock a row `X`, a transaction must first acquire an intent lock (IS or IX) on the table containing it, and on the database containing that table, walking from coarse to fine. Now, a transaction that wants to lock the *entire table* with a plain S or X lock only has to check the table-level lock state — if any other transaction holds an IX there, it knows *immediately*, without scanning individual rows, that a conflicting row-level lock exists somewhere below.

**Compatibility matrix** (rows = held, columns = requested):

| | IS | IX | S | SIX | X |
|---|---|---|---|---|---|
| **IS** | ✅ | ✅ | ✅ | ✅ | ❌ |
| **IX** | ✅ | ✅ | ❌ | ❌ | ❌ |
| **S**  | ✅ | ❌ | ✅ | ❌ | ❌ |
| **SIX**| ✅ | ❌ | ❌ | ❌ | ❌ |
| **X**  | ❌ | ❌ | ❌ | ❌ | ❌ |

**Practical Postgres implementation:** this is exactly what Postgres's table-level lock modes are, running automatically underneath every row-level lock you take. When you `UPDATE` a row, Postgres silently takes a `ROW EXCLUSIVE` (an IX-style intent lock) on the *table* in addition to the exclusive lock on the *row* — so an operation like `LOCK TABLE accounts` (a table-wide X-style lock) can immediately see it conflicts, without scanning rows. Inspect this directly:

```sql
SELECT locktype, relation::regclass, mode, granted FROM pg_locks WHERE pid = pg_backend_pid();
-- you'll see both a row-level lock AND a table-level "ROW EXCLUSIVE" intent-style lock from one UPDATE
```

## 4.10 Isolation Levels in Practice

Tying 4.5's anomalies directly to what you actually configure:

| Level | Prevents | Allows | Trade-off |
|---|---|---|---|
| **Read Uncommitted** | Nothing | Dirty, unrepeatable, phantom reads | Fastest, essentially unused in practice (Postgres treats this the same as Read Committed) |
| **Read Committed** (default in Postgres) | Dirty reads | Unrepeatable reads, phantom reads | Good default for most apps |
| **Repeatable Read** | Dirty + unrepeatable reads | Phantom reads (Postgres's implementation actually also prevents phantoms via MVCC snapshots) | Row values stay stable for the whole transaction |
| **Serializable** | All of the above | Nothing | Safest, most contention/retries |

```sql
SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;
```

Postgres implements `SERIALIZABLE` using **Serializable Snapshot Isolation (SSI)** — rather than locking everything pessimistically, it lets transactions run concurrently against a snapshot and detects at commit time whether the interleaving *would* have violated serializability (4.3/4.4); if so, it aborts one transaction with a serialization error your application must retry. This gets you the correctness of full serializability with less locking overhead than naive 2PL, at the cost of needing retry logic in your app for the (usually rare) abort case.

## 4.11 System Recovery — Types of Failure

Before covering *how* recovery works, know *what* it's recovering from:

- **Transaction failure** — a logical error (a `CHECK` constraint violation) or a system-detected error (a deadlock victim, 4.8) aborts one transaction. The rest of the system is fine; that one transaction just needs to roll back.
- **System crash** — the whole database process (or OS) crashes: power loss, kernel panic, container killed. Everything in **volatile storage** (RAM — including in-progress transactions and anything not yet flushed to disk) is lost, but **non-volatile storage** (disk) is intact.
- **Disk failure** — non-volatile storage itself is damaged or lost. Recovering from this requires a *separate* copy of the data — backups or replicas — not just the log on the same failed disk. This is why "stable storage" in recovery theory really means storage that's been replicated (e.g., RAID, or a physically separate backup) so it's treated as effectively never lost.

## 4.12 Recovery and Atomicity — Undo and Redo

To guarantee Atomicity (4.2) and Durability survive a crash, the database needs a way to reconstruct "what should the data look like" after an unplanned restart. Two operations do this:

- **Undo** — roll back the changes of a transaction that was active (not yet committed) when the crash happened, since it violated Atomicity by being left half-done.
- **Redo** — reapply the changes of a transaction that *had* committed, but whose changes hadn't yet been physically written from memory to the actual data files on disk when the crash happened — otherwise a committed transaction's effects would be lost, violating Durability.

Which one a given transaction needs depends on when its changes actually touched the on-disk data files relative to when it committed — which is exactly what the logging strategy in 4.13 is designed to make cheap and unambiguous to determine.

## 4.13 Log-Based Recovery

**The core idea:** before any change touches the actual data files, write a record of that change to an append-only **log** on stable storage first — this is **Write-Ahead Logging (WAL)**, already introduced in 4.2's Durability section. If the system crashes, the log is the ground truth for reconstructing exactly what happened.

**Log record format**, conceptually:

```
<Ti, start>
<Ti, data_item, old_value, new_value>
<Ti, commit>          -- or <Ti, abort>
```

Each write logs the transaction ID, what changed, and both the old and new value — the old value is what makes **undo** possible, the new value is what makes **redo** possible.

**Two strategies for when the actual data file gets modified:**

- **Deferred database modification** — don't touch the actual data files until *after* the transaction commits; only the log is written during execution. Recovery only ever needs **redo** (replay committed transactions the crash caught before their writes hit disk) — never undo, since nothing uncommitted ever reached the data files in the first place.
- **Immediate database modification** — write to the actual data files *during* the transaction, before commit (this is what real systems do, for performance — you can't defer everything). Recovery needs **both undo** (for transactions active at crash time) **and redo** (for committed transactions not yet flushed).

**Checkpoints** — without them, recovery would have to replay the *entire* log from the beginning of time on every restart. A checkpoint periodically records "everything up to this point in the log is already safely reflected in the data files," so recovery only needs to start from the most recent checkpoint:

```sql
CHECKPOINT;   -- Postgres also does this automatically on a schedule (checkpoint_timeout, max_wal_size)
```

**The recovery algorithm on restart** (a simplified version of the industry-standard ARIES algorithm), scanning forward from the last checkpoint:

1. **Analysis pass** — scan the log from the last checkpoint to determine which transactions were active (uncommitted) at crash time, and which had committed.
2. **Redo pass** — reapply *every* logged change since the checkpoint, for both committed and uncommitted transactions — this repeats history to reconstruct the exact state right before the crash, including uncommitted work (deliberately — it's undone next).
3. **Undo pass** — roll back the changes of transactions that were still active (not committed) at crash time, using the old values logged for each.

**Practical Postgres implementation:**

```bash
# WAL files live here — this IS the log described above
ls /var/lib/postgresql/data/pg_wal/

# Inspect WAL contents directly
pg_waldump /var/lib/postgresql/data/pg_wal/000000010000000000000001
```

```sql
-- Checkpoint tuning
SHOW checkpoint_timeout;   -- how often Postgres checkpoints on a schedule
SHOW max_wal_size;         -- checkpoint sooner if WAL grows past this
```

When Postgres restarts after a crash, it automatically performs exactly the analysis/redo/undo sequence above using the WAL — this is why a Postgres container that gets killed `-9` and restarted comes back with all committed data intact and no half-finished transactions, with no manual intervention required.

## 4.14 Two-Phase Commit Protocol (2PC)

Everything so far assumes one database. **Distributed transactions** span multiple databases or services that each need to commit (or all abort) together — e.g., debiting an account in one database while writing an audit record in another. 2PC is the classical protocol for this:

```mermaid
sequenceDiagram
    participant C as Coordinator
    participant P1 as Participant 1
    participant P2 as Participant 2
    Note over C,P2: Phase 1 — Voting / Prepare
    C->>P1: PREPARE
    C->>P2: PREPARE
    P1-->>C: Yes (logged, ready to commit)
    P2-->>C: Yes (logged, ready to commit)
    Note over C,P2: Phase 2 — Commit
    C->>P1: COMMIT
    C->>P2: COMMIT
    P1-->>C: Acknowledged
    P2-->>C: Acknowledged
```

- **Phase 1 (voting/prepare)** — the coordinator asks every participant "can you commit?" Each participant does the actual work, writes it to its *own* log, and replies yes ("prepared") or no — but does **not** commit yet, even after voting yes.
- **Phase 2 (commit)** — if *every* participant voted yes, the coordinator tells all of them to commit for real; if *any* voted no, it tells all of them to abort.

**The blocking problem:** a participant that voted "yes" and is waiting for the coordinator's final decision is now in an "in-doubt" state — it's holding locks and can't unilaterally decide to commit or abort, because the coordinator might have told other participants the opposite. If the coordinator crashes at exactly this moment, that participant *must block* — holding its locks — until the coordinator recovers. This is 2PC's fundamental weakness, and it's exactly why it's largely avoided at large scale in favor of eventual-consistency patterns (Sagas: a sequence of local transactions each with a compensating "undo" action, rather than one distributed atomic commit).

**Practical Postgres implementation** — `PREPARE TRANSACTION` / `COMMIT PREPARED` implement exactly this two-phase protocol for distributed transactions coordinated across multiple Postgres instances:

```sql
BEGIN;
UPDATE accounts SET balance = balance - 100 WHERE id = 1;
PREPARE TRANSACTION 'transfer_123';   -- Phase 1: prepared, but not committed

-- later, once every participant has confirmed it's prepared:
COMMIT PREPARED 'transfer_123';       -- Phase 2: commit
-- or: ROLLBACK PREPARED 'transfer_123';
```

`max_prepared_transactions` must be configured above 0 to use this, and a prepared-but-never-resolved transaction holds its locks indefinitely (the blocking problem above) — so a real coordinator process must reliably follow through with commit or rollback.

## 4.15 Practical Postgres Recovery

Putting 4.13 to use for real disaster recovery, not just crash recovery:

```bash
# Logical backup (a SQL dump — good for smaller databases, migrations)
pg_dump -U postgres mydb > backup.sql

# Physical backup (a full copy of the data directory — needed for PITR below)
pg_basebackup -U postgres -D /backup/base -Fp -Xs -P
```

**Point-in-Time Recovery (PITR)** uses a base backup plus *archived* WAL files to restore the database to any exact moment — not just the moment of the last backup — by replaying WAL forward from the base backup up to (but not including) a bad transaction, e.g. right before an accidental `DELETE FROM accounts;` with no `WHERE` clause:

```
# postgresql.conf
archive_mode = on
archive_command = 'cp %p /archive/%f'
```

```
# recovery target, when restoring
recovery_target_time = '2026-09-05 14:29:00'
```

This is the practical payoff of everything in Part 4: because every committed change was written to the WAL *before* it touched the data files (4.13), that log is a complete, replayable history — which is what makes "restore to exactly one second before someone ran a destructive query" possible at all.

---
# PART 5 — Security: Access Control and Encryption

A correct, fast, crash-safe database is still a liability if anyone can read or write anything. This Part covers who can touch what, and how to protect data even if someone gets past that.

## 5.1 Access Control Models

| Model | How it decides | Typical implementation |
|---|---|---|
| **DAC** (Discretionary Access Control) | The owner of an object grants/revokes privileges to others, at their discretion | `GRANT`/`REVOKE` — Postgres's default model |
| **MAC** (Mandatory Access Control) | The *system* enforces access based on fixed security labels/clearance levels — no owner discretion | OS-level (SELinux), rare in typical app databases |
| **RBAC** (Role-Based Access Control) | Permissions attach to *roles*; users are assigned roles, not individual permissions | `CREATE ROLE`, `GRANT role TO user` |
| **ABAC** (Attribute-Based Access Control) | Rules evaluate attributes of the user, the data, and the context at query time | Row-Level Security policies using session context |

**DAC/RBAC practical implementation in Postgres:**

```sql
-- Create roles for job functions, not individual people
CREATE ROLE readonly;
CREATE ROLE app_writer;

GRANT SELECT ON ALL TABLES IN SCHEMA public TO readonly;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO app_writer;

-- Create actual login users and assign them roles
CREATE ROLE reporting_bot LOGIN PASSWORD 'secret';
GRANT readonly TO reporting_bot;

CREATE ROLE api_service LOGIN PASSWORD 'secret';
GRANT app_writer TO api_service;

-- Revoke is the mirror image
REVOKE INSERT ON orders FROM app_writer;
```

**Column-level privileges** — grant access to specific columns only, e.g. letting a support role see customer names but not emails:

```sql
GRANT SELECT (id, name) ON customers TO support_role;
```

**Practical rule — principle of least privilege:** never use a superuser/admin credential in application code. Create a dedicated role per service with exactly the permissions that service needs — a reporting dashboard gets `readonly`, the API server gets `app_writer` scoped to only the tables it touches, and nobody but a human operator with a separate credential gets `DROP TABLE` rights.

## 5.2 Row-Level Security (ABAC in Practice)

RBAC controls *which tables/columns* a role can touch; **Row-Level Security (RLS)** controls *which rows* — the practical implementation of ABAC, since the policy can depend on attributes of the current session:

```sql
ALTER TABLE orders ENABLE ROW LEVEL SECURITY;

CREATE POLICY customer_sees_own_orders ON orders
    USING (customer_id = current_setting('app.current_customer_id')::INT);
```

Once enabled, every query against `orders` — even a plain `SELECT * FROM orders` — is automatically filtered to only the rows the policy allows, based on a session variable the application sets after authenticating the user:

```sql
SET app.current_customer_id = '1';
SELECT * FROM orders;   -- only returns customer 1's orders, automatically
```

This is the practical technique for true multi-tenant isolation at the database layer, rather than trusting every application query to remember a `WHERE customer_id = ?` clause.

## 5.3 Encryption

Three distinct layers, each protecting against a different threat:

- **Encryption in transit** — protects data moving between the application and the database over the network from interception. Enable TLS:
  ```
  # postgresql.conf
  ssl = on
  ssl_cert_file = 'server.crt'
  ssl_key_file = 'server.key'
  ```
  ```go
  // connection string — require and verify TLS
  "postgres://user:pass@host:5432/db?sslmode=verify-full"
  ```
  `sslmode=require` encrypts but doesn't verify the server's identity (vulnerable to a man-in-the-middle); `verify-full` also checks the certificate matches the hostname — use `verify-full` for anything that isn't strictly local development.

- **Encryption at rest** — protects the data files on disk if the physical disk or a backup is stolen. Postgres doesn't encrypt data files natively by default; this is handled at the filesystem/volume level (LUKS, or a cloud provider's KMS-managed encrypted volume) or with a dedicated extension.

- **Column-level (application-visible) encryption** — encrypts specific sensitive fields so even someone with raw table access can't read them, using the `pgcrypto` extension:
  ```sql
  CREATE EXTENSION pgcrypto;

  CREATE TABLE customers (
      id SERIAL PRIMARY KEY,
      name VARCHAR(100),
      ssn_encrypted BYTEA   -- store ciphertext, never plaintext
  );

  INSERT INTO customers (name, ssn_encrypted)
  VALUES ('Alice', pgp_sym_encrypt('123-45-6789', 'encryption-key-from-env-var'));

  SELECT name, pgp_sym_decrypt(ssn_encrypted, 'encryption-key-from-env-var') AS ssn
  FROM customers WHERE id = 1;
  ```

## 5.4 Password Hashing vs. Encryption — a critical distinction

**Encryption is reversible** (with the key); **hashing is one-way**. Passwords must be **hashed**, never encrypted — even the database administrator should never be able to recover a user's actual password, only verify a guess against it.

```go
import "golang.org/x/crypto/bcrypt"

// Storing a password
hashed, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
// store `hashed` in the database — never the plaintext, never a reversibly-encrypted version

// Verifying a login attempt
err := bcrypt.CompareHashAndPassword(hashedFromDB, []byte(attemptedPassword))
if err == nil {
    // password correct
}
```

`bcrypt` (or `argon2`) automatically **salts** each hash, so two users with the same password get different stored hashes — this defeats precomputed rainbow-table attacks. Never write your own hashing scheme, and never use a fast general-purpose hash like plain SHA-256 for passwords — it's *designed* to be fast, which is exactly the wrong property for something an attacker will try to brute-force.

## 5.5 Key Management Basics

- Never hardcode encryption keys, database passwords, or API credentials in source code — pull them from environment variables or a dedicated secrets manager (AWS Secrets Manager, HashiCorp Vault, etc.).
- Rotate credentials periodically, and immediately after any suspected exposure.
- Separate keys by environment (dev/staging/prod) so a leaked development key can't touch production data.

## 5.6 Security in Redis and MongoDB

**Redis:**
```bash
# redis.conf
requirepass a-strong-password-here
```
```bash
redis-cli -a a-strong-password-here
```
Redis 6+ also supports fine-grained **ACLs** (a practical RBAC implementation for Redis):
```
ACL SETUSER readonly_app on >password ~cache:* +get +exists
```
This creates a user restricted to `GET`/`EXISTS` commands on keys matching `cache:*` only — the same least-privilege principle as 5.1, applied to Redis's command/key model. Enable TLS (`tls-port`) for encryption in transit, same rationale as 5.3.

**MongoDB:**
```javascript
// Enable authentication, create a role-scoped user
db.createUser({
  user: "app_writer",
  pwd: "secret",
  roles: [{ role: "readWrite", db: "myapp" }]   // RBAC, same principle as 5.1
})
```
MongoDB uses **SCRAM** authentication by default (never sends the plaintext password over the wire), supports TLS for in-transit encryption, and offers **field-level encryption** for the column-level-encryption equivalent from 5.3 — specific document fields stay encrypted even in memory on the server, decryptable only by a client holding the key.

---
# PART 6 — Redis: Building and Querying

Redis ("**Re**mote **Di**ctionary **S**erver") is an in-memory data store that keeps its entire dataset in RAM, which is why reads and writes typically complete in well under a millisecond — no disk seek in the hot path. It's usually reached for as a cache, a session store, a queue, or a leaderboard, not as a primary source of truth. This guide reuses the same `customers`/`orders` running example to show where Redis sits *alongside* Postgres, not instead of it.

```bash
docker exec -it redis redis-cli
```

## 6.1 What Redis Actually Is

Unlike SQL or MongoDB, Redis has no tables, no schema, and (mostly) no query language — just **keys mapped to values**, where the value can be one of several structures (string, list, hash, set, sorted set...). Everything lives in memory, which means:

- Speed is exceptional, but **your dataset can't exceed available RAM** without extra configuration.
- Data is volatile by default — a crash loses everything **unless** persistence is configured (6.9).
- Redis is **single-threaded** for command execution. This is actually why it's so predictable: no locks, no race conditions between commands, every command runs to completion before the next one starts (its own tiny, built-in version of Atomicity from 4.2). The trade-off: one slow command (scanning a huge dataset) blocks every other client until it finishes.

**Basic commands, usable on any key regardless of type:**

```bash
SET session:1 "active"       # create/overwrite a key
GET session:1                 # read it back
EXISTS session:1              # 1 if it exists, 0 if not
DEL session:1                 # delete it
EXPIRE session:1 3600         # auto-delete after 3600 seconds
TTL session:1                 # seconds remaining before expiry (-1 = no expiry, -2 = doesn't exist)
```

Avoid `KEYS *` in production — it scans the entire keyspace and blocks the single thread while doing it. Use `SCAN` instead, which walks the keyspace incrementally without blocking:

```bash
SCAN 0
```

## 6.2 Strings

The simplest structure — text, a number, or serialized JSON. The basis for simple caching and counters.

```bash
SET customer:1:name "Alice"
GET customer:1:name                 # "Alice"

SET page_views 0
INCR page_views                     # atomic increment, now 1
INCRBY page_views 10                # now 11
DECR page_views                     # now 10

SETEX cache:customer:1 300 "{...}"  # set with a 300-second expiry in one command
MSET a 1 b 2 c 3                    # set multiple keys at once
MGET a b c                          # get multiple keys at once
```

`INCR`/`DECR` are **atomic** — safe for concurrent counters (page views, rate limiters) even with many clients hitting the same key at once, no race condition (a direct consequence of 6.1's single-threaded execution).

## 6.3 Lists — Ordered Collections

A linked list of values, efficient for pushing/popping from either end — the natural structure for queues and recent-activity feeds.

```bash
LPUSH recent_orders:1 "order:104"   # push to the left (front)
RPUSH recent_orders:1 "order:105"   # push to the right (back)
LRANGE recent_orders:1 0 -1         # get the whole list (0 to -1 = all)
LPOP recent_orders:1                # remove and return the leftmost item
LLEN recent_orders:1                # length of the list
```

**Use case — a simple job queue.** A worker blocks and waits for new work instead of polling:

```bash
LPUSH job_queue "process_order:104"
BRPOP job_queue 0                    # blocks until an item is available, 0 = wait forever
```

## 6.4 Hashes — Objects Within a Key

A hash maps field names to values *within a single key* — useful for representing a record like a customer without needing a separate `SET` key per field.

```bash
HSET customer:1 name "Alice" city "Boston" total_spent "70.00"
HGET customer:1 name                  # "Alice"
HGETALL customer:1                    # all fields and values
HDEL customer:1 city                  # remove one field
HINCRBY customer:1 total_spent 20     # atomically add to a numeric field
```

This is closer to the "row" concept from SQL/MongoDB than a plain string — one key, multiple named fields, updated independently.

## 6.5 Sets — Unique, Unordered Collections

A set stores unique values with no order, and supports fast membership checks plus mathematical set operations.

```bash
SADD customer:1:tags "vip" "boston" "repeat_buyer"
SISMEMBER customer:1:tags "vip"       # 1 (true) or 0 (false)
SMEMBERS customer:1:tags               # all tags
SREM customer:1:tags "vip"             # remove a tag

SADD boston_customers 1 3
SADD vip_customers 1 2
SINTER boston_customers vip_customers  # customers who are both: {1}
SUNION boston_customers vip_customers  # all of either: {1, 2, 3}
SDIFF boston_customers vip_customers   # in Boston but not VIP: {3}
```

**Use case:** deduplicating unique visitors, tagging, "friends in common" style intersections.

## 6.6 Sorted Sets — Ranked Collections

Like a set, but every member has a **score**, and Redis keeps members ordered by that score automatically — Redis's native structure for leaderboards and anything ranked.

```bash
ZADD leaderboard 70 "Alice"
ZADD leaderboard 75 "Bob"
ZADD leaderboard 10 "Carla"

ZRANGE leaderboard 0 -1 WITHSCORES     # ascending: Carla, Alice, Bob
ZREVRANGE leaderboard 0 -1 WITHSCORES  # descending: Bob, Alice, Carla
ZSCORE leaderboard "Alice"             # 70
ZRANK leaderboard "Alice"              # 0-indexed position, ascending
ZINCRBY leaderboard 20 "Carla"         # Carla's score becomes 30
ZRANGEBYSCORE leaderboard 50 100       # everyone scoring 50–100
```

This maps directly onto a query like "top 5 customers by total spend" — something that needs `GROUP BY` + `ORDER BY` + `LIMIT` in SQL happens here as a single structure kept pre-sorted at all times.

## 6.7 Expiration — Making Data Temporary

Because Redis is often used as a cache, letting keys expire automatically is core functionality — it avoids manually cleaning up stale sessions, cache entries, or rate-limit counters.

```bash
EXPIRE customer:1 3600      # this key disappears in 1 hour
PEXPIRE customer:1 3600000  # same thing, in milliseconds
TTL customer:1               # check seconds remaining
PERSIST customer:1           # cancel the expiration, key becomes permanent again

SET cache:top_customers "[...]" EX 60   # set with a 60-second expiry, one command
```

## 6.8 Pub/Sub and Streams — Messaging

**Pub/Sub** is fire-and-forget: publishers send messages to a channel; subscribers listening *at that moment* receive them in real time. No storage, no replay — if no one's listening when a message is published, it's gone.

```bash
# Client A (subscriber)
SUBSCRIBE order_updates

# Client B (publisher)
PUBLISH order_updates "order:104 shipped"
```

`PSUBSCRIBE order_*` subscribes to any channel matching a pattern. Good for live notifications, not for anything that must be durable.

**Streams** persist every message and let multiple consumer groups read (and re-read) history — conceptually a lightweight Kafka topic living inside Redis. Use this instead of Pub/Sub when a message must not be missed:

```bash
XADD orders_stream '*' customer_id 1 amount 50.00   # '*' auto-generates an ID
XRANGE orders_stream - +                              # read all entries, oldest to newest
```

## 6.9 Persistence — Making Redis Durable

By default, an unconfigured Redis loses everything on restart (6.1) — fine for a pure cache, not fine if losing that data would actually hurt. Two options, echoing the Durability discussion in 4.2:

- **RDB (snapshotting)** — periodically saves a full point-in-time snapshot to disk. Compact, fast to restart from, but you can lose everything since the last snapshot.
- **AOF (Append-Only File)** — logs every write operation as it happens, similar in spirit to the WAL from 4.13. More durable (configurable to fsync every write, every second, or never), larger on disk.

```
# redis.conf
appendonly yes
appendfsync everysec
```

**Practical rule:** if data loss would genuinely hurt, that data belongs in Postgres as the source of truth, with Redis just caching it — don't lean on Redis persistence as your only durability story.

## 6.10 Practical Patterns in Go

**Pattern 1 — Cache-aside** (the most common real use): check Redis first, fall back to the real database on a miss, then populate the cache.

```go
func getCustomerName(ctx context.Context, rdb *redis.Client, db *pgxpool.Pool, id string) (string, error) {
    cached, err := rdb.Get(ctx, "customer:"+id).Result()
    if err == nil {
        return cached, nil // cache hit
    }

    var name string
    db.QueryRow(ctx, "SELECT name FROM customers WHERE id = $1", id).Scan(&name)

    rdb.Set(ctx, "customer:"+id, name, time.Hour) // populate cache for next time
    return name, nil
}
```

**Pattern 2 — Rate limiting:**

```go
func isAllowed(ctx context.Context, rdb *redis.Client, key string, limit int64) bool {
    count, _ := rdb.Incr(ctx, key).Result()
    if count == 1 {
        rdb.Expire(ctx, key, time.Minute) // start a 1-minute window on first hit
    }
    return count <= limit
}
```

**Pattern 3 — Pub/Sub:**

```go
rdb.Publish(ctx, "notifications", "new order created")

sub := rdb.Subscribe(ctx, "notifications")
for msg := range sub.Channel() {
    fmt.Println("received:", msg.Payload)
}
```

## 6.11 Connecting from Go

```bash
go get github.com/redis/go-redis/v9
```

```go
import "github.com/redis/go-redis/v9"

rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
ctx := context.Background()
rdb.Set(ctx, "key", "value", 0)
val, _ := rdb.Get(ctx, "key").Result()
```

## 6.12 Practical Gotchas

- Redis is in-memory — without persistence (6.9), a restart loses everything. Fine for caches/sessions, not fine as your only copy of important data.
- Don't store huge objects in Redis — it's optimized for many small, fast lookups, not as a blob store.
- Secure it (5.6) — an unauthenticated Redis instance exposed to the network is a common, serious real-world breach.

---
# PART 7 — MongoDB: Building and Querying

MongoDB stores data as **documents** (JSON-like objects) in **collections** (like tables, but with no fixed schema). This guide reuses the same `customers`/`orders` example, translated into document form, to show exactly where the relational modeling from Part 1 changes shape.

```bash
docker exec -it mongo mongosh
use myapp
```

## 7.1 Core CRUD

```javascript
// Create
db.customers.insertOne({ name: "Alice", email: "alice@example.com", city: "Boston" })

// Read
db.customers.find({ email: "alice@example.com" })

// Update
db.customers.updateOne({ email: "alice@example.com" }, { $set: { name: "Alice Smith" } })

// Delete
db.customers.deleteOne({ email: "alice@example.com" })
```

There's no `CREATE TABLE` step — collections and fields are created implicitly on first insert. This is MongoDB's core trade-off: fast to start, flexible to change, but the database won't stop inconsistent shapes from creeping in (a `name` field that's a string on one document and a number on another) — that discipline has to come from your application code (or `$jsonSchema` validation rules, if you want the database to enforce it).

## 7.2 Modeling Data as Documents

Recall the relational version from Part 1 — `customers` and `orders`, joined by `customer_id`. The document equivalent can nest the related data directly, avoiding a join for the common read pattern:

```json
{
  "_id": 1,
  "name": "Alice",
  "city": "Boston",
  "orders": [
    { "amount": 50.00, "order_date": "2026-01-05" },
    { "amount": 20.00, "order_date": "2026-02-10" }
  ]
}
```

This is the core NoSQL trade-off: **embedding** (nesting) trades some redundancy and update complexity for read speed — you get a customer and all their orders in a single lookup, no join required.

## 7.3 Embedding vs. Referencing — the decision you'll make constantly

This is the document-database equivalent of a normalization decision from Part 1.

**Embed when:**
- The nested data is only ever read together with the parent (an order's line items).
- The nested data doesn't grow unbounded (a customer's 3 addresses, not their 10,000 orders).
- The nested data rarely changes independently of the parent.

**Reference when:**
- The related data is large, unbounded, or queried on its own (all orders across all customers).
- Multiple parents need to point to the same record (a product referenced by many orders — embedding it everywhere would duplicate that product's data into every order).

**Referencing example** — the same idea as a foreign key, resolved manually or via `$lookup` (7.6):

```json
// customers collection
{ "_id": 1, "name": "Alice", "city": "Boston" }

// orders collection
{ "_id": 101, "customer_id": 1, "amount": 50.00, "order_date": "2026-01-05" }
{ "_id": 102, "customer_id": 1, "amount": 20.00, "order_date": "2026-02-10" }
```

A common real-world pattern mixes both: embed a *summary* (name, last-order date) for fast reads, reference the full record for anything deeper.

**Common shapes this decision takes:**

- **One-to-few** (a customer's 2–3 addresses) → embed as an array.
- **One-to-many, bounded** (an order's line items) → embed.
- **One-to-many, unbounded** (a customer's orders over years) → reference, as a separate collection queried by `customer_id`.
- **Many-to-many** (products appearing in many orders) → reference both ways, or embed a small denormalized snapshot (product name/price *at time of purchase* — the same denormalization trade-off as 1.5) so historical orders don't change if the product is later edited.

## 7.4 CRUD Against Embedded and Referenced Data

```javascript
// CREATE
db.customers.insertOne({ _id: 4, name: "Dana", city: "Seattle" })

// READ
db.customers.find({ city: "Boston" })
db.customers.findOne({ _id: 1 })

// UPDATE
db.customers.updateOne({ _id: 1 }, { $set: { city: "Cambridge" } })

// UPDATE (push into an embedded array)
db.customers.updateOne(
  { _id: 1 },
  { $push: { orders: { amount: 30.00, order_date: "2026-04-01" } } }
)

// DELETE
db.customers.deleteOne({ _id: 4 })
```

## 7.5 Querying and Filtering

The document equivalent of `WHERE`, `ORDER BY`, and a `SELECT` column list — apply the "How to Write Any Query" method from 2.2 here too.

```javascript
// equivalent of WHERE amount > 50
db.orders.find({ amount: { $gt: 50 } })

// equivalent of WHERE city IN ('Boston', 'Denver')
db.customers.find({ city: { $in: ["Boston", "Denver"] } })

// equivalent of WHERE amount > 20 AND city = 'Boston'
db.orders.find({ amount: { $gt: 20 }, city: "Boston" })

// sorting and limiting — equivalent of ORDER BY / LIMIT
db.orders.find().sort({ amount: -1 }).limit(2)

// projection — equivalent of SELECT name, city
db.customers.find({}, { name: 1, city: 1, _id: 0 })

// querying inside embedded arrays
db.customers.find({ "orders.amount": { $gt: 40 } })
```

## 7.6 Aggregation Pipeline

The NoSQL equivalent of `GROUP BY`, `JOIN`, and multi-step transformations — a pipeline of stages, each passing its output to the next:

```javascript
db.orders.aggregate([
  { $match: { amount: { $gt: 0 } } },                              // like WHERE
  { $group: { _id: "$customer_id", total: { $sum: "$amount" } } }, // like GROUP BY + SUM
  { $sort: { total: -1 } },                                        // like ORDER BY
  { $limit: 5 }                                                    // like LIMIT
])
```

**`$lookup`** performs a join-like operation across collections when you chose to reference rather than embed:

```javascript
db.orders.aggregate([
  {
    $lookup: {
      from: "customers",
      localField: "customer_id",
      foreignField: "_id",
      as: "customer_info"
    }
  }
])
```

Note the performance trade-off: `$lookup` is slower than reading an already-embedded field — exactly why embedding is preferred for hot read paths (7.3).

## 7.7 Indexing

Same goal as Postgres (2.11/3.3) — speed up lookups — but you index document fields, including nested ones:

```javascript
db.orders.createIndex({ customer_id: 1 })
db.customers.createIndex({ email: 1 }, { unique: true })

// compound index, mirrors ORDER BY city, name
db.customers.createIndex({ city: 1, name: -1 })

db.customers.find({ email: "alice@example.com" }).explain("executionStats")
```

`.explain()` is Mongo's `EXPLAIN` (3.2) — it shows whether the query used the index (`IXSCAN`) or scanned everything (`COLLSCAN`). Index any field you regularly filter or sort by; more indexes speed up reads but slow down writes, same trade-off as Postgres.

## 7.8 Transactions in MongoDB

Modern MongoDB (4.0+) supports genuine **multi-document ACID transactions** (4.2) — necessary any time an operation spans more than one document or collection and must be all-or-nothing, the same motivating problem as 4.2's account transfer example:

```javascript
const session = client.startSession();
try {
  session.startTransaction();
  db.accounts.updateOne({ _id: 1 }, { $inc: { balance: -100 } }, { session });
  db.accounts.updateOne({ _id: 2 }, { $inc: { balance: 100 } }, { session });
  session.commitTransaction();
} catch (error) {
  session.abortTransaction();
} finally {
  session.endSession();
}
```

This is genuinely useful, but reach for it deliberately: the entire appeal of embedding (7.3) is that a *single-document* update is already atomic in MongoDB without any transaction wrapper at all — needing a multi-document transaction is often a sign the data would model more naturally as embedded, or that this specific operation is a rare enough exception to justify the extra machinery.

## 7.9 Connecting from Go

```bash
go get go.mongodb.org/mongo-driver/mongo
```

```go
package main

import (
    "context"
    "fmt"
    "log"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Customer struct {
    Name  string `bson:"name"`
    Email string `bson:"email"`
}

func main() {
    client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(context.Background())

    collection := client.Database("myapp").Collection("customers")

    collection.InsertOne(context.Background(), Customer{Name: "Bob", Email: "bob@example.com"})

    var result Customer
    collection.FindOne(context.Background(), bson.M{"email": "bob@example.com"}).Decode(&result)
    fmt.Println(result)
}
```

`bson:"name"` tags work just like `json:"name"` tags in your Gin structs — they control how your Go struct maps to Mongo's document fields. `bson.M{...}` is Mongo's equivalent of Gin's `gin.H{...}` — a shorthand map for building queries/filters.

## 7.10 Practical Gotchas

- MongoDB **can** do joins (`$lookup`), but it's less natural and generally slower than SQL joins — if your data is heavily relational, that's a sign you probably want Postgres instead (Part 10).
- Schema flexibility is a double-edged sword: nothing stops inconsistent documents from creeping in over time without discipline (validation rules, or enforcing it in your Go structs).
- Don't reach for MongoDB just because it's popular — reach for it when your data genuinely has a document shape (deeply nested, varying structure) that would be awkward across many relational tables.

---
# PART 8 — The Rest of the NoSQL Landscape

MongoDB (document) and Redis (key-value) are the two NoSQL types this guide builds in depth, since they cover the vast majority of real projects. Two more types exist and are worth recognizing even if you don't reach for them often:

| Type | Stores data as | Example products | Best for |
|---|---|---|---|
| Document | JSON-like documents | MongoDB, Couchbase | Flexible, nested records (profiles, orders) |
| Key-Value | A key mapped to an opaque value | Redis, DynamoDB | Caching, sessions, fast lookups by ID |
| Wide-Column | Rows with dynamic, sparse columns | Cassandra, HBase | Massive write volume, time-series data |
| Graph | Nodes and relationships | Neo4j, Amazon Neptune | Highly connected data (social networks, recommendations) |

## 8.1 Wide-Column Stores (Cassandra)

Built for very high write throughput across many machines, with tables designed around how you'll *read* the data, not ad-hoc filtering:

```sql
CREATE TABLE orders_by_customer (
    customer_id INT,
    order_date  TIMESTAMP,
    amount      DECIMAL,
    PRIMARY KEY (customer_id, order_date)
);

SELECT * FROM orders_by_customer WHERE customer_id = 1;
```

The key design difference from SQL: you model the table around the *query* you'll run (here, "get a customer's orders sorted by date"), not around normalized entities. A different query pattern (e.g. "get orders by date across all customers") typically means creating a second table with a different key structure, duplicating the data on purpose.

## 8.2 Graph Databases (Neo4j)

For when relationships *between* records matter as much as the records themselves — traversing many hops is slow with SQL joins but native here:

```cypher
CREATE (a:Customer {name: "Alice"})
CREATE (b:Customer {name: "Bob"})
CREATE (a)-[:REFERRED]->(b)

// find everyone Alice referred, directly or indirectly
MATCH (a:Customer {name: "Alice"})-[:REFERRED*1..3]->(referred)
RETURN referred.name;
```

Common uses: social networks, fraud detection, recommendation engines — anywhere "friends of friends" or multi-hop paths are the actual question.

---

# PART 9 — Distributed Systems: CAP, Consistency, and Scaling

## 9.1 The CAP Theorem

This is the distributed-systems analog of Part 4's Isolation/Durability trade-offs — what you give up in exchange for scaling across many machines. A distributed database can only fully guarantee two of three properties at once:

- **Consistency** — every read gets the most recent write.
- **Availability** — every request gets a response, even during a network failure.
- **Partition tolerance** — the system keeps working despite network splits between nodes.

Since network partitions happen regardless of what you'd prefer, the real choice in practice is **C vs. A** during a partition:

| Choice | Example databases | Behavior during a partition |
|---|---|---|
| CP (Consistency + Partition tolerance) | MongoDB (default), HBase | Some nodes refuse requests rather than risk stale data |
| AP (Availability + Partition tolerance) | Cassandra, DynamoDB | All nodes keep responding, but may briefly return stale data |

Most document and key-value stores also offer **eventual consistency** as a tunable setting — writes propagate to all replicas *eventually*, not instantly:

```javascript
// MongoDB: read from a secondary replica, may be slightly stale
db.orders.find().readPref("secondaryPreferred")
```

## 9.2 Scaling: Replication and Sharding

The mechanisms that let databases handle far more data/traffic than a single server — much of why you'd choose a distributed system at all.

**Replication** — copy data across multiple nodes for redundancy and read scaling (the same idea as Postgres read replicas in 3.8):

```javascript
rs.status();   // MongoDB replica set status — writes go to the primary, reads can spread across secondaries
```

**Sharding** — split data across nodes so no single machine holds it all, based on a shard key:

```javascript
sh.shardCollection("mydb.orders", { customer_id: 1 });
```

Choosing a good shard key matters enormously — a poor choice (e.g. a field with few distinct values) creates "hot" shards that get overloaded while others sit idle, the distributed-systems equivalent of the lock-contention hot spots from Part 4.

## 9.3 Quick Reference: SQL ↔ NoSQL Terminology

| Concept | SQL equivalent | NoSQL approach |
|---|---|---|
| Row | Table row | Document / item / node |
| Relationship | Foreign key + `JOIN` | Embedding, reference + `$lookup`, or graph edge |
| `WHERE` | Filter clause | Query filter object |
| `GROUP BY` | Aggregate clause | Aggregation pipeline `$group` |
| `INDEX` | Speeds up lookups | Same concept, indexes on any field |
| Schema enforcement | Enforced by table definition | Optional/flexible (schema-on-read) |
| Transaction guarantees | ACID by default | Varies — many now support multi-document ACID (7.8), but tune consistency explicitly |

---
# PART 10 — When to Use What: The Practical Decision Framework

You've now seen how Postgres, Redis, and MongoDB actually work internally — their keys, their locking, their query engines, their durability guarantees. This section turns that into concrete, mechanism-level implementation decisions, not hand-wavy "it depends."

## 10.1 Decide by the mechanism you actually need

| If your requirement is... | The mechanism that satisfies it | Which database, and why |
|---|---|---|
| Multiple related writes must never partially apply (money, inventory, contracts) | Atomicity + Consistency + a real locking/MVCC engine (Part 4) | **Postgres** — `CHECK`/`FOREIGN KEY` constraints plus `BEGIN`/`COMMIT` give you enforcement the database itself guarantees, not just application-code discipline |
| Data has real, stable relationships you'll query from multiple angles | Foreign keys + `JOIN` (2.5), normalized to 3NF (1.4) | **Postgres** — a relational engine is *built* for "give me customers whose orders in the last month exceed $X," which needs a join and a filter, not a document shape |
| Sub-millisecond reads of small, simple values by exact key | O(1) key lookup, everything in RAM (6.1) | **Redis** — as a cache *in front of* Postgres (cache-aside, 6.10), never as the only copy of data that matters |
| Data expires naturally (sessions, rate-limit counters, short-lived tokens) | Built-in `EXPIRE`/TTL (6.7) | **Redis** — the TTL mechanism itself does the cleanup work you'd otherwise have to build with a cron job in Postgres |
| A ranked, always-sorted view (leaderboards, "top N") | Sorted set, kept ordered on every write (6.6) | **Redis** — what would be `GROUP BY` + `ORDER BY` + `LIMIT` recomputed on every read in SQL is a pre-sorted structure in Redis |
| Records whose shape genuinely varies per instance (arbitrary webhook payloads, activity logs, evolving product attributes across wildly different categories) | Schema-on-read, no `ALTER TABLE` needed per shape (7.1) | **MongoDB** — forcing this into rigid columns means either a sprawling nullable-column table or an EAV anti-pattern; a document collection models it directly |
| A record and all its child data are almost always read together, and the child data is bounded | Embedding (7.3) | **MongoDB** — one document read replaces a join, at the cost of some redundancy if the child data ever needs to be queried independently |
| Your schema will change shape weekly during early prototyping, before requirements settle | Schema-on-read, no migration per change (7.1) | **MongoDB** — cheaper to iterate on early, *if* you accept the app-code discipline cost from 7.10's gotcha once the shape does stabilize |
| Data that's naturally a graph — many-hop relationship traversal is the actual query | Native graph traversal (8.2) | **Neo4j** or similar — SQL joins get slow and awkward past 2–3 hops of "friends of friends" |
| Enormous write volume across many machines, queried by a fixed access pattern | Wide-column, table-per-query-shape (8.1) | **Cassandra**/DynamoDB — accept modeling around the query up front in exchange for write throughput no single-node database can match |

## 10.2 A worked decision, start to finish

Take a concrete feature: **"show each user their 5 most recent orders instantly, and give them a live notification when an order ships."**

1. **Recent orders, correctly, with real relationships to customers/products** → this needs joins and durable correctness → **Postgres**, using the schema from Part 1, indexed on `(customer_id, order_date)` per Part 3 so the query is an Index Scan, not a Seq Scan.
2. **"Instantly"** → even a well-indexed Postgres query is a network round trip plus a query plan; if this is hit on every page load, cache the last-5-orders result **per customer** in Redis with a short TTL (cache-aside, 6.10), invalidated (or just left to expire) whenever a new order is placed.
3. **Live notification on shipment** → this is exactly Redis Pub/Sub's use case (6.8) — not durable, not something you'd want as the source of truth for "did the order ship" (that fact lives in Postgres, written durably in the same transaction that updates order status), but perfect for "push a real-time UI update to whoever's currently looking."

The result uses all three engines' actual mechanisms deliberately: Postgres as the durable, correct source of truth; Redis as a cache for the hot read path *and* as the transport for a real-time signal that doesn't need to be durable. This is the normal shape of a real production system — not a sign of over-engineering.

## 10.3 Quick Decision Cheat Sheet

| You need... | Use |
|---|---|
| Structured data with real relationships (users, orders, payments) | **PostgreSQL** |
| Strong consistency, financial-style data, complex queries/joins | **PostgreSQL** |
| Fast, in-memory caching to speed up repeated reads | **Redis** (in front of Postgres, not instead of it) |
| Session storage, rate limiting, short-lived data | **Redis** |
| Real-time counters, leaderboards, pub/sub messaging | **Redis** |
| Deeply nested / wildly varying document shapes | **MongoDB** |
| Rapid prototyping where the schema will change a lot early on | **MongoDB** |
| A typical web app backend (the default choice) | **PostgreSQL** + **Redis** for caching — covers the vast majority of real projects |

---

# PART 11 — Putting It All Together: A Polyglot Architecture

**Most real production systems use more than one of these together.** A typical shape:

```mermaid
flowchart TD
    App["Go API server (Gin)"]
    App -->|source of truth, ACID writes| PG["PostgreSQL<br/>customers, orders, products<br/>(Parts 1, 2, 4, 5)"]
    App -->|cache-aside reads, sessions, rate limits| Redis["Redis<br/>(Part 6)"]
    App -->|flexible/variable-shape data| Mongo["MongoDB<br/>webhook payloads, activity logs<br/>(Part 7)"]
    PG -.->|WAL replication, PITR| PGReplica["Read replica / backups<br/>(Part 4.15)"]
```

- **Postgres** holds the core transactional data — orders, payments, anything where a `CHECK` constraint or a foreign key actually needs to be enforced, and where the recovery guarantees from Part 4 matter.
- **Redis** caches hot data, holds sessions, and handles rate-limiting/pub-sub — all disposable if it's lost, because Postgres is the real source of truth behind it.
- **MongoDB** handles the specific subsystem that genuinely needs flexible documents — arbitrary webhook payloads, activity feeds, or a product catalog with wildly varying attributes across categories.

Don't feel pressure to pick just one database for an entire project — pick the right mechanism (Part 10.1) for each piece of the system.

---

# PART 12 — Master Quick-Reference Cheat Sheets

### SQL syntax

| Clause / Feature | Purpose |
|---|---|
| `SELECT ... WHERE` | Retrieve and filter rows |
| `ORDER BY` / `LIMIT` | Sort and cap results |
| `JOIN` | Combine rows from related tables |
| `GROUP BY` / `HAVING` | Aggregate rows into summary groups |
| Subquery / CTE | Nest one query's result inside another |
| `UNION` / `INTERSECT` / `EXCEPT` | Combine or compare separate result sets |
| `INSERT`/`UPDATE`/`DELETE` | Modify data |
| `CREATE TABLE` + constraints | Define structure and enforce validity |
| `INDEX` | Speed up lookups on specific columns |
| `VIEW` | Save and reuse a query as a virtual table |
| `BEGIN`/`COMMIT`/`ROLLBACK` | Group statements into one atomic transaction |
| Window functions | Per-row calculations across related rows |

### Transactions, concurrency, recovery (Part 4)

| Concept | One-line definition |
|---|---|
| ACID | Atomicity, Consistency, Isolation, Durability — the four transaction guarantees |
| Serializability | A concurrent schedule that behaves as if it ran serially |
| Precedence graph | Cycle = not serializable; test via conflicting operations between transactions |
| Lost update / dirty read / unrepeatable read / phantom read | The four concurrency anomalies isolation levels prevent |
| S / X locks | Shared (read) / Exclusive (write) locks on a data item |
| 2PL | Growing phase (acquire only) then shrinking phase (release only) — guarantees serializability |
| Deadlock | Circular wait; detected via wait-for graph, avoided via consistent lock ordering |
| Intent locks (IS/IX/SIX) | Coarse-grained "heads up" locks that make multi-granularity locking checkable in O(1) |
| WAL / log-based recovery | Log the change before applying it; undo uncommitted, redo committed-but-unflushed on crash |
| 2PC | Prepare (vote) then commit — atomic commit across multiple databases, at the cost of a blocking failure mode |

### Security (Part 5)

| Concept | One-line definition |
|---|---|
| DAC / RBAC / ABAC | Owner-granted / role-based / attribute-and-policy-based access control |
| Row-Level Security | Per-row filtering enforced by the database itself, based on session context |
| Encryption in transit / at rest / column-level | TLS between app and DB / disk-level / `pgcrypto` on specific fields |
| Hashing ≠ encryption | Passwords are hashed (one-way, salted, e.g. bcrypt) — never encrypted (reversible) |

### Redis command families

| Type | Core commands |
|---|---|
| String | `SET` `GET` `INCR` `EXPIRE` |
| List | `LPUSH` `RPUSH` `LPOP` `LRANGE` |
| Hash | `HSET` `HGET` `HGETALL` |
| Set | `SADD` `SISMEMBER` `SINTER` |
| Sorted Set | `ZADD` `ZRANGE` `ZSCORE` |
| Pub/Sub | `SUBSCRIBE` `PUBLISH` |
| Streams | `XADD` `XRANGE` |

### MongoDB operator families

| SQL | MongoDB |
|---|---|
| `WHERE x > 5` | `{ x: { $gt: 5 } }` |
| `WHERE x IN (...)` | `{ x: { $in: [...] } }` |
| `ORDER BY x DESC LIMIT n` | `.sort({x:-1}).limit(n)` |
| `GROUP BY` + `SUM` | `$group` + `$sum` in `.aggregate()` |
| `JOIN` | `$lookup` in `.aggregate()`, or embed instead |

---

This guide, plus the four short videos linked at the top, covers everything needed to go from "what is a database" to designing, building, querying, optimizing, securing, and recovering a real Postgres + Redis + MongoDB system in Go — as an industry-grade practitioner, not just someone who can write a `SELECT`.
