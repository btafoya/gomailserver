# Add pagination to admin List endpoints

## Overview

Implement pagination for domain, user, and alias List handlers using the existing PaginatedResponse middleware. Currently these handlers return all records without pagination.

## Rationale

The PaginatedResponse structure exists in middleware/responses.go and is ready to use. The webmail handlers already implement pagination (page/limit/offset pattern). The DomainService even has a TODO comment about pagination support.

---
*This spec was created from ideation and is pending detailed specification.*
