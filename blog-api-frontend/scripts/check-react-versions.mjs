import { createRequire } from 'node:module';
import console from 'node:console';
import process from 'node:process';

const require = createRequire(import.meta.url);
const versions = {};

for (const name of ['react', 'react-dom']) {
  try {
    versions[name] = require(`${name}/package.json`).version;
  } catch {
    versions[name] = 'not installed';
  }
}

if (versions.react === 'not installed' || versions['react-dom'] === 'not installed' ||
    versions.react !== versions['react-dom']) {
  console.error(`React version check failed: react=${versions.react}, react-dom=${versions['react-dom']}.`);
  console.error('Install react and react-dom at the exact same version, update package-lock.json, then run npm ci.');
  process.exitCode = 1;
} else {
  console.log(`React versions match: ${versions.react}`);
}
