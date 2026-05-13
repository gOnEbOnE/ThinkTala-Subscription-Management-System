# Graph Report - subscription  (2026-04-20)

## Corpus Check
- Corpus is ~39,671 words - fits in a single context window. You may not need a graph.

## Summary
- 320 nodes · 537 edges · 24 communities detected
- Extraction: 59% EXTRACTED · 41% INFERRED · 0% AMBIGUOUS · INFERRED: 218 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Createpackage Deletepackage|Createpackage Deletepackage]]
- [[_COMMUNITY_Pkcs Manager|Pkcs Manager]]
- [[_COMMUNITY_Controller Welcome|Controller Welcome]]
- [[_COMMUNITY_Config Shutdown|Config Shutdown]]
- [[_COMMUNITY_Newcircuitbreaker Enhancedcircuitbreaker|Newcircuitbreaker Enhancedcircuitbreaker]]
- [[_COMMUNITY_Inactive Form|Inactive Form]]
- [[_COMMUNITY_Repository Main|Repository Main]]
- [[_COMMUNITY_Dispatchandwait Writeheader|Dispatchandwait Writeheader]]
- [[_COMMUNITY_Form Secureform|Form Secureform]]
- [[_COMMUNITY_Service Createjwt|Service Createjwt]]
- [[_COMMUNITY_Newuuid Randomstring|Newuuid Randomstring]]
- [[_COMMUNITY_Generaterandomnumber Generaterandomstring|Generaterandomnumber Generaterandomstring]]
- [[_COMMUNITY_Loginpayload Loginresult|Loginpayload Loginresult]]
- [[_COMMUNITY_Pngcompression Processimage|Pngcompression Processimage]]
- [[_COMMUNITY_Isvalidemail Newsmtpclient|Isvalidemail Newsmtpclient]]
- [[_COMMUNITY_Validator Validatorerror|Validator Validatorerror]]
- [[_COMMUNITY_Types Userlistdto|Types Userlistdto]]
- [[_COMMUNITY_Createpackagedto Package|Createpackagedto Package]]
- [[_COMMUNITY_Theme Inittheme|Theme Inittheme]]
- [[_COMMUNITY_Navigate Logoutuser|Navigate Logoutuser]]
- [[_COMMUNITY_Client Auth|Client Auth]]
- [[_COMMUNITY_Multi Section|Multi Section]]
- [[_COMMUNITY_Time Formatdateid|Time Formatdateid]]
- [[_COMMUNITY_Toast Injecttoasthtml|Toast Injecttoasthtml]]

## God Nodes (most connected - your core abstractions)
1. `New()` - 28 edges
2. `Get()` - 18 edges
3. `main()` - 14 edges
4. `packageService` - 13 edges
5. `SecureForm` - 13 edges
6. `Set()` - 12 edges
7. `GetEnv()` - 11 edges
8. `Dispatcher` - 10 edges
9. `NewController()` - 9 edges
10. `packageRepo` - 9 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `InitJWTLoadKeys()`  [INFERRED]
  subscription\main.go → subscription\core\utils\jwt.go
- `New()` --calls--> `Connect()`  [INFERRED]
  subscription\core\app.go → subscription\core\database\connection.go
- `init()` --calls--> `New()`  [INFERRED]
  subscription\core\utils\validator.go → subscription\core\app.go
- `NewSMTPClient()` --calls--> `GetEnv()`  [INFERRED]
  subscription\core\utils\mail.go → subscription\core\utils\env.go
- `main()` --calls--> `LoadEnv()`  [INFERRED]
  subscription\main.go → subscription\core\utils\env.go

## Hyperedges (group relationships)
- **Account User Management Views** — account_wrapper_page_admin_shell, account_users_active_page_user_table, account_users_inactive_page_user_table, account_users_banned_page_blocked_users, account_users_page_user_table [INFERRED 0.78]
- **Authentication UI Flow** — login_page_login_form, login1_page_legacy_login_form, register_page_registration_form, reset_page_password_reset_form [INFERRED 0.83]
- **MacroQuant and ThinkArah Concept Cluster** — landing_index_macroquant_dashboard, landing_macrodata_macro_dashboard, landing_index2_macroquant_spa, landing_thinkarah_project_blueprint [INFERRED 0.74]

## Communities

### Community 0 - "Createpackage Deletepackage"
Cohesion: 0.08
Nodes (9): Connect(), Config, DBWrapper, dummyScanner, LoadEnv(), packageRepo, packageService, Service (+1 more)

### Community 1 - "Pkcs Manager"
Cohesion: 0.1
Nodes (33): GetAgent(), GetClientIP(), GetUserAgent(), New(), Decrypt(), Encrypt(), PKCS5Padding(), PKCS5UnPadding() (+25 more)

### Community 2 - "Controller Welcome"
Cohesion: 0.08
Nodes (10): NewController(), Controller, JSONResponse, ResponseHelper, Controller, Controller, Controller, Controller (+2 more)

### Community 3 - "Config Shutdown"
Cohesion: 0.09
Nodes (16): Config, Dispatcher, HandlerFunc, Job, JobResult, JobType, SlidingWindowRateLimiter, App (+8 more)

### Community 4 - "Newcircuitbreaker Enhancedcircuitbreaker"
Cohesion: 0.11
Nodes (8): NewCircuitBreaker(), EnhancedCircuitBreaker, EnhancedMetrics, WorkerPool, NewDispatcher(), NewEnhancedMetrics(), NewSlidingWindowRateLimiter(), newWorkerPool()

### Community 5 - "Inactive Form"
Cohesion: 0.1
Nodes (20): Account System Summary Dashboard, Active Users Table View, Banned Users View, Users List API (/api/users), Inactive Users Data API (/account/users/data/inactive), User Detail API (/account/users/detail/:id), User Update API (/account/users/update/:id), Inactive Users Table View (+12 more)

### Community 6 - "Repository Main"
Cohesion: 0.13
Nodes (8): ToInt(), Repository, userRepo, main(), MigrateAndSeed(), Repository, NewRepository(), Repository

### Community 7 - "Dispatchandwait Writeheader"
Cohesion: 0.29
Nodes (2): Controller, Init()

### Community 8 - "Form Secureform"
Cohesion: 0.25
Nodes (1): SecureForm

### Community 9 - "Service Createjwt"
Cohesion: 0.24
Nodes (4): CreateJWT(), Service, Service, NewService()

### Community 10 - "Newuuid Randomstring"
Cohesion: 0.32
Nodes (4): RandomString(), isAllowedExtension(), UploadFile(), validateMimeType()

### Community 11 - "Generaterandomnumber Generaterandomstring"
Cohesion: 0.38
Nodes (3): GenerateRandomNumber(), GenerateRandomString(), secureRandom()

### Community 12 - "Loginpayload Loginresult"
Cohesion: 0.4
Nodes (4): Login, LoginPayload, LoginResult, User

### Community 13 - "Pngcompression Processimage"
Cohesion: 0.6
Nodes (4): pngCompression(), ProcessImage(), saveImage(), ImagePreset

### Community 14 - "Isvalidemail Newsmtpclient"
Cohesion: 0.4
Nodes (2): NewSMTPClient(), SMTPClient

### Community 15 - "Validator Validatorerror"
Cohesion: 0.5
Nodes (4): ValidatorError, init(), msgForTag(), ValidateStruct()

### Community 16 - "Types Userlistdto"
Cohesion: 0.5
Nodes (3): User, UserListDTO, UserUpdateDTO

### Community 17 - "Createpackagedto Package"
Cohesion: 0.5
Nodes (3): CreatePackageDTO, Package, UpdatePackageDTO

### Community 18 - "Theme Inittheme"
Cohesion: 0.67
Nodes (2): initTheme(), updateThemeUI()

### Community 19 - "Navigate Logoutuser"
Cohesion: 0.5
Nodes (0): 

### Community 20 - "Client Auth"
Cohesion: 0.5
Nodes (4): Auth Logout API, Client Dashboard Page, Client KYC Route, Client Subscription Route

### Community 21 - "Multi Section"
Cohesion: 0.5
Nodes (4): MacroQuant Multi-Section SPA, MacroQuant Dashboard Landing, Macro Data Landing Dashboard, ThinkArah Project Blueprint

### Community 22 - "Time Formatdateid"
Cohesion: 1.0
Nodes (2): FormatDateID(), TimeAgo()

### Community 23 - "Toast Injecttoasthtml"
Cohesion: 1.0
Nodes (0): 

## Ambiguous Edges - Review These
- `Account Wrapper Admin Shell` → `Legacy Login Form`  [AMBIGUOUS]
  subscription/public/views/login1/page.html · relation: references

## Knowledge Gaps
- **46 isolated node(s):** `Repository`, `UserListDTO`, `UserUpdateDTO`, `User`, `Repository` (+41 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **Thin community `Toast Injecttoasthtml`** (2 nodes): `za-toast.js`, `injectToastHTML()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **What is the exact relationship between `Account Wrapper Admin Shell` and `Legacy Login Form`?**
  _Edge tagged AMBIGUOUS (relation: references) - confidence is low._
- **Why does `New()` connect `Pkcs Manager` to `Createpackage Deletepackage`, `Controller Welcome`, `Config Shutdown`, `Newcircuitbreaker Enhancedcircuitbreaker`, `Repository Main`, `Service Createjwt`, `Newuuid Randomstring`, `Validator Validatorerror`?**
  _High betweenness centrality (0.310) - this node is a cross-community bridge._
- **Why does `main()` connect `Repository Main` to `Createpackage Deletepackage`, `Pkcs Manager`, `Controller Welcome`, `Config Shutdown`, `Form Secureform`, `Service Createjwt`?**
  _High betweenness centrality (0.156) - this node is a cross-community bridge._
- **Why does `GetEnv()` connect `Pkcs Manager` to `Createpackage Deletepackage`, `Config Shutdown`, `Repository Main`, `Isvalidemail Newsmtpclient`?**
  _High betweenness centrality (0.054) - this node is a cross-community bridge._
- **Are the 27 inferred relationships involving `New()` (e.g. with `main()` and `.CreatePackage()`) actually correct?**
  _`New()` has 27 INFERRED edges - model-reasoned connections that need verification._
- **Are the 13 inferred relationships involving `Get()` (e.g. with `.GetInactiveUsers()` and `.Logout()`) actually correct?**
  _`Get()` has 13 INFERRED edges - model-reasoned connections that need verification._
- **Are the 13 inferred relationships involving `main()` (e.g. with `LoadEnv()` and `InitJWTLoadKeys()`) actually correct?**
  _`main()` has 13 INFERRED edges - model-reasoned connections that need verification._