# Graph Report - users  (2026-04-20)

## Corpus Check
- 73 files · ~50,607 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 388 nodes · 735 edges · 28 communities detected
- Extraction: 54% EXTRACTED · 45% INFERRED · 0% AMBIGUOUS · INFERRED: 332 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Repository Createuser|Repository Createuser]]
- [[_COMMUNITY_Controller Submit|Controller Submit]]
- [[_COMMUNITY_App Pkcs|App Pkcs]]
- [[_COMMUNITY_Dispatcher Shutdown|Dispatcher Shutdown]]
- [[_COMMUNITY_Service Manager|Service Manager]]
- [[_COMMUNITY_Getcsrftoken Jsonresponse|Getcsrftoken Jsonresponse]]
- [[_COMMUNITY_Endpoint Otp|Endpoint Otp]]
- [[_COMMUNITY_Endpoint Inactive|Endpoint Inactive]]
- [[_COMMUNITY_Createuserinput Createuserresult|Createuserinput Createuserresult]]
- [[_COMMUNITY_Kycdetailresult Kyclistitem|Kycdetailresult Kyclistitem]]
- [[_COMMUNITY_Assumeroleinput Assumeroleresult|Assumeroleinput Assumeroleresult]]
- [[_COMMUNITY_Randomstring Slugify|Randomstring Slugify]]
- [[_COMMUNITY_Generaterandomnumber Generaterandomstring|Generaterandomnumber Generaterandomstring]]
- [[_COMMUNITY_Database Concept|Database Concept]]
- [[_COMMUNITY_Otprecord Otpverifyinput|Otprecord Otpverifyinput]]
- [[_COMMUNITY_Pngcompression Processimage|Pngcompression Processimage]]
- [[_COMMUNITY_Validator Validatorerror|Validator Validatorerror]]
- [[_COMMUNITY_Greed Rationale|Greed Rationale]]
- [[_COMMUNITY_Theme Inittheme|Theme Inittheme]]
- [[_COMMUNITY_Navigate Logoutuser|Navigate Logoutuser]]
- [[_COMMUNITY_Hashpwd Migrateandseed|Hashpwd Migrateandseed]]
- [[_COMMUNITY_Formatdateid Timeago|Formatdateid Timeago]]
- [[_COMMUNITY_Table Active|Table Active]]
- [[_COMMUNITY_Full Suite|Full Suite]]
- [[_COMMUNITY_Toast Injecttoasthtml|Toast Injecttoasthtml]]
- [[_COMMUNITY_Core Cluster|Core Cluster]]
- [[_COMMUNITY_Assume Role|Assume Role]]
- [[_COMMUNITY_Logout Cluster|Logout Cluster]]

## God Nodes (most connected - your core abstractions)
1. `Get()` - 35 edges
2. `New()` - 24 edges
3. `ApiJSON()` - 24 edges
4. `main()` - 15 edges
5. `GetEnv()` - 14 edges
6. `adminRepo` - 13 edges
7. `SecureForm` - 13 edges
8. `dispatchNotification()` - 12 edges
9. `kycRepo` - 12 edges
10. `Set()` - 12 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `NewAdminController()`  [INFERRED]
  users\main.go → users\app\modules\kyc\admin_controller.go
- `New()` --calls--> `Connect()`  [INFERRED]
  users\core\app.go → users\core\database\connection.go
- `init()` --calls--> `New()`  [INFERRED]
  users\core\utils\validator.go → users\core\app.go
- `main()` --calls--> `InitJWTLoadKeys()`  [INFERRED]
  users\main.go → users\core\utils\jwt.go
- `main()` --calls--> `ToInt()`  [INFERRED]
  users\main.go → users\core\utils\converts.go

## Hyperedges (group relationships)
- **Account User Administration Flow** — account_wrapper_navigation_shell, account_users_inactive_page, account_users_detail_api, account_users_update_api [INFERRED 0.86]
- **Authentication and Onboarding Flow** — auth_login_modern_page, auth_register_page, auth_verify_otp_page, auth_login_api_modern, auth_register_api, auth_verify_otp_api [INFERRED 0.89]
- **MacroQuant Analytics Experience Cluster** — landing_macroquant_dashboard, landing_macroquant_full_suite, landing_macroquant_macrodata_page, landing_thinkarah_blueprint [INFERRED 0.82]

## Communities

### Community 0 - "Repository Createuser"
Cohesion: 0.06
Nodes (17): adminRepo, Repository, Service, Connect(), Config, DBWrapper, dummyScanner, kycRepo (+9 more)

### Community 1 - "Controller Submit"
Cohesion: 0.09
Nodes (21): Controller, extractKYCID(), extractKYCIDFromAction(), NewAdminController(), GetAgent(), GetClientIP(), GetUserAgent(), NewController() (+13 more)

### Community 2 - "App Pkcs"
Cohesion: 0.08
Nodes (35): New(), HandlerFunc, ToInt(), App, Config, Decrypt(), Encrypt(), PKCS5Padding() (+27 more)

### Community 3 - "Dispatcher Shutdown"
Cohesion: 0.07
Nodes (14): NewCircuitBreaker(), Config, Dispatcher, EnhancedCircuitBreaker, EnhancedMetrics, Job, JobResult, JobType (+6 more)

### Community 4 - "Service Manager"
Cohesion: 0.11
Nodes (18): saveSessionToRedis(), Service, NewSMTPClient(), Init(), GetRedisClient(), IsRedisEnabled(), PublishNotificationEvent(), RedisSet() (+10 more)

### Community 5 - "Getcsrftoken Jsonresponse"
Cohesion: 0.13
Nodes (3): JSONResponse, ResponseHelper, SecureForm

### Community 6 - "Endpoint Otp"
Cohesion: 0.18
Nodes (13): Legacy Login Auth API Endpoint, Modern Login Auth API Endpoint, Legacy Login Page, Modern Login Page, Account Login Route, Account Reset Route, Register API Endpoint, Multi-Step Registration Page (+5 more)

### Community 7 - "Endpoint Inactive"
Cohesion: 0.22
Nodes (9): Account Dashboard Summary View, Banned Users Table View, User Detail API Endpoint, Inactive Users Data API Endpoint, Inactive Users Management View, User Update API Endpoint, Account Wrapper Navigation Shell, Account Settings Password View (+1 more)

### Community 8 - "Createuserinput Createuserresult"
Cohesion: 0.25
Nodes (7): CreateUserInput, CreateUserResult, EditUserInput, GetUsersParams, GetUsersResponse, UserDetail, UserListItem

### Community 9 - "Kycdetailresult Kyclistitem"
Cohesion: 0.25
Nodes (7): KYCDetailResult, KYCListItem, KYCReviewPayload, KYCStatusResult, KYCSubmission, KYCSubmitPayload, KYCSubmitResult

### Community 10 - "Assumeroleinput Assumeroleresult"
Cohesion: 0.25
Nodes (7): AssumeRoleInput, AssumeRoleResult, Login, LoginPayload, LoginResult, RoleInfo, User

### Community 11 - "Randomstring Slugify"
Cohesion: 0.38
Nodes (4): RandomString(), isAllowedExtension(), UploadFile(), validateMimeType()

### Community 12 - "Generaterandomnumber Generaterandomstring"
Cohesion: 0.38
Nodes (3): GenerateRandomNumber(), GenerateRandomString(), secureRandom()

### Community 13 - "Database Concept"
Cohesion: 0.33
Nodes (6): Database Bypass Mode Concept, Rationale: Keep App Running Without Database, Rationale: Background Jobs Avoid Blocking User Responses, Worker Pool Concurrency Concept, ZAFramework Documentation, Pink AI Marketing Landing Page

### Community 14 - "Otprecord Otpverifyinput"
Cohesion: 0.4
Nodes (4): OTPRecord, OTPVerifyInput, RegisterInput, RegisterResult

### Community 15 - "Pngcompression Processimage"
Cohesion: 0.6
Nodes (4): pngCompression(), ProcessImage(), saveImage(), ImagePreset

### Community 16 - "Validator Validatorerror"
Cohesion: 0.5
Nodes (4): ValidatorError, init(), msgForTag(), ValidateStruct()

### Community 17 - "Greed Rationale"
Cohesion: 0.4
Nodes (5): ThinkArah Project Blueprint, Confluence 3 Judges Logic, Fear and Greed Risk Brake, Rationale: Block Buys in Extreme Greed, Rationale: Require Multi-Algorithm Agreement

### Community 18 - "Theme Inittheme"
Cohesion: 0.67
Nodes (2): initTheme(), updateThemeUI()

### Community 19 - "Navigate Logoutuser"
Cohesion: 0.5
Nodes (0): 

### Community 20 - "Hashpwd Migrateandseed"
Cohesion: 1.0
Nodes (2): hashPwd(), MigrateAndSeed()

### Community 21 - "Formatdateid Timeago"
Cohesion: 1.0
Nodes (2): FormatDateID(), TimeAgo()

### Community 22 - "Table Active"
Cohesion: 1.0
Nodes (3): Active Users Table View, Users Listing API Endpoint, Generic Users Table View

### Community 23 - "Full Suite"
Cohesion: 0.67
Nodes (3): MacroQuant Dashboard Page, MacroQuant Full Suite Single-Page Dashboard, MacroQuant Macro Data Page

### Community 24 - "Toast Injecttoasthtml"
Cohesion: 1.0
Nodes (0): 

### Community 25 - "Core Cluster"
Cohesion: 1.0
Nodes (0): 

### Community 26 - "Assume Role"
Cohesion: 1.0
Nodes (0): 

### Community 27 - "Logout Cluster"
Cohesion: 1.0
Nodes (0): 

## Ambiguous Edges - Review These
- `Account Dashboard Summary View` → `Account Wrapper Navigation Shell`  [AMBIGUOUS]
  users/public/views/account/wrapper/page.html · relation: references
- `Legacy Login Page` → `Account Register Route`  [AMBIGUOUS]
  users/public/views/login1/page.html · relation: references
- `Legacy Login Page` → `Account Reset Route`  [AMBIGUOUS]
  users/public/views/login1/page.html · relation: references

## Knowledge Gaps
- **63 isolated node(s):** `Repository`, `CreateUserInput`, `CreateUserResult`, `UserListItem`, `GetUsersParams` (+58 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **Thin community `Toast Injecttoasthtml`** (2 nodes): `za-toast.js`, `injectToastHTML()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Core Cluster`** (1 nodes): `api.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Assume Role`** (1 nodes): `assume_role.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Logout Cluster`** (1 nodes): `logout.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **What is the exact relationship between `Account Dashboard Summary View` and `Account Wrapper Navigation Shell`?**
  _Edge tagged AMBIGUOUS (relation: references) - confidence is low._
- **What is the exact relationship between `Legacy Login Page` and `Account Register Route`?**
  _Edge tagged AMBIGUOUS (relation: references) - confidence is low._
- **What is the exact relationship between `Legacy Login Page` and `Account Reset Route`?**
  _Edge tagged AMBIGUOUS (relation: references) - confidence is low._
- **Why does `New()` connect `App Pkcs` to `Repository Createuser`, `Controller Submit`, `Dispatcher Shutdown`, `Getcsrftoken Jsonresponse`, `Randomstring Slugify`, `Validator Validatorerror`?**
  _High betweenness centrality (0.223) - this node is a cross-community bridge._
- **Why does `main()` connect `App Pkcs` to `Repository Createuser`, `Controller Submit`, `Service Manager`, `Getcsrftoken Jsonresponse`, `Hashpwd Migrateandseed`?**
  _High betweenness centrality (0.116) - this node is a cross-community bridge._
- **Why does `GetEnv()` connect `App Pkcs` to `Controller Submit`, `Service Manager`?**
  _High betweenness centrality (0.066) - this node is a cross-community bridge._
- **Are the 30 inferred relationships involving `Get()` (e.g. with `.CreateUser()` and `.GetUsers()`) actually correct?**
  _`Get()` has 30 INFERRED edges - model-reasoned connections that need verification._