---
name: check-isa-steps
description: Analyzes BDD feature files to find missing step definitions and creates an implementation plan. Use when working with .isa.feature files, implementing BDD scenarios, or when the user asks to "check missing steps" or "implement feature steps".
allowed-tools: Read, Glob, Grep, TodoWrite, Edit, Write
---

# Execution Steps

## 1. 讀取 feature 檔案

讀取用戶指定的 .isa.feature 檔案，提取所有 Given/When/Then/And 步驟。

## 2. 建立索引

讀取 `tests/bdd/steps/register.go`，提取所有 `Register*Steps(sc)` 呼叫，建立檔案路徑映射。

範例: `database.RegisterUserSteps(sc)` → `tests/bdd/steps/database/user.go`

## 3. 匹配步驟

對每個步驟：

1. 使用 [reference.md](reference.md) 的匹配規則推斷候選檔案（1-2 個）
2. 用 Grep 搜尋 `sc.Step(` 提取該檔案中的正則表達式
3. 匹配步驟 → 標記已實作或缺少

## 4. 判斷實作位置

對缺少的步驟，參考 [reference.md](reference.md) 判斷：

- 現有 entity → 加入現有檔案
- 新 entity → 建議新檔案
- 確定函式名稱（camelCase）和正則表達式

## 5. 建立 TodoWrite 計畫

為每個缺少的步驟建立 todo：

- **content**: `"實作 '{步驟}' 於 {檔案路徑} (函式: {函式名稱})"`
- **activeForm**: `"實作 '{步驟}'"`
- **status**: `"pending"`

輸出簡短摘要：找到 X 個缺少步驟，等待 lgtm。

## 6. 實作步驟（用戶確認後）

逐一實作每個 todo，包括：

1. 寫入或編輯步驟檔案
2. 在對應的 `Register*Steps` 函式中註冊步驟
3. 必要時更新 `register.go` 新增 `Register*Steps` 呼叫
