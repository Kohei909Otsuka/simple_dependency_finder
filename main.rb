require 'bundler/setup'
require 'optparse'
require 'rubrowser' # TODO: あとで必要なところだけcopyしてくる
require 'csv'

opt = OptionParser.new

params = {
  path: []
}

opt.on('-path', '--path', Array, 'abs path to search comma separated') do |v|
  params[:path] = v
end
opt.parse(ARGV)

data = Rubrowser::Data.new(params[:path])

# ex: {"FirstClass": 1, "SecondClass": 2}
node_hash = {}
data.definitions.each_with_index do |definition, index|
  node_hash[definition.to_s] = index + 1
end

# write nodes.csv
CSV.open('neo4j/import/nodes_ruby.csv', 'wb') do |csv|
  csv << %w(id name path)
  data.definitions.each_with_index do |definition, index|
    csv << [ index + 1, definition.to_s, definition.file]
  end
end

# write relations.csv
CSV.open('neo4j/import/relations_ruby.csv', 'wb') do |csv|
  csv << %w(from to)
  data.relations.each do |relation|
    from = node_hash[relation.caller_namespace.to_s]
    to = node_hash[relation.namespace.to_s]
    next if from.nil? || to.nil?
    csv << [from, to]
  end
end
