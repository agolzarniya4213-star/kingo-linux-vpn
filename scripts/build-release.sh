#!/bin/bash
set -e

# پیکربندی
VERSION=$(cat VERSION 2>/dev/null || echo "1.0.0")
BUILD_DIR="build-release"
INSTALL_DIR="$BUILD_DIR/install"

# پاک‌سازی ساخت قبلی
rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"

# ساخت با CMake
cd "$BUILD_DIR"
cmake .. -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=/usr
make -j$(nproc)

# نصب در مسیر موقت
make DESTDIR="$INSTALL_DIR" install

# ایجاد بسته‌های مختلف
echo "Building packages for version $VERSION..."

# AppImage
./scripts/build-appimage.sh "$VERSION"

# .deb package (برای Ubuntu/Debian)
./scripts/build-deb.sh "$VERSION" "$INSTALL_DIR"

# .rpm package (برای Fedora/RHEL)
./scripts/build-rpm.sh "$VERSION" "$INSTALL_DIR"

echo "Build completed successfully!"
ls -la "$BUILD_DIR"/*.AppImage "$BUILD_DIR"/*.deb "$BUILD_DIR"/*.rpm 2>/dev/null || true
