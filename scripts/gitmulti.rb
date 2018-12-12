#!/usr/bin/env ruby

open(ARGV[0]).read.split("\n").each do |line|
    puts `hub issue create -m "#{line}"`
end