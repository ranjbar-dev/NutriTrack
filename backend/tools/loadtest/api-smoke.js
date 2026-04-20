const baseUrl = process.env.NUTRITRACK_BASE_URL || 'http://localhost:8080'
const requestsPerScenario = Number(process.env.NUTRITRACK_REQUESTS || 20)

const scenarios = [
  { name: 'health', path: '/api/health' },
  { name: 'auth-me', path: '/api/auth/me' },
]

async function timeRequest(url, cookie) {
  const started = performance.now()
  const response = await fetch(url, {
    headers: cookie ? { Cookie: cookie } : {},
  })
  return {
    ok: response.ok,
    status: response.status,
    durationMs: Number((performance.now() - started).toFixed(2)),
  }
}

function percentile(values, p) {
  if (!values.length) return 0
  const sorted = [...values].sort((a, b) => a - b)
  const index = Math.min(sorted.length - 1, Math.ceil((p / 100) * sorted.length) - 1)
  return sorted[index]
}

async function run() {
  const cookie = process.env.NUTRITRACK_COOKIE || ''
  const results = []

  for (const scenario of scenarios) {
    for (let i = 0; i < requestsPerScenario; i += 1) {
      results.push({
        scenario: scenario.name,
        ...(await timeRequest(`${baseUrl}${scenario.path}`, cookie)),
      })
    }
  }

  const durations = results.map(result => result.durationMs)
  const p95 = percentile(durations, 95)
  const failed = results.filter(result => !result.ok)

  console.log(JSON.stringify({
    totalRequests: results.length,
    p95,
    avg: Number((durations.reduce((sum, value) => sum + value, 0) / durations.length).toFixed(2)),
    failed: failed.length,
    targetP95Ms: 200,
    pass: p95 <= 200 && failed.length === 0,
    results,
  }, null, 2))
}

run().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
