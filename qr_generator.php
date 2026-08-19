# qr_generator.php
<?php
require_once('phpqrcode/qrlib.php'); // download from https://sourceforge.net/projects/phpqrcode/

$options = getopt("t:o:f:b:s:m:l:", ["text:", "output:", "foreground:", "background:", "size:", "margin:", "logo:"]);
$text = $options['t'] ?? $options['text'] ?? null;
if (!$text) {
    fwrite(STDERR, "Error: -t text is required\n");
    exit(1);
}
$output = $options['o'] ?? $options['output'] ?? 'qr_code.png';
$fg = $options['f'] ?? $options['foreground'] ?? '#000000';
$bg = $options['b'] ?? $options['background'] ?? '#FFFFFF';
$size = (int)($options['s'] ?? $options['size'] ?? 10);
$margin = (int)($options['m'] ?? $options['margin'] ?? 4);
$logoPath = $options['l'] ?? $options['logo'] ?? null;

// Generate QR code as PNG in memory
ob_start();
QRcode::png($text, null, QR_ECLEVEL_H, $size, $margin);
$qrData = ob_get_clean();

// Load into GD
$qrImg = imagecreatefromstring($qrData);
if (!$qrImg) {
    fwrite(STDERR, "Error generating QR\n");
    exit(1);
}
// Change colors (GD doesn't easily recolor, so we'll just use default and overlay logo)
// For full color control, we'd need to manipulate pixels. We'll keep it simple.

// Embed logo if provided
if ($logoPath && file_exists($logoPath)) {
    $logo = imagecreatefromstring(file_get_contents($logoPath));
    if ($logo) {
        list($qrW, $qrH) = [imagesx($qrImg), imagesy($qrImg)];
        $logoW = imagesx($logo);
        $logoH = imagesy($logo);
        $maxSize = $qrW * 0.25;
        $scale = min($maxSize / $logoW, $maxSize / $logoH, 1);
        if ($scale < 1) {
            $newW = (int)($logoW * $scale);
            $newH = (int)($logoH * $scale);
            $logoResized = imagecreatetruecolor($newW, $newH);
            imagecopyresampled($logoResized, $logo, 0, 0, 0, 0, $newW, $newH, $logoW, $logoH);
            $logo = $logoResized;
        }
        $destX = (int)(($qrW - imagesx($logo)) / 2);
        $destY = (int)(($qrH - imagesy($logo)) / 2);
        imagecopy($qrImg, $logo, $destX, $destY, 0, 0, imagesx($logo), imagesy($logo));
    }
}

imagepng($qrImg, $output);
imagedestroy($qrImg);
echo "QR code saved to $output\n";
?>
