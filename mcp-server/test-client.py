#!/usr/bin/env python3
"""
Simple test client for the DCT MCP Server
"""

import json
import subprocess
import sys


def send_message(proc, message):
    """Send a JSON-RPC message to the server"""
    data = json.dumps(message) + "\n"
    proc.stdin.write(data.encode())
    proc.stdin.flush()


def read_response(proc, timeout=5):
    """Read a JSON-RPC response from the server"""
    import select
    
    # Use select to timeout if no data is available
    if select.select([proc.stdout], [], [], timeout)[0]:
        line = proc.stdout.readline().decode().strip()
        if line:
            return json.loads(line)
    return None


def test_dct_mcp_server(dct_path):
    """Test the DCT MCP server"""
    cmd = ["./dct-mcp-server", dct_path]

    with subprocess.Popen(
        cmd, stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE
    ) as proc:
        try:
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

            # List tools
            list_msg = {"jsonrpc": "2.0", "id": 2, "method": "tools/list", "params": {}}
            send_message(proc, list_msg)
            response = read_response(proc)
            if response:
                print("\\nTools list response:", json.dumps(response, indent=2))
            else:
                print("\\nNo response to tools/list")

            # Test data_peek tool
            peek_msg = {
                "jsonrpc": "2.0",
                "id": 3,
                "method": "tools/call",
                "params": {
                    "name": "data_peek",
                    "arguments": {"file_path": "../examples/left.parquet", "lines": 3},
                },
            }
            send_message(proc, peek_msg)
            response = read_response(proc)
            if response:
                print("\\nPeek tool response:", json.dumps(response, indent=2))
            else:
                print("\\nNo response to data_peek")

        except Exception as e:
            print(f"Error: {e}")
            assert proc.stderr is not None
            stderr = proc.stderr.read().decode()
            if stderr:
                print(f"Server stderr: {stderr}")


if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python3 test-client.py <path-to-dct-binary>")
        sys.exit(1)

    test_dct_mcp_server(sys.argv[1])

