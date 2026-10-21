import { execSync } from 'child_process';
import fs from 'fs';
import path from 'path';

function getArgs() {
  const args = process.argv.slice(2);
  let bump = 'auto';
  let json = false;
  for (let i = 0; i < args.length; i++) {
    if (args[i] === '--bump' && args[i + 1]) {
      bump = args[++i];
    } else if (args[i] === '--json') {
      json = true;
    }
  }
  return { bump, json };
}

function getLatestTag() {
  try {
    const stdout = execSync('git describe --tags --abbrev=0', { encoding: 'utf8' }).trim();
    return stdout;
  } catch {
    return 'v1.0.0-beta';
  }
}

function getCommitMessagesSince(tag) {
  try {
    const stdout = execSync(`git log ${tag}..HEAD --oneline`, { encoding: 'utf8' }).trim();
    return stdout ? stdout.split('\n') : [];
  } catch {
    try {
      const stdout = execSync('git log -n 10 --oneline', { encoding: 'utf8' }).trim();
      return stdout.split('\n');
    } catch {
      return [];
    }
  }
}

function parseSemver(tag) {
  const clean = tag.replace(/^v/, '');
  const [base, prerelease] = clean.split('-');
  const parts = base.split('.').map((p) => parseInt(p, 10) || 0);
  while (parts.length < 3) parts.push(0);
  return {
    major: parts[0],
    minor: parts[1],
    patch: parts[2],
    prerelease: prerelease || ''
  };
}

function determineBump(commits, requestedBump) {
  if (['major', 'minor', 'patch'].includes(requestedBump)) {
    return requestedBump;
  }

  let hasBreaking = false;
  let hasFeat = false;
  let hasFix = false;

  for (const c of commits) {
    const msg = c.toLowerCase();
    if (msg.includes('breaking change') || msg.includes('breaking:')) {
      hasBreaking = true;
    } else if (msg.includes('feat(') || msg.includes('feat:')) {
      hasFeat = true;
    } else if (msg.includes('fix(') || msg.includes('fix:')) {
      hasFix = true;
    }
  }

  if (hasBreaking) return 'major';
  if (hasFeat) return 'minor';
  if (hasFix) return 'patch';
  return 'patch';
}

function applyBump(sem, type) {
  if (type === 'major') {
    return `${sem.major + 1}.0.0`;
  }
  if (type === 'minor') {
    return `${sem.major}.${sem.minor + 1}.0`;
  }
  return `${sem.major}.${sem.minor}.${sem.patch + 1}`;
}

const { bump: requestedBump, json: isJson } = getArgs();
const currentRef = process.env.GITHUB_REF || '';

let result;

if (currentRef.startsWith('refs/tags/v')) {
  const tag = currentRef.replace('refs/tags/', '');
  const version = tag.replace(/^v/, '');
  result = {
    release: true,
    bump: 'tag',
    previous: version,
    previousTag: tag,
    version: version,
    tag: tag,
    reason: `Release triggered directly by tag ${tag}`
  };
} else {
  const prevTag = getLatestTag();
  const sem = parseSemver(prevTag);
  const commits = getCommitMessagesSince(prevTag);
  const bumpType = determineBump(commits, requestedBump);
  const nextVer = applyBump(sem, bumpType);
  const nextTag = `v${nextVer}`;

  const shouldRelease = commits.length > 0 || requestedBump !== 'auto';

  result = {
    release: shouldRelease,
    bump: bumpType,
    previous: prevTag.replace(/^v/, ''),
    previousTag: prevTag,
    version: nextVer,
    tag: nextTag,
    reason: shouldRelease
      ? `SemVer bump (${bumpType}) based on ${commits.length} commit(s)`
      : 'No release-triggering commits found since last tag'
  };
}

if (isJson) {
  console.log(JSON.stringify(result, null, 2));
} else {
  console.log(`Version decision: ${result.tag} (${result.reason})`);
}
