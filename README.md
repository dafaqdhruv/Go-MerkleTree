# Go-Test

A simple implementation of Merkle Tree in Go.  
Takes an array of strings as input and returns the Merkle root of the tree.


### How to run
 1. Clone the repo
 2. go run test.go

### Test run result
Providing the input of the following array with sha256 encryption :  
```["hello", "this", "is", "a", "merkle", "tree"]```  
The resulting merkle looks a bit like this   


```mermaid 
graph 
A(c393...eaaa) --> B(c06c...f5cd)
A --> C(487e...d094)
B --> D(8815...0d52)
B --> E(1bc4...70e2)
C --> F(487e...d094)
D --> G(2cf2...9824)  --> M(hello)
D --> H(1eb7...8408)  --> N(this)
E --> I(fa51...57f6)  --> O(is)
E --> J(ca97...48bb)  --> P(a)
F --> K(7975...5590)  --> Q(merkle)
F --> L(dc9c...0622)  --> R(tree)
```