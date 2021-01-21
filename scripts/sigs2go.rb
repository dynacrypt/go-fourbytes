#!/usr/bin/env ruby

require 'json'

case ARGV.size
when 0
  $stderr.puts "Syntax: #{File.basename($0)} DIR"
  exit 0
when 1
  dir = ARGV[0]
else
  $stderr.puts "Error: too many arguments!"
  exit 1
end

sigs = {}
Dir["#{dir}/signatures/????????"].each do |f|
  sig = File.basename(f)
  argfile = File.join(dir, 'with_parameter_names', sig)
  funcs = File.read(f).strip.split(';').map { |s| s.strip } - [ '' ]
  funcargs = (File.read(argfile) rescue '').split(';').map { |s| s.strip } - [ '' ]
  sigs[sig] = {
    'calls' => funcs,
    'callsWithArgs' => funcargs,
  }
end

calls = [
  "var sig2calls = map[string][]string{"
]
callsWithArgs = [
  "var sig2callsWithArgs = map[string][]string{"
]

sigs.keys.sort.each do |sig|
  f = sigs[sig]['calls'].map { |s| '"'+s+'"' }.join(', ')
  calls << "  \"#{sig}\":  []string{ #{f} },"

  if (sigs[sig]['callsWithArgs'] || []).any?
    f = sigs[sig]['callsWithArgs'].map { |s| '"'+s+'"' }.join(', ')
    callsWithArgs << "  \"#{sig}\":  []string{ #{f} },"
  end
end

puts calls, "}", ""
puts callsWithArgs, "}", ""
