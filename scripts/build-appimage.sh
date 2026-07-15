#!/bin/bash
set -e

VERSION=$1
BUILD_DIR="build-release"
APPDIR="$BUILD_DIR/AppDir"

# ایجاد ساختار AppDir
mkdir -p "$APPDIR/usr/bin"
mkdir -p "$APPDIR/usr/lib"
mkdir -p "$APPDIR/usr/share/applications"
mkdir -p "$APPDIR/usr/share/icons/hicolor/256x256/apps"

# کپی فایل‌های اجرایی
cp "$BUILD_DIR/qt/kingo-tray" "$APPDIR/usr/bin/"
cp "$BUILD_DIR/cmd/kingo-linux-vpn" "$APPDIR/usr/bin/"

# کپی منابع
cp -r "$BUILD_DIR/ui/resources" "$APPDIR/usr/share/kingo-vpn/"

# کپی فایل دسکتاپ
cp packaging/kingo-linux-vpn.desktop "$APPDIR/usr/share/applications/"

# دانلود آیکون
wget -q "https://raw.githubusercontent.com/agolzarniya4213-star/kingo-linux-vpn/main/qt/icons/connected-symbolic.svg" \
  -O "$APPDIR/usr/share/icons/hicolor/256x256/apps/kingo-vpn.svg" || \
  cp qt/icons/connected-symbolic.svg "$APPDIR/usr/share/icons/hicolor/256x256/apps/kingo-vpn.svg"

# ایجاد AppImage
wget -q "https://github.com/AppImage/AppImageKit/releases/download/continuous/appimagetool-x86_64.AppImage" \
  -O "$BUILD_DIR/appimagetool"
chmod +x "$BUILD_DIR/appimagetool"

"$BUILD_DIR/appimagetool" "$APPDIR" "Kingo-Linux-VPN-$VERSION-x86_64.AppImage"

echo "AppImage created: Kingo-Linux-VPN-$VERSION-x86_64.AppImage"
