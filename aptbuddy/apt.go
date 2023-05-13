package aptbuddy

// http://ftp.debian.org/debian/dists/stable/main/binary-amd64/Packages.gz
/*
Package: 0ad
Source: 0ad (0.0.23.1-5)
Version: 0.0.23.1-5+b1
Installed-Size: 20268
Maintainer: Debian Games Team <pkg-games-devel@lists.alioth.debian.org>
Architecture: amd64
Depends: 0ad-data (>= 0.0.23.1), 0ad-data (<= 0.0.23.1-5), 0ad-data-common (>= 0.0.23.1), 0ad-data-common (<= 0.0.23.1-5), libboost-filesystem1.74.0 (>= 1.74.0), libc6 (>= 2.29), libcurl3-gnutls (>= 7.16.2), libenet7, libgcc-s1 (>= 3.4), libgl1, libgloox18 (>= 1.0.24), libicu67 (>= 67.1-1~), libminiupnpc17 (>= 1.9.20140610), libnspr4 (>= 2:4.9.2), libnvtt2, libopenal1 (>= 1.14), libpng16-16 (>= 1.6.2-1), libsdl2-2.0-0 (>= 2.0.12+dfsg1), libsodium23 (>= 1.0.14), libstdc++6 (>= 9), libvorbisfile3 (>= 1.1.2), libwxbase3.0-0v5 (>= 3.0.5.1+dfsg), libwxgtk3.0-gtk3-0v5 (>= 3.0.5.1+dfsg), libx11-6, libxcursor1 (>> 1.1.2), libxml2 (>= 2.9.0), zlib1g (>= 1:1.2.0)
Pre-Depends: dpkg (>= 1.15.6~)
Description: Real-time strategy game of ancient warfare
Homepage: http://play0ad.com/
Description-md5: d943033bedada21853d2ae54a2578a7b
Tag: game::strategy, interface::graphical, interface::x11, role::program,
 uitoolkit::sdl, uitoolkit::wxwidgets, use::gameplaying,
 x11::application
Section: games
Priority: optional
Filename: pool/main/0/0ad/0ad_0.0.23.1-5+b1_amd64.deb
Size: 5588508
MD5sum: 35412374733ae00cbbc7260596e1d78c
SHA256: 610e9f9c41be18af516dd64a6dc1316dbfe1bb8989c52bafa556de9e381d3e29

Package: 0ad-data
...
*/
// http://ftp.debian.org/debian/dists/stable/main/source/Sources.gz
/*
Package: 0ad
Binary: 0ad
Version: 0.0.23.1-5
Maintainer: Debian Games Team <pkg-games-devel@lists.alioth.debian.org>
Uploaders: Vincent Cheng <vcheng@debian.org>, Ludovic Rousseau <rousseau@debian.org>
Build-Depends: autoconf, debhelper-compat (= 12), dpkg-dev (>= 1.15.5), libboost-dev, libboost-filesystem-dev, libcurl4-gnutls-dev | libcurl4-dev, libenet-dev (>= 1.3), libgloox-dev (>= 1.0.10), libicu-dev, libminiupnpc-dev (>= 1.6), libnspr4-dev, libnvtt-dev (>= 2.0.8-1+dfsg-4~), libogg-dev, libopenal-dev, libpng-dev, libsdl2-dev (>= 2.0.2), libsodium-dev (>= 1.0.14), libvorbis-dev, libwxgtk3.0-gtk3-dev, libxcursor-dev, libxml2-dev, pkg-config, python2, python3, zlib1g-dev
Architecture: amd64 arm64 armhf i386 kfreebsd-amd64 kfreebsd-i386
Standards-Version: 4.5.0
Format: 3.0 (quilt)
Files:
 f8edafc49f74ae8eccfafd8613b97015 2438 0ad_0.0.23.1-5.dsc
 4fa111410ea55de7a013406ac1013668 31922812 0ad_0.0.23.1.orig.tar.xz
 43a5bf77192a8eebdbe763cdd1d72fa3 73620 0ad_0.0.23.1-5.debian.tar.xz
Vcs-Browser: https://salsa.debian.org/games-team/0ad
Vcs-Git: https://salsa.debian.org/games-team/0ad.git
Checksums-Sha256:
 f55d001ac0abbcc636e12f4b8d9df269c5fce93178287878b382ae89bc41b9ba 2438 0ad_0.0.23.1-5.dsc
 01bff7641ee08cac896c54d518d7e4b01752513105558f212e3199d747512a37 31922812 0ad_0.0.23.1.orig.tar.xz
 aff899c0b6a0c2ff746e051504a3e3ac7bb6070c21eb5a5ef5fb55d55391b0e0 73620 0ad_0.0.23.1-5.debian.tar.xz
Homepage: http://play0ad.com/
Package-List:
 0ad deb games optional arch=amd64,arm64,armhf,i386,kfreebsd-amd64,kfreebsd-i386
Directory: pool/main/0/0ad
Priority: source
Section: games

Package: 0ad-data
*/

// todo
// where do LONG descriptions stem from, which `apt-cache show` provides?
//  --> http://ftp.fr.debian.org/debian/dists/stable/main/i18n/Translation-en.bz2
// get Source: name - unique?
