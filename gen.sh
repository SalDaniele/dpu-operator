openssl req -newkey rsa:4096 -nodes -sha256 -keyout /opt/registry/certs/domain.key -x509 -days 365 -out /opt/registry/certs/domain.crt -subj "/CN=wsfd-advnetlab239.anl.eng.bos2.dc.redhat.com" -addext "subjectAltName = DNS:wsfd-advnetlab239.anl.eng.bos2.dc.redhat.com"
podman run --name myregistry -p 5000:5000 -v /opt/registry/data:/var/lib/registry:z -v /opt/registry/auth:/auth:z -v /opt/registry/certs:/certs:z -e "REGISTRY_HTTP_TLS_CERTIFICATE=/certs/domain.crt" -e "REGISTRY_HTTP_TLS_KEY=/certs/domain.key" -e REGISTRY_COMPATIBILITY_SCHEMA1_ENABLED=true -d docker.io/library/registry:latest
oc create secret tls registry-cas --cert=/opt/registry/certs/domain.crt --key=/opt/registry/certs/domain.key
oc patch image.config.openshift.io/cluster --patch '{"spec":{"additionalTrustedCA":{"name":"registry-cas"}}}' --type=merge
openssl s_client -connect wsfd-advnetlab239.anl.eng.bos2.dc.redhat.com:5000
podman push c4b7878d65c8 localhost:5000/dpu-operator:001-x86_64
make undeploy; IMG=wsfd-advnetlab239.anl.eng.bos2.dc.redhat.com:5000/dpu-operator:001-x86_64 make deploy; oc create -f ~/a.yaml
