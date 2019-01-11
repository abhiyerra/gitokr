#!/usr/bin/env ruby

require "commonmarker"
require "octokit"

repo = ARGV[0]
doc = CommonMarker.render_doc(open(ARGV[1]).read)

outline = Hash.new { "" }

doc.walk do |node|
  if node.type == :header and node.header_level == 1
    header = node.to_plaintext.strip
    n = node.next
    outline[header] = n.walk.map do |subnode|
      if subnode.type == :list || subnode.type == :list_item
        ""
      else
        subnode.string_content rescue "\n"
      end
    end.join if n
  end
end

client = Octokit::Client.new(:access_token => ENV["GITHUB_AUTH_TOKEN"])

existing_issues = client.list_issues(repo)
outline.each do |k, v|
  existing = existing_issues.select { |i| i.title == k }.first
  if existing
    client.update_issue(repo, existing.number, k, v)
  else
    client.create_issue(repo, k, v)
  end
end
