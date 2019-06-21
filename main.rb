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

# pp data.definitions[0..20]
# pp data.relations

# write nodes.csv
CSV.open('neo4j/import/nodes_ruby.csv', 'wb') do |csv|
  csv << %w(id name path)
  data.definitions.each_with_index do |definition, index|
    csv << [ index + 1, definition.to_s, definition.file]
  end
end
