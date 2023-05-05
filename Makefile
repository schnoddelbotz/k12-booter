
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
	./sinebeep -duration 200 ^cdefgg:,aaaag,aaaag,ffffee,ggggc
	./sinebeep -duration 200 ^^cdefgg:,aaaag,aaaag,ffffee,XXXXggggc,::::c
	# Now looking for DAISY notes ...?
	# https://tedgioia.substack.com/p/how-an-ibm-computer-learned-to-sing

clean:
	rm -f k12-booter iso3166 flagit sinebeep
