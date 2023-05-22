
K12_VERSION := $(shell git describe --tags | cut -dv -f2)
LDFLAGS := -X main.AppVersion=$(K12_VERSION)

all: k12-booter

k12-booter: *.go
	go build -ldflags='$(LDFLAGS)'

internationalization/countrydata.go: iso3166
	./iso3166 > internationalization/countrydata.go 
	go fmt internationalization/countrydata.go 

iso3166: internationalization/iso3166/*.go
	cd internationalization/iso3166 && go build -o ../../iso3166

flagit: internationalization/flagit/*.go
	cd internationalization/flagit && go build -o ../../flagit
	banner -w 40 k12 | ./flagit 
	cowsay k12-booter | lolcat

sinebeep: sounds/sinebeep/*.go 
	cd sounds/sinebeep && go build -o ../../sinebeep

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

apt-bleve-experiment: k12-booter install-bleve-cli
	./k12-booter -apt
	bleve query aptbuddy_en.bleve 'Section:math^100 Description:transitional^-100 Package:/common/^-100' --fields

install-bleve-cli:
	go install github.com/blevesearch/bleve/v2/cmd/bleve@latest

clean:
	rm -f k12-booter iso3166 flagit sinebeep
	# rm -rf aptbuddy_en.bleve/
