// The consumer fixture installs without a lockfile, so mise, the root package,
// and the fixture must name the same pnpm before `pnpm install` runs.
import { readFileSync } from "node:fs";

const packageManager = (path) => JSON.parse(readFileSync(path, "utf8")).packageManager;
const expected = `pnpm@${process.argv[2]}`;
const root = packageManager(new URL("../../../package.json", import.meta.url));
const fixture = packageManager(new URL("./package.json", import.meta.url));

if (root !== expected) {
  console.error(`package.json packageManager ${root} differs from mise ${expected}`);
  process.exit(1);
}
if (fixture !== root) {
  console.error(`TypeScript consumer packageManager ${fixture} differs from package.json ${root}`);
  process.exit(1);
}
