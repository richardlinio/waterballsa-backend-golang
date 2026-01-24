---
name: analyze-isa-feature
description: Analyzes ISA (Integration System Architecture) feature files to create comprehensive implementation plans for backend API endpoints. Use when given an .isa.feature file to understand what needs to be implemented and create a detailed, step-by-step implementation guide.
allowed-tools: Read, Glob, Grep, Task, Write, WebFetch, WebSearch
---

# Purpose

This skill systematically analyzes ISA BDD feature files and produces detailed implementation plans for backend API features. It bridges the gap between high-level test specifications and concrete code implementation by:

1. Understanding what the feature tests expect
2. Analyzing the current codebase state
3. Identifying what's missing
4. Creating a complete, actionable implementation plan

# When to Use

- User provides an `.isa.feature` file and asks "how to implement this"
- User asks to "analyze this ISA feature" or "create implementation plan"
- User wants to understand what code is needed for BDD tests to pass
- User references swagger.yaml and db-schema.dbml for API implementation

# Execution Steps

## Phase 1: Feature Analysis

### Step 1.1: Read the ISA Feature File

Read the specified `.isa.feature` file to understand:

- **Scenarios**: What behaviors are being tested?
- **API endpoints**: Which HTTP methods and paths? (GET, POST, PUT, DELETE)
- **Request/Response formats**: What JSON structures are expected?
- **Business logic**: What are the key rules? (auto-completion, validation, status transitions)
- **Data setup requirements**: What database records are needed?

**Key questions to answer:**

- What endpoints need to be implemented?
- What are the request/response DTOs?
- What business rules must the service layer enforce?
- What database queries are needed?

### Step 1.2: Read API Specification (swagger.yaml)

Location: `/docs/api-docs/swagger.yaml` and `/docs/api-docs/openapi/paths/*.yaml`

Extract:

- **Endpoint paths**: Full URL patterns with path parameters
- **Request schemas**: Required/optional fields, validation rules
- **Response schemas**: Structure and field types
- **Status codes**: Success (200, 201) and error codes (400, 401, 403, 404, 409, 500)
- **Authentication requirements**: Bearer token, user authorization
- **API documentation**: Chinese descriptions of behavior

**Example findings:**

```
GET /users/{userId}/missions/{missionId}/progress
  - Auth: JWT required
  - Response: { missionId, status, watchPositionSeconds }
  - Status: UNCOMPLETED, COMPLETED, DELIVERED
```

### Step 1.3: Read Database Schema (db-schema.dbml)

Location: `/docs/db-schema.dbml`

Identify:

- **Relevant tables**: Which tables store the domain data?
- **Column definitions**: Types, constraints, defaults
- **Relationships**: Foreign keys and joins
- **Indexes**: Performance optimization hints
- **Enums**: Valid status values and types

**Example findings:**

```
Table user_mission_progress {
  user_id, mission_id (composite unique key)
  status (enum: UNCOMPLETED, COMPLETED, DELIVERED)
  watch_position_seconds (integer, default 0)
}
```

## Phase 2: Codebase Exploration

### Step 2.1: Launch Explore Agent (Parallel Search)

Use the **Task tool with subagent_type=Explore** to search the codebase efficiently.

**Search focus areas:**

1. **Existing handlers**: Search for similar endpoint patterns
2. **Services**: Look for related business logic
3. **Repositories**: Find data access patterns
4. **Models**: Check domain model definitions
5. **DTOs**: Examine request/response structures
6. **Error handling**: Understand error codes and patterns
7. **Middleware**: Authentication and authorization patterns
8. **Database migrations**: Verify table schemas match documentation

**Key files to find:**

```
internal/handler/*_handler.go
internal/service/*_service.go
internal/repository/*_repository.go
internal/model/*.go
internal/dto/*.go
internal/apperror/*.go
internal/middleware/*.go
internal/router/router.go
internal/app/app.go
migrations/*.sql
```

**Agent prompt template:**

```
Explore the codebase to understand the current state of [FEATURE] implementation.
Find:
1. Existing handlers, services, repositories related to [DOMAIN]
2. Database schema and models for [TABLES]
3. Any existing DTOs or API endpoints for [PATH PATTERN]
4. Authentication/authorization patterns
5. Error handling patterns

Focus on files in: internal/{handler,service,repository,dto,model,middleware}/
Provide a summary of what exists and what's missing.
```

### Step 2.2: Analyze Findings

Categorize discoveries:

**✅ What Exists:**

- Database tables and migrations
- Related domain models
- Authentication middleware
- Error handling framework
- Architectural patterns
- Similar implementations to reference

**❌ What's Missing:**

- New domain models
- DTOs for requests/responses
- SQLc queries
- Repository methods
- Service business logic
- Handler endpoints
- Error codes
- Route registrations
- Dependency wiring

## Phase 3: Architecture Analysis

### Step 3.1: Understand Project Patterns

Read [CLAUDE.md](../../CLAUDE.md) to understand:

- **Layer architecture**: Handler → Service → Repository → Database
- **Configuration patterns**: Pass-by-value vs pass-by-pointer
- **Error handling**: AppError system, centralized middleware
- **Naming conventions**: No abbreviations (progressRepository not progressRepo)
- **Code organization**: Where files should be created
- **Testing patterns**: BDD test structure

### Step 3.2: Map to Layer Architecture

For each missing piece, determine its layer:

```
Presentation Layer (Handler):
  - Extract path params and query params
  - Validate authorization (user can only access own resources)
  - Bind JSON request bodies
  - Call service layer
  - Convert domain models to DTOs
  - Return HTTP responses

Business Logic Layer (Service):
  - Validate business rules
  - Coordinate repository calls
  - Calculate derived values (e.g., auto-completion status)
  - Handle complex workflows
  - Wrap errors with AppError

Data Access Layer (Repository):
  - Execute SQLc queries
  - Convert between sqlc models and domain models
  - Handle pgx.ErrNoRows gracefully
  - Manage transactions if needed

Data Layer (SQLc Queries):
  - Write type-safe SQL queries
  - Define query parameters
  - Return structs matching database schema
```

### Step 3.3: Identify Dependencies

For each new component, list what it depends on:

**Example:**

```
ProgressHandler depends on:
  - ProgressService (constructor injection)
  - Logger (*slog.Logger)
  - Timeout (time.Duration)

ProgressService depends on:
  - ProgressRepository (constructor injection)
  - Logger (*slog.Logger)
  - Timeout (time.Duration)

ProgressRepository depends on:
  - Database pool (*pgxpool.Pool)
  - SQLc Querier (auto-generated)
```

## Phase 4: Plan Creation

### Step 4.1: Structure the Implementation Plan

Create a plan file with these sections:

1. **Overview**: Brief description of the feature
2. **References**: Links to ISA feature, API spec, DB schema, project guide
3. **Context Summary**: What exists vs what's missing
4. **API Requirements**: Endpoint specifications with examples
5. **Implementation Steps**: Detailed, ordered steps
6. **Critical Files**: Create new vs modify existing
7. **Key Business Rules**: Important logic to implement correctly
8. **Verification Plan**: How to test the implementation
9. **Architecture Alignment**: Patterns to follow
10. **Success Criteria**: Definition of done

### Step 4.2: Write Detailed Implementation Steps

For each step, include:

**Step template:**

```markdown
### Step X: [Action Title]

**File:** [path/to/file.go] (create new / modify existing)

**What to do:**

- Bullet point 1 with specific code pattern
- Bullet point 2 with specific method signature
- Bullet point 3 with error handling approach

**Code structure:**
[Show example struct, function signature, or code snippet]

**Dependencies:**

- What this depends on
- What needs to be done first

**Notes:**

- Important considerations
- Common pitfalls to avoid
```

### Step 4.3: Sequence Steps Correctly

Order steps to minimize dependency issues:

1. **Foundation first**: Models, DTOs, constants
2. **Database layer**: SQLc queries (run `make sqlc` after)
3. **Data access**: Repository layer
4. **Business logic**: Service layer
5. **Presentation**: Handler layer
6. **Infrastructure**: Error codes, routes, wiring

**Correct order example:**

```
Step 1: Domain Model (no dependencies)
Step 2: DTOs (depends on domain model)
Step 3: SQLc Queries (no dependencies)
Step 4: Repository (depends on Step 3 auto-generated code)
Step 5: Service (depends on repository)
Step 6: Handler (depends on service)
Step 7: Error codes (no dependencies, can be earlier)
Step 8: Routes (depends on handler)
Step 9: App wiring (depends on all components)
```

### Step 4.4: Include Business Logic Details

For service layer, be explicit about:

- **Validation rules**: What to check and when
- **Calculations**: Formulas and conditions
- **State transitions**: Valid status changes
- **Edge cases**: How to handle nulls, not found, etc.

**Example:**

```markdown
### Calculate Completion Status

Logic:

- Get video duration from mission_resources table
- If watchPositionSeconds >= duration → status = "COMPLETED"
- If watchPositionSeconds < duration → status = "UNCOMPLETED"
- Allow COMPLETED to remain COMPLETED (re-watching)

Edge cases:

- No video resource → return error
- watchPositionSeconds negative → return validation error
- watchPositionSeconds > duration → allowed, set to COMPLETED
```

### Step 4.5: Specify Error Handling

For each error scenario, define:

- **Error code**: ERR_SOMETHING
- **Error message**: Chinese user-facing message
- **HTTP status**: 400, 401, 403, 404, 409, 500
- **When to throw**: Specific conditions

**Example:**

```markdown
Add error codes:

1. CodeUnauthorizedAccess
   - Message: "無法存取其他使用者的資源"
   - Status: 403 Forbidden
   - When: userId in path doesn't match JWT user

2. CodeInvalidWatchPosition
   - Message: "無效的觀看位置"
   - Status: 400 Bad Request
   - When: watchPositionSeconds < 0
```

### Step 4.6: Define Verification Steps

Create a multi-phase verification plan:

**Phase 1: Build Verification**

```bash
make sqlc  # Generate code
make fmt   # Format
make lint  # Check code quality
make build # Compile
```

**Phase 2: Test Execution**

```bash
cd tests/bdd
go test -v -tags=isa ./...
```

**Phase 3: Scenario Coverage**
List each ISA scenario with expected behavior:

```
1. ✅ Scenario: 首次觀看影片應從頭開始
   - GET /progress returns watchPositionSeconds=0

2. ✅ Scenario: 觀看完成後應自動標記為已完成
   - PUT with position >= duration sets status=COMPLETED
```

**Phase 4: Manual Testing (Optional)**
Provide curl commands for manual API testing.

## Phase 5: Plan Review

### Step 5.1: Cross-Reference Documentation

Verify plan matches:

- ✅ ISA feature scenarios
- ✅ Swagger API specifications
- ✅ Database schema definitions
- ✅ Project coding standards (CLAUDE.md)

### Step 5.2: Check Completeness

Ensure plan includes:

- ✅ All files that need creation
- ✅ All files that need modification
- ✅ All error cases handled
- ✅ Authorization checks in place
- ✅ SQLc generation command
- ✅ Dependency wiring complete
- ✅ Routes registered correctly

### Step 5.3: Validate Architecture Alignment

Confirm plan follows:

- ✅ Layer-based architecture pattern
- ✅ Config pass-by-value pattern
- ✅ Centralized error handling
- ✅ JWT authentication flow
- ✅ Naming conventions (no abbreviations)
- ✅ Existing code patterns

## Output Format

The final plan should be a markdown file with:

```markdown
# Implementation Plan: [Feature Name]

## Overview

[Brief description]

## References

- ISA Feature: [link]
- API Spec: [link]
- DB Schema: [link]
- Project Guide: [link]

## Context Summary

### What Exists

- ✅ [List existing infrastructure]

### What's Missing

- ❌ [List gaps to fill]

## API Requirements

### [Endpoint 1]

**Purpose:** [What it does]
**Request:** [JSON example]
**Response:** [JSON example]
**Business Logic:** [Key rules]

## Implementation Steps

### Step 1: [Action]

**File:** [path] (create/modify)
[Detailed instructions]

### Step 2: [Action]

...

## Critical Files to Modify/Create

[Organized list with file paths]

## Key Business Rules

[Numbered list of important logic]

## Verification Plan

[Multi-phase testing approach]

## Architecture Alignment

[Checklist of patterns to follow]

## Success Criteria

[Definition of done with checkboxes]
```

# Key Principles

1. **Be Specific**: Don't say "implement the handler" - say "create GetProgress method that extracts userId from c.Param('userId')"

2. **Show Code Structure**: Provide function signatures, struct definitions, SQL queries

3. **Explain Dependencies**: Make clear what must be done before what

4. **Include Examples**: Show JSON request/response, code snippets, SQL queries

5. **Handle Errors**: Don't forget error codes, validation, authorization

6. **Follow Patterns**: Reference existing code that follows the same pattern

7. **Be Actionable**: Each step should be clear enough to implement without guessing

8. **Verify Completeness**: Cross-reference all documentation sources

# Common Pitfalls to Avoid

❌ **Don't:**

- Skip authorization checks (userId validation)
- Forget to register routes in router.go
- Miss dependency wiring in app.go
- Ignore error handling for database queries
- Use abbreviations in variable names
- Create pointer configs (use values)
- Forget to run `make sqlc` after adding queries

✅ **Do:**

- Validate user can only access own resources
- Wire all dependencies through constructors
- Handle pgx.ErrNoRows gracefully
- Use full names (progressRepository not progressRepo)
- Follow pass-by-value for config
- Run code generation commands
- Test all BDD scenarios

# Success Indicators

A good implementation plan should:

- ✅ Be implementable by another developer without asking questions
- ✅ Cover all scenarios in the ISA feature file
- ✅ Match the API specification exactly
- ✅ Follow all project architectural patterns
- ✅ Include proper error handling
- ✅ Have clear verification steps
- ✅ Reference specific file paths and line numbers where helpful

# Example Usage

**User asks:**

```
Given this ISA feature file: tests/bdd/features/isa/progress/watch-video-and-track-progress.isa.feature
Please analyze and create an implementation plan.
```

**Skill executes:**

1. Read the feature file
2. Read swagger.yaml for progress endpoints
3. Read db-schema.dbml for user_mission_progress table
4. Launch Explore agent to find existing progress-related code
5. Analyze findings and identify gaps
6. Create detailed implementation plan with 10 steps
7. Include business rules (auto-completion logic)
8. Define verification with 9 test scenarios
9. Write plan to file

**Output:**
A comprehensive markdown plan that enables implementation of the complete feature with all endpoints, business logic, error handling, and tests passing.
