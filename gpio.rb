require 'epoll'
require 'net/http'

def watch pin, on:
  File.binwrite "/sys/class/gpio/export", pin.to_s

  retries = 0
  begin
    File.binwrite "/sys/class/gpio/gpio#{pin}/edge", on
  rescue
    raise if retries > 3
    sleep 0.1
    retries += 1
    retry
  end

  fd = File.open "/sys/class/gpio/gpio#{pin}/value", 'r'
  yield fd.read.chomp

  epoll = Epoll.create
  epoll.add fd, Epoll::PRI

  loop do
    fd.seek 0, IO::SEEK_SET
    epoll.wait
    yield fd.read.chomp
  end
ensure
  File.binwrite "/sys/class/gpio/unexport", pin.to_s
end

pin  = ENV['KICKER_SENSOR_PIN']
team = ENV['KICKER_TEAM']
uri  = URI('http://hkick:3000/goals')

watch pin, on: :both do |value|
  p value
  next unless value.to_i == 0

  res = Net::HTTP.post_form(uri, 'team' => team)
  p res
end
