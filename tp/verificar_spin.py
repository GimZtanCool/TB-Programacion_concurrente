"""Exhaustive safety/liveness checks and negative controls; requires Spin and GCC."""
from pathlib import Path
import csv
import hashlib
import re
import shutil
import subprocess

ROOT = Path(__file__).resolve().parent
EVIDENCE = ROOT / "evidencias"
BUILD = ROOT / ".build" / "spin"


def command(args, cwd, log):
    result = subprocess.run(args, cwd=cwd, text=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
    log.append("$ " + " ".join(args) + "\n" + result.stdout + f"\nexit_code={result.returncode}\n")
    if result.returncode != 0:
        raise RuntimeError(log[-1])
    return result.stdout


def main():
    EVIDENCE.mkdir(exist_ok=True)
    versions = []
    for args in (["spin", "-V"], ["gcc", "--version"], ["python3", "--version"], ["uname", "-a"]):
        command(args, ROOT, versions)
    versions.append("modelo_tp.pml SHA256=" + hashlib.sha256((ROOT / "modelo_tp.pml").read_bytes()).hexdigest())
    (EVIDENCE / "spin_entorno.txt").write_text("\n".join(versions), encoding="utf-8")
    rows = []
    cases = [(f"w{w}_{kind}", w, kind, None) for w in (1, 2, 3) for kind in ("safety", "liveness")]
    cases += [("control_escritura", 2, "safety", "BAD_WRITER"), ("control_espera", 2, "safety", "BAD_WAIT")]
    for name, workers, kind, mutation in cases:
        folder = BUILD / name
        folder.mkdir(parents=True, exist_ok=True)
        shutil.copyfile(ROOT / "modelo_tp.pml", folder / "modelo_tp.pml")
        log = []
        gen = ["spin", f"-DWORKERS={workers}"]
        if mutation:
            gen += [f"-D{mutation}"]
        gen += ["-a", "modelo_tp.pml"]
        try:
            command(gen, folder, log)
            compile_args = ["gcc", "-O2"]
            if kind == "safety":
                compile_args += ["-DSAFETY", "-DNOCLAIM"]
            compile_args += ["-o", "pan", "pan.c"]
            command(compile_args, folder, log)
            check = ["./pan", "-m100000", "-w22"]
            if kind == "liveness":
                check += ["-a"]
            output = command(check, folder, log)
            match = re.search(r"errors:\s*(\d+)", output)
            if not match:
                raise RuntimeError("Missing verification result")
            errors = int(match.group(1))
            incomplete = "Search not completed" in output or "max search depth too small" in output
            if mutation:
                if errors == 0:
                    raise RuntimeError("Negative control did not detect an error")
                trail = folder / "modelo_tp.pml.trail"
                if trail.exists():
                    shutil.copyfile(trail, EVIDENCE / f"{name}.trail")
                    replay = ["spin", f"-DWORKERS={workers}", f"-D{mutation}", "-t", "-p", "-g", "modelo_tp.pml"]
                    command(replay, folder, log)
            elif errors or incomplete:
                raise RuntimeError("Verification failed or was incomplete")
            states = re.search(r"([\d.e+]+) states, stored", output)
            depth = re.search(r"depth reached (\d+)", output)
            rows.append([name, workers, kind, errors, states.group(1) if states else "", depth.group(1) if depth else "", "counterexample_expected" if mutation else "complete"])
            print(name, "errors=", errors, "states=", rows[-1][4], flush=True)
        finally:
            (EVIDENCE / f"spin_{name}.txt").write_text("\n".join(log), encoding="utf-8")
    with (EVIDENCE / "spin_resumen.csv").open("w", encoding="utf-8", newline="") as f:
        writer = csv.writer(f)
        writer.writerow(["case", "workers", "check", "errors", "states_stored", "depth", "status"])
        writer.writerows(rows)


if __name__ == "__main__":
    main()
