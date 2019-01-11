#!/usr/bin/env ruby

require "commonmarker"
require "octokit"
require "pp"

repo = ARGV[0]
doc = CommonMarker.render_doc(open(ARGV[1]).read)

outline = Hash.new { "" }

current_header = ""
doc.each do |node|
  if node.type == :header and node.header_level == 1
    current_header = node.to_commonmark.gsub(/^# /, "").strip # .string_content
  elsif node.type == :list
    outline[current_header] += node.to_commonmark
  else
    outline[current_header] += node.to_commonmark
  end
end

client = Octokit::Client.new(:access_token => ENV["GITHUB_AUTH_TOKEN"])

existing_issues = client.list_issues(repo)
outline.each do |k, v|
  next if k == ""
  existing = existing_issues.select { |i| i.title == k }.first
  if existing
    client.update_issue(repo, existing.number, k, v)
  else
    client.create_issue(repo, k, v)
  end
end
