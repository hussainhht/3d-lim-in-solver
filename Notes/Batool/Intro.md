
**Project Overview for the Team**

This project is basically about building an **Ant Farm** made of rooms connected by paths, and then **simulating ants moving through it**.

Before worrying about how ants move, we first need to model the **ant farm** itself.

---

### Graph Analogy

If you’ve used **Obsidian**, you’ve seen those nodes and links between them. That’s a **graph**.

- **Nodes**: entities (rooms, people, documents, etc.)
    
- **Edges**: relationships or connections between entities (paths between rooms, links between notes).
    

---

### Types of Graphs

- **Directed graph**: edges have a direction. Example:  
    `NoteA → NoteB` means you can go from A to B, but not the other way unless there’s also `NoteB → NoteA`.
    
- **Undirected graph**: edges go both ways by default. Example:  
    `A ↔ B` means A is connected to B and B is connected to A.
    

Obsidian’s graph view can be **directed** (arrows on) or **undirected** (arrows off).  
In our ant farm, **ants can go both ways**, so we’ll use an **undirected graph**.

---

### Graphs in Go

Go doesn’t have a built-in graph data structure (unlike Python’s NetworkX or C++’s Boost).  
Instead, we’ll represent our graph with a **map**.

Two common approaches:

1. **Direct map**: `map[string][]string` — where each key is a room, and the value is the list of connected rooms.
    
2. **Graph struct**:
    
    ```go
    type Graph struct {
        adj map[string][]string
    }
    ```
    
    This keeps all graph logic together and makes it easier to add helper methods.
    

---

### Building an Undirected Graph in Go

In an undirected graph, if A is connected to B:

```go
A → B  
B → A
```

We have to store both directions explicitly in our map.

---

**So far:**  
We’re just **building the ant farm’s structure** (rooms + paths).  
We haven’t started the simulation yet — that’s the next step once the farm is ready.

---

