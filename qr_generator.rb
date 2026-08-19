# qr_generator.rb
require 'rqrcode'
require 'chunky_png'
require 'optparse'
require 'oily_png' # optional for speed

options = {}
OptionParser.new do |opts|
  opts.banner = "Usage: ruby qr_generator.rb -t TEXT [options]"
  opts.on("-t TEXT", "--text TEXT", "Text to encode") { |v| options[:text] = v }
  opts.on("-o FILE", "--output FILE", "Output file name", default: "qr_code.png") { |v| options[:output] = v }
  opts.on("-f COLOR", "--foreground COLOR", "Foreground color (hex or name)", default: "#000000") { |v| options[:fg] = v }
  opts.on("-b COLOR", "--background COLOR", "Background color", default: "#FFFFFF") { |v| options[:bg] = v }
  opts.on("-s SIZE", "--size SIZE", Integer, "Size per module (pixels)", default: 10) { |v| options[:size] = v }
  opts.on("-m MARGIN", "--margin MARGIN", Integer, "Margin in modules", default: 4) { |v| options[:margin] = v }
  opts.on("-l LOGO", "--logo LOGO", "Path to logo image") { |v| options[:logo] = v }
end.parse!

unless options[:text]
  $stderr.puts "Error: -t text is required"
  exit 1
end

qr = RQRCode::QRCode.new(options[:text], level: :h)
png = ChunkyPNG::Image.new(qr.modules.size + 2 * options[:margin], qr.modules.size + 2 * options[:margin], ChunkyPNG::Color.from_hex(options[:bg]))

fg_color = ChunkyPNG::Color.from_hex(options[:fg])
size = options[:size]
margin = options[:margin]

qr.modules.each_index do |y|
  qr.modules.each_index do |x|
    if qr.modules[y][x]
      (0...size).each do |dy|
        (0...size).each do |dx|
          png[ (margin + x) * size + dx, (margin + y) * size + dy ] = fg_color
        end
      end
    end
  end
end

# Embed logo if provided
if options[:logo] && File.exist?(options[:logo])
  begin
    logo = ChunkyPNG::Image.from_file(options[:logo])
    logo_width = logo.width
    logo_height = logo.height
    max_logo_size = (qr.modules.size * size) * 0.25
    scale = [max_logo_size / logo_width, max_logo_size / logo_height, 1].min
    if scale < 1
      logo = logo.resample_bilinear((logo_width * scale).round, (logo_height * scale).round)
    end
    x_offset = (png.width - logo.width) / 2
    y_offset = (png.height - logo.height) / 2
    png.compose!(logo, x_offset, y_offset)
  rescue => e
    $stderr.puts "Warning: could not embed logo: #{e}"
  end
end

png.save(options[:output])
puts "QR code saved to #{options[:output]}"
