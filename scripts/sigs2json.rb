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
    'call' => funcs,
    'callWithArgs' => funcargs,
  }
end

puts JSON.generate(sigs)
