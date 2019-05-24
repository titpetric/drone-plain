docker:
	docker build --rm --no-cache -t titpetric/drone-plain .

push:
	docker push titpetric/drone-plain

.PHONY: docker push