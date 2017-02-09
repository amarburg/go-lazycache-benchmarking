
task :build do
  %w( frame_stress random_walk ).each { |dir|
    chdir( dir ) {
      sh "go build"
    }
  }
end
