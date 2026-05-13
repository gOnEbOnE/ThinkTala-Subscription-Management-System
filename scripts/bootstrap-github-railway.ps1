[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [string]$Owner,

    [Parameter(Mandatory = $true)]
    [string]$Repo,

    [ValidateSet("public", "private")]
    [string]$Visibility = "public",

    [string]$Pat,

    [string]$PrimaryBranch = "main",

    [switch]$PushCurrentBranch = $true
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function Write-Step {
    param([string]$Message)
    Write-Host "[bootstrap] $Message"
}

function New-GitHubHeaders {
    param([string]$Token)
    $headers = @{
        "Accept" = "application/vnd.github+json"
        "X-GitHub-Api-Version" = "2022-11-28"
    }

    if ($Token) {
        $headers["Authorization"] = "Bearer $Token"
    }

    return $headers
}

function Test-GitHubRepoExists {
    param(
        [string]$RepoOwner,
        [string]$RepoName,
        [hashtable]$Headers
    )

    $uri = "https://api.github.com/repos/$RepoOwner/$RepoName"
    try {
        Invoke-RestMethod -Method Get -Uri $uri -Headers $Headers | Out-Null
        return $true
    }
    catch {
        return $false
    }
}

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
Set-Location $repoRoot

Write-Step "Working directory: $repoRoot"

if (-not (Test-Path ".git")) {
    throw "Current directory is not a git repository: $repoRoot"
}

$headers = New-GitHubHeaders -Token $Pat
$repoExists = Test-GitHubRepoExists -RepoOwner $Owner -RepoName $Repo -Headers $headers

if (-not $repoExists) {
    if (-not $Pat) {
        throw "GitHub repo $Owner/$Repo does not exist and no PAT was provided. Create it first at https://github.com/new, then rerun this script."
    }

    Write-Step "Repository not found. Creating GitHub repository $Owner/$Repo ($Visibility)."

    $isPrivate = $Visibility -eq "private"
    $body = @{
        name = $Repo
        private = $isPrivate
        auto_init = $false
    } | ConvertTo-Json

    $user = Invoke-RestMethod -Method Get -Uri "https://api.github.com/user" -Headers $headers

    if ($user.login -eq $Owner) {
        Invoke-RestMethod -Method Post -Uri "https://api.github.com/user/repos" -Headers $headers -Body $body -ContentType "application/json" | Out-Null
    }
    else {
        Invoke-RestMethod -Method Post -Uri "https://api.github.com/orgs/$Owner/repos" -Headers $headers -Body $body -ContentType "application/json" | Out-Null
    }

    Write-Step "Repository created."
}
else {
    Write-Step "Repository already exists: $Owner/$Repo"
}

$remoteUrl = "https://github.com/$Owner/$Repo.git"
$existingRemotes = git remote

if ($existingRemotes -contains "github") {
    Write-Step "Updating existing remote 'github' => $remoteUrl"
    git remote set-url github $remoteUrl
}
else {
    Write-Step "Adding remote 'github' => $remoteUrl"
    git remote add github $remoteUrl
}

$currentBranch = (git branch --show-current).Trim()
if (-not $currentBranch) {
    throw "Unable to detect current branch."
}

$localBranchNames = @(git branch --format="%(refname:short)")
if (-not ($localBranchNames -contains $PrimaryBranch)) {
    Write-Step "Primary branch '$PrimaryBranch' not found locally. Creating from current HEAD."
    git branch $PrimaryBranch
}

Write-Step "Pushing primary branch '$PrimaryBranch' to remote 'github'."
git push -u github $PrimaryBranch

if ($PushCurrentBranch -and $currentBranch -ne $PrimaryBranch) {
    Write-Step "Pushing current branch '$currentBranch' to remote 'github'."
    git push -u github $currentBranch
}

Write-Step "Done. Next required setup in GitHub Actions secrets/variables:"
Write-Host "  - Secret: RAILWAY_TOKEN"
Write-Host "  - Secret: RAILWAY_PROJECT_ID"
Write-Host "  - Variable: RAILWAY_ENVIRONMENT=production"
