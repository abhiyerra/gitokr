#!/usr/bin/env ruby

require "octokit"
require "tmpdir"
require "byebug"

client = Octokit::Client.new(:access_token => ENV["GITHUB_AUTH_TOKEN"])

dest_repo = ARGV[0]
src_dir = File.absolute_path(".")
dest_path = `git rev-parse --show-prefix` # Get prefix from git root

new_repo = client.repository(dest_repo)
clone_url = new_repo[:clone_url]

Dir.mktmpdir do |dest_dir|
  puts `git clone #{clone_url} #{dest_dir}`
  puts `cd #{dest_dir} && git checkout -b gitdownstrem-#{Time.now.to_i}`
  puts `rsync -a #{src_dir}/ #{dest_dir}/#{dest_path}`
  puts `cd #{dest_dir} && git remote add origin #{clone_url}`
  puts `cd #{dest_dir} && git add .`
  puts `cd #{dest_dir} && git commit -m "initial commit"`
  puts `cd #{dest_dir} && git push --set-upstream origin master`
end
