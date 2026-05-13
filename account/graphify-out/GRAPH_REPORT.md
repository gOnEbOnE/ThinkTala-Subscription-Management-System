# Graph Report - account  (2026-04-20)

## Corpus Check
- Corpus is ~34,479 words - fits in a single context window. You may not need a graph.

## Summary
- 242 nodes · 368 edges · 26 communities detected
- Extraction: 66% EXTRACTED · 34% INFERRED · 0% AMBIGUOUS · INFERRED: 125 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_App Config|App Config]]
- [[_COMMUNITY_Main Pkcs|Main Pkcs]]
- [[_COMMUNITY_Circuit Breaker|Circuit Breaker]]
- [[_COMMUNITY_Router Logout|Router Logout]]
- [[_COMMUNITY_Repository Connection|Repository Connection]]
- [[_COMMUNITY_Form Secureform|Form Secureform]]
- [[_COMMUNITY_Endpoint Table|Endpoint Table]]
- [[_COMMUNITY_Controller Response|Controller Response]]
- [[_COMMUNITY_Manager Redis|Manager Redis]]
- [[_COMMUNITY_Overview Calendar|Overview Calendar]]
- [[_COMMUNITY_String Generaterandomnumber|String Generaterandomnumber]]
- [[_COMMUNITY_Generator Uploader|Generator Uploader]]
- [[_COMMUNITY_Image Pngcompression|Image Pngcompression]]
- [[_COMMUNITY_Mail Isvalidemail|Mail Isvalidemail]]
- [[_COMMUNITY_Validator Validatorerror|Validator Validatorerror]]
- [[_COMMUNITY_Legacy Tawk|Legacy Tawk]]
- [[_COMMUNITY_Types Currentuser|Types Currentuser]]
- [[_COMMUNITY_Theme Inittheme|Theme Inittheme]]
- [[_COMMUNITY_Navigate Logoutuser|Navigate Logoutuser]]
- [[_COMMUNITY_Service Newservice|Service Newservice]]
- [[_COMMUNITY_Time Formatdateid|Time Formatdateid]]
- [[_COMMUNITY_Zaframework Documentation|Zaframework Documentation]]
- [[_COMMUNITY_Uuid Createuuid|Uuid Createuuid]]
- [[_COMMUNITY_Toast Injecttoasthtml|Toast Injecttoasthtml]]
- [[_COMMUNITY_Pinkai Marketing|Pinkai Marketing]]
- [[_COMMUNITY_System Summary|System Summary]]

## God Nodes (most connected - your core abstractions)
1. `New()` - 24 edges
2. `Get()` - 15 edges
3. `SecureForm` - 13 edges
4. `GetEnv()` - 12 edges
5. `main()` - 11 edges
6. `Set()` - 11 edges
7. `Dispatcher` - 10 edges
8. `HandlerFunc` - 10 edges
9. `AuthMiddleware()` - 9 edges
10. `Authorize()` - 8 edges

## Surprising Connections (you probably didn't know these)
- `New()` --calls--> `Connect()`  [INFERRED]
  account\core\app.go → account\core\database\connection.go
- `init()` --calls--> `New()`  [INFERRED]
  account\core\utils\validator.go → account\core\app.go
- `NewSMTPClient()` --calls--> `GetEnv()`  [INFERRED]
  account\core\utils\mail.go → account\core\utils\env.go
- `main()` --calls--> `LoadEnv()`  [INFERRED]
  account\main.go → account\core\utils\env.go
- `main()` --calls--> `InitJWTLoadKeys()`  [INFERRED]
  account\main.go → account\core\utils\jwt.go

## Hyperedges (group relationships)
- **Authentication Form Flow** — login_page_modern_login, login1_page_legacy_login, register_page_signup, reset_page_password_reset, auth_route_login [EXTRACTED 0.92]
- **User Administration Table Views** — users_page_user_management, active_page_users_table, banned_page_users_table, inactive_page_user_management, active_api_users_endpoint [EXTRACTED 0.90]
- **MacroQuant Analysis Stack** — index_macroquant_overview, index2_macroquant_suite, macrodata_macro_calendar, macroquant_neuromac_model, macroquant_central_bank_bias, macroquant_economic_calendar [INFERRED 0.84]

## Communities

### Community 0 - "App Config"
Cohesion: 0.09
Nodes (16): Config, Dispatcher, HandlerFunc, Job, JobResult, JobType, SlidingWindowRateLimiter, App (+8 more)

### Community 1 - "Main Pkcs"
Cohesion: 0.12
Nodes (20): New(), ToInt(), Decrypt(), Encrypt(), PKCS5Padding(), PKCS5UnPadding(), GetEnv(), CreateJWT() (+12 more)

### Community 2 - "Circuit Breaker"
Cohesion: 0.11
Nodes (8): NewCircuitBreaker(), EnhancedCircuitBreaker, EnhancedMetrics, WorkerPool, NewDispatcher(), NewEnhancedMetrics(), NewSlidingWindowRateLimiter(), newWorkerPool()

### Community 3 - "Router Logout"
Cohesion: 0.21
Nodes (12): GetAgent(), GetClientIP(), GetUserAgent(), Destroy(), extractToken(), Get(), getManager(), Set() (+4 more)

### Community 4 - "Repository Connection"
Cohesion: 0.14
Nodes (7): Repository, userRepo, Connect(), Config, DBWrapper, dummyScanner, LoadEnv()

### Community 5 - "Form Secureform"
Cohesion: 0.27
Nodes (1): SecureForm

### Community 6 - "Endpoint Table"
Cohesion: 0.18
Nodes (13): Admin Panel Shell, Users API Endpoint, Active Users Table View, Login Route, Banned Users Table View, User Detail Endpoint, User Update Endpoint, Inactive Users Management (+5 more)

### Community 7 - "Controller Response"
Cohesion: 0.18
Nodes (5): Controller, NewController(), Controller, JSONResponse, ResponseHelper

### Community 8 - "Manager Redis"
Cohesion: 0.22
Nodes (7): Init(), GetRedisClient(), IsRedisEnabled(), RedisGet(), RedisSet(), Config, Manager

### Community 9 - "Overview Calendar"
Cohesion: 0.25
Nodes (9): MacroQuant Full Suite Landing, MacroQuant Overview Landing (Copy 2), MacroQuant Overview Landing, Macro Data Calendar View, Central Bank Bias Module, AI Economic Calendar Module, NeuroMac v4.2 Model, ThinkArah AI Engine (+1 more)

### Community 10 - "String Generaterandomnumber"
Cohesion: 0.38
Nodes (3): GenerateRandomNumber(), GenerateRandomString(), secureRandom()

### Community 11 - "Generator Uploader"
Cohesion: 0.38
Nodes (4): RandomString(), isAllowedExtension(), UploadFile(), validateMimeType()

### Community 12 - "Image Pngcompression"
Cohesion: 0.6
Nodes (4): pngCompression(), ProcessImage(), saveImage(), ImagePreset

### Community 13 - "Mail Isvalidemail"
Cohesion: 0.4
Nodes (2): NewSMTPClient(), SMTPClient

### Community 14 - "Validator Validatorerror"
Cohesion: 0.5
Nodes (4): ValidatorError, init(), msgForTag(), ValidateStruct()

### Community 15 - "Legacy Tawk"
Cohesion: 0.5
Nodes (5): Legacy Login Page, Tawk.to Chat Widget, Modern Login Page, Register Route, Reset Route

### Community 16 - "Types Currentuser"
Cohesion: 0.5
Nodes (3): CurrentUser, Login, User

### Community 17 - "Theme Inittheme"
Cohesion: 0.67
Nodes (2): initTheme(), updateThemeUI()

### Community 18 - "Navigate Logoutuser"
Cohesion: 0.5
Nodes (0): 

### Community 19 - "Service Newservice"
Cohesion: 0.67
Nodes (1): Service

### Community 20 - "Time Formatdateid"
Cohesion: 1.0
Nodes (2): FormatDateID(), TimeAgo()

### Community 21 - "Zaframework Documentation"
Cohesion: 0.67
Nodes (3): ZAFramework Documentation, Routing System (http.ServeMux Wrapper), Session Manager (Redis/Cookie/JWT)

### Community 22 - "Uuid Createuuid"
Cohesion: 1.0
Nodes (0): 

### Community 23 - "Toast Injecttoasthtml"
Cohesion: 1.0
Nodes (0): 

### Community 24 - "Pinkai Marketing"
Cohesion: 1.0
Nodes (2): PinkAI Marketing Landing, PinkAI SDK Multi-language Concept

### Community 25 - "System Summary"
Cohesion: 1.0
Nodes (1): System Summary Dashboard

## Knowledge Gaps
- **32 isolated node(s):** `Repository`, `Service`, `Login`, `User`, `CurrentUser` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **Thin community `Uuid Createuuid`** (2 nodes): `uuid.go`, `CreateUUID()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Toast Injecttoasthtml`** (2 nodes): `za-toast.js`, `injectToastHTML()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Pinkai Marketing`** (2 nodes): `PinkAI Marketing Landing`, `PinkAI SDK Multi-language Concept`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `System Summary`** (1 nodes): `System Summary Dashboard`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `New()` connect `Main Pkcs` to `App Config`, `Circuit Breaker`, `Router Logout`, `Repository Connection`, `Controller Response`, `Generator Uploader`, `Validator Validatorerror`?**
  _High betweenness centrality (0.250) - this node is a cross-community bridge._
- **Why does `main()` connect `Main Pkcs` to `App Config`, `Repository Connection`, `Form Secureform`, `Controller Response`?**
  _High betweenness centrality (0.118) - this node is a cross-community bridge._
- **Why does `GetEnv()` connect `Main Pkcs` to `App Config`, `Router Logout`, `Repository Connection`, `Mail Isvalidemail`?**
  _High betweenness centrality (0.066) - this node is a cross-community bridge._
- **Are the 23 inferred relationships involving `New()` (e.g. with `main()` and `GetEnv()`) actually correct?**
  _`New()` has 23 INFERRED edges - model-reasoned connections that need verification._
- **Are the 10 inferred relationships involving `Get()` (e.g. with `.Logout()` and `CORSMiddleware()`) actually correct?**
  _`Get()` has 10 INFERRED edges - model-reasoned connections that need verification._
- **Are the 11 inferred relationships involving `GetEnv()` (e.g. with `main()` and `.Logout()`) actually correct?**
  _`GetEnv()` has 11 INFERRED edges - model-reasoned connections that need verification._
- **Are the 10 inferred relationships involving `main()` (e.g. with `LoadEnv()` and `InitJWTLoadKeys()`) actually correct?**
  _`main()` has 10 INFERRED edges - model-reasoned connections that need verification._