okr:
		gitokr  ./OKR.yml | dot -Tsvg > OKR.svg
		open -a "/Applications/Google Chrome.app" ./OKR.svg
