
all: k12-booter

k12-booter: *.go */*.go
	go build

# internationalization/countrydata.go: iso3166
#	./iso3166 > internationalization/countrydata.go 
#	go fmt internationalization/countrydata.go 

iso3166: internationalization/iso3166/*.go
	cd internationalization/iso3166 && go build -o ../../iso3166

clean:
	rm -f k12-booter iso3166
