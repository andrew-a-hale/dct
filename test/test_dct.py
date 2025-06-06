#!/bin/python3
import subprocess
import pytest
import os


class BuildError(Exception):
    def __init__(self):
        super().__init__("failed to build dct")


PEEK_SUPPORTED_FILE_TYPES = ["csv", "json", "ndjson", "parquet"]
PROFILE_SUPPORTED_FILE_TYPES = ["csv", "json", "ndjson", "parquet"]

b = subprocess.run(["go", "build"], capture_output=True)
if b.stderr:
    raise BuildError()

subprocess.run(["chmod", "+x", "./dct"])


def helper_peek_default(filetype: str):
    out = subprocess.run(
        ["./dct", "peek", f"./test/resources/left.{filetype}"],
        capture_output=True,
    )

    assert out.stdout == open("./test/expected/test_peek_default.txt", mode="rb").read()


def helper_peek_5lines(filetype: str):
    out = subprocess.run(
        ["./dct", "peek", f"./test/resources/left.{filetype}", "-n", "5"],
        capture_output=True,
    )

    assert out.stdout == open("./test/expected/test_peek_5lines.txt", mode="rb").read()


def helper_peek_output(filetype: str):
    subprocess.run(
        [
            "./dct",
            "peek",
            f"./test/resources/left.{filetype}",
            "-o",
            "tmp_test_peek_output.csv",
        ],
    )

    assert (
        open("./tmp_test_peek_output.csv", mode="rb").read()
        == open("./test/expected/test_peek_output.csv", mode="rb").read()
    )

    os.remove("./tmp_test_peek_output.csv")


@pytest.mark.parametrize("type", PEEK_SUPPORTED_FILE_TYPES)
def test_peek_default(type: str):
    helper_peek_default(type)


@pytest.mark.parametrize("type", PEEK_SUPPORTED_FILE_TYPES)
def test_peek_5lines(type: str):
    helper_peek_5lines(type)


@pytest.mark.parametrize("type", PEEK_SUPPORTED_FILE_TYPES)
def test_peek_output(type: str):
    helper_peek_output(type)


def test_diff_not_equal():
    out = subprocess.run(
        [
            "./dct",
            "diff",
            "a",
            "./test/resources/left.csv",
            "./test/resources/right.csv",
        ],
        capture_output=True,
    )

    assert (
        out.stdout == open("./test/expected/test_diff_not_equal.txt", mode="rb").read()
    )


def test_diff_empty():
    out = subprocess.run(
        [
            "./dct",
            "diff",
            "a",
            "./test/resources/left.csv",
            "./test/resources/empty.csv",
        ],
        capture_output=True,
    )

    assert str(out.stderr, "utf-8").endswith(
        "attempted to diff when least one of the files have no data\n"
    )


def test_diff_equal():
    out = subprocess.run(
        [
            "./dct",
            "diff",
            "a",
            "./test/resources/left.csv",
            "./test/resources/left.csv",
        ],
        capture_output=True,
    )

    assert out.stdout == open("./test/expected/test_diff_equal.txt", mode="rb").read()


def test_diff_key():
    out = subprocess.run(
        [
            "./dct",
            "diff",
            "a=b",
            "./test/resources/left.csv",
            "./test/resources/right.csv",
        ],
        capture_output=True,
    )

    assert out.stdout == open("./test/expected/test_diff_key.txt", mode="rb").read()


def test_diff_keys():
    out = subprocess.run(
        [
            "./dct",
            "diff",
            "a,b",
            "./test/resources/left.csv",
            "./test/resources/right.csv",
        ],
        capture_output=True,
    )

    assert out.stdout == open("./test/expected/test_diff_keys.txt", mode="rb").read()


# dct diff a test/resources/left.csv test/resources/right.csv -m '[{"agg":"mean","left":"b","right":"b"},{"agg":"count_distinct","left":"c","right":"c"}]'
def test_diff_metric_string():
    out = subprocess.run(
        [
            "./dct",
            "diff",
            "a",
            "./test/resources/left.csv",
            "./test/resources/right.csv",
            "-m",
            """[{"agg":"mean","left":"b","right":"b"},{"agg":"count_distinct","left":"c","right":"c"}]""",
        ],
        capture_output=True,
    )

    assert (
        out.stdout
        == open("./test/expected/test_diff_metric_string.txt", mode="rb").read()
    )


# dct diff a test/resources/left.csv test/resources/right.csv -m test/resources/metrics.json
def test_diff_metric_file():
    out = subprocess.run(
        [
            "./dct",
            "diff",
            "a",
            "./test/resources/left.csv",
            "./test/resources/right.csv",
            "-m",
            "./test/resources/metrics.json",
        ],
        capture_output=True,
    )

    assert (
        out.stdout
        == open("./test/expected/test_diff_metric_file.txt", mode="rb").read()
    )


# dct diff a test/resources/left.csv test/resources/left.csv -m test/resources/metrics.json -a
def test_diff_metric_file_all():
    out = subprocess.run(
        [
            "./dct",
            "diff",
            "a",
            "./test/resources/left.csv",
            "./test/resources/left.csv",
            "-m",
            "./test/resources/metrics.json",
            "-a",
        ],
        capture_output=True,
    )

    assert (
        out.stdout
        == open("./test/expected/test_diff_metric_file_all.txt", mode="rb").read()
    )


# dct diff a test/resources/left.csv test/resources/right.csv -m test/resources/metrics.json
def test_diff_output():
    subprocess.run(
        [
            "./dct",
            "diff",
            "a",
            "./test/resources/left.csv",
            "./test/resources/right.csv",
            "-m",
            "./test/resources/metrics.json",
            "-o",
            "./tmp_test_diff_output.csv",
        ],
    )

    assert (
        open("./tmp_test_diff_output.csv", mode="rb").read()
        == open("./test/expected/test_diff_output.csv", mode="rb").read()
    )

    os.remove("./tmp_test_diff_output.csv")


def test_chart():
    subprocess.run(
        [
            "./dct",
            "chart",
            "-w",
            "50",
            "./test/resources/left.csv",
            "1",
            "count",
        ],
        capture_output=True,
    )

    # skip
    assert True


def test_version():
    out = subprocess.run(
        ["./dct", "version"],
        capture_output=True,
    )

    assert out.stderr == b""


def test_generator():
    out = subprocess.run(
        [
            "./dct",
            "gen",
            "test/resources/generator-schema.json",
        ],
        capture_output=True
    )

    # can't test the random data? need to implement rng seed
    header = out.stdout.decode().splitlines()[0]
    expected_header = open("test/expected/test_generator.csv", mode="r").read().strip()
    assert header == expected_header
    assert out.stderr == b""


def test_flattify_ndjson():
    out = subprocess.run(
        [
            "./dct",
            "flattify",
            "test/resources/flattify.ndjson",
        ],
        capture_output=True,
    )

    # map isn't sorted
    assert out.stderr == b""

def test_flattify_json_file():
    out = subprocess.run(
        [
            "./dct",
            "flattify",
            "test/resources/flattify.json",
        ],
        capture_output=True,
    )

    # map isn't sorted
    assert out.stderr == b""
    
def test_flattify_json():
    out = subprocess.run(
        [
            "./dct",
            "flattify",
            """{
                "a": 1,
                "b": {
                    "a": 1
                }
            }"""
        ],
        capture_output=True,
    )

    # map isn't sorted
    assert out.stderr == b""

def test_flattify_json_array():
    out = subprocess.run(
        [
            "./dct",
            "flattify",
            """[ 1, 2 ]"""
        ],
        capture_output=True,
    )

    # map isn't sorted
    assert out.stderr == b""

def test_flattify_json_array_sql():
    out = subprocess.run(
        [
            "./dct",
            "flattify",
            "-s",
            """[ 1, 2 ]"""
        ],
        capture_output=True,
    )

    # map isn't sorted
    assert out.stderr == b""
    assert out.stdout == open("./test/expected/test_flattify_json_array.sql", mode="rb").read()

def test_flattify_json_object_digit_key():
    out = subprocess.run(
        [
            "./dct",
            "flattify",
            """{"0": [ 1, 2 ]}"""
        ],
        capture_output=True,
    )

    # map isn't sorted
    assert out.stderr == b""

def test_flattify_json_object_digit_key_sql():
    out = subprocess.run(
        [
            "./dct",
            "flattify",
            "-s",
            """{"0": [ 1, 2 ]}"""
        ],
        capture_output=True,
    )

    # map isn't sorted
    assert out.stderr == b""
    assert out.stdout == open("./test/expected/test_flattify_json_object_digit_key.sql", mode="rb").read()


def test_flattify_ndjson_select():
    out = subprocess.run(
        [
            "./dct",
            "flattify",
            "-s",
            "test/resources/flattify.ndjson",
        ],
        capture_output=True,
    )

    assert out.stderr == b""
    assert out.stdout == open("./test/expected/test_flattify_ndjson_select.sql", mode="rb").read()

def test_flattify_json_select():
    out = subprocess.run(
        [
            "./dct",
            "flattify",
            "-s",
            "test/resources/flattify.json",
        ],
        capture_output=True,
    )

    assert out.stderr == b""
    assert out.stdout == open("./test/expected/test_flattify_json_select.sql", mode="rb").read()


@pytest.mark.parametrize("filetype", PROFILE_SUPPORTED_FILE_TYPES)
def test_prof(filetype: str):
    out = subprocess.run(
        ["./dct", "prof", f"./test/resources/left.{filetype}"],
        capture_output=True,
    )

    # output is not ordered!
    assert out.stderr == b""
    assert out.stdout != b""