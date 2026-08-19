// QRGenerator.cs
using System;
using System.Drawing;
using System.Drawing.Imaging;
using QRCoder;
using SixLabors.ImageSharp;
using SixLabors.ImageSharp.PixelFormats;
using SixLabors.ImageSharp.Processing;
using SixLabors.ImageSharp.Drawing.Processing;

class QRGenerator
{
    static int Main(string[] args)
    {
        string text = null, output = "qr_code.png", fg = "#000000", bg = "#FFFFFF", logo = null;
        int size = 10, margin = 4;

        for (int i = 0; i < args.Length; i++)
        {
            switch (args[i])
            {
                case "-t": case "--text": text = args[++i]; break;
                case "-o": case "--output": output = args[++i]; break;
                case "-f": case "--foreground": fg = args[++i]; break;
                case "-b": case "--background": bg = args[++i]; break;
                case "-s": case "--size": size = int.Parse(args[++i]); break;
                case "-m": case "--margin": margin = int.Parse(args[++i]); break;
                case "-l": case "--logo": logo = args[++i]; break;
            }
        }
        if (text == null)
        {
            Console.Error.WriteLine("Error: -t text is required");
            return 1;
        }

        // Generate QR using QRCoder
        var generator = new QRCodeGenerator();
        var qrData = generator.CreateQrCode(text, QRCodeGenerator.ECCLevel.H);
        var qrCode = new QRCode(qrData);
        var bitmap = qrCode.GetGraphic(size, Color.FromName(fg), Color.FromName(bg), true, margin: margin);

        // Embed logo if provided
        if (logo != null && System.IO.File.Exists(logo))
        {
            try
            {
                using var logoImg = Image.Load<Rgba32>(logo);
                var qrImg = Image.Load<Rgba32>(bitmap.ToBytes());
                int logoSize = (int)(Math.Min(qrImg.Width, qrImg.Height) * 0.25);
                logoImg.Mutate(x => x.Resize(logoSize, logoSize));
                qrImg.Mutate(x => x.DrawImage(logoImg, new Point((qrImg.Width - logoSize) / 2, (qrImg.Height - logoSize) / 2), 1f));
                qrImg.Save(output);
            }
            catch (Exception ex)
            {
                Console.Error.WriteLine($"Warning: could not embed logo: {ex.Message}");
                bitmap.Save(output, ImageFormat.Png);
            }
        }
        else
        {
            bitmap.Save(output, ImageFormat.Png);
        }
        Console.WriteLine($"QR code saved to {output}");
        return 0;
    }
}
