#!/usr/bin/env ruby

require "octokit"
require "tmpdir"

client = Octokit::Client.new(:access_token => ENV["GITHUB_AUTH_TOKEN"])

dest_repo = ARGV[0]
src_dir = File.absolute_path(".")
dest_path = `git rev-parse --show-prefix` # Get prefix from git root

new_repo = client.create_repository(dest_repo, {
  :homepage => "https://paper.dropbox.com/doc/Coding-Guide--AVRNhdoQFtz012lQFoOuqZ_gAg-aDM0OWpacNt0UTGFfuRCr",
  :description => "Look at the Issues for what needs to be done, look at the link for coding guide",
  :private => true,
  :has_wiki => false,
  :has_downloads => false,
})

clone_url = new_repo[:clone_url]

puts `git remote add #{dest_repo} #{clone_url}`

Dir.mktmpdir do |dest_dir|
  puts `cd #{dest_dir} && git init`
  puts `rsync -a #{src_dir}/ #{dest_dir}/#{dest_path}`
  puts `cd #{dest_dir} && git remote add origin #{clone_url}`
  puts `cd #{dest_dir} && git add .`
  puts `cd #{dest_dir} && git commit -m "initial commit"`
  puts `cd #{dest_dir} && git push --set-upstream origin master`
end
