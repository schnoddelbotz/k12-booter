
K12_VERSION := $(shell git describe --tags | cut -dv -f2)
LDFLAGS := -X schnoddelbotz/k12-booter/cmd/k12-booter/cmd.AppVersion=$(K12_VERSION)

all: k12-booter aptbuddy sinebeep

k12-booter: bin/k12-booter
bin/k12-booter: cmd/k12-booter/main.go cmd/k12-booter/cmd/*.go
	cd bin && go build -ldflags='$(LDFLAGS)' ../cmd/k12-booter

aptbuddy: bin/aptbuddy
bin/aptbuddy: cmd/aptbuddy/main.go cmd/aptbuddy/cmd/*.go
	cd bin && go build -ldflags='$(LDFLAGS)' ../cmd/aptbuddy

sinebeep: bin/sinebeep
bin/sinebeep: cmd/sinebeep/main.go
	cd bin && go build ../cmd/sinebeep

internationalization/countrydata.go: iso3166
	./iso3166 > internationalization/countrydata.go 
	go fmt internationalization/countrydata.go 

iso3166: internationalization/iso3166/*.go
	cd internationalization/iso3166 && go build -o ../../iso3166

flagit: internationalization/flagit/*.go
	cd internationalization/flagit && go build -o ../../flagit
	banner -w 40 k12 | ./flagit 
	cowsay k12-booter | lolcat

#

playsine: sinebeep
	./sinebeep -duration 150 cCdDefFgGaAb^:cCdDefFgGaAb^:cCdDefFgGaAbbAaGgFfeDdCcVbAaGgFfeDdCcVXbAaGgFfeDdCc
	./sinebeep -duration 200 cdefgg,aaaag,::aaaag,ffffee,ggggc
	./sinebeep -duration 200 ^^cdefgg:,aaaag,aaaag,ffffee,XXXXggggc,::::c
	./sinebeep -duration 200 ^^^^cVcVcVcVcVcVcVcVc
	# Now looking for DAISY notes ...?
	# https://tedgioia.substack.com/p/how-an-ibm-computer-learned-to-sing
	# Trying here: Love story - sorry for the buggy implementation!
	# https://en.wikipedia.org/wiki/Francis_Lai
	./sinebeep -duration 200 ^^cVee^cc,Vee^ccVefedddbb,ddbbdedcccaa,ccaacdcVbbb^aVb^aa
	# polyphony via shell jobs - 2 voices
	./sinebeep -duration 200 ^cVee^cc,Vee^ccVefedddbb & \
	./sinebeep -duration 200 ^^cVee^cc,Vee^ccVefedddbb
	# 3 voices
	./sinebeep -duration 200 ^cVee^cc,Vee^ccVefedddbb,ddbbdedcccaa,ccaacdcVbbb^aVb^aa & \
	./sinebeep -duration 200 ^^cVee^cc,Vee^ccVefedddbb,ddbbdedcccaa,ccaacdcVbbb^aVb^aa & \
	./sinebeep -duration 200 ^^^cVee^cc,Vee^ccVefedddbb,ddbbdedcccaa,ccaacdcVbbb^aVb^aa

apt-bleve-experiment: aptbuddy install-bleve-cli
	./bin/aptbuddy index
	bleve query aptbuddy_en.bleve 'Section:math^100 Description:transitional^-100 Package:/common/^-100' --fields

install-bleve-cli:
	go install github.com/blevesearch/bleve/v2/cmd/bleve@latest

clean:
	cd bin && rm -f k12-booter iso3166 flagit sinebeep aptbuddy

realclean: clean 
	rm -rf aptbuddy_en.bleve/ Packages.gz Translation-??.bz2
