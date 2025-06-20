import subprocess
import json
import select


class BuildError(Exception):
    def __init__(self):
        super().__init__("failed to build mcp-server")


b = subprocess.run(["go", "build", "-C", "mcp-server"], capture_output=True)
if b.stderr:
    raise BuildError()

subprocess.run(["chmod", "+x", "./mcp-server/mcp-server"])


def send_message(proc, message):
    """Send a JSON-RPC message to the server"""
    data = json.dumps(message) + "\n"
    proc.stdin.write(data.encode())
    proc.stdin.flush()


def read_response(proc, timeout=5):
    """Read a JSON-RPC response from the server"""
    # Use select to timeout if no data is available
    if select.select([proc.stdout], [], [], timeout)[0]:
        line = proc.stdout.readline().decode().strip()
        if line:
            return json.loads(line)
    return None


def setup_server():
    """Start the MCP server and initialize it"""
    cmd = ["./mcp-server/mcp-server", "./dct"]
    proc = subprocess.Popen(
        cmd, stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE
    )

    # Initialize
    init_msg = {
        "jsonrpc": "2.0",
        "id": 1,
        "method": "initialize",
        "params": {
            "protocolVersion": "2024-11-05",
            "capabilities": {},
            "clientInfo": {"name": "test-client", "version": "1.0.0"},
        },
    }
    send_message(proc, init_msg)
    response = read_response(proc)
    if response:
        print("Initialize response:", json.dumps(response, indent=2))
    else:
        print("No response to initialize - server may not be responding")

    return proc


def test_tools_list():
    """Test listing available tools"""
    with setup_server() as proc:
        try:
            list_msg = {"jsonrpc": "2.0", "id": 2, "method": "tools/list", "params": {}}
            send_message(proc, list_msg)
            response = read_response(proc)
            if response:
                print("Tools list response:", json.dumps(response, indent=2))
                assert "tools" in response.get("result", {})
                tools = response["result"]["tools"]
                tool_names = [tool["name"] for tool in tools]
                expected_tools = [
                    "data_peek",
                    "data_infer",
                    "data_diff",
                    "data_chart",
                    "data_generate",
                    "data_flattify",
                    "data_js2sql",
                    "data_profile",
                ]
                for tool_name in expected_tools:
                    assert tool_name in tool_names, (
                        f"Tool {tool_name} not found in tools list"
                    )
            else:
                print("No response to tools/list")
                assert False, "No response to tools/list"

        except Exception as e:
            print(f"Error: {e}")
            if proc.stderr:
                stderr = proc.stderr.read().decode()
                if stderr:
                    print(f"Server stderr: {stderr}")
            raise


def test_data_peek():
    """Test data_peek tool"""
    with setup_server() as proc:
        try:
            peek_msg = {
                "jsonrpc": "2.0",
                "id": 3,
                "method": "tools/call",
                "params": {
                    "name": "data_peek",
                    "arguments": {"file_path": "examples/left.csv", "lines": 3},
                },
            }
            send_message(proc, peek_msg)
            response = read_response(proc)
            if response:
                print("data_peek response:", json.dumps(response, indent=2))
                assert "result" in response
            else:
                print("No response to data_peek")
                assert False, "No response to data_peek"

        except Exception as e:
            print(f"Error: {e}")
            if proc.stderr:
                stderr = proc.stderr.read().decode()
                if stderr:
                    print(f"Server stderr: {stderr}")
            raise


def test_data_infer():
    """Test data_infer tool"""
    with setup_server() as proc:
        try:
            infer_msg = {
                "jsonrpc": "2.0",
                "id": 4,
                "method": "tools/call",
                "params": {
                    "name": "data_infer",
                    "arguments": {
                        "file_path": "examples/left.csv",
                        "table": "test_table",
                    },
                },
            }
            send_message(proc, infer_msg)
            response = read_response(proc)
            if response:
                print("data_infer response:", json.dumps(response, indent=2))
                assert "result" in response
            else:
                print("No response to data_infer")
                assert False, "No response to data_infer"

        except Exception as e:
            print(f"Error: {e}")
            if proc.stderr:
                stderr = proc.stderr.read().decode()
                if stderr:
                    print(f"Server stderr: {stderr}")
            raise


def test_data_diff():
    """Test data_diff tool"""
    with setup_server() as proc:
        try:
            diff_msg = {
                "jsonrpc": "2.0",
                "id": 5,
                "method": "tools/call",
                "params": {
                    "name": "data_diff",
                    "arguments": {
                        "keys": "id",
                        "file1": "examples/left.csv",
                        "file2": "examples/right.csv",
                    },
                },
            }
            send_message(proc, diff_msg)
            response = read_response(proc)
            if response:
                print("data_diff response:", json.dumps(response, indent=2))
                assert "result" in response
            else:
                print("No response to data_diff")
                assert False, "No response to data_diff"

        except Exception as e:
            print(f"Error: {e}")
            if proc.stderr:
                stderr = proc.stderr.read().decode()
                if stderr:
                    print(f"Server stderr: {stderr}")
            raise


def test_data_chart():
    """Test data_chart tool"""
    with setup_server() as proc:
        try:
            chart_msg = {
                "jsonrpc": "2.0",
                "id": 6,
                "method": "tools/call",
                "params": {
                    "name": "data_chart",
                    "arguments": {
                        "file_path": "examples/chart.csv",
                        "column_index": 0,
                        "width": 50,
                    },
                },
            }
            send_message(proc, chart_msg)
            response = read_response(proc)
            if response:
                print("data_chart response:", json.dumps(response, indent=2))
                assert "result" in response
            else:
                print("No response to data_chart")
                assert False, "No response to data_chart"

        except Exception as e:
            print(f"Error: {e}")
            if proc.stderr:
                stderr = proc.stderr.read().decode()
                if stderr:
                    print(f"Server stderr: {stderr}")
            raise


def test_data_generate():
    """Test data_generate tool"""
    with setup_server() as proc:
        try:
            test_schema = '{"type":"object","properties":{"id":{"type":"integer"},"name":{"type":"string"},"age":{"type":"integer"}}}'
            generate_msg = {
                "jsonrpc": "2.0",
                "id": 7,
                "method": "tools/call",
                "params": {
                    "name": "data_generate",
                    "arguments": {"schema": test_schema, "lines": 2, "format": "csv"},
                },
            }
            send_message(proc, generate_msg)
            response = read_response(proc)
            if response:
                print("data_generate response:", json.dumps(response, indent=2))
                assert "result" in response
            else:
                print("No response to data_generate")
                assert False, "No response to data_generate"

        except Exception as e:
            print(f"Error: {e}")
            if proc.stderr:
                stderr = proc.stderr.read().decode()
                if stderr:
                    print(f"Server stderr: {stderr}")
            raise


def test_data_flattify():
    """Test data_flattify tool"""
    with setup_server() as proc:
        try:
            test_json = '{"user":{"profile":{"name":"John","age":30},"settings":{"theme":"dark"}}}'
            flattify_msg = {
                "jsonrpc": "2.0",
                "id": 8,
                "method": "tools/call",
                "params": {
                    "name": "data_flattify",
                    "arguments": {"input": test_json, "sql": False},
                },
            }
            send_message(proc, flattify_msg)
            response = read_response(proc)
            if response:
                print("data_flattify response:", json.dumps(response, indent=2))
                assert "result" in response
            else:
                print("No response to data_flattify")
                assert False, "No response to data_flattify"

        except Exception as e:
            print(f"Error: {e}")
            if proc.stderr:
                stderr = proc.stderr.read().decode()
                if stderr:
                    print(f"Server stderr: {stderr}")
            raise


def test_data_js2sql():
    """Test data_js2sql tool"""
    with setup_server() as proc:
        try:
            js2sql_msg = {
                "jsonrpc": "2.0",
                "id": 9,
                "method": "tools/call",
                "params": {
                    "name": "data_js2sql",
                    "arguments": {
                        "schema_file": "examples/generator-schema.json",
                        "table_name": "users",
                    },
                },
            }
            send_message(proc, js2sql_msg)
            response = read_response(proc)
            if response:
                print("data_js2sql response:", json.dumps(response, indent=2))
                assert "result" in response
            else:
                print("No response to data_js2sql")
                assert False, "No response to data_js2sql"

        except Exception as e:
            print(f"Error: {e}")
            if proc.stderr:
                stderr = proc.stderr.read().decode()
                if stderr:
                    print(f"Server stderr: {stderr}")
            raise


def test_data_profile():
    """Test data_profile tool"""
    with setup_server() as proc:
        try:
            profile_msg = {
                "jsonrpc": "2.0",
                "id": 10,
                "method": "tools/call",
                "params": {
                    "name": "data_profile",
                    "arguments": {"file_path": "examples/left.csv"},
                },
            }
            send_message(proc, profile_msg)
            response = read_response(proc)
            if response:
                print("data_profile response:", json.dumps(response, indent=2))
                assert "result" in response
            else:
                print("No response to data_profile")
                assert False, "No response to data_profile"

        except Exception as e:
            print(f"Error: {e}")
            if proc.stderr:
                stderr = proc.stderr.read().decode()
                if stderr:
                    print(f"Server stderr: {stderr}")
            raise
