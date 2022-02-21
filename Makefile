run:
	kubectl apply -f k8s/deployment.yml
	kubectl apply -f k8s/ingress.yml
	kubectl apply -f k8s/service.yml
	kubectl apply -f k8s/issuer.yml