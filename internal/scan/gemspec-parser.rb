require 'bundler'

deps = Gem::Specification.load(ARGV[0]).dependencies
deps.each do |dep|
  if dep.runtime?()
    puts "#{dep.name}"
  end
end
