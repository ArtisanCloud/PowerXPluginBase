#!/usr/bin/env node

import { promises as fs } from "node:fs";
import path from "node:path";
import process from "node:process";
import { execFile } from "node:child_process";
import { promisify } from "node:util";

import Ajv from "ajv";
import yaml from "js-yaml";
import { diffString } from "json-diff";

const execFileAsync = promisify(execFile);

const DEFAULT_CONFIG = "contracts/compatibility.yaml";
const DEFAULT_OUTPUT = "build/compat";

function parseArgs(argv) {
  const args = { config: DEFAULT_CONFIG, outDir: DEFAULT_OUTPUT };
  for (let i = 2; i < argv.length; i += 1) {
    const arg = argv[i];
    switch (arg) {
      case "--config":
        args.config = argv[++i];
        break;
      case "--out":
      case "--output":
        args.outDir = argv[++i];
        break;
      case "--help":
      case "-h":
        printHelp();
        process.exit(0);
        break;
      default:
        console.warn(`[compat] Unknown argument: ${arg}`);
    }
  }
  return args;
}

function printHelp() {
  console.log(`Usage: node scripts/check-compatibility.mjs [--config <path>] [--out <dir>]

Validates capability/schema compatibility using the instructions found in contracts/compatibility.yaml.

Options:
  --config <path>   Compatibility config file (default: ${DEFAULT_CONFIG})
  --out <dir>       Directory for generated reports (default: ${DEFAULT_OUTPUT})
`);
}

async function readConfig(configPath) {
  const raw = await fs.readFile(configPath, "utf8");
  const doc = yaml.load(raw);
  if (!doc || !Array.isArray(doc.items)) {
    throw new Error(
      `[compat] Invalid config structure. Expected 'items' array in ${configPath}`
    );
  }
  return doc.items;
}

async function ensureDir(dir) {
  await fs.mkdir(dir, { recursive: true });
}

async function loadJSONLike(filePath) {
  const content = await fs.readFile(filePath, "utf8");
  if (filePath.endsWith(".yaml") || filePath.endsWith(".yml")) {
    return yaml.load(content);
  }
  return JSON.parse(content);
}

async function validateSchema(schemaPath) {
  try {
    const schema = JSON.parse(await fs.readFile(schemaPath, "utf8"));
    const ajv = new Ajv({ strict: false });
    ajv.compile(schema);
    return { ok: true };
  } catch (err) {
    return { ok: false, error: err };
  }
}

async function generateStructuredDiff(baselinePath, currentPath) {
  const [baseline, current] = await Promise.all([
    loadJSONLike(baselinePath),
    loadJSONLike(currentPath),
  ]);
  const diff = diffString(baseline, current) ?? "No differences";
  return diff;
}

async function runOpenAPIDiff(root, baseline, current, safeId, outDir) {
  const binName = process.platform === "win32" ? "openapi-diff.cmd" : "openapi-diff";
  const binPath = path.join(root, "scripts", "node_modules", ".bin", binName);
  if (!(await exists(binPath))) {
    throw new Error(
      `openapi-diff binary not found at ${binPath}. Run npm install --prefix scripts first.`
    );
  }
  const { stdout } = await execFileAsync(binPath, ["--json", baseline, current], {
    cwd: root,
    maxBuffer: 10 * 1024 * 1024,
  });
  const reportPath = path.join(outDir, `${safeId}.openapi.diff.json`);
  await fs.writeFile(reportPath, stdout, "utf8");
  return reportPath;
}

function resolvePath(root, p) {
  if (!p) return null;
  return path.resolve(root, p);
}

async function main() {
  try {
    const { config, outDir } = parseArgs(process.argv);
    const repoRoot = process.cwd();
    const resolvedConfig = resolvePath(repoRoot, config);
    const resolvedOut = resolvePath(repoRoot, outDir);

    await ensureDir(resolvedOut);
    const items = await readConfig(resolvedConfig);

    const summary = [];

    for (const item of items) {
      const kind = item.kind || "capability";
      const id = item.id || item.path || "unknown";
      const current = resolvePath(repoRoot, item.current);
      const baseline = resolvePath(repoRoot, item.baseline);
      const changeType = item.change_type || "unspecified";

      const entry = {
        kind,
        id,
        changeType,
        status: "ok",
        notes: [],
      };

      if (!current) {
        entry.status = "error";
        entry.notes.push("Missing 'current' path in config item");
        summary.push(entry);
        continue;
      }
      if (!(await exists(current))) {
        entry.status = "error";
        entry.notes.push(`Current file not found: ${current}`);
        summary.push(entry);
        continue;
      }

      if (kind === "schema") {
        const result = await validateSchema(current);
        if (!result.ok) {
          entry.status = "error";
          entry.notes.push(`Invalid JSON Schema: ${result.error.message}`);
        } else {
          entry.notes.push("JSON Schema compiled successfully");
        }
      }

      if (!baseline || !(await exists(baseline))) {
        entry.status = entry.status === "ok" ? "warn" : entry.status;
        entry.notes.push(
          baseline
            ? `Baseline file missing: ${baseline}`
            : "No baseline specified; skipping diff"
        );
        summary.push(entry);
        continue;
      }

      try {
        const safeId = id.replace(/[^a-zA-Z0-9._-]/g, "_");
        if (kind === "openapi") {
          const reportPath = await runOpenAPIDiff(
            repoRoot,
            baseline,
            current,
            safeId,
            resolvedOut
          );
          entry.notes.push(`Diff written to ${path.relative(repoRoot, reportPath)}`);
        } else {
          const diff = await generateStructuredDiff(baseline, current);
          const reportPath = path.join(
            resolvedOut,
            `${safeId}.${kind}.diff.txt`
          );
          await fs.writeFile(reportPath, diff, "utf8");
          entry.notes.push(`Diff written to ${path.relative(repoRoot, reportPath)}`);
        }
      } catch (err) {
        entry.status = "error";
        entry.notes.push(`Diff failed: ${err.message}`);
      }

      summary.push(entry);
    }

    const summaryPath = path.join(resolvedOut, "report.json");
    await fs.writeFile(summaryPath, JSON.stringify(summary, null, 2), "utf8");
    console.log(
      `[compat] Completed compatibility analysis. Summary: ${path.relative(
        repoRoot,
        summaryPath
      )}`
    );
  } catch (err) {
    console.error(`[compat] ${err.message}`);
    process.exit(1);
  }
}

async function exists(p) {
  if (!p) return false;
  try {
    await fs.access(p);
    return true;
  } catch {
    return false;
  }
}

await main();
