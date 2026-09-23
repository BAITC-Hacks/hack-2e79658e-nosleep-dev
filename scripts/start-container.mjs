import { spawn } from "node:child_process";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const appRoot = dirname(fileURLToPath(import.meta.url));

const processes = [
  spawn(join(appRoot, "api"), { stdio: "inherit" }),
  spawn(process.execPath, [join(appRoot, "apps/web/server.js")], { stdio: "inherit" }),
];

let stopping = false;
let exitCode = 0;

function stop(signal = "SIGTERM") {
  if (stopping) return;
  stopping = true;
  for (const process of processes) {
    if (process.exitCode === null && process.signalCode === null) process.kill(signal);
  }
  setTimeout(() => {
    for (const process of processes) {
      if (process.exitCode === null && process.signalCode === null) process.kill("SIGKILL");
    }
  }, 12_000).unref();
}

for (const signal of ["SIGINT", "SIGTERM"]) {
  process.on(signal, () => stop(signal));
}

for (const child of processes) {
  child.on("error", (error) => {
    console.error(`Could not start a service: ${error.message}`);
    exitCode = 1;
    stop();
  });
  child.on("exit", (code) => {
    if (!stopping) {
      console.error(`A service exited unexpectedly with code ${code ?? "unknown"}.`);
      exitCode = code || 1;
      stop();
    }
    if (processes.every((process) => process.exitCode !== null || process.signalCode !== null)) {
      process.exitCode = exitCode;
    }
  });
}
