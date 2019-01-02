okr:
		gitokr  ./OKR.yml | dot -Tsvg > OKR.svg
		open -a "/Applications/Google Chrome.app" ./OKR.svg

release:
		docker build -t abhiyerra/gitokr  .
		docker tag abhiyerra/gitokr abhiyerra/gitokr
		docker push abhiyerra/gitokr

build:
		cd gitokr && dep ensure && go build -o gitokr && mv gitokr /bin/gitokr
		cd gitcanvas && dep ensure && go build -o gitcanvas && mv gitcanvas /bin/gitcanvas
		cd gitcron && dep ensure &&  go build -o gitcron && mv gitcron /bin/gitcron
		# cd gitsop && dep ensure && go install
