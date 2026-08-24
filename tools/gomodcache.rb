#!/usr/bin/env ruby
# gomodcache.rb — Resolve installed Go module versions and inspect their API.
#
# Ends the "which copy of the module is the real official one installed on this
# machine" confusion by operating on the actual GOMODCACHE (the location `go`
# resolves dependencies from), NOT stray mirrors like ~/.sandbox.
#
# Usage:
#   gomodcache.rb clihelp                     # list cached versions
#   gomodcache.rb clihelp --latest            # show latest cached version
#   gomodcache.rb clihelp --api v0.2.14       # dump exported API of one version
#   gomodcache.rb clihelp --api v0.2.14 --name Execute --methods
#   gomodcache.rb --root                       # show resolved GOMODCACHE path
#
# Flags:
#   --root         Print the resolved GOMODCACHE path and exit
#   --latest       Print only the highest cached semver
#   --api VER      Inspect the exported API of the given cached version
#   --name SUBSTR  Substring filter on symbols (with --api)
#   --funcs/--types/--methods   Kind filters (with --api)
#   --json         JSON output where applicable
#
# Exit 0 on success; 1 when nothing found; 2 on usage error.

require 'optparse'
require 'json'

def gomodcache_path
  out = `go env GOMODCACHE 2>/dev/null`.strip
  return out unless out.empty?

  go_path = `go env GOPATH 2>/dev/null`.strip
  File.join(go_path, 'pkg', 'mod')
end

# Derive the module-cache directory name for a module path. Go escapes
# uppercase letters with a trailing '!' to keep file systems case-correct.
def cache_dir_name(module_name)
  module_name.split('/').map { |seg| seg.match?(/[A-Z]/) ? "#{seg}!" : seg }
end

# Returns absolute directory paths for every cached version of the module.
# Accepts either a full module path (e.g. "github.com/a/b") or just the leaf
# package name (e.g. "b"), in which case a recursive cache search is used.
def cached_module_dirs(root, module_name)
  if module_name.include?('/')
    segs = cache_dir_name(module_name)
    leaf = segs[-1]
    parent = segs[0...-1].join('/')
    return Dir.glob(File.join(root, parent, "#{leaf}@v*")).select { |d| File.directory?(d) }
  end
  Dir.glob(File.join(root, '**', "#{module_name}@v*")).select { |d| File.directory?(d) }
end

def version_of(dir)
  m = dir.match(/@(v\d+\.\d+(?:\.\d+)?(?:-[0-9A-Za-z.-]+)?)$/)
  m && m[1]
end

def sort_versions(dirs)
  dirs.map { |d| version_of(d) }.compact.sort_by { |v| v.split(/[.-]/).map { |p| p.to_i } }
end

def print_api(root, version, module_name, opts)
  dirs = cached_module_dirs(root, module_name)
  dir = dirs.find { |d| version_of(d) == version } || dirs.find { |d| version_of(d).to_s.start_with?(version) }
  unless dir
    warn "gomodcache.rb: version #{version} of #{module_name} not cached"
    return 1
  end
  api_script = File.expand_path('goapi.rb', __dir__)
  cmd = [api_script, '--path', dir]
  cmd << '--funcs' if opts[:funcs]
  cmd << '--types' if opts[:types]
  cmd << '--methods' if opts[:methods]
  cmd << '--name' << opts[:name] if opts[:name]
  cmd << '--json' if opts[:json]
  system(*cmd)
  $?.exitstatus || 0
end

def main(argv)
  opts = { name: nil, funcs: false, types: false, methods: false, json: false }
  parser = OptionParser.new do |o|
    o.on('--root', 'Print GOMODCACHE path and exit') { opts[:root] = true }
    o.on('--latest', 'Print only the highest cached version') { opts[:latest] = true }
    o.on('--api VER', 'Inspect API of the given version') { |v| opts[:api] = v }
    o.on('--name SUBSTR', 'Symbol substring filter (with --api)') { |v| opts[:name] = v }
    o.on('--funcs', 'Only functions (with --api)') { opts[:funcs] = true }
    o.on('--types', 'Only types (with --api)') { opts[:types] = true }
    o.on('--methods', 'Only methods (with --api)') { opts[:methods] = true }
    o.on('--json', 'JSON output') { opts[:json] = true }
    o.on('-h', '--help', 'Show help') do
      puts o.banner
      exit 0
    end
  end
  parser.parse!(argv)
  module_name = argv[0]
  root = gomodcache_path

  if opts[:root]
    puts root
    return 0
  end

  unless module_name
    warn 'gomodcache.rb: module name required (see --help)'
    return 2
  end

  versions = sort_versions(cached_module_dirs(root, module_name))
  if versions.empty?
    warn "gomodcache.rb: no cached versions of #{module_name} in #{root}"
    return 1
  end

  return print_api(root, opts[:api], module_name, opts) if opts[:api]

  if opts[:latest]
    puts versions[-1]
    return 0
  end

  if opts[:json]
    versions.each { |v| puts JSON.generate(version: v, module: module_name, cache: root) }
  else
    puts "Cached versions of #{module_name} in #{root}:"
    versions.each { |v| puts "  #{v}" }
  end
  0
end

exit(main(ARGV.dup))
