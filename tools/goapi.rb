#!/usr/bin/env ruby
# goapi.rb — Explore the exported API surface of Go packages (stdlib or module cache).
#
# Use when you need the exact signature of a function, type, or method in a
# specific Go version / installed module WITHOUT writing throwaway probe
# programs. Parses .go source files directly.
#
# Usage:
#   goapi.rb --path /usr/lib/go-1.26/src/io --name Write
#   goapi.rb --path ~/.go/pkg/mod/github.com/sarielhp/clihelp@v0.2.14 --methods
#   goapi.rb --path /usr/lib/go-1.26/src/strings --name Sort
#   goapi.rb --stdlib --name Exists
#
# Flags:
#   --path DIR     Directory tree to scan (repeatable)
#   --stdlib       Shortcut for --path /usr/lib/go-1.26/src
#   --name SUBSTR  Substring filter on symbol names
#   --funcs        Only show top-level funcs
#   --types        Only show types (struct/interface)
#   --methods      Only show methods (receiver funcs)
#   --json         Emit JSON lines
#   --help         Show this help
#
# Exit 0 on success (even with 0 hits); 2 on usage error.

require 'optparse'
require 'json'

def collect_go_files(dir)
  heap = [File.expand_path(dir)]
  files = []
  until heap.empty?
    cur = heap.pop
    next unless File.directory?(cur)

    Dir.each_child(cur) do |entry|
      next if %w[testdata vendor].include?(entry) || entry.start_with?('.')

      full = File.join(cur, entry)
      if File.directory?(full)
        heap << full
      elsif entry.end_with?('.go') && !entry.end_with?('_test.go')
        files << full
      end
    end
  end
  files
end

# Returns array of [name, kind, signature] where kind is :func, :type, or :method.
def parse_signatures(file)
  out = []
  text = File.read(file, mode: 'r')
  text = text.gsub(%r{//.*$}, '')
  pos = 0
  while (m = text.match(/^(func|type)\b/, pos))
    start = m.begin(0)
    nl = text.index("\n", start)
    nl ||= text.length
    line = text[start...nl]
    pos = nl + 1

    if m[1] == 'func'
      fm = line.match(/^func\s+(\([^)]*\)\s+)?([A-Z]\w*)\s*\(([^)]*)\)/)
      next unless fm

      recv = fm[1] ? "#{fm[1].gsub(/\s+/, ' ').strip} " : ''
      name = fm[2]
      params = fm[3].to_s.strip
      sig = "func #{recv}#{name}(#{params})"
      out << [name, recv.empty? ? :func : :method, sig]
    elsif m[1] == 'type'
      tm = line.match(/^type\s+(\w+)\s+(struct|interface)\b/)
      out << [tm[1], :type, "type #{tm[1]} #{tm[2]}"] if tm
    end
  end
  out
end

def main(argv)
  opts = { paths: [], name: nil, funcs: false, types: false, methods: false, json: false }
  parser = OptionParser.new do |o|
    o.on('--path DIR', 'Go source directory to scan (repeatable)') { |v| opts[:paths] << v }
    o.on('--stdlib', 'Scan /usr/lib/go-1.26/src') { opts[:paths] << '/usr/lib/go-1.26/src' }
    o.on('--name SUBSTR', 'Filter on symbol substring') { |v| opts[:name] = v }
    o.on('--funcs', 'Only top-level funcs') { opts[:funcs] = true }
    o.on('--types', 'Only types') { opts[:types] = true }
    o.on('--methods', 'Only methods') { opts[:methods] = true }
    o.on('--json', 'Emit JSON lines') { opts[:json] = true }
    o.on('-h', '--help', 'Show help') do
      puts o.banner
      exit 0
    end
  end
  parser.parse!(argv)

  if opts[:paths].empty?
    warn 'goapi.rb: no --path/--stdlib given (see --help)'
    return 2
  end

  any_kind = opts[:funcs] || opts[:types] || opts[:methods]
  wanted = lambda do |kind|
    return true unless any_kind

    { func: opts[:funcs], type: opts[:types], method: opts[:methods] }[kind]
  end

  files = opts[:paths].flat_map { |p| collect_go_files(p) }.uniq.sort
  results = []
  files.each do |f|
    cwd = Dir.pwd + '/'
    rel = f.start_with?(cwd) ? f[/#{Regexp.escape(cwd)}(.*)/, 1] : f
    parse_signatures(f).each do |name, kind, sig|
      next unless wanted.call(kind)
      next if opts[:name] && !name.include?(opts[:name])

      results << [name, kind.to_s, sig, rel]
    end
  end

  results.sort_by! { |n, _, _, f| [f, n] }.uniq!

  if opts[:json]
    results.each do |n, k, sig, f|
      puts JSON.generate(name: n, kind: k, signature: sig, file: f)
    end
  else
    results.each do |n, k, _, f|
      printf("%-34s %-7s %s\n", n, k, f)
    end
  end
  0
end

exit(main(ARGV.dup))
