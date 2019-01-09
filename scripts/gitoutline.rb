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
    n.walk.each { |subnode| outline[header] += subnode.to_commonmark } if n
  end
end

client = Octokit::Client.new(:access_token => ENV["GITHUB_AUTH_TOKEN"])
outline.each do |k, v|
  client.create_issue(repo, k, v)
end
