# L8QL - Layer 8 API Query Language

[![Go Version](https://img.shields.io/badge/go-1.23.8-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-GPL%20v3-blue.svg)](LICENSE)

L8QL (Layer 8 API Query Language) is an alteration and facade of the SQL language designed for querying Graph Models at runtime. It provides a single, simple, and common API to query graph model data of any Go struct, eliminating the need for complex API integrations.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage Examples](#usage-examples)
- [Query Syntax](#query-syntax)
- [API Reference](#api-reference)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

## Overview

We keep inventing the wheel, over and over again, trying to create APIs for our services/products and spending enormous time and money trying to integrate different products. Most software engineers consider infrastructure components like Kafka, NATS, DB, ETCD, etc., as the "Wheel" and rush to implement and use those infrastructure components. However, they are doing the complete opposite.

While creating those from scratch is a nice challenge, it isn't as expensive as maintaining APIs and integrations over time. Using infrastructure components is a very easy task that can take a month or even weeks, while building a stable API and integrating with different products might take years, huge amounts of money, and constant costly maintenance over time.

If we do an analogy to Language, the infrastructure components are the alphabet, while the API is the actual Languages. Just as two persons, each knowing a different language but with the same alphabet, cannot speak to each other, two products built with the same infrastructure cannot communicate with each other and require very expensive, highly maintenance integrations.

**L8QL comes to ease the language/API challenge by presenting a single, simple & common API to query the graph model & data of a product at runtime.**

## Features

- **SQL-like Query Language**: Familiar syntax for developers
- **Graph Model Support**: Native support for complex nested structures
- **Runtime Introspection**: Automatic model discovery and schema generation
- **Filter and Match**: Powerful filtering capabilities with pattern matching
- **Sorting and Pagination**: Built-in support for `sort-by`, `limit`, `page`, `ascending`/`descending`
- **Case Sensitivity Control**: Optional `match-case` functionality
- **Deep Path Navigation**: Query nested objects and collections seamlessly
- **Hash Support**: Built-in MD5 hash generation for data integrity and caching
- **Advanced Sorting**: Sort by value with support for complex data types
- **Type Safety**: Strong typing with Go's reflection system
- **Zero Dependencies**: Lightweight design with minimal external dependencies

## Architecture

L8QL consists of several key components:

### 1. Introspector
The Model Introspector accepts a Go struct and introspects it by drilling down to discover its attributes and sub-objects. From this data, it creates:
- **Internal Schema**: Mapping struct→table and attribute→column
- **Graph Schema**: Mapping relations between the root struct and its sub-structs

### 2. Parser
The parser parses the query string and validates syntax correctness. It divides the query into:
- **Requested Columns**: Selected fields
- **Tables**: Target struct types
- **Criteria**: Divided into Expressions, Comparators & Conditions

### 3. Interpreter
The Interpreter takes a syntax-valid parsed query and validates it via the Introspector schema. It matches string representations of attributes to discovered Columns & tables. The outcome is an Interpreter Query instance for filtering elements.

### 4. Instance & Attribute System
- **Instance**: String representation of an instance inside the model (e.g., `"Employee[key]"`)
- **Attribute**: String representation of a struct attribute (e.g., `"Employee[key].Addresses.Address[0].Line2"`)

## Installation

```bash
go get github.com/saichler/l8ql/go
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/saichler/l8ql/go/gsql/interpreter"
    "github.com/saichler/l8types/go/ifs"
)

// Define your model
type Employee struct {
    Name      string
    Age       int32
    Addresses []Address
}

type Address struct {
    Line1   string
    Line2   string
    Zip     string
    Country string
}

func main() {
    // Create a query
    query, err := interpreter.NewQuery(
        "select name, age from employee where name='John' or age>25", 
        resources, // your IResources implementation
    )
    if err != nil {
        panic(err)
    }
    
    // Test with your data
    employees := []*Employee{
        {Name: "John", Age: 30},
        {Name: "Jane", Age: 25},
        {Name: "Bob", Age: 35},
    }
    
    // Filter matching elements
    for _, emp := range employees {
        if query.Match(emp) {
            fmt.Printf("Matched: %+v\n", emp)
        }
    }
}
```

## Usage Examples

### Basic Queries

```sql
-- Select all fields
select * from employee

-- Select specific fields
select name, age from employee where age > 25

-- Complex conditions with parentheses
select name from employee where (age > 25 and name = 'John') or country = 'USA'
```

### Deep Path Navigation

```sql
-- Query nested objects
select name from employee where addresses.country = 'USA'

-- Query array elements
select name from employee where addresses[0].zip = '12345'

-- Query map values
select name from employee where metadata.department = 'Engineering'
```

### Advanced Features

```sql
-- Sorting and pagination
select * from employee where age > 25 sort-by age descending limit 10 page 1

-- Case-sensitive matching
select * from employee where name = 'john' match-case

-- Pattern matching
select * from employee where name = 'J*'
```

## Query Syntax

L8QL supports the following SQL-like syntax:

### Basic Structure
```sql
select <columns> from <table> [where <conditions>] [sort-by <column>] [ascending|descending] [limit <number>] [page <number>] [match-case]
```

### Supported Comparators
- `=` - Equal
- `!=` - Not Equal  
- `<` - Less Than
- `<=` - Less Than or Equal
- `>` - Greater Than
- `>=` - Greater Than or Equal
- `in` - In (for arrays/collections)
- `not-in` - Not In

### Logical Operators
- `and` - Logical AND
- `or` - Logical OR
- `()` - Parentheses for grouping

### Special Features
- `*` - Wildcard for selecting all columns
- `sort-by <column>` - Sort results by specified column
- `ascending`/`descending` - Sort order
- `limit <n>` - Limit results to n items (max 1000)
- `page <n>` - Page number for pagination
- `match-case` - Enable case-sensitive string matching

## API Reference

### Core Interfaces

#### Query Interface
```go
type Query interface {
    Match(any interface{}) bool
    Filter(list []interface{}, onlySelectedColumns bool) []interface{}
    String() string
    // ... other methods
}
```

#### Creating Queries
```go
// From SQL string
query, err := interpreter.NewQuery(sqlString, resources)

// From parsed query object
query, err := interpreter.NewFromQuery(parsedQuery, resources)
```

### Key Methods

- `Match(any interface{}) bool` - Test if an object matches the query criteria
- `Filter([]interface{}, bool) []interface{}` - Filter a slice of objects
- `Properties() []ifs.IProperty` - Get selected properties
- `Criteria() ifs.IExpression` - Get the where clause expression

## Testing

The project includes comprehensive test suites:

### Running Tests
```bash
# Run all tests
go test ./...

# Run with coverage
go test -v -coverpkg=./gsql/... -coverprofile=cover.html ./...

# Run the test script (includes security checks)
./test.sh
```

### Test Categories
- **Parser Tests**: Validate SQL parsing and syntax validation
- **Interpreter Tests**: Test query execution and matching
- **Integration Tests**: End-to-end query scenarios

## Project Structure

```
l8ql/
├── go/
│   ├── gsql/
│   │   ├── interpreter/          # Query interpretation and execution
│   │   │   ├── Query.go
│   │   │   ├── Expression.go
│   │   │   ├── Condition.go
│   │   │   ├── Comparator.go
│   │   │   └── comparators/      # Comparison operators
│   │   └── parser/               # SQL parsing
│   │       ├── Query.go
│   │       ├── Expression.go
│   │       ├── Condition.go
│   │       └── Comparator.go
│   ├── tests/                    # Test suites
│   ├── test.sh                   # Test runner script
│   ├── go.mod                    # Go module definition
│   └── go.sum                    # Dependency checksums
├── LICENSE                       # GPL v3 License
└── README.md                     # This file
```

## Dependencies

- **github.com/saichler/l8types** - Core type definitions and interfaces
- **github.com/saichler/reflect** - Enhanced reflection utilities
- **github.com/saichler/l8test** - Testing infrastructure
- **Standard Go libraries** - No external runtime dependencies

## Performance Considerations

- L8QL uses Go's reflection system for type introspection
- Query parsing is done once and can be reused
- Filtering is performed in-memory
- Suitable for moderate-sized datasets (thousands to tens of thousands of objects)
- For large datasets, consider implementing custom optimizations

## Limitations

- **In-Memory Processing**: All filtering happens in memory
- **Limit Cap**: Maximum limit is 1000 items per query
- **Go Structs Only**: Currently supports Go structs only
- **No Joins**: No support for SQL-style joins between different types

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Run the test suite (`./test.sh`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Development Guidelines

- Follow Go coding conventions
- Add comprehensive tests for new features
- Update documentation for API changes
- Ensure all tests pass before submitting
- Use meaningful commit messages

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

## Support

For questions, issues, or contributions:

- **Issues**: [GitHub Issues](https://github.com/saichler/l8ql/issues)
- **Discussions**: Use GitHub Discussions for questions and ideas
- **Documentation**: Check the code documentation and tests for detailed examples

## Roadmap

- [ ] Performance optimizations for large datasets
- [ ] Support for additional data types
- [ ] Query optimization and caching
- [ ] Integration with popular Go ORMs
- [ ] Support for other programming languages
- [ ] Advanced aggregation functions
- [ ] Query plan visualization tools

---

**L8QL** - Simplifying graph data querying with familiar SQL syntax.