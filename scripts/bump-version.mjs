import fs from 'fs';
import path from 'path';

const version = process.argv[2];
if (!version) {
  console.error('Usage: node scripts/bump-version.mjs <version>');
  process.exit(1);
}

const cleanVersion = version.replace(/^v/, '');
console.log(`Bumping version to ${cleanVersion}...`);

// 1. Update internal/platform/version/version.go
const goVersionPath = path.resolve('internal/platform/version/version.go');
if (fs.existsSync(goVersionPath)) {
  let goContent = fs.readFileSync(goVersionPath, 'utf8');
  goContent = goContent.replace(/const Current = "[^"]+"/, `const Current = "${cleanVersion}"`);
  fs.writeFileSync(goVersionPath, goContent, 'utf8');
  console.log(`Updated ${goVersionPath}`);
}

// 2. Update wails.json
const wailsPath = path.resolve('wails.json');
if (fs.existsSync(wailsPath)) {
  const wailsConfig = JSON.parse(fs.readFileSync(wailsPath, 'utf8'));
  if (wailsConfig.info) {
    wailsConfig.info.productVersion = cleanVersion;
  }
  fs.writeFileSync(wailsPath, JSON.stringify(wailsConfig, null, 2) + '\n', 'utf8');
  console.log(`Updated ${wailsPath}`);
}

// 3. Update frontend/package.json
const pkgPath = path.resolve('frontend/package.json');
if (fs.existsSync(pkgPath)) {
  const pkg = JSON.parse(fs.readFileSync(pkgPath, 'utf8'));
  pkg.version = cleanVersion;
  fs.writeFileSync(pkgPath, JSON.stringify(pkg, null, 2) + '\n', 'utf8');
  console.log(`Updated ${pkgPath}`);
}

console.log('Version bump complete.');
