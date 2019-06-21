require 'optparse'

opt = OptionParser.new

params = {
  path: []
}

opt.on('-path', '--path', Array, 'abs path to search comma separated') do |v|
  params[:path] = v
end

opt.parse(ARGV)
p params
