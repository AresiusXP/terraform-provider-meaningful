package meaningful

import (
	"github.com/hashicorp/terraform/helper/schema"

	"fmt"
	"github.com/imroc/req"
	"encoding/json"
	"bytes"
	"errors"
)

func resourceMeaningfulName() *schema.Resource {
	return &schema.Resource{
		Create: resourceMeaningfulNameCreate,
		Read: 	resourceMeaningfulNameRead,
		Delete: resourceMeaningfulNameDelete,

		Schema: map[string]*schema.Schema{
			"tenant_id": &schema.Schema{
				Type:			schema.TypeString,
				Required:		true,
				Sensitive:		true,
				ForceNew:		true,
				Description:	"Tenant ID for requesting Token",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"client_id": &schema.Schema{
				Type:			schema.TypeString,
				Required:		true,
				Sensitive:		true,
				ForceNew:		true,
				Description:	"SPN client Id.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"client_secret": &schema.Schema{
				Type:			schema.TypeString,
				Required:		true,
				Sensitive:		true,
				ForceNew:		true,
				Description:	"SPN client secret.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"meaningful_env": &schema.Schema{
				Type:			schema.TypeString,
				Required:		true,
				Sensitive:		true,
				ForceNew:		true,
				Description:	"QA or Prod Meaningful app",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"resource_type": &schema.Schema{
				Type:			schema.TypeString,
				Required:		true,
				ForceNew:		true,
				Description:	"Type of resource to deploy in Azure",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"deployment_id": &schema.Schema{
				Type:			schema.TypeString,
				Required:		true,
				Sensitive:		true,
				ForceNew:		true,
				Description:	"Eg RPA001, TASEMB, ADVCAT",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"location": &schema.Schema{
				Type:			schema.TypeString,
				Required:		true,
				Sensitive:		true,
				ForceNew:		true,
				Description:	"Location where deployment will happen. E.g. westeurope",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"environment": &schema.Schema{
				Type:			schema.TypeString,
				Required:		true,
				ForceNew:		true,
				Description:	"E.g. Development, Production, etc.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceMeaningfulNameCreate(data *schema.ResourceData, meta interface{}) error {

	clientId := data.Get("client_id").(string)
	clientSecret := data.Get("client_secret").(string)
	azureTenantId := data.Get("tenant_id").(string)
	deploymentId := data.Get("deployment_id").(string)

	var mnfClientIdResource, mnfUri string

	switch data.Get("meaningful_env") {
	case "QA":
		mnfClientIdResource = "5cb2e0ee-79f0-42a5-8bc5-c290f098b890" // QA
		mnfUri = "https://uksqmngweb001.azurewebsites.net/api/Meaning" // QA	
	case "Prod":
		mnfClientIdResource = "d127f0cc-fb77-4c4b-99e5-09b7c41979bf" // Prod
		mnfUri = "https://ukspmngweb001.azurewebsites.net/api/Meaning" // Prod
	}

	tokenUri := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/token", azureTenantId)

	resourceType := data.Get("resource_type").(string)
	location := data.Get("location").(string)
	environmentName := data.Get("environment").(string)

	tokenResp := getToken(tokenUri, clientId, clientSecret, mnfClientIdResource)
	token := tokenResp["access_token"].(string)

	name, err := getName(mnfUri, token, resourceType, deploymentId, location, environmentName)

	if err != nil {
		return err
	}

	data.Set("name", name)
	data.SetId(name)

	return nil
}

func resourceMeaningfulNameRead(data *schema.ResourceData, meta interface{}) error {

	clientId := data.Get("client_id").(string)
	clientSecret := data.Get("client_secret").(string)
	azureTenantId := data.Get("tenant_id").(string)
	deploymentId := data.Get("deployment_id").(string)
	mnfName := data.Get("name").(string)

	var mnfClientIdResource, mnfUri string

	switch data.Get("meaningful_env") {
	case "QA":
		mnfClientIdResource = "5cb2e0ee-79f0-42a5-8bc5-c290f098b890" // QA
		mnfUri = "https://uksqmngweb001.azurewebsites.net/api/Meaning/Generated" // QA	
	case "Prod":
		mnfClientIdResource = "d127f0cc-fb77-4c4b-99e5-09b7c41979bf" // Prod
		mnfUri = "https://ukspmngweb001.azurewebsites.net/api/Meaning/Generated" // Prod
	}

	tokenUri := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/token", azureTenantId)

	resourceType := data.Get("resource_type").(string)
	location := data.Get("location").(string)
	environmentName := data.Get("environment").(string)

	tokenResp := getToken(tokenUri, clientId, clientSecret, mnfClientIdResource)
	token := tokenResp["access_token"].(string)

	nameExists, err := checkName(mnfUri, token, resourceType, deploymentId, location, environmentName, mnfName)

	if err != nil {
		return err
	} 

	if nameExists == false {
		data.SetId("")
	}

	return nil
}

func resourceMeaningfulNameDelete(data *schema.ResourceData, meta interface{}) error {

	clientId := data.Get("client_id").(string)
	clientSecret := data.Get("client_secret").(string)
	azureTenantId := data.Get("tenant_id").(string)
	deploymentId := data.Get("deployment_id").(string)
	resourceType := data.Get("resource_type").(string)
	location := data.Get("location").(string)
	environmentName := data.Get("environment").(string)

	var mnfClientIdResource, mnfUri string

	switch data.Get("meaningful_env") {
	case "QA":
		mnfClientIdResource = "5cb2e0ee-79f0-42a5-8bc5-c290f098b890" // QA
		mnfUri = "https://uksqmngweb001.azurewebsites.net/api/Meaning/Reset" // QA	
	case "Prod":
		mnfClientIdResource = "d127f0cc-fb77-4c4b-99e5-09b7c41979bf" // Prod
		mnfUri = "https://ukspmngweb001.azurewebsites.net/api/Meaning/Reset" // Prod
	}

	tokenUri := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/token", azureTenantId)

	tokenResp := getToken(tokenUri, clientId, clientSecret, mnfClientIdResource)
	token := tokenResp["access_token"].(string)

	httpCode, err := deleteName(mnfUri, token, resourceType, deploymentId, location, environmentName)

	if err != nil {
		return err
	} 

	if httpCode == 204 {
		data.SetId("")
		return nil
	} else {
		return errors.New(fmt.Sprintf("Unexpected result. HTTP return code: %d", httpCode))
	}
}

func getToken(url string, clientId string, clientSecret string, mnfClientIdResource string) map[string]interface{} {
	tokenHeader := req.Header{
		"Content-type": "application/x-www-form-urlencoded",
	}
	
	tokenBody := req.Param{
		"client_id":		clientId,
		"client_secret":	clientSecret,
		"scope":			"openid profile offline_access",
		"grant_type":		"client_credentials",
		"resource":			mnfClientIdResource,
	}
	
	response, err := req.Post(url, tokenHeader, tokenBody)
	if err != nil {
		panic(err)
	}
	
	// io.ReadCloser to string
	buf := new(bytes.Buffer)
    buf.ReadFrom(response.Response().Body)
	newStr := buf.String()
	
	//string to JSON
	in := []byte(newStr)
	var jsonResp map[string]interface{}
	json.Unmarshal(in, &jsonResp)

	return jsonResp
}

func getName (url string, token string, resourceType string, deploymentId string, location string, environmentName string) (string, error) {

	mnfHeader := req.Header{
		"Authorization": fmt.Sprintf("Bearer %s", token),
		"Content-Type": "application/json",
	}

	mnfBody := map[string]string{
		"productName": deploymentId[0:3],
		"resourceTypeName": resourceType,
		"projectName": deploymentId[3:6],
		"regionName": location,
		"environmentName": environmentName,
	}
	var returnName string

	response, err := req.Post(url, mnfHeader, req.BodyJSON(&mnfBody))
	if err != nil {
		return returnName, err
	}

	// io.ReadCloser to string
	buf := new(bytes.Buffer)
    buf.ReadFrom(response.Response().Body)
	newStr := buf.String()

	//string to JSON
	in := []byte(newStr)
	var arrayJson = []*MeaningfulResponse{}
	json.Unmarshal(in, &arrayJson)

	returnName = arrayJson[0].Name

	return returnName, err
}

func checkName (url string, token string, resourceType string, deploymentId string, location string, environmentName string, mnfName string) (bool, error) {

	mnfHeader := req.Header{
		"Authorization": fmt.Sprintf("Bearer %s", token),
		"Content-Type": "application/json",
	}

	mnfBody := map[string]string{
		"productName": deploymentId[0:3],
		"resourceTypeName": resourceType,
		"projectName": deploymentId[3:6],
		"regionName": location,
		"environmentName": environmentName,
	}

	response, err := req.Post(url, mnfHeader, req.BodyJSON(&mnfBody))

	// io.ReadCloser to string
	buf := new(bytes.Buffer)
    buf.ReadFrom(response.Response().Body)
	newStr := buf.String()

	//string to JSON
	in := []byte(newStr)
	var arrayJson []MeaningfulResponse
	json.Unmarshal(in, &arrayJson)

	nameExists := false
	
	for i := range arrayJson {
		if arrayJson[i].Name == mnfName {
			nameExists = true
		}
	}

	return nameExists, err
}

func deleteName (url string, token string, resourceType string, deploymentId string, location string, environmentName string) (int, error) {

	mnfHeader := req.Header{
		"Authorization": fmt.Sprintf("Bearer %s", token),
		"Content-Type": "application/json",
	}

	mnfBody := map[string]string{
		"productName": deploymentId[0:3],
		"resourceTypeName": resourceType,
		"projectName": deploymentId[3:6],
		"regionName": location,
		"environmentName": environmentName,
	}

	response, err := req.Post(url, mnfHeader, req.BodyJSON(&mnfBody))
	
	httpCode := response.Response().StatusCode

	return httpCode, err
}

type MeaningfulResponse struct {
	Name string
}