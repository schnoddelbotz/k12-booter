package aptbuddy

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
	"regexp"
	"schnoddelbotz/k12-booter/utility"

	"github.com/blevesearch/bleve/v2"
)

/*

Assumptions made here:
1. The debian repositories are public.
2. k12-booter's intension is to increase OSS adoption, be it on debian or not.
3. k12-booter states package descriptions origin and does not try to monetarise.

If you see any issue here, please ping me. Thank you. Peace.

https://www.debian.org/mirror/list -> CDN, https -> https://deb.debian.org/debian/

*/

type APTPackage struct {
	Package         string
	Description     string
	Homepage        string
	Tags            []string
	Section         string
	LongDescription string // depends on I18N ...
}

const (
	PackagesGZ = "Packages.gz"
)

func FetchFiles(translation string) error {
	packagesURL := fmt.Sprintf("https://deb.debian.org/debian/dists/stable/main/binary-amd64/%s", PackagesGZ)
	utility.Fetch(packagesURL, PackagesGZ)
	i18nURL := fmt.Sprintf("https://deb.debian.org/debian/dists/stable/main/i18n/%s", translationFilename(translation))
	utility.Fetch(i18nURL, translationFilename(translation))
	return nil
}

func (b *Buddy) Debian2Bleve(translation string) error {
	// UGLY + BUGGY( check eg. `gzcat Packages.gz | grep '^Section: math$' | wc -l` vs 345 here).
	// And todo:
	// add LongDescriptions
	// add https://popcon.debian.org/source/by_inst for scoring? how to use in query?
	packGZ, err := os.Open(PackagesGZ)
	utility.Fatal(err)
	zr, err := gzip.NewReader(packGZ)
	utility.Fatal(err)
	fileScanner := bufio.NewScanner(zr)
	fileScanner.Split(bufio.ScanLines)

	var packagesRegEx = regexp.MustCompile(`^(\S+): (.*)`)

	var p APTPackage
	var pn string
	var cs string

	batchSize := 10000
	packagesCount := 0
	batch := b.index.NewBatch()
	for fileScanner.Scan() {
		line := fileScanner.Text()
		if line == "" {
			//fmt.Printf("flush current package - batch index %s\n", pn)
			packagesCount++
			batch.Index(pn, p)
			p = APTPackage{}
			pn = ""
			if packagesCount%batchSize == 0 {
				fmt.Printf("Indexing batch (%d docs)...\n", packagesCount)
				err := b.index.Batch(batch)
				if err != nil {
					return err
				}
				batch = b.index.NewBatch()
			}
			continue
		}

		pMatch := packagesRegEx.FindStringSubmatch(line)
		if len(pMatch) == 3 {
			//fmt.Printf("'%s' = '%s'\n", pMatch[1], pMatch[2])
			cs = pMatch[1]
		}
		switch cs {
		case "Package":
			p.Package = pMatch[2]
			pn = pMatch[2]
		case "Description":
			p.Description = pMatch[2]
		case "Section":
			p.Section = pMatch[2]
		case "Homepage":
			p.Homepage = pMatch[2]
		case "Tag":
			p.Tags = append(p.Tags, line) // BS FIXME; explode ,
		}
	}

	fmt.Printf("Indexing batch (%d docs)...\n", packagesCount)
	err = b.index.Batch(batch)
	if err != nil {
		return err
	}

	zr.Close()
	packGZ.Close()
	return nil
}

func (b *Buddy) Experiments() {
	r := b.Search("0ad")
	fmt.Printf("\n%+v\n", r)
	fmt.Printf("%+v\n", r.MaxScore)

	b.FieldDict("Description", 1000)
	b.FieldDict("Section", 1)

	r = b.Search("Section:math")
	fmt.Printf("\n%+v\n", r)

	b.FieldDict("Tags", 1)
	r = b.Search(`Tags:"lang:ada"`)
	fmt.Printf("\n%+v\n", r)

	b.FacetQueryExperiment()
}

func (b *Buddy) FacetQueryExperiment() {
	i := b.index
	facet := bleve.NewFacetRequest("Section", 10)
	query := bleve.NewMatchAllQuery()
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.AddFacet("facet name", facet)
	searchResults, err := i.Search(searchRequest)
	if err != nil {
		panic(err)
	}

	// total number of terms
	fmt.Println(searchResults.Facets["facet name"].Total)
	// numer of docs with no value for this field
	fmt.Println(searchResults.Facets["facet name"].Missing)
	// term with highest occurrences in field name
	fmt.Println(searchResults.Facets["facet name"].Terms.Terms()[0].Term)
	fmt.Printf("%+v", searchResults.Facets["facet name"].Terms)
}

func translationFilename(translation string) string {
	return fmt.Sprintf("Translation-%s.bz2", translation)
}

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
