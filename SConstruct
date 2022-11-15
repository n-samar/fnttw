from pathlib import Path
import logging
import os
import subprocess
import sys

clang_version = "clang++-10"
cpp_version = "c++20"

build_dir = Path("build")
src_dir = Path("src")
test_dir = Path("test")
benchmark_dir = Path("benchmark")
target_dir = Path("target")

debug = int(ARGUMENTS.get("debug", 0))
clang = str(Path("/usr/bin/") / clang_version)

env = Environment(CXX=clang, ENV=os.environ)
env.Tool("compilation_db")
env.CompilationDatabase()

logging.basicConfig(
    format="{levelname}[{lineno}]: {message}",
    style="{",
    level=logging.DEBUG,
)

env.Append(CPPPATH=src_dir)
env.VariantDir(build_dir, src_dir, duplicate=0)


env.Append(
    CPPFLAGS=[
        "-Werror",
        "-fno-omit-frame-pointer",
        "-g",
        "-ggdb",
        f"-std={cpp_version}",
        "-Wall",
    ]
)

if debug == 1:
    env.Append(CPPFLAGS=["-O0", "-fsanitize=address"])
    # nsamar: Linker incantations required for LIBASAN to work
    # https://github.com/google/sanitizers/issues/856
    env.Append(LINKFLAGS=["-fsanitize=address", "-no-pie", "-shared-libasan"])
elif debug == 0:
    env.Append(CPPFLAGS="-O3")

env.Append(
    LIBS=[
        "stdc++fs",
        "pthread",
        "glog",
        "gtest_main",
        "gtest",
        "gflags",
        "pfm",
        File("/data/sanchez/tools/nsamar/18.04/janncy/x86_64-linux-gnu/libbenchmark.a"),
    ]
)

env.Append(LINKFLAGS = ["-fopenmp"])
env.Append(CCFLAGS = ["-fopenmp"])

def program(name, source, extra_sources = []):
    target = str(build_dir / name)

    if isinstance(source, Path):
        source = str(source)

    sources = [source, Glob(src_dir / "*.cc") + extra_sources]
    return env.Program(target, sources)

program("tests", Glob(test_dir / "*.cc"))
program("benchmark", Glob(benchmark_dir / "*.cc"))
program("main", Glob(target_dir / "main.cc"))
