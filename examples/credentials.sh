for ARGUMENT in "$@"
do

    KEY=$(echo $ARGUMENT | cut -f1 -d=)
    VALUE=$(echo $ARGUMENT | cut -f2 -d=)   

    case "$KEY" in
            aws)        AWS_FILE=${VALUE} ;;
            gcp)        GCP_FILE=${VALUE} ;;
            azure)      AZURE_FILE=${VALUE} ;;
            alibaba)    ALIBABA_FILE=${VALUE} ;;
            tencent)    TENCENT_FILE=${VALUE} ;;
            openstack)  OPENSTACK_FILE=${VALUE} ;;
            ibm)        IBM_FILE=${VALUE} ;;
            cloudit)    CLOUDIT_FILE=${VALUE} ;;
            *)   
    esac    


done
if [ "${AWS_FILE}" != "" ]; then 
	#FILE="${HOME}/.aws/credentials"
	export AWS_SECRET_ID="$(head -n 3 ${AWS_FILE} | tail -n 1 | sed  '/^$/d; s/\r//; s/aws_access_key_id = //g')"
	export AWS_SECRET_KEY="$(head -n 2 ${AWS_FILE} | tail -n 1 | sed  '/^$/d; s/\r//; s/aws_secret_access_key = //g')"
	echo "export AWS_SECRET_ID=\"${AWS_SECRET_ID}\""
	echo "export AWS_SECRET_KEY=\"${AWS_SECRET_KEY}\""
fi

if [ "${GCP_FILE}" != "" ]; then 
	#FILE="${HOME}/.ssh/google-credential-cloudbarista.json"
	export GCP_PROJECT=$(cat ${GCP_FILE} | jq ".project_id" | tr -d '"')
	export GCP_PKEY=$(cat ${GCP_FILE} | jq ".private_key" | tr -d '"')
	export GCP_SA=$(cat ${GCP_FILE} | jq ".client_email" | tr -d '"')
	echo "export GCP_PROJECT=\"${GCP_PROJECT}\""
	echo "export GCP_PKEY=\"${GCP_PKEY}\""
	echo "export GCP_SA=\"${GCP_SA}\""
fi
if [ "${AZURE_FILE}" != "" ]; then 
	#FILE="${HOME}/.azure/azure-credential-cloudbarista.json"
	export AZURE_SUBSCRIPTION_ID="$(cat ${AZURE_FILE} | jq '.subscriptionId' | sed  '/^$/d; s/\r//; s/"//g')"
	export AZURE_TENANT_ID="$(cat ${AZURE_FILE} | jq '.tenantId' | sed  '/^$/d; s/\r//; s/"//g')"
	export AZURE_CLIENT_ID="$(cat ${AZURE_FILE} | jq '.clientId' | sed  '/^$/d; s/\r//; s/"//g')"
	export AZURE_CLIENT_SECRET="$(cat ${AZURE_FILE} | jq '.clientSecret' | sed  '/^$/d; s/\r//; s/"//g')"
	echo "export AZURE_SUBSCRIPTION_ID=\"${AZURE_SUBSCRIPTION_ID}\""
	echo "export AZURE_TENANT_ID=\"${AZURE_TENANT_ID}\""
	echo "export AZURE_CLIENT_ID=\"${AZURE_CLIENT_ID}\""
	echo "export AZURE_CLIENT_SECRET=\"${AZURE_CLIENT_SECRET}\""
fi
if [ "${ALIBABA_FILE}" != "" ]; then 
	#FILE="${HOME}/.ssh/alibaba_accesskey.csv"
	export ALIBABA_SECRET_ID="$(cat ${ALIBABA_FILE} | awk 'FNR==2' | cut -f1 -d ',')"
	export ALIBABA_SECRET_KEY="$(cat ${ALIBABA_FILE} | awk 'FNR==2' | cut -f2 -d ',')"
	echo "export ALIBABA_SECRET_ID=\"${ALIBABA_SECRET_ID}\""
	echo "export ALIBABA_SECRET_KEY=\"${ALIBABA_SECRET_KEY}\""
fi
if [ "${TENCENT_FILE}" != "" ]; then 
	#FILE="${HOME}/.tccli/default.credential"
	export TENCENT_SECRET_ID="$(cat ${TENCENT_FILE} | jq '.secretId' | sed  '/^$/d; s/\r//; s/"//g')"
	export TENCENT_SECRET_KEY="$(cat ${TENCENT_FILE} | jq '.secretKey' | sed  '/^$/d; s/\r//; s/"//g')"
	echo "export TENCENT_SECRET_ID=\"${TENCENT_SECRET_ID}\""
	echo "export TENCENT_SECRET_KEY=\"${TENCENT_SECRET_KEY}\""
fi
if [ "${IBM_FILE}" != "" ]; then 
	#FILE="${HOME}/.ssh/ibm-apikey.json"
	export IBM_API_KEY="$(cat ${IBM_FILE} | jq '.apikey' | sed  '/^$/d; s/\r//; s/"//g')"
	echo "export IBM_API_KEY=\"${IBM_API_KEY}\""
fi
if [ "${OPENSTACK_FILE}" != "" ]; then 
    #FILE="${HOME}/.ssh/openstack-openrc.sh"
	source "${OPENSTACK_FILE}"
    echo "export OS_AUTH_URL=\"${OS_AUTH_URL}\""
    echo "export OS_USERNAME=\"${OS_USERNAME}\""
    echo "export OS_PASSWORD=\"${OS_PASSWORD}\""
    echo "export OS_USER_DOMAIN_NAME=\"${OS_USER_DOMAIN_NAME}\""
    echo "export OS_PROJECT_ID=\"${OS_PROJECT_ID}\""
    echo "export OS_REGION_NAME=\"${OS_REGION_NAME}\""
fi
if [ "${CLOUDIT_FILE}" != "" ]; then 
    #FILE="${HOME}/.ssh/cloudit-credential.sh"
	source "${CLOUDIT_FILE}"
    echo "export CLOUDIT_ENDPOINT=\"${CLOUDIT_ENDPOINT}\""
    echo "export CLOUDIT_TOKEN=\"${CLOUDIT_TOKEN}\""
    echo "export CLOUDIT_USERNAME=\"${CLOUDIT_USERNAME}\""
    echo "export CLOUDIT_PASSWORD=\"${CLOUDIT_PASSWORD}\""
    echo "export CLOUDIT_TENANT_ID=\"${CLOUDIT_TENANT_ID}\""
    echo "export CLOUDIT_REGION=\"${CLOUDIT_REGION}\""
    echo "export CLOUDIT_ZONE=\"${CLOUDIT_ZONE}\""
fi
