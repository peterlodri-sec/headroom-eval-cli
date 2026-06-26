"""Fix the evals import in app.py — robust version."""
c = open('app.py').read()

# Check if broken
if 'from evals.headroom_runner import' in c:
    c = c.replace(
        'from evals.headroom_runner import execute_swe_trajectory, TrajectoryMetrics, SWE_TASKS',
        'from headroom_runner import execute_swe_trajectory, TrajectoryMetrics, SWE_TASKS'
    )
    open('app.py', 'w').write(c)
    print('✅ import fixed: evals → headroom_runner')
elif 'from headroom_runner import' in c:
    print('✓ import already fixed')
else:
    print('⚠️ import line not found — check app.py')
