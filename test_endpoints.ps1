# test_endpoints.ps1

$BASE_URL = "http://localhost:8080"

Write-Host "1. Obteniendo token JWT..." -ForegroundColor Cyan
$tokenResp = Invoke-RestMethod -Method POST -Uri "$BASE_URL/auth/test-token"
$TOKEN = $tokenResp.token
Write-Host "Token: $TOKEN"
Write-Host ""

$Headers = @{ Authorization = "Bearer $TOKEN" }

Write-Host "2. GET /tasks (sin filtros)" -ForegroundColor Cyan
Invoke-RestMethod -Method GET -Uri "$BASE_URL/tasks" -Headers $Headers | ConvertTo-Json -Depth 8
Write-Host ""

Write-Host "3. GET /tasks?priority=high" -ForegroundColor Cyan
Invoke-RestMethod -Method GET -Uri "$BASE_URL/tasks?priority=high" -Headers $Headers | ConvertTo-Json -Depth 8
Write-Host ""

Write-Host "4. GET /tasks?status=todo,in-progress" -ForegroundColor Cyan
Invoke-RestMethod -Method GET -Uri "$BASE_URL/tasks?status=todo,in-progress" -Headers $Headers | ConvertTo-Json -Depth 8
Write-Host ""

Write-Host "5. GET /tasks/search?q=proyecto" -ForegroundColor Cyan
Invoke-RestMethod -Method GET -Uri "$BASE_URL/tasks/search?q=proyecto" -Headers $Headers | ConvertTo-Json -Depth 8
Write-Host ""

Write-Host "6. GET /tasks/overdue" -ForegroundColor Cyan
Invoke-RestMethod -Method GET -Uri "$BASE_URL/tasks/overdue" -Headers $Headers | ConvertTo-Json -Depth 8
Write-Host ""

Write-Host "7. GET /tasks/completed" -ForegroundColor Cyan
Invoke-RestMethod -Method GET -Uri "$BASE_URL/tasks/completed" -Headers $Headers | ConvertTo-Json -Depth 8
Write-Host ""

Write-Host "8. GET /tasks/by-subject/subject-ic-6821" -ForegroundColor Cyan
Invoke-RestMethod -Method GET -Uri "$BASE_URL/tasks/by-subject/subject-ic-6821" -Headers $Headers | ConvertTo-Json -Depth 8
Write-Host ""

Write-Host "9. GET /tasks/by-period/period-2025-1" -ForegroundColor Cyan
Invoke-RestMethod -Method GET -Uri "$BASE_URL/tasks/by-period/period-2025-1" -Headers $Headers | ConvertTo-Json -Depth 8
Write-Host ""

Write-Host "✅ Tests completados" -ForegroundColor Green
