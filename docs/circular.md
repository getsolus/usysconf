# Circular Dependencies

## Example 1

A -> A

```
g.HasCircular
    g.circular("A", [])
        visited = ["A"]
        for: dep = "A"
            for: v = "A"
                if: dep == v: true
                    found: ["A", "A"]
    found = ["A","A"]
    print: "Circular dependency: A -> A"
```

## Example 2

A -> B
B -> A

```
g.HasCircular
    visited = []
    g.circular("A", visited)
        visited = ["A"]
        for: dep = "B"
            for: v = "A"
                if: dep == v: false
            if: len(found) == 0
                g.circular("B", visited)
                    visited = ["A", "B"]
                    for: dep = "A"
                        for: v = "A"
                            if: dep == v: true
                                found = ["B", "A"]
                found = ["B", "A"]
                if: len(found) !=0
                    found = ["A", "B", "A"]
    found = ["A", "B", "A"]
    print: "Circular dependency: A -> B -> A"
```

## Example 3

A -> B
B -> C
C -> D
D -> B

```
g.HasCircular
    visited = []
    g.circular("A", visited)
        visited = ["A"]
        for: dep = "B"
            for: v = "A"
                if: dep == v: false
            if: len(found) == 0: true
                g.circular("B", visited)
                    visited = ["A", "B"]
                    for: dep = "C"
                        for: v = "A"
                            if: dep == v: false
                        for: v = "B"
                            if: dep == v: false
                    if: len(found) == 0: true
                        g.circular("C". visited)
                            visited = ["A", "B", "C"]
                            for: dep = "D"
                                for: v = "A"
                                    if: dep == v: false
                                for: v = "B"
                                    if: dep == v: false
                                for: v = "C"
                                    if: dep == v: false
                            if: len(found) == 0: true
                                g.circular("D". visited)
                                    visited = ["A", "B", "C", "D"]
                                    for: dep = "B"
                                        for: v = "A"
                                            if: dep == v: false
                                        for: v = "B"
                                            if: dep == v: true
                                                found = ["D", "B"]
                                found = ["D", "B"]
                                if len(found) != 0: true
                                    found = ["C", "D", "B"]
                        found = ["C", "D", "B"]
                        if len(found) != 0: true
                            found = ["B", "C", "D", "B"]
                found = ["B", "C", "D", "B"]
                    if len(found) != 0: true
                        found = ["A", "B", "C", "D", "B"]            
    found = ["A", "B", "C", "D", "B"]
    print: "Circular dependency: A -> B -> C -> D -> B"
```
