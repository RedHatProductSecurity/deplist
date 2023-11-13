require 'bundler'


def get_lockfile_deps(lockfile, deps = nil)
    lockfileParser = Bundler::LockfileParser.new(lockfile)
    rvalue = []

    lockfileParser.specs.each do |spec|
      if deps
        if deps.include?(spec.name)
          rvalue.append(spec)
        end
      else
        rvalue.append(spec)
      end
    end

    return rvalue
end

def get_paths()
  dir = "."
  if ARGV.length > 0
    dir = ARGV[0]
  end

  gemfilePath = File.join(dir, "Gemfile")
  lockfilePath = File.join(dir, "Gemfile.lock")

  return gemfilePath, lockfilePath
end


def get_gemfiles(gemfilePath, lockfilePath)
  gemfile = nil
  lockfile = nil
  if File.exist? File.expand_path gemfilePath
    gemfile = Bundler.read_file(gemfilePath)
  end
  if File.exist? File.expand_path lockfilePath
    lockfile = Bundler.read_file(lockfilePath)
  end

  return gemfile, lockfile
end

def print_spec_recursive(dep, all_specs)
    all_specs.each do |spec|
      if spec.name == dep.name
        puts "#{spec.name} #{spec.version}"
        spec.dependencies.each do |spec2|
          print_spec_recursive(spec2, all_specs)
        end
      end
    end
end


gemfilePath, lockfilePath = get_paths()
gemfile, lockfile = get_gemfiles(gemfilePath, lockfilePath)

exclude_groups = [
  :development,
  :guard,
  :packaging,
  :release,
  :system_tests,
  :test,
]

if gemfile
  if lockfile
    gemfileParser = Bundler::Definition.build(gemfilePath, lockfilePath, nil)
  else
    gemfileParser = Bundler::Definition.build(gemfilePath, '', nil)
  end

  # First get runtime deps from Gemfile
  runtime_deps = []
  gemfileParser.dependencies.each do |dep|
    # puts dep
    if dep.should_include?
      if !exclude_groups.any? { |g| dep.groups.include?(g) }
        # puts "#{dep.name} #{dep.groups}"
        runtime_deps.append(dep.name)
      end
    end
  end

  # Then compare with lockfile to get more specific version info and transitive deps
  if lockfile
    runtime_specs = get_lockfile_deps(lockfile, runtime_deps)
    all_specs = get_lockfile_deps(lockfile, nil)
    runtime_specs.each do |spec|
      puts "#{spec.name} #{spec.version}"
      spec.dependencies.each do |dep2|
        print_spec_recursive(dep2, all_specs)
      end
    end
  else
    # if no lockfile then just print deps without version info
    runtime_deps.each do |r|
      puts "#{r}"
    end
  end
  
end
