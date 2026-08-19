// qr_generator.cpp
#include <iostream>
#include <fstream>
#include <vector>
#include <string>
#include <cstring>
#include <cmath>
#include <getopt.h>
#include "QrCode.hpp" // Nayuki's QR-Code-generator: https://github.com/nayuki/QR-Code-generator
#define STB_IMAGE_WRITE_IMPLEMENTATION
#include "stb_image_write.h"

using namespace std;
using namespace qrcodegen;

int main(int argc, char **argv) {
    string text, output = "qr_code.png", fg = "#000000", bg = "#FFFFFF", logoPath;
    int size = 10, margin = 4;

    static struct option long_options[] = {
        {"text", required_argument, 0, 't'},
        {"output", required_argument, 0, 'o'},
        {"foreground", required_argument, 0, 'f'},
        {"background", required_argument, 0, 'b'},
        {"size", required_argument, 0, 's'},
        {"margin", required_argument, 0, 'm'},
        {"logo", required_argument, 0, 'l'},
        {0,0,0,0}
    };
    int c;
    while ((c = getopt_long(argc, argv, "t:o:f:b:s:m:l:", long_options, nullptr)) != -1) {
        switch (c) {
            case 't': text = optarg; break;
            case 'o': output = optarg; break;
            case 'f': fg = optarg; break;
            case 'b': bg = optarg; break;
            case 's': size = stoi(optarg); break;
            case 'm': margin = stoi(optarg); break;
            case 'l': logoPath = optarg; break;
        }
    }
    if (text.empty()) {
        cerr << "Error: -t text is required\n";
        return 1;
    }

    // Generate QR
    QrCode qr = QrCode::encodeText(text.c_str(), QrCode::Ecc::HIGH);
    int qrSize = qr.getSize();
    int imgSize = (qrSize + 2 * margin) * size;
    vector<unsigned char> img(imgSize * imgSize * 4, 0); // RGBA

    // Parse colors (simplistic: accept #RRGGBB only)
    auto parseColor = [](const string& hex) -> tuple<unsigned char,unsigned char,unsigned char> {
        unsigned int r,g,b;
        if (hex.size()==7 && hex[0]=='#') {
            sscanf(hex.c_str(), "#%02x%02x%02x", &r, &g, &b);
        } else {
            r=0; g=0; b=0;
        }
        return {r,g,b};
    };
    auto [fr,fgc,fb] = parseColor(fg);
    auto [br,bgc,bb] = parseColor(bg);

    for (int y = 0; y < imgSize; y++) {
        for (int x = 0; x < imgSize; x++) {
            int qx = (x / size) - margin;
            int qy = (y / size) - margin;
            bool pixel = (qx >= 0 && qx < qrSize && qy >= 0 && qy < qrSize && qr.getModule(qx, qy));
            int idx = (y * imgSize + x) * 4;
            if (pixel) {
                img[idx] = fr; img[idx+1] = fgc; img[idx+2] = fb; img[idx+3] = 255;
            } else {
                img[idx] = br; img[idx+1] = bgc; img[idx+2] = bb; img[idx+3] = 255;
            }
        }
    }

    // Save PNG
    if (!stbi_write_png(output.c_str(), imgSize, imgSize, 4, img.data(), imgSize*4)) {
        cerr << "Error writing PNG\n";
        return 1;
    }
    cout << "QR code saved to " << output << "\n";
    return 0;
}
