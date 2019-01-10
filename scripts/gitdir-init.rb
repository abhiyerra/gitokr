#!/usr/bin/env ruby

require "octokit"
require "tmpdir"
require "byebug"

client = Octokit::Client.new(:access_token => ENV["GITHUB_AUTH_TOKEN"])

dest_repo = ARGV[0]
src_dir = File.absolute_path(".")
dest_path = `git rev-parse --show-prefix` # Get prefix from git root

is_new_repo = false
new_repo = nil

begin
  new_repo = client.repository(dest_repo)
rescue Octokit::UnprocessableEntity, Octokit::NotFound
  is_new_repo = true
  new_repo = client.create_repository(dest_repo.split('/')[1], {
    :homepage => "https://paper.dropbox.com/doc/Coding-Guide--AVRNhdoQFtz012lQFoOuqZ_gAg-aDM0OWpacNt0UTGFfuRCr",
    :description => "Look at the Issues for what needs to be done, look at the link for coding guide",
    :private => true,
    :has_wiki => false,
    :has_downloads => false,
  })
end

clone_url = new_repo[:clone_url]

puts `git remote add #{dest_repo} #{clone_url}`

Dir.mktmpdir do |dest_dir|
  if is_new_repo
    puts `cd #{dest_dir} && git init`
  else
    puts `git clone #{clone_url} #{dest_dir}`
  end

  puts `cd #{dest_dir} && mkdir -p #{dest_path}`
  puts `rsync -a #{src_dir}/ #{dest_dir}/#{dest_path}`
  puts `cd #{dest_dir} && git remote add origin #{clone_url}`
  puts `cd #{dest_dir} && git add .`
  puts `cd #{dest_dir} && git commit -m "initial commit"`
  puts `cd #{dest_dir} && git push --set-upstream origin master`
end
