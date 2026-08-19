# qr_generator.py
import argparse
import qrcode
from PIL import Image, ImageDraw
import sys

def generate_qr(text, output, fg, bg, size, margin, logo):
    qr = qrcode.QRCode(
        version=1,
        error_correction=qrcode.constants.ERROR_CORRECT_H,
        box_size=size,
        border=margin,
    )
    qr.add_data(text)
    qr.make(fit=True)
    img = qr.make_image(fill_color=fg, back_color=bg).convert('RGB')
    if logo:
        try:
            logo_img = Image.open(logo)
            # resize logo to fit 1/4 of QR size
            qr_width, qr_height = img.size
            logo_size = int(min(qr_width, qr_height) * 0.25)
            logo_img = logo_img.resize((logo_size, logo_size), Image.Resampling.LANCZOS)
            pos = ((qr_width - logo_size) // 2, (qr_height - logo_size) // 2)
            img.paste(logo_img, pos, mask=logo_img)
        except Exception as e:
            print(f"Warning: could not embed logo: {e}", file=sys.stderr)
    img.save(output)
    print(f"QR code saved to {output}")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="QR Badge Generator")
    parser.add_argument("-t", "--text", required=True, help="Text to encode")
    parser.add_argument("-o", "--output", default="qr_code.png", help="Output file name")
    parser.add_argument("-f", "--foreground", default="#000000", help="Foreground color (hex or name)")
    parser.add_argument("-b", "--background", default="#FFFFFF", help="Background color")
    parser.add_argument("-s", "--size", type=int, default=10, help="Size per module (pixels)")
    parser.add_argument("-m", "--margin", type=int, default=4, help="Margin in modules")
    parser.add_argument("-l", "--logo", help="Path to logo image")
    args = parser.parse_args()
    generate_qr(args.text, args.output, args.foreground, args.background, args.size, args.margin, args.logo)
