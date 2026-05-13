# Graph Report - management  (2026-04-21)

## Corpus Check
- Corpus is ~5,366 words - fits in a single context window. You may not need a graph.

## Summary
- 68 nodes · 132 edges · 6 communities detected
- Extraction: 100% EXTRACTED · 0% INFERRED · 0% AMBIGUOUS
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Community 0|Community 0]]
- [[_COMMUNITY_Community 1|Community 1]]
- [[_COMMUNITY_Community 2|Community 2]]
- [[_COMMUNITY_Community 3|Community 3]]
- [[_COMMUNITY_Community 4|Community 4]]
- [[_COMMUNITY_Community 5|Community 5]]

## God Nodes (most connected - your core abstractions)
1. `buildDummyDashboardPayload()` - 16 edges
2. `buildDummyPackageDashboardPayload()` - 12 edges
3. `buildDummyPackageDetailPayload()` - 10 edges
4. `setAuditError()` - 7 edges
5. `main()` - 6 edges
6. `round2()` - 6 edges
7. `getDashboardCustomers()` - 5 edges
8. `getDashboardPackages()` - 5 edges
9. `compareRate()` - 5 edges
10. `getDashboardPackageDetail()` - 4 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `requirePackageDashboardRole()`  [EXTRACTED]
  management\main.go → management\main.go  _Bridges community 4 → community 3_
- `getDashboardCustomers()` --calls--> `buildDummyDashboardPayload()`  [EXTRACTED]
  management\main.go → management\main.go  _Bridges community 3 → community 2_
- `getDashboardPackages()` --calls--> `buildDummyPackageDashboardPayload()`  [EXTRACTED]
  management\main.go → management\main.go  _Bridges community 3 → community 1_
- `buildDummyPackageDashboardPayload()` --calls--> `calcTotalPages()`  [EXTRACTED]
  management\main.go → management\main.go  _Bridges community 1 → community 2_

## Communities

### Community 0 - "Community 0"
Cohesion: 0.08
Nodes (17): churnedCustomerItem, churnRatePoint, customerDetailPayload, dashboardCustomerPayload, dashboardPackagePayload, dummyCustomer, dummyPackageTransaction, loyalCustomerItem (+9 more)

### Community 1 - "Community 1"
Cohesion: 0.27
Nodes (13): aggregatePackageMetrics(), buildDummyPackageDashboardPayload(), buildDummyPackageDetailPayload(), buildMonthBuckets(), calculateRate(), compareFloatRate(), compareRate(), getDummyPackageTransactions() (+5 more)

### Community 2 - "Community 2"
Cohesion: 0.23
Nodes (12): buildDummyCustomerDetail(), buildDummyDashboardPayload(), calcTotalPages(), countDummyActiveAtDate(), countDummyActiveInRange(), countDummyChurnedInRange(), dummyDurationMonths(), getDashboardCustomerDetail() (+4 more)

### Community 3 - "Community 3"
Cohesion: 0.36
Nodes (8): getDashboardCustomers(), getDashboardPackageDetail(), getDashboardPackages(), handleServerError(), parseIntDefault(), parsePeriod(), requirePackageDashboardRole(), setAuditError()

### Community 4 - "Community 4"
Cohesion: 0.33
Nodes (6): auditMiddleware(), ensureManagementSchema(), getEnv(), main(), requireManagementRole(), saveAuditLog()

### Community 5 - "Community 5"
Cohesion: 0.67
Nodes (3): fetchMonthlyChurnCustomers(), fetchMonthlyNewCustomers(), monthKey()

## Knowledge Gaps
- **17 isolated node(s):** `monthlyPoint`, `churnRatePoint`, `loyalCustomerItem`, `churnedCustomerItem`, `dashboardCustomerPayload` (+12 more)
  These have ≤1 connection - possible missing edges or undocumented components.