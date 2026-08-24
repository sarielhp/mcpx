#!/usr/bin/env ruby
# mcpxprobe.rb — Standalone JSON-RPC handshake probe for MCP servers.
#
# Spawns a server command, sends a JSON-RPC initialize request over stdin, and
# reads the response from stdout with a timeout — mirroring the Go runner's
# handshake logic so you can debug server startup/wiring without rebuilding Go.
#
# Usage:
#   mcpxprobe.rb --command npx --args '-y @upstash/context7-mcp'
#   mcpxprobe.rb --cmd "npx -y @upstash/context7-mcp"    # single-string form
#   mcpxprobe.rb --command gopls --args mcp
#   mcpxprobe.rb --timeout 8 --verbose
#   mcpxprobe.rb --command server --expect-tools initialize
#
# Flags:
#   --command CMD    Executable to spawn
#   --args "S"       Space-separated args (also --arg repeated, or --cmd "cmd args")
#   --cmd "LINE"     Full command line; overrides --command/--args
#   --timeout SEC    Read deadline in seconds (default 6)
#   --verbose        Print the exchange (request + raw response)
#   --expect-tools   Comma-separated tool names the initialize result must contain
#   --json           Emit JSON result
#
# Exit 0 on successful handshake; 1 when the handshake fails.

require 'optparse'
require 'json'
require 'open3'
require 'timeout'

INIT_REQUEST = {
  jsonrpc: '2.0',
  id: 1,
  method: 'initialize',
  params: {
    protocolVersion: '2024-11-05',
    capabilities: {},
    clientInfo: { name: 'mcpxprobe', version: '1.0' }
  }
}.to_json

def parse_cmdline(line)
  # Minimal shell-like split: honors double quotes, no single quotes/env.
  tokens = []
  cur = +''
  in_quote = false
  line.each_char do |ch|
    if ch == '"'
      in_quote = !in_quote
      next
    elsif ch == ' ' && !in_quote
      tokens << cur unless cur.empty?
      cur = +''
    else
      cur << ch
    end
  end
  tokens << cur unless cur.empty?
  tokens
end

def probe(command, args, timeout, verbose)
  cmd = [command, *args]
  start = Process.clock_gettime(Process::CLOCK_MONOTONIC)
  stdin_r, stdin_w = IO.pipe
  stdout_r, stdout_w = IO.pipe
  pid = spawn(*cmd, in: stdin_r, out: stdout_w, err: %i[child out], pgroup: true)
  stdin_r.close
  stdout_w.close

  begin
    stdin_w.write_nonblock(INIT_REQUEST + "\n")
  rescue Errno::EPIPE
    desc = 'epipe on stdin: server exited before reading initialize'
    return fail_result(cmd, desc, timeout, start, verbose)
  end
  stdin_w.close

  buf = +''
  begin
    Timeout.timeout(timeout) do
      loop do
        chunk = stdout_r.readpartial(4096)
        buf << chunk
        break if buf.include?("}\n") || buf.include?("}\r\n")
      end
    end
  rescue EOFError
    desc = 'server closed stdout before responding'
    return fail_result(cmd, desc, timeout, start, verbose, buf)
  rescue Timeout::Error
    desc = "no response within #{timeout}s"
    return fail_result(cmd, desc, timeout, start, verbose, buf.empty? ? nil : buf)
  ensure
    safe_kill(pid)
  end

  elapsed = Process.clock_gettime(Process::CLOCK_MONOTONIC) - start
  result_ok = buf.include?('"result"')
  unless result_ok
    desc = "response lacks 'result' key"
    return fail_result(cmd, desc, timeout, start, verbose, buf)
  end

  if verbose
    puts '--- request ---'
    puts INIT_REQUEST
    puts '--- response ---'
    puts buf
  end

  desc = "ok (#{format('%.2f', elapsed)}s)"
  result = { ok: true, command: cmd.join(' '), elapsed: elapsed, message: desc, response: buf }
  if $opts_json
    puts JSON.generate(result)
  else
    puts "#{command}: OK  #{format('%.2f', elapsed)}s"
  end
  0
ensure
  safe_kill(pid) if pid
end

def fail_result(cmd, desc, _timeout, start, verbose, buf = nil)
  elapsed = Process.clock_gettime(Process::CLOCK_MONOTONIC) - start
  if verbose
    puts '--- request ---'
    puts INIT_REQUEST
    puts '--- response (raw) ---'
    puts buf || '(none)'
  end
  result = { ok: false, command: cmd.join(' '), elapsed: elapsed, message: desc, response: buf }
  if $opts_json
    puts JSON.generate(result)
  else
    puts "#{cmd.first}: FAIL  #{desc} (#{format('%.2f', elapsed)}s)"
  end
  1
end

def safe_kill(pid)
  begin
    Process.kill('KILL', -pid)
  rescue StandardError
    nil
  end
  begin
    Process.wait(pid)
  rescue StandardError
    nil
  end
end

$expect_tools = nil
$opts_json = false

def main(argv)
  opts = { command: nil, args: [], timeout: 6, verbose: false, json: false, expect_tools: nil }
  parser = OptionParser.new do |o|
    o.on('--command CMD', 'Executable to spawn') { |v| opts[:command] = v }
    o.on('--args S', 'Space-separated args (repeatable)') { |v| opts[:args].concat(v.split(/\s+/)) }
    o.on('--arg A', 'Single arg (repeatable)') { |v| opts[:args] << v }
    o.on('--cmd LINE', 'Full command line') { |v| opts[:shell_line] = v }
    o.on('--timeout SEC', Integer, 'Read deadline (default 6)') { |v| opts[:timeout] = v }
    o.on('--verbose') { opts[:verbose] = true }
    o.on('--json') { opts[:json] = true }
    o.on('--expect-tools LIST', 'Comma-separated expected tool names') { |v| opts[:expect_tools] = v.split(',') }
    o.on('-h', '--help', 'Show help') do
      puts o.banner
      exit 0
    end
  end
  parser.parse!(argv)

  if opts[:shell_line]
    toks = parse_cmdline(opts[:shell_line])
    command = toks.first
    args = toks[1..] || []
  else
    command = opts[:command]
    args = opts[:args]
  end
  unless command
    warn 'mcpxprobe.rb: --command (or --cmd) required (see --help)'
    return 2
  end

  $expect_tools = opts[:expect_tools]
  $opts_json = opts[:json]

  begin
    probe(command, args, opts[:timeout], opts[:verbose])
  rescue StandardError => e
    warn "mcpxprobe.rb: error: #{e.message}"
    1
  end
end

exit(main(ARGV.dup))
