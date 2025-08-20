Sure! Here’s a clean **Obsidian-style note** for this concept:

---

# Handling Maps in Go – Pointers vs. Values

## Issue

When you have a map pointer:

```go
var occupied *map[string]bool
occupied[r.name] = true // ❌ invalid operation
```

Go gives an error:

```
invalid operation: cannot index occupied (variable of type *map[string]bool)
```

---

## Why

* `occupied` is a pointer to a map (`*map[string]bool`), not the map itself.
* Indexing a pointer directly is **not allowed** in Go.

---

## Solutions

### 1️⃣ Dereference the pointer

```go
(*occupied)[r.name] = true
```

* Explicitly accesses the map the pointer refers to.

### 2️⃣ Use a normal map (preferred)

```go
occupied := make(map[string]bool)
```

* Pass `occupied` as a value, not a pointer:

```go
func BFS(currentRoom *Room, endRoom *Room, occupied map[string]bool) []*Room
```

* Then indexing works naturally:

```go
occupied[r.name] = true
```

---

## ✅ Summary

* Use a pointer only if you really need to share the **map reference** across functions in a way that might reassign the map itself.
* Most of the time, passing a map **by value** works perfectly because maps are reference types in Go.

---

You can save this in Obsidian as `Go_Map_Pointers.md` or just paste it directly into your notes.

If you want, I can also make a **mini visual diagram** showing map vs. pointer behavior—it’s really handy for quick recall. Do you want me to do that?
