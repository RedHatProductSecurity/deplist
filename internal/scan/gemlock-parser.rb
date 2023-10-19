require 'bundler'

# Read the Gemfile.lock
lockfile = Bundler.read_file(ARGV[0])

# Create a new LockfileParser object
parser = Bundler::LockfileParser.new(lockfile)

# Print all the dependencies
parser.specs.each do |spec|
  puts "#{spec.name} #{spec.version}"
end
