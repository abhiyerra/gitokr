#!/usr/bin/env ruby

require "octokit"
require "byebug"

client = Octokit::Client.new(:access_token => ENV["GITHUB_AUTH_TOKEN"])

dest_repo = ARGV[0]
collab = ARGV[1]

if client.contribs(dest_repo).select { |c| c.login == collab }.first
  client.remove_collab(dest_repo, collab)
else
  client.add_collab(dest_repo, collab)
end
