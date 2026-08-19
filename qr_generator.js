// qr_generator.js
const QRCode = require('qrcode');
const sharp = require('sharp');
const fs = require('fs');
const path = require('path');

const args = require('minimist')(process.argv.slice(2), {
    string: ['t', 'o', 'f', 'b', 'l'],
    default: { o: 'qr_code.png', f: '#000000', b: '#FFFFFF', s: 10, m: 4 },
    alias: { t: 'text', o: 'output', f: 'foreground', b: 'background', s: 'size', m: 'margin', l: 'logo' }
});

if (!args.t) {
    console.error('Error: -t text is required');
    process.exit(1);
}

const text = args.t;
const output = args.o;
const fg = args.f;
const bg = args.b;
const size = parseInt(args.s) || 10;
const margin = parseInt(args.m) || 4;
const logoPath = args.l;

// Generate QR as SVG (to get high-quality vector) then render to PNG
QRCode.toString(text, {
    type: 'svg',
    color: { dark: fg, light: bg },
    margin: margin,
    width: size * 40, // approximate width
    errorCorrectionLevel: 'H'
}, (err, svgString) => {
    if (err) {
        console.error('Error generating QR:', err);
        process.exit(1);
    }
    // Convert SVG to PNG using sharp
    const svgBuffer = Buffer.from(svgString);
    sharp(svgBuffer)
        .png()
        .toBuffer()
        .then(async (pngBuffer) => {
            if (logoPath && fs.existsSync(logoPath)) {
                // Overlay logo
                const qrImage = sharp(pngBuffer);
                const metadata = await qrImage.metadata();
                const logoSize = Math.min(metadata.width, metadata.height) * 0.25;
                const logoBuffer = await sharp(logoPath)
                    .resize({ width: Math.round(logoSize), height: Math.round(logoSize), fit: 'contain' })
                    .toBuffer();
                const composited = await sharp({
                    create: { width: metadata.width, height: metadata.height, channels: 4, background: { r: 0, g: 0, b: 0, alpha: 0 } }
                })
                .composite([
                    { input: pngBuffer, gravity: 'center' },
                    { input: logoBuffer, gravity: 'center' }
                ])
                .png()
                .toBuffer();
                fs.writeFileSync(output, composited);
            } else {
                fs.writeFileSync(output, pngBuffer);
            }
            console.log(`QR code saved to ${output}`);
        })
        .catch(err => {
            console.error('Error processing image:', err);
            process.exit(1);
        });
});
