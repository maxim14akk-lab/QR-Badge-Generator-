// QRGenerator.java
import com.google.zxing.*;
import com.google.zxing.common.BitMatrix;
import com.google.zxing.qrcode.QRCodeWriter;
import com.google.zxing.qrcode.decoder.ErrorCorrectionLevel;

import javax.imageio.ImageIO;
import java.awt.*;
import java.awt.image.BufferedImage;
import java.io.File;
import java.io.IOException;

public class QRGenerator {
    public static void main(String[] args) throws Exception {
        String text = null, output = "qr_code.png", fg = "#000000", bg = "#FFFFFF", logo = null;
        int size = 10, margin = 4;

        for (int i = 0; i < args.length; i++) {
            switch (args[i]) {
                case "-t": case "--text": text = args[++i]; break;
                case "-o": case "--output": output = args[++i]; break;
                case "-f": case "--foreground": fg = args[++i]; break;
                case "-b": case "--background": bg = args[++i]; break;
                case "-s": case "--size": size = Integer.parseInt(args[++i]); break;
                case "-m": case "--margin": margin = Integer.parseInt(args[++i]); break;
                case "-l": case "--logo": logo = args[++i]; break;
            }
        }
        if (text == null) {
            System.err.println("Error: -t text is required");
            System.exit(1);
        }

        QRCodeWriter writer = new QRCodeWriter();
        BitMatrix matrix = writer.encode(text, BarcodeFormat.QR_CODE, size * 40, size * 40,
                java.util.Map.of(EncodeHintType.ERROR_CORRECTION, ErrorCorrectionLevel.H,
                                 EncodeHintType.MARGIN, margin));

        BufferedImage image = new BufferedImage(matrix.getWidth(), matrix.getHeight(), BufferedImage.TYPE_INT_RGB);
        Color fgColor = Color.decode(fg);
        Color bgColor = Color.decode(bg);
        for (int x = 0; x < matrix.getWidth(); x++) {
            for (int y = 0; y < matrix.getHeight(); y++) {
                image.setRGB(x, y, matrix.get(x, y) ? fgColor.getRGB() : bgColor.getRGB());
            }
        }

        // Embed logo if provided
        if (logo != null) {
            try {
                BufferedImage logoImg = ImageIO.read(new File(logo));
                int logoSize = (int) (Math.min(matrix.getWidth(), matrix.getHeight()) * 0.25);
                Image scaledLogo = logoImg.getScaledInstance(logoSize, logoSize, Image.SCALE_SMOOTH);
                Graphics2D g = image.createGraphics();
                int x = (image.getWidth() - logoSize) / 2;
                int y = (image.getHeight() - logoSize) / 2;
                g.drawImage(scaledLogo, x, y, null);
                g.dispose();
            } catch (Exception e) {
                System.err.println("Warning: could not embed logo: " + e.getMessage());
            }
        }

        ImageIO.write(image, "png", new File(output));
        System.out.println("QR code saved to " + output);
    }
}
