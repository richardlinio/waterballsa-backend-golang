# Slice 初始化風格

在初始化 Slice 時，推薦使用 `append` 模式，並預先分配容量以提升效能。

### ✅ 推薦模式

```go
// 使用 append 模式 (推薦)
items := make([]T, 0, len(source))
for _, item := range source {
    items = append(items, T{...})
}
```

### 原則

- 使用 `make([]T, 0, capacity)` 預分配容量以避免多次重新分配記憶體。
- 透過 `append` 逐一加入元素，程式碼更為簡潔且不易出錯。
- 避免使用索引賦值 `items[i] = ...`，因為這需要手動管理索引，且在初始化空 Slice 時會導致 panic。
