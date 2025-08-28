
---

## **Depth-First Search (DFS)**

**Definition:**  
Depth-First Search (DFS) is a graph traversal algorithm used to explore nodes and edges in a graph or tree structure. It works by starting at a given node and exploring as far as possible along each branch before backtracking.

---

### **How It Works**

1. **Start** at the root (or a given starting node).
    
2. **Visit** the node and mark it as visited.
    
3. **Explore** the first unvisited neighbor and repeat the process.
    
4. If a node has no unvisited neighbors, **backtrack** to the most recent node that still has unexplored neighbors.
    
5. Continue until **all nodes** reachable from the start have been visited.
    

This process can be implemented **recursively** or **iteratively** (using a stack).

---

### **Key Characteristics**

- **Traversal Order:** Goes deep into one path before moving to another.
    
- **Data Structure Used:** Commonly uses a **stack** (explicit stack in iterative form, or the call stack in recursion).
    
- **Complexity:**
    
    - **Time:** `O(V + E)` where `V` = number of vertices, `E` = number of edges.
        
    - **Space:** `O(V)` in the worst case (due to recursion depth or stack size).
        
- **Type:** Non-weighted, does not guarantee the shortest path.
    

---

### **Common Use Cases**

- **Pathfinding** in mazes or puzzle solving (when finding _a_ path is enough, not necessarily the shortest).
    
- **Topological sorting** of a Directed Acyclic Graph (DAG).
    
- **Cycle detection** in graphs.
    
- **Connected components** detection in an undirected graph.
    
- **AI search problems** and decision trees.
    

---

### **Pseudocode (Recursive)**

```plaintext
procedure DFS(node):
    mark node as visited
    for each neighbor in node's adjacency list:
        if neighbor is not visited:
            DFS(neighbor)
```

---

### **Example (Iterative in Pseudocode)**

```plaintext
procedure DFS_Iterative(start):
    create empty stack
    push start to stack

    while stack is not empty:
        node = pop from stack
        if node is not visited:
            mark node as visited
            push all unvisited neighbors to stack
```

---

### **Illustrative Example**

If we have the following graph:

```
A - B - D
|   |
C   E
```

**DFS starting at A** might visit nodes in the order:

```
A → B → D → E → C
```

_(Exact order can vary depending on adjacency ordering.)_

---

Do you want me to also make a **small visual diagram** showing how DFS explores step-by-step so it’s even clearer in your notes? That would make your project’s explanation pop visually.